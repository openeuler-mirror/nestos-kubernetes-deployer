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
	"context"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/infraasset"
	"os"
	"testing"
	"time"
)

func TestExtend(t *testing.T) {
	globalYaml := "global_config.yaml"
	opts.Opts.RootOptDir = "../data"
	cmd := NewExtendCommand()
	args := []string{"--num", "0", "--cluster-id", "k8s-007"}
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Logf("Failed to execute command: %v", err)
	}

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
	}

	t.Run("ExtendCmd Fail", func(t *testing.T) {
		if err := runExtendCmd(cmd, args); err == nil {
			t.Log("Expected error, got nil")
		}
		// Clean up
		if err := os.RemoveAll("logs"); err != nil {
			t.Logf("Failed to remove logs folder: %v", err)
		}

		if _, err := os.Stat(globalYaml); os.IsNotExist(err) {
			t.Logf("Expected global_config.yaml to be created, but it does not exist")
		}

		if err := os.Remove(globalYaml); err != nil {
			t.Logf("Failed to remove global_config.yaml: %v", err)
		}
		// Clean up
		if err := os.RemoveAll("logs"); err != nil {
			t.Logf("Failed to remove logs folder: %v", err)
		}

		if _, err := os.Stat(globalYaml); os.IsNotExist(err) {
			t.Logf("Expected global_config.yaml to be created, but it does not exist")
		}

		if err := os.Remove(globalYaml); err != nil {
			t.Logf("Failed to remove global_config.yaml: %v", err)
		}
	})
	t.Run("extendCluster Fail", func(t *testing.T) {
		err := extendCluster(cc, 1)
		if err != nil {
			t.Log(err)
		}
	})
	t.Run("extendArray Fail", func(t *testing.T) {
		err := extendArray(cc, 1)
		if err != nil {
			t.Log(err)
		}
	})
	t.Run("checkNodesReady Fail", func(t *testing.T) {
		err := checkNodesReady(context.Background(), cc, 1)
		if err != nil {
			t.Log(err)
		}
	})
	//t.Run("extendCluster Fail", func(t *testing.T) {
	//	_, err := getReadyNodesCount(context.Background(), nil)
	//	if err != nil {
	//		t.Log(err)
	//	}
	//})
	t.Run("waitForMinimumReadyNodes Fail", func(t *testing.T) {
		err := waitForMinimumReadyNodes(context.Background(), nil, 1, time.Second)
		if err != nil {
			t.Log(err)
		}
	})
}
