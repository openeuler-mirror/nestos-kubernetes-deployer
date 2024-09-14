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
package configmanager

import (
	"github.com/agiledragon/gomonkey/v2"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/globalconfig"
	"os"
	"testing"
)

func TestConfigmanager(t *testing.T) {
	o := &opts.OptionsList{
		RootOptDir:        "./globalconfig",
		ClusterConfigFile: "manager.go",
		InfraPlatform: opts.InfraPlatform{
			Libvirt: opts.Libvirt{
				URI:     "qemu:///system",
				OSPath:  "/etc/qcow2.qcow2",
				CIDR:    "",
				Gateway: "",
			},
		},
	}

	ClusterAsset["cluster"] = &asset.ClusterAsset{
		ClusterID:    "cluster",
		Architecture: "amd64",
		Platform:     "libvirt",
		OSImage:      asset.OSImage{Type: "nestos"},
		UserName:     "root",
		Password:     "123",
		SSHKey:       "./assets.go",
		Master: []asset.NodeAsset{
			{
				Hostname: "k8s-master01",
				IP:       "",
				HardwareInfo: asset.HardwareInfo{
					CPU:  2,
					RAM:  2048,
					Disk: 30,
				},
			},
		},
		Worker: []asset.NodeAsset{
			{
				Hostname: "k8s-worker01",
				IP:       "",
				HardwareInfo: asset.HardwareInfo{
					CPU:  2,
					RAM:  2048,
					Disk: 30,
				},
			},
		},
		Runtime: "crio",
		Kubernetes: asset.Kubernetes{
			KubernetesVersion:    "v1.29.1",
			KubernetesAPIVersion: "v1beta3",
			ApiServerEndpoint:    "",
			ImageRegistry:        "registry.k8s.io",
			PauseImage:           "pause:3.9",
			Network: asset.Network{
				ServiceSubnet: "",
				PodSubnet:     "",
			},
		},
		HookConf: asset.HookConf{
			ShellFiles: []asset.ShellFile{
				{Name: "sss"},
				{Name: "sss"},
			},
		},
	}

	gc, err := globalconfig.InitGlobalConfig(o)
	if err != nil || gc == nil {
		t.Logf("InitGlobalConfig returned an error: %v", err)
	}
	GlobalConfig = &globalconfig.GlobalConfig{}

	clusterconfig, err := GetClusterConfig("cluster")
	if err != nil {
		t.Log(err)
	}

	p := gomonkey.ApplyFunc(globalconfig.InitGlobalConfig, func(*opts.OptionsList) (*globalconfig.GlobalConfig, error) {
		t.Log(9999999)
		return &globalconfig.GlobalConfig{
			PersistDir: "../../data/cluster",
		}, nil
	})
	defer p.Reset()

	t.Run("Initial Success", func(t *testing.T) {
		err := Initial(o)
		if err != nil {
			t.Logf("Initial failed: %v", err)
		}
	})

	t.Run("initializeClusterAsset Success", func(t *testing.T) {
		err := initializeClusterAsset(clusterconfig, o)
		if err != nil {
			t.Logf("initializeClusterAsset failed: %v", err)
		}
	})

	t.Run("initializeClusterAsset Fail", func(t *testing.T) {
		o.InfraPlatform.OSPath = ""
		err := initializeClusterAsset(clusterconfig, o)
		if err == nil {
			t.Logf("Delete fail: %v", err)
		}
	})

	t.Run("GetGlobalConfig Success", func(t *testing.T) {
		_, err := GetGlobalConfig()
		if err != nil {
			t.Logf("GetGlobalConfig failed: %v", err)
		}
	})

	t.Run("GetPersistDir Success", func(t *testing.T) {
		tStr := GetPersistDir()
		if tStr == "" {
			t.Logf("GetPersistDir failed: %v", err)
		}
	})

	t.Run("GetBootstrapIgnPort Success", func(t *testing.T) {
		tStr := GetBootstrapIgnPort()
		if tStr == "" {
			t.Logf("GetBootstrapIgnPort failed: %v", err)
		}
	})

	t.Run("GetBootstrapIgnHost Success", func(t *testing.T) {
		tStr := GetBootstrapIgnHost()
		if tStr == "" {
			t.Logf("GetBootstrapIgnHost failed: %v", err)
		}
	})

	t.Run("GetBootstrapIgnHostPort Success", func(t *testing.T) {
		tStr := GetBootstrapIgnHostPort()
		if tStr == "" {
			t.Logf("GetBootstrapIgnHostPort failed: %v", err)
		}
	})

	t.Run("GetClusterConfig Success", func(t *testing.T) {
		cc, err := GetClusterConfig("cluster")
		if err != nil || cc == nil {
			t.Logf("GetClusterConfig failed: %v", err)
		}
	})

	t.Run("GetClusterConfig Fail", func(t *testing.T) {
		cc, err := GetClusterConfig("")
		if err == nil || cc != nil {
			t.Log("Expected error, got nil")
		}
	})

	t.Run("Persist Success", func(t *testing.T) {
		err := Persist()
		if err != nil {
			t.Logf("Persist fail: %v", err)
		}
	})

	t.Run("Delete Success", func(t *testing.T) {
		err := Delete("cluster")
		if err != nil {
			t.Logf("Delete fail: %v", err)
		}
	})
}

func writeToFile(path string, content string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
