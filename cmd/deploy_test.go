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
package cmd

import (
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/infraasset"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestDeploy(t *testing.T) {
	err := os.Chdir("../data")
	if err != nil {
		t.Fatal(err)
	}
	opts.Opts.RootOptDir = "./"

	cc := &asset.ClusterAsset{
		ClusterID:    "cluster",
		Architecture: "amd64",
		Platform:     "pxe",
		InfraPlatform: infraasset.PXEAsset{
			IP:             "",
			HTTPServerPort: "10",
			HTTPRootDir:    "./",
			TFTPServerPort: "20",
			TFTPRootDir:    "./",
		},
		OSImage:  asset.OSImage{Type: "nestos"},
		UserName: "root",
		SSHKey:   "./test.yaml",
		Password: "123",
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
	}

	clusterData, err := yaml.Marshal(cc)
	if err != nil {
		return
	}
	if err := os.WriteFile("test.yaml", clusterData, 0644); err != nil {
		return
	}

	cmd := NewDeployCommand()
	args := []string{"--file", "test.yaml"}
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("Failed to execute command: %v", err)
	}

	t.Run("DeployCmd Fail", func(t *testing.T) {
		if err := runDeployCmd(cmd, args); err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("clusterCreatePost Fail", func(t *testing.T) {
		if err := clusterCreatePost(cc); err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("deployHousekeeper Fail", func(t *testing.T) {
		err := deployHousekeeper(nil, "./test.yaml")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("applyNetworkPlugin Fail", func(t *testing.T) {
		err := applyNetworkPlugin("./", true)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
