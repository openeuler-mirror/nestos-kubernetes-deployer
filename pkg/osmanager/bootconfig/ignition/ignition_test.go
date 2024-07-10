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
package ignition

import (
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/globalconfig"
	"nestos-kubernetes-deployer/pkg/constants"
	"os"
	"testing"
)

func TestIgnition(t *testing.T) {
	err := os.Chdir("../../../../data")
	if err != nil {
		t.Fatal(err)
	}
	clusterAsset := &asset.ClusterAsset{
		ClusterID: "cluster",
		Platform:  "libvirt",
		Master: []asset.NodeAsset{
			{Hostname: "k8s-master01", IP: "192.168.1.1"},
			{Hostname: "k8s-master02", IP: "192.168.1.2"},
		},
		Worker: []asset.NodeAsset{
			{Hostname: "k8s-worker01", IP: "192.168.1.6"},
		},
		SSHKey: "./assets.go",
		HookConf: asset.HookConf{
			ShellFiles: []asset.ShellFile{
				{Name: "sss"},
				{Name: "sss"},
			},
		},
	}
	bootstrapBaseurl := "a/b"
	configmanager.GlobalConfig = &globalconfig.GlobalConfig{
		PersistDir: "./",
		BootstrapUrl: globalconfig.BootstrapUrl{
			BootstrapIgnHost: "127.0.0.1",
			BootstrapIgnPort: "9080",
		},
	}

	ci := NewIgnition(clusterAsset, bootstrapBaseurl)

	t.Run("GenerateBootConfig", func(t *testing.T) {
		clusterAsset.Runtime = constants.Crio
		err := ci.GenerateBootConfig()
		if err != nil {
			t.Log("test fail", err)
			return
		}
		t.Log("success")
	})

	t.Run("GenerateBootConfig_fail", func(t *testing.T) {
		clusterAsset.Runtime = "podman"
		err := ci.GenerateBootConfig()
		if err == nil {
			t.Log("expected failure, got success")
		}
	})
}
