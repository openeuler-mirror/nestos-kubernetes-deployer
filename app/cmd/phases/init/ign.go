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
package phases

import (
	"bytes"
	"fmt"
	"io"
	"nestos-kubernetes-deployer/app/apis/nkd"
	"nestos-kubernetes-deployer/app/cmd/phases/workflow"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/clarketm/json"
	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/sirupsen/logrus"
	"github.com/vincent-petithory/dataurl"
)

func NewGenerateIgnCmd() workflow.Phase {
	return workflow.Phase{
		Name:  "ign",
		Short: "run ign to genertate ignition file",
		Run:   runGenerateIgnConfig,
	}
}

type commonTemplateData struct {
	SSHKey          string
	APIServerURL    string
	Hsip            string //HostName + IP
	ImageRegistry   string
	PodSandboxImage string
	KubeVersion     string
	ServiceSubnet   string
	PodSubnet       string
	Token           string
	NodeType        string
	NodeName        string
}

var (
	enabledServices = []string{
		"kubelet.service",
		"set-kernel-para.service",
		"disable-selinux.service",
		"init-cluster.service",
		"install-cni-plugin.service",
		"join-cluster.service",
	}
)

func runGenerateIgnConfig(r workflow.RunData, node string) error {
	data, ok := r.(InitData)
	if !ok {
		panic(fmt.Sprintf("Expect to fetch configuration data but got a %T", r))
	}
	var (
		hsip        string
		hostName    string
		oneNodeName string
		nodeCount   int
	)
	if node == "master" {
		nodeCount = data.MasterCfg().System.Count
		hostName = data.MasterCfg().System.HostName
		for i := 0; i < nodeCount; i++ {
			oneNodeName = fmt.Sprintf("%s%02d", hostName, i+1)
			temp := data.MasterCfg().System.Ips[i] + oneNodeName + "\n"
			hsip = hsip + temp
		}
		for j := 0; j < nodeCount; j++ {
			ctd := getMasterTmplData(data.MasterCfg(), j+1)
			if err := generateConfig(ctd); err != nil {
				return err
			}
		}
	} else {
		nodeCount = data.WorkerCfg().System.Count
		for j := 0; j < nodeCount; j++ {
			ctd := getWorkerTmplData(data.WorkerCfg(), j+1)
			if err := generateConfig(ctd); err != nil {
				return err
			}
		}
	}
	return nil
}

func getMasterTmplData(nkdConfig *nkd.Master, count int) *commonTemplateData {
	oneNodeName := fmt.Sprintf("%s%d", nkdConfig.System.HostName, count)
	return &commonTemplateData{
		SSHKey:          nkdConfig.System.SSHKey,
		APIServerURL:    "",
		ImageRegistry:   nkdConfig.Repo.Registry,
		PodSandboxImage: "",
		ServiceSubnet:   nkdConfig.Kubeadm.Networking.ServiceSubnet,
		PodSubnet:       nkdConfig.Kubeadm.Networking.PodSubnet,
		Token:           "",
		NodeName:        oneNodeName,
		NodeType:        "master",
	}
}

func getWorkerTmplData(nkdConfig *nkd.Worker, count int) *commonTemplateData {
	oneNodeName := fmt.Sprintf("%s%d", nkdConfig.System.HostName, count)
	return &commonTemplateData{
		SSHKey:          nkdConfig.System.SSHKey,
		APIServerURL:    "",
		ImageRegistry:   nkdConfig.Repo.Registry,
		PodSandboxImage: "",
		Token:           nkdConfig.Worker.Discovery.TlsBootstrapToken,
		NodeName:        oneNodeName,
		NodeType:        "worker",
	}
}

func generateConfig(ctd *commonTemplateData) error {
	config := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Passwd: igntypes.Passwd{
			Users: []igntypes.PasswdUser{
				{
					Name: "root",
					SSHAuthorizedKeys: []igntypes.SSHAuthorizedKey{
						igntypes.SSHAuthorizedKey(ctd.SSHKey),
					},
					Groups: []igntypes.Group{
						igntypes.Group("adm"),
						igntypes.Group("sudo"),
						igntypes.Group("systemd-journal"),
						igntypes.Group("wheel"),
					},
				},
			},
		},
		Storage: igntypes.Storage{
			Links: []igntypes.Link{
				{
					Node: igntypes.Node{
						Path: "/etc/local/time",
					},
					LinkEmbedded1: igntypes.LinkEmbedded1{
						Target: "/usr/share/zoneinfo/Asia/Shanghai",
					},
				},
			},
		},
	}
	nodeFilesPath := fmt.Sprintf("data/ignition/%s/files", ctd.NodeType)
	if err := AddStorageFiles(&config, "/", nodeFilesPath, ctd); err != nil {
		logrus.Errorf("failed to add files to a ignition config: %v", err)
		return err
	}
	nodeUnitPath := fmt.Sprintf("data/ignition/%s/systemd/", ctd.NodeType)
	if err := AddSystemdUnits(&config, nodeUnitPath, ctd, enabledServices); err != nil {
		logrus.Errorf("failed to add systemd units to a ignition config: %v", err)
		return err
	}
	ignName := fmt.Sprintf("%s%s", ctd.NodeName, ".ign")
	if err := generateFile(&config, "./", ignName); err != nil {
		logrus.Errorf("failed to generate ignition file: %v", err)
		return err
	}
	return nil
}

/*
Add files to a ignition config
Parameters:
config - the ignition config to be modified
tmplData - struct to used to render templates
*/
func AddStorageFiles(config *igntypes.Config, base string, uri string, tmplData interface{}) error {
	file, err := os.Open(uri)
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
			err = AddStorageFiles(config, path.Join(base, name), path.Join(uri, name), tmplData)
			if err != nil {
				return err
			}
		}
		return nil
	}
	_, data, err := readFile(info.Name(), file, tmplData)
	if err != nil {
		return err
	}
	ign := fileFromBytes(strings.TrimSuffix(base, ".template"), 0644, data)
	config.Storage.Files = appendFiles(config.Storage.Files, ign)
	return nil
}

/*
Add systemd units to a ignition config
Parameters:
config - the ignition config to be modified
uri - path under data/ignition specifying the systemd units files to be included
tmplData - struct to used to render templates
enabledServices - a list of systemd units to be enabled by default
*/
func AddSystemdUnits(config *igntypes.Config, uri string, tmplData interface{}, enabledServices []string) error {
	enabled := make(map[string]struct{}, len(enabledServices))
	for _, s := range enabledServices {
		enabled[s] = struct{}{}
	}

	dir, err := os.Open(uri)
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
		file, err := os.Open(dir)
		if err != nil {
			return err
		}
		defer file.Close()
		name, contents, err := readFile(childInfo.Name(), file, tmplData)
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

// Read data from the file
func readFile(name string, file io.Reader, tmplData interface{}) (realName string, data []byte, err error) {
	data, err = io.ReadAll(file)
	if err != nil {
		return "", nil, err
	}
	if filepath.Ext(name) == ".template" {
		name = strings.TrimSuffix(name, ".template")
		tmpl := template.New(name)
		tmpl, err := tmpl.Parse(string(data))
		if err != nil {
			return "", nil, err
		}
		stringData := applyTmplData(tmpl, tmplData)
		data = []byte(stringData)
	}

	return name, data, nil
}

func applyTmplData(tmpl *template.Template, data interface{}) string {
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, data); err != nil {
		panic(err)
	}
	return buf.String()
}

func fileFromBytes(path string, mode int, contents []byte) igntypes.File {
	return igntypes.File{
		Node: igntypes.Node{
			Path:      path,
			Overwrite: ignutil.BoolToPtr(true),
		},
		FileEmbedded1: igntypes.FileEmbedded1{
			Mode: &mode,
			Contents: igntypes.Resource{
				Source: ignutil.StrToPtr(dataurl.EncodeBytes(contents)),
			},
		},
	}
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
Generate a ignition config
Parameters:
config - the ignition config to be saved
filePath - the path to save the file
fileName - the name to save the file
*/
func generateFile(config *igntypes.Config, filePath string, fileName string) error {
	data, err := json.Marshal(config)
	if err != nil {
		logrus.Errorf("failed to Marshal ignition config: %v", err)
		return err
	}
	path := filepath.Join(filePath, fileName)
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		logrus.Errorf("failed to Mkdir: %v", err)
		return err
	}
	if err := os.WriteFile(path, data, 0640); err != nil {
		logrus.Errorf("failed to save ignition file: %v", err)
		return err
	}
	return nil
}
