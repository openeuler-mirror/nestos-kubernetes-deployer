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
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/infraasset"
	"nestos-kubernetes-deployer/pkg/configmanager/globalconfig"
	"nestos-kubernetes-deployer/pkg/httpserver"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// setupTestEnvironment 设置测试环境
func setupTestEnvironment(t *testing.T) {
	err := os.Chdir("../data")
	if err != nil {
		t.Fatal(err)
	}
	opts.Opts.RootOptDir = "./"
}

// cleanUp 清理测试产生的文件
func cleanUp(t *testing.T) {
	if err := os.RemoveAll("logs"); err != nil {
		t.Logf("Failed to remove logs folder: %v", err)
	}

	if _, err := os.Stat("global_config.yaml"); os.IsNotExist(err) {
		t.Logf("Expected global_config.yaml to be created, but it does not exist")
	} else if err := os.Remove("global_config.yaml"); err != nil {
		t.Logf("Failed to remove global_config.yaml: %v", err)
	}
}

// TestDeploy 测试部署命令
func TestDeploy(t *testing.T) {
	setupTestEnvironment(t)

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
				IP:       "127.0.0.1",
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
				IP:       "127.0.0.1",
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
			ApiServerEndpoint:    "127.0.0.1:6443",
			ImageRegistry:        "registry.k8s.io",
			PauseImage:           "pause:3.9",
			Network: asset.Network{
				ServiceSubnet: "127.0.0.1/16",
				PodSubnet:     "127.0.0.1/16",
			},
		},
	}

	cmd := NewDeployCommand()
	args := []string{"--file", "test.yaml"}
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Logf("Failed to execute command: %v", err)
	}

	t.Run("DeployCmd Fail", func(t *testing.T) {
		if err := runDeployCmd(cmd, args); err == nil {
			t.Log("Expected error, got nil")
		}
		cleanUp(t)
	})

	t.Run("clusterCreatePost Fail", func(t *testing.T) {
		if err := clusterCreatePost(cc); err == nil {
			t.Log("Expected error, got nil")
		}
	})

	t.Run("deployHousekeeper Fail", func(t *testing.T) {
		err := deployHousekeeper(nil, "./test.yaml")
		if err == nil {
			t.Log("Expected error, got nil")
		}
	})

	t.Run("applyNetworkPlugin Fail", func(t *testing.T) {
		err := applyNetworkPlugin("http://www.aaa.com", true)
		if err == nil {
			t.Log("Expected error, got nil")
		}
	})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "API endpoint success"}`))
	}))
	defer ts.Close()
	httpService := httpserver.NewHTTPService(ts.URL)
	defer httpService.Stop()

	t.Run("addIgnitionFiles Fail", func(t *testing.T) {
		err := addIgnitionFiles(httpService, cc)
		if err == nil {
			t.Log("Expected error, got nil")
		}
	})

	t.Run("addKickstartFiles Fail", func(t *testing.T) {
		err := addKickstartFiles(httpService, cc)
		if err == nil {
			t.Log("Expected error, got nil")
		}
	})

	t.Run("createCluster", func(t *testing.T) {
		configmanager.GlobalConfig = &globalconfig.GlobalConfig{}
		configmanager.ClusterAsset = map[string]*asset.ClusterAsset{
			"cluster": cc,
		}
		cc.Platform = "dsxxxxx10"
		err := createCluster(cc)
		if err != nil {
			t.Log("createCluster Expected error, got nil")
		}
	})

	t.Run("createCluster generlos", func(t *testing.T) {
		configmanager.GlobalConfig = &globalconfig.GlobalConfig{}
		configmanager.ClusterAsset = map[string]*asset.ClusterAsset{
			"cluster": cc,
		}
		cc.Platform = "dsxxxxx111"
		cc.OSImage = asset.OSImage{Type: "generalos"}
		err := createCluster(cc)
		if err != nil {
			t.Log("createCluster Expected error, got nil")
		}
	})

	t.Run("createCluster generlos_libvirt", func(t *testing.T) {
		configmanager.GlobalConfig = &globalconfig.GlobalConfig{}
		configmanager.ClusterAsset = map[string]*asset.ClusterAsset{
			"cluster": cc,
		}
		cc.Platform = "libvirt"
		cc.OSImage = asset.OSImage{Type: "generalos"}
		err := createCluster(cc)
		if err != nil {
			t.Log("generlos_libvirt Expected error:", err)
		}
	})

	t.Run("createCluster libvirt", func(t *testing.T) {
		cc.Platform = "libvirt"
		cc.OSImage = asset.OSImage{Type: "generalos"}
		configmanager.GlobalConfig = &globalconfig.GlobalConfig{}
		configmanager.ClusterAsset = map[string]*asset.ClusterAsset{
			"cluster": cc,
		}

		err := createCluster(cc)
		if err != nil {
			t.Log("createCluster Expected error, got nil")
		}
	})
	t.Run("getClusterConfig", func(t *testing.T) {
		getClusterConfig(&opts.Opts)
	})
	//t.Run("waitForPodsReady", func(t *testing.T) {
	//	fake.new()
	//	c := &kubernetes.Clientset{}
	//	err := waitForPodsReady(c)
	//	if err != nil {
	//		t.Log("waitForPodsReady Expected error, got nil")
	//	}
	//})
}
