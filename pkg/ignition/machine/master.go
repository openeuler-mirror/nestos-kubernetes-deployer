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
package machine

import (
	"nestos-kubernetes-deployer/pkg/configmanager/asset/cluster"
	"nestos-kubernetes-deployer/pkg/ignition"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/sirupsen/logrus"
)

type masterTmplData struct {
	APIServerURL    string
	Hsip            string //HostName + IP
	ImageRegistry   string
	PauseImage      string
	KubeVersion     string
	ServiceSubnet   string
	PodSubnet       string
	Token           string
	CorednsImageTag string
	IpSegment       string
	ReleaseImageURl string
	CertificateKey  string
}

var (
	enabledServices = []string{
		"kubelet.service",
		"set-kernel-para.service",
		"disable-selinux.service",
		"init-cluster.service",
		"install-cni-plugin.service",
		"join-master.service",
		"release-image-pivot.service",
	}
)

type Master struct {
	ClusterAsset cluster.ClusterAsset
	CertFiles    []CertFile
	IgnFiles     []IgnFile
}

type CertFile struct {
	Path    string
	Mode    int
	Content []byte
}

type IgnFile struct {
	Data []byte
}

func (m *Master) GenerateFiles() error {
	mtd := getTmplData(m.ClusterAsset)
	generateFile := ignition.Common{
		UserName:        m.ClusterAsset.NodeAsset[0].UserName,
		SSHKey:          m.ClusterAsset.NodeAsset[0].SSHKey,
		PassWord:        m.ClusterAsset.NodeAsset[0].PassWord,
		NodeType:        "controlplane",
		TmplData:        mtd,
		EnabledServices: enabledServices,
		Config:          &igntypes.Config{},
	}
	if err := generateFile.Generate(); err != nil {
		logrus.Errorf("failed to generate %s ignition file: %v", m.ClusterAsset.NodeAsset[0].UserName, err)
		return err
	}
	for _, file := range m.CertFiles {
		ignFile := ignition.FileWithContents(file.Path, file.Mode, file.Content)
		generateFile.Config.Storage.Files = ignition.AppendFiles(generateFile.Config.Storage.Files, ignFile)
	}
	data, err := ignition.Marshal(generateFile.Config)
	if err != nil {
		logrus.Errorf("failed to Marshal ignition config: %v", err)
		return err
	}
	appendData(m, data)
	for i := 1; i < m.ClusterAsset.Master.Count; i++ {
		generateFile.UserName = m.ClusterAsset.NodeAsset[i].UserName
		generateFile.SSHKey = m.ClusterAsset.NodeAsset[i].SSHKey
		generateFile.PassWord = m.ClusterAsset.NodeAsset[i].PassWord
		generateFile.NodeType = "master"
		if err := generateFile.Generate(); err != nil {
			logrus.Errorf("failed to generate %s ignition file: %v", m.ClusterAsset.NodeAsset[i].UserName, err)
			return err
		}
		data, err := ignition.Marshal(generateFile.Config)
		if err != nil {
			logrus.Errorf("failed to Marshal ignition config: %v", err)
			return err
		}
		appendData(m, data)
	}

	return nil
}

func getTmplData(c cluster.ClusterAsset) *masterTmplData {
	return &masterTmplData{
		APIServerURL:    c.Kubernetes.ApiServer_Endpoint,
		ImageRegistry:   c.Kubernetes.Insecure_Registry,
		PauseImage:      c.Kubernetes.Pause_Image,
		KubeVersion:     c.Kubernetes.Kubernetes_Version,
		ServiceSubnet:   c.Network.Service_Subnet,
		PodSubnet:       c.Network.Pod_Subnet,
		Token:           c.Kubernetes.Token,
		CorednsImageTag: c.Network.CoreDNS_Image_Version,
		ReleaseImageURl: c.Kubernetes.Release_Image_URL,
	}
}

func appendData(master *Master, data []byte) {
	ignFile := IgnFile{
		Data: data,
	}
	master.IgnFiles = append(master.IgnFiles, ignFile)
}
