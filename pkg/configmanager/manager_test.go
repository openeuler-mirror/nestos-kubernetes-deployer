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
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/globalconfig"
	"testing"
)

func TestConfigmanager(t *testing.T) {
	opts := &opts.OptionsList{
		RootOptDir: "/tmp",
		InfraPlatform: opts.InfraPlatform{
			Libvirt: opts.Libvirt{
				URI:     "qemu:///system",
				OSPath:  "/etc/qcow2.qcow2",
				CIDR:    "192.168.132.0/24",
				Gateway: "192.168.132.1",
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

	gc, err := globalconfig.InitGlobalConfig(opts)
	if err != nil || gc == nil {
		t.Fatalf("InitGlobalConfig returned an error: %v", err)
	}

	clusterconfig, err := GetClusterConfig("cluster")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Initial Success", func(t *testing.T) {
		err := Initial(opts)
		if err != nil {
			t.Fatalf("Initial failed: %v", err)
		}
	})

	t.Run("initializeClusterAsset Success", func(t *testing.T) {
		err := initializeClusterAsset(clusterconfig, opts)
		if err != nil {
			t.Fatalf("initializeClusterAsset failed: %v", err)
		}
	})

	t.Run("initializeClusterAsset Fail", func(t *testing.T) {
		opts.InfraPlatform.OSPath = ""
		err := initializeClusterAsset(clusterconfig, opts)
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
