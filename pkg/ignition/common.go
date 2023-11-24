/*
Copyright 2023 KylinSoft  Co., Ltd.

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
	"nestos-kubernetes-deployer/data"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/cluster"
	"nestos-kubernetes-deployer/pkg/utils"
	"os"
	"path"
	"strings"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/sirupsen/logrus"
)

var (
	enabledServices = []string{
		"kubelet.service",
		"set-kernel-para.service",
		"disable-selinux.service",
		"init-cluster.service",
		"install-cni-plugin.service",
		"join-master.service",
		"join-worker.service",
		"release-image-pivot.service",
	}
)

type tmplData struct {
	SSHKey          string
	APIServerURL    string
	Hsip            string //HostName + IP
	ImageRegistry   string
	PauseImageTag   string
	KubeVersion     string
	ServiceSubnet   string
	PodSubnet       string
	Token           string
	NodeType        string
	NodeName        string
	CorednsImageTag string
	IpSegment       string
	ReleaseImageURl string
	PasswordHash    string
	CertificateKey  string
}

type Common struct {
	Config       *igntypes.Config
	ClusterAsset cluster.ClusterAsset
	Files        []File
}

type File struct {
	Path    string
	Mode    int
	Content []byte
}

func (c *Common) GenerateFile() error {
	c.Config = &igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Passwd: igntypes.Passwd{
			Users: []igntypes.PasswdUser{
				{
					Name: "root",
					SSHAuthorizedKeys: []igntypes.SSHAuthorizedKey{
						igntypes.SSHAuthorizedKey("/*SSHKEY*/"),
					},
					PasswordHash: nil, /*PasswordHasH*/
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
	//get template data
	td := GetTmplData(c.ClusterAsset)

	//todo：对配置项参数解析，生成不同的Ignition文件

	if err := AppendStorageFiles(c.Config, "/", "", td); err != nil {
		logrus.Errorf("failed to add files to a ignition config: %v", err)
		return err
	}
	if err := AppendSystemdUnits(c.Config, "", td, enabledServices); err != nil {
		logrus.Errorf("failed to add systemd units to a ignition config: %v", err)
		return err
	}

	for _, file := range c.Files {
		ignFile := FileWithContents(file.Path, file.Mode, file.Content)
		c.Config.Storage.Files = appendFiles(c.Config.Storage.Files, ignFile)
	}
	return nil
}

func (c *Common) SaveFile(filename string) error {
	data, err := Marshal(c.Config)
	if err != nil {
		logrus.Errorf("failed to Marshal ignition config: %v", err)
		return err
	}
	if err := os.WriteFile(filename, data, 0640); err != nil {
		logrus.Errorf("failed to save ignition file: %v", err)
		return err
	}
	return nil
}

func GetTmplData(c cluster.ClusterAsset) *tmplData {
	return &tmplData{
		KubeVersion: c.KubernetesVersion,
	}
}

/*
AppendStorageFiles add files to a ignition config
Parameters:
  - config: the ignition config to be modified
  - tmplData: struct to used to render templates
*/
func AppendStorageFiles(config *igntypes.Config, base string, uri string, tmplData interface{}) error {
	file, err := data.Assets.Open(uri)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		children, err := file.Readdir(0)
		if err != nil {
			return err
		}
		if err = file.Close(); err != nil {
			return err
		}

		for _, childInfo := range children {
			name := childInfo.Name()
			err = AppendStorageFiles(config, path.Join(base, name), path.Join(uri, name), tmplData)
			if err != nil {
				return err
			}
		}
		return nil
	}
	_, data, err := utils.GetCompleteFile(info.Name(), file, tmplData)
	if err != nil {
		return err
	}
	ignFile := FileWithContents(strings.TrimSuffix(base, ".template"), 0755, data)
	config.Storage.Files = appendFiles(config.Storage.Files, ignFile)
	return nil
}

func appendFiles(files []igntypes.File, file igntypes.File) []igntypes.File {
	for i, f := range files {
		if f.Node.Path == file.Node.Path {
			files[i] = file
			return files
		}
	}
	files = append(files, file)
	return files
}

/*
Add systemd units to a ignition config
Parameters:
  - config: the ignition config to be modified
  - uri: path under data/ignition specifying the systemd units files to be included
  - tmplData: struct to used to render templates
  - enabledServices: a list of systemd units to be enabled by default
*/
func AppendSystemdUnits(config *igntypes.Config, uri string, tmplData interface{}, enabledServices []string) error {
	enabled := make(map[string]struct{}, len(enabledServices))
	for _, s := range enabledServices {
		enabled[s] = struct{}{}
	}

	dir, err := data.Assets.Open(uri)
	if err != nil {
		return err
	}
	defer dir.Close()

	child, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, childInfo := range child {
		dir := path.Join(uri, childInfo.Name())
		file, err := data.Assets.Open(dir)
		if err != nil {
			return err
		}
		defer file.Close()
		name, contents, err := utils.GetCompleteFile(childInfo.Name(), file, tmplData)
		if err != nil {
			return err
		}
		unit := igntypes.Unit{
			Name:     name,
			Contents: ignutil.StrToPtr(string(contents)),
		}
		if _, ok := enabled[name]; ok {
			unit.Enabled = ignutil.BoolToPtr(true)
		}
		config.Systemd.Units = append(config.Systemd.Units, unit)
	}
	return nil
}
