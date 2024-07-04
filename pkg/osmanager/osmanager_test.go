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
package osmanager

import (
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/infraasset"
	"nestos-kubernetes-deployer/pkg/configmanager/globalconfig"
	"os"
	"testing"
)

func TestOsmanager(t *testing.T) {
	err := os.Chdir("../../data")
	if err != nil {
		t.Fatal(err)
	}

	configmanager.ClusterAsset["cluster"] = &asset.ClusterAsset{
		ClusterID:    "cluster",
		Architecture: "amd64",
		Platform:     "libvirt",
		InfraPlatform: &infraasset.LibvirtAsset{
			URI:     "qemu:///system",
			OSPath:  "/etc/qcow2.qcow2",
			CIDR:    "192.168.132.0/24",
			Gateway: "192.168.132.1",
		},
		OSImage:  asset.OSImage{Type: "nestos"},
		UserName: "root",
		Password: "123",
		SSHKey:   "./assets.go",
		Master: []asset.NodeAsset{
			{
				Hostname: "k8s-master01",
				IP:       "192.168.132.11",
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
				IP:       "192.168.132.12",
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
			ApiServerEndpoint:    "192.168.132.11:6443",
			ImageRegistry:        "registry.k8s.io",
			PauseImage:           "pause:3.9",
			Network: asset.Network{
				ServiceSubnet: "10.96.0.0/16",
				PodSubnet:     "10.244.0.0/16",
			},
		},
		HookConf: asset.HookConf{
			ShellFiles: []asset.ShellFile{
				{Name: "sss"},
				{Name: "sss"},
			},
		},
	}

	configmanager.GlobalConfig = &globalconfig.GlobalConfig{
		PersistDir: "/",
		BootstrapUrl: globalconfig.BootstrapUrl{
			BootstrapIgnHost: "127.0.0.1",
			BootstrapIgnPort: "9080",
		},
	}

	clusterconfig, err := configmanager.GetClusterConfig("cluster")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("NewOSManager_Success", func(t *testing.T) {
		ns := NewOSManager(clusterconfig)
		if ns == nil {
			t.Error("Expected non-nil NewOSManager instance")
		}
	})

	t.Run("GenerateOSConfig_Success", func(t *testing.T) {
		ns := NewOSManager(clusterconfig)
		if ns == nil {
			t.Error("Expected non-nil NewOSManager instance")
		}
		err := ns.GenerateOSConfig()
		if err != nil {
			t.Fatalf("Failed to create NewOSManager instance: %v", err)
		}
	})

	t.Run("GenerateOSConfig_Fail", func(t *testing.T) {
		ns := NewOSManager(clusterconfig)
		if ns == nil {
			t.Error("Expected non-nil NewOSManager instance")
		}
		configmanager.ClusterAsset["cluster"].ServiceSubnet = ""
		err = ns.GenerateOSConfig()
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
