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
package bootconfig

import (
	"errors"
	"fmt"
	"nestos-kubernetes-deployer/data"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/runtime"
	"nestos-kubernetes-deployer/pkg/constants"
	"nestos-kubernetes-deployer/pkg/utils"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/clarketm/json"
	ignutil "github.com/coreos/ignition/v2/config/util"
	"gopkg.in/yaml.v2"
)

type TmplData struct {
	NodeName          string
	APIServerURL      string
	ImageRegistry     string
	Runtime           string
	CriSocket         string
	PauseImage        string
	KubeVersion       string
	ServiceSubnet     string
	PodSubnet         string
	Token             string
	CaCertHash        string
	ReleaseImageURl   string
	CertificateKey    string
	Hsip              string //HostName + IP
	KubeadmApiVersion string
	HookFilesPath     string
	CertsUrl          string
	IsControlPlane    bool
	IsDocker          bool
	IsIsulad          bool
	IsCrio            bool
}

func GetTmplData(c *asset.ClusterAsset) (*TmplData, error) {
	var hsipStrings []string
	for _, master := range c.Master {
		hsipStrings = append(hsipStrings, master.IP+" "+master.Hostname)
	}
	hsip := strings.Join(hsipStrings, "\n")
	engine, err := runtime.GetRuntime(c.Runtime)
	if err != nil {
		return nil, err
	}

	return &TmplData{
		APIServerURL:      c.Kubernetes.ApiServerEndpoint,
		ImageRegistry:     c.Kubernetes.ImageRegistry,
		Runtime:           c.Runtime,
		PauseImage:        c.Kubernetes.PauseImage,
		KubeVersion:       c.Kubernetes.KubernetesVersion,
		KubeadmApiVersion: c.Kubernetes.KubernetesAPIVersion,
		ServiceSubnet:     c.Network.ServiceSubnet,
		PodSubnet:         c.Network.PodSubnet,
		Token:             c.Kubernetes.Token,
		CaCertHash:        c.Kubernetes.CaCertHash,
		ReleaseImageURl:   c.Kubernetes.ReleaseImageURL,
		CertificateKey:    c.Kubernetes.CertificateKey,
		Hsip:              hsip,
		HookFilesPath:     constants.HookFilesPath,
		IsDocker:          runtime.IsDocker(engine),
		IsIsulad:          runtime.IsIsulad(engine),
		IsCrio:            runtime.IsCrio(engine),
	}, nil
}

/*
AppendStorageFiles: 向提供的切片中追加存储文件的信息。
参数：
  - config：指向 File 结构切片的指针，文件信息将被追加到其中。
  - base：要从中开始遍历目录的基本路径。
  - uri：要处理的文件或目录的 URI。
  - tmplData：用于模板渲染的数据。
  - enabledFiles：已启用处理的文件名列表。
*/
func AppendStorageFiles(config *[]File, base string, uri string, tmplData interface{}, enabledFiles []string) error {
	enabled := make(map[string]struct{}, len(enabledFiles))
	for _, s := range enabledFiles {
		enabled[s] = struct{}{}
	}

	file, err := data.Assets.Open(uri)
	if err != nil {
		fmt.Printf("err: %v\n", err)
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
			err = AppendStorageFiles(config, path.Join(base, name), path.Join(uri, name), tmplData, enabledFiles)
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
	if _, ok := enabled[base]; ok {
		ignFile := fileWithContents(strings.TrimSuffix(base, ".template"), constants.StorageFilesMode, data)
		*config = appendFiles(*config, ignFile)
	}
	return nil
}

/*
AppendSystemdUnits: 向 Systemd 结构中追加信息
参数：
  - config：指向 Systemd 结构的指针，其中将添加Systemd信息。
  - uri：要打开的目录的 URI。
  - tmplData：用于模板渲染的数据。
  - enabledServices：已启用处理的服务名列表
*/

func AppendSystemdUnits(config *Systemd, uri string, tmplData interface{}, enabledServices []string) error {
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

		if _, ok := enabled[name]; ok {
			unit := Unit{
				Name:     name,
				Contents: string(contents),
			}
			unit.Enabled = ignutil.BoolToPtr(true)
			config.Units = append(config.Units, unit)
		}
	}
	return nil
}

func GetSavePath(clusterID string) string {
	return filepath.Join(configmanager.GetPersistDir(), clusterID, constants.BootConfigSaveDir)
}

func SaveYAML(data interface{}, filePath string, fileName string, header string) error {
	return saveFile(data, filePath, fileName, header, yaml.Marshal)
}

func SaveJSON(data interface{}, filePath string, fileName string) error {
	return saveFile(data, filePath, fileName, "", json.Marshal)
}

func Marshal(input interface{}) ([]byte, error) {
	return json.Marshal(input)
}

func SaveFile(data []byte, filePath, fileName string) error {
	if data == nil {
		return errors.New("data is nil")
	}

	return saveDataToFile(string(data), filePath, fileName)
}

func saveFile(data interface{}, filePath, fileName, header string, marshalFunc func(interface{}) ([]byte, error)) error {
	if data == nil {
		return errors.New("data is nil")
	}

	dataBytes, err := marshalFunc(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	dataString := header + string(dataBytes)

	return saveDataToFile(dataString, filePath, fileName)
}

func saveDataToFile(data, filePath, fileName string) error {
	fullPath := filepath.Join(filePath, fileName)
	if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	if err := os.WriteFile(fullPath, []byte(data), os.ModePerm); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func fileWithContents(path string, mode os.FileMode, contents []byte) File {
	return File{
		Node: Node{
			Path:      path,
			Overwrite: ignutil.BoolToPtr(true),
		},
		FileEmbedded1: FileEmbedded1{
			Mode: mode,
			Contents: Resource{
				Source: contents,
			},
		},
	}
}

func appendFiles(files []File, file File) []File {
	for i, f := range files {
		if f.Node.Path == file.Node.Path {
			files[i] = file
			return files
		}
	}
	files = append(files, file)
	return files
}
