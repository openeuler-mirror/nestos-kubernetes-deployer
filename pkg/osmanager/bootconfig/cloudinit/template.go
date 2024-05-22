/*
Copyright 2024 KylinSoft  Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cloudinit

import (
	"encoding/base64"
	"fmt"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/runtime"
	"nestos-kubernetes-deployer/pkg/constants"
	"nestos-kubernetes-deployer/pkg/osmanager/bootconfig"
	"nestos-kubernetes-deployer/pkg/utils"
	"os"

	"github.com/sirupsen/logrus"
)

type template struct {
	clusterAsset    *asset.ClusterAsset
	config          *CloudinitConfig
	enabledServices []string
	enabledFiles    []string
}

func newTemplate(clusterAsset *asset.ClusterAsset, config *CloudinitConfig, enabledServices, enabledFiles []string) *template {
	return &template{
		clusterAsset:    clusterAsset,
		config:          config,
		enabledServices: enabledServices,
		enabledFiles:    enabledFiles,
	}
}

func (t *template) GenerateBootConfig(url string, nodeType string) error {
	var (
		bootConfigFiles   []bootconfig.File
		bootConfigSystemd bootconfig.Systemd
	)
	sshkeyContent, err := os.ReadFile(t.clusterAsset.SSHKey)
	if err != nil {
		logrus.Debug("Failed to read sshkey content:", err)
		return err
	}
	tmplData, err := bootconfig.GetTmplData(t.clusterAsset)
	if err != nil {
		return err
	}
	if nodeType == constants.Controlplane {
		tmplData.IsControlPlane = true
	}
	tmplData.CertsUrl = utils.ConstructURL(url, constants.CertsFiles)

	//set container engine config
	engine, err := runtime.GetRuntime(t.clusterAsset.Runtime)
	if err != nil {
		return err
	}
	tmplData.CriSocket = engine.GetRuntimeCriSocket()
	if runtime.IsIsulad(engine) {
		t.enabledFiles = append(t.enabledFiles, constants.IsuladConfig)
	} else if runtime.IsDocker(engine) {
		t.enabledFiles = append(t.enabledFiles, constants.DockerConfig)
	}

	if err := bootconfig.AppendStorageFiles(&bootConfigFiles, "/", constants.BootConfigFilesPath, tmplData, t.enabledFiles); err != nil {
		logrus.Errorf("failed to add files to an cloudinit config: %v", err)
		return err
	}

	if err := bootconfig.AppendSystemdUnits(&bootConfigSystemd, constants.BootConfigSystemdPath, tmplData, t.enabledServices); err != nil {
		logrus.Errorf("failed to add systemd units to an cloudinit config: %v", err)
		return err
	}

	t.config = &CloudinitConfig{
		SSHPasswordAuth: true,
		SSHAuthorizedKeys: []string{
			string(sshkeyContent),
		},

		Chpasswd: ChpasswdConfig{
			List:   t.clusterAsset.UserName + ":" + t.clusterAsset.Password,
			Expire: false,
		},
	}

	for _, f := range bootConfigFiles {
		cf := WriteFile{
			EnCoding:    "b64",
			Content:     base64.StdEncoding.EncodeToString(f.Contents.Source),
			Path:        f.Path,
			Permissions: f.Mode,
		}
		t.config.WriteFiles = append(t.config.WriteFiles, cf)
	}
	for _, u := range bootConfigSystemd.Units {
		cf := WriteFile{
			Content:     u.Contents,
			Path:        fmt.Sprintf("/etc/systemd/system/%s", u.Name),
			Permissions: constants.SystemdServiceMode,
		}
		t.config.WriteFiles = append(t.config.WriteFiles, cf)
		t.config.RunCmds = append(t.config.RunCmds, "systemctl enable "+u.Name)
		t.config.RunCmds = append(t.config.RunCmds, "systemctl start "+u.Name)
	}
	return nil
}
