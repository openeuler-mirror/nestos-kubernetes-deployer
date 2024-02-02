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
	"fmt"
	"nestos-kubernetes-deployer/data"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/utils"
	"path"
	"strings"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/sirupsen/logrus"
)

var (
	EnabledServices = []string{
		"kubelet.service",
		"set-kernel-para.service",
		"disable-selinux.service",
		"init-cluster.service",
		"join-master.service",
		"release-image-pivot.service",
		"join-worker.service",
		"set-hostname.service",
	}
)

type TmplData struct {
	NodeName        string
	APIServerURL    string
	ImageRegistry   string
	PauseImage      string
	KubeVersion     string
	ServiceSubnet   string
	PodSubnet       string
	Token           string
	CorednsImageTag string
	ReleaseImageURl string
	CertificateKey  string
	Hsip            string //HostName + IP
}

type Common struct {
	UserName        string
	SSHKey          string
	PassWord        string
	NodeType        string
	TmplData        interface{}
	EnabledServices []string
	Config          *igntypes.Config
}

func (c *Common) Generate() error {
	c.Config = &igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Passwd: igntypes.Passwd{
			Users: []igntypes.PasswdUser{
				{
					Name: c.UserName,
					SSHAuthorizedKeys: []igntypes.SSHAuthorizedKey{
						igntypes.SSHAuthorizedKey(c.SSHKey),
					},
					PasswordHash: &c.PassWord,
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

	nodeFilesPath := fmt.Sprintf("ignition/%s/files", c.NodeType)
	if err := appendStorageFiles(c.Config, "/", nodeFilesPath, c.TmplData); err != nil {
		logrus.Errorf("failed to add files to a ignition config: %v", err)
		return err
	}
	nodeUnitPath := fmt.Sprintf("ignition/%s/systemd/", c.NodeType)
	if err := appendSystemdUnits(c.Config, nodeUnitPath, c.TmplData, c.EnabledServices); err != nil {
		logrus.Errorf("failed to add systemd units to a ignition config: %v", err)
		return err
	}

	return nil
}

/*
AppendStorageFiles add files to a ignition config
Parameters:
  - config: the ignition config to be modified
  - tmplData: struct to used to render templates
*/
func appendStorageFiles(config *igntypes.Config, base string, uri string, tmplData interface{}) error {
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
			err = appendStorageFiles(config, path.Join(base, name), path.Join(uri, name), tmplData)
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
	config.Storage.Files = AppendFiles(config.Storage.Files, ignFile)
	return nil
}

/*
Add systemd units to a ignition config
Parameters:
  - config: the ignition config to be modified
  - uri: path under data/ignition specifying the systemd units files to be included
  - tmplData: struct to used to render templates
  - enabledServices: a list of systemd units to be enabled by default
*/
func appendSystemdUnits(config *igntypes.Config, uri string, tmplData interface{}, enabledServices []string) error {
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

func GetTmplData(c *asset.ClusterAsset) *TmplData {
	var hsip string
	for i := 0; i < len(c.Master); i++ {
		temp := c.Master[i].IP + " " + c.Master[i].Hostname + "\n"
		hsip = hsip + temp
	}

	return &TmplData{
		APIServerURL:    c.Kubernetes.ApiServer_Endpoint,
		ImageRegistry:   c.Kubernetes.Image_Registry,
		PauseImage:      c.Kubernetes.Pause_Image,
		KubeVersion:     c.Kubernetes.Kubernetes_Version,
		ServiceSubnet:   c.Network.Service_Subnet,
		PodSubnet:       c.Network.Pod_Subnet,
		Token:           c.Kubernetes.Token,
		CorednsImageTag: c.Network.CoreDNS_Image_Version,
		ReleaseImageURl: c.Kubernetes.Release_Image_URL,
		CertificateKey:  c.Kubernetes.CertificateKey,
		Hsip:            hsip,
	}
}
