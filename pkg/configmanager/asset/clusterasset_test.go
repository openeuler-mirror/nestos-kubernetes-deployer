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
package asset

import (
	"nestos-kubernetes-deployer/cmd/command/opts"
	"testing"
)

func TestClusterasset(t *testing.T) {
	opts := &opts.OptionsList{
		RootOptDir: "/tmp",
		InfraPlatform: opts.InfraPlatform{
			Libvirt: opts.Libvirt{
				URI:     "qemu:///system",
				OSPath:  "/etc/qcow2.qcow2",
				CIDR:    "127.0.0.1/24",
				Gateway: "127.0.0.1",
			},
		},
	}

	cc := &ClusterAsset{
		ClusterID:    "cluster",
		Architecture: "amd64",
		Platform:     "libvirt",
		OSImage:      OSImage{Type: "nestos"},
		UserName:     "root",
		Password:     "123",
		SSHKey:       "./assets.go",
		Master: []NodeAsset{
			{
				Hostname: "k8s-master01",
				IP:       "127.0.0.1",
				HardwareInfo: HardwareInfo{
					CPU:  2,
					RAM:  2048,
					Disk: 30,
				},
			},
		},
		Worker: []NodeAsset{
			{
				Hostname: "k8s-worker01",
				IP:       "127.0.0.1",
				HardwareInfo: HardwareInfo{
					CPU:  2,
					RAM:  2048,
					Disk: 30,
				},
			},
		},
		Runtime: "crio",
		Kubernetes: Kubernetes{
			KubernetesVersion:    "v1.29.1",
			KubernetesAPIVersion: "v1beta3",
			ApiServerEndpoint:    "127.0.0.1:6443",
			ImageRegistry:        "registry.k8s.io",
			PauseImage:           "pause:3.9",
			Network: Network{
				ServiceSubnet: "127.0.0.1/16",
				PodSubnet:     "127.0.0.1/16",
			},
		},
		HookConf: HookConf{
			ShellFiles: []ShellFile{
				{Name: "sss"},
				{Name: "sss"},
			},
		},
	}

	t.Run("CheckStringValue Success", func(t *testing.T) {
		err := CheckStringValue(&cc.ClusterID, opts.ClusterID, "cluster")
		if err != nil {
			t.Logf("CheckStringValue failed: %v", err)
		}
	})

	t.Run("InitClusterAsset Success", func(t *testing.T) {
		clusterConfig, err := cc.InitClusterAsset(opts)
		if err != nil || clusterConfig == nil {
			t.Errorf("InitClusterAsset failed: %v", err)
		}
	})

	t.Run("Delete Success", func(t *testing.T) {
		err := cc.Delete("sss")
		if err != nil {
			t.Errorf("Delete failed: %v", err)
		}
	})

	t.Run("Persist Success", func(t *testing.T) {
		err := cc.Persist("/tmp")
		if err != nil {
			t.Errorf("Persist failed: %v", err)
		}
	})

	t.Run("InitClusterAsset KubernetesAPIVersion Fail", func(t *testing.T) {
		opts.KubernetesAPIVersion = 21
		clusterConfig, err := cc.InitClusterAsset(opts)
		if err == nil || clusterConfig != nil {
			t.Log("Expected error, got nil")
		}
	})
}
