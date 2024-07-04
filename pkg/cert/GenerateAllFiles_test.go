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
package cert

import (
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/globalconfig"
	"testing"
)

func TestOsmanager(t *testing.T) {
	configmanager.ClusterAsset["cluster"] = &asset.ClusterAsset{
		ClusterID:    "cluster",
		Architecture: "amd64",
		Platform:     "libvirt",
		OSImage:      asset.OSImage{Type: "nestos"},
		UserName:     "root",
		Password:     "123",
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
		Runtime: "crio",
		Kubernetes: asset.Kubernetes{
			Network: asset.Network{
				ServiceSubnet: "10.96.0.0/16",
				PodSubnet:     "10.244.0.0/16",
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

	t.Run("NewCertGenerator_Success", func(t *testing.T) {
		ns := NewCertGenerator(clusterconfig.ClusterID, &clusterconfig.Master[0])
		if ns == nil {
			t.Error("Expected non-nil NewCertGenerator instance")
		}
	})

	ns := NewCertGenerator(clusterconfig.ClusterID, &clusterconfig.Master[0])
	if ns == nil {
		t.Error("Expected non-nil NewOSManager instance")
	}

	t.Run("GenerateAllFiles_Success", func(t *testing.T) {
		if err := ns.GenerateAllFiles(); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("GenerateAllFiles_Fail", func(t *testing.T) {
		configmanager.ClusterAsset["cluster"].ServiceSubnet = ""
		err = ns.GenerateAllFiles()
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
