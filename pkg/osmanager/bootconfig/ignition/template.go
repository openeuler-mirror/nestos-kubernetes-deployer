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
package ignition

import (
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/runtime"
	"nestos-kubernetes-deployer/pkg/constants"
	"nestos-kubernetes-deployer/pkg/osmanager/bootconfig"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/vincent-petithory/dataurl"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
)

type template struct {
	clusterAsset    *asset.ClusterAsset
	config          *igntypes.Config
	enabledServices []string
	enabledFiles    []string
}

func newTemplate(clusterAsset *asset.ClusterAsset, config *igntypes.Config, enabledServices, enabledFiles []string) *template {
	return &template{
		clusterAsset:    clusterAsset,
		config:          config,
		enabledServices: enabledServices,
		enabledFiles:    enabledFiles,
	}
}

func (t *template) GenerateBootConfig() error {
	var (
		files         []bootconfig.File
		systemdConfig bootconfig.Systemd
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

	//set container engine config
	containerRuntime, err := runtime.GetRuntime(t.clusterAsset.Runtime)
	if err != nil {
		return err
	}
	tmplData.CriSocket = containerRuntime.GetRuntimeCriSocket()
	if constants.Isulad == containerRuntime.GetRuntimeClient() {
		t.enabledFiles = append(t.enabledFiles, constants.IsuladConfig)
	}

	if err := bootconfig.AppendStorageFiles(&files, "/", constants.BootConfigFilesPath, tmplData, t.enabledFiles); err != nil {
		logrus.Errorf("failed to add files to an Ignition config: %v", err)
		return err
	}
	if err := bootconfig.AppendSystemdUnits(&systemdConfig, constants.BootConfigSystemdPath, tmplData, t.enabledServices); err != nil {
		logrus.Errorf("failed to add systemd units to an Ignition config: %v", err)
		return err
	}

	t.config = &igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Passwd: igntypes.Passwd{
			Users: []igntypes.PasswdUser{
				{
					Name: t.clusterAsset.UserName,
					SSHAuthorizedKeys: []igntypes.SSHAuthorizedKey{
						igntypes.SSHAuthorizedKey(sshkeyContent),
					},
					PasswordHash: &t.clusterAsset.Password,
				},
			},
		},
		Storage: igntypes.Storage{
			Links: []igntypes.Link{
				{
					Node: igntypes.Node{Path: "/etc/local/time"},
					LinkEmbedded1: igntypes.LinkEmbedded1{
						Target: "/usr/share/zoneinfo/Asia/Shanghai",
					},
				},
			},
		},
	}

	for _, f := range files {
		str := dataurl.EncodeBytes(f.Contents.Source)
		mode := int(f.Mode.Perm())
		ignFile := igntypes.File{
			Node: igntypes.Node{
				Path:      f.Path,
				Overwrite: f.Overwrite,
			},
			FileEmbedded1: igntypes.FileEmbedded1{
				Mode: &mode,
				Contents: igntypes.Resource{
					Source: &str,
				},
			},
		}
		t.config.Storage.Files = append(t.config.Storage.Files, ignFile)
	}

	if len(t.clusterAsset.HookConf.ShellFiles) > 0 {
		for _, sf := range t.clusterAsset.HookConf.ShellFiles {
			ignHookFile := igntypes.File{
				Node: igntypes.Node{
					Path:      constants.HookFilesPath,
					Overwrite: ignutil.BoolToPtr(true),
				},
				FileEmbedded1: igntypes.FileEmbedded1{
					Mode: &sf.Mode,
					Contents: igntypes.Resource{
						Source: ignutil.StrToPtr(dataurl.EncodeBytes(sf.Content)),
					},
				},
			}
			t.config.Storage.Files = append(t.config.Storage.Files, ignHookFile)
		}
	}

	for _, u := range systemdConfig.Units {
		unit := igntypes.Unit{
			Name:     u.Name,
			Contents: u.Contents,
			Enabled:  u.Enabled,
		}
		t.config.Systemd.Units = append(t.config.Systemd.Units, unit)
	}

	return nil
}
