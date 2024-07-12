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
package cloudinit

import (
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/globalconfig"
	"nestos-kubernetes-deployer/pkg/constants"
	"os"
	"testing"
)

func TestCloudinit(t *testing.T) {
	err := os.Chdir("../../../../data")
	if err != nil {
		t.Fatal(err)
	}
	clusterAsset := &asset.ClusterAsset{
		ClusterID: "cluster",
		Master: []asset.NodeAsset{
			{Hostname: "master"},
			{Hostname: "master02"},
		},
		Worker: []asset.NodeAsset{
			{Hostname: "worker"},
		},
		SSHKey: "./assets.go",
		HookConf: asset.HookConf{
			ShellFiles: []asset.ShellFile{
				{Name: "ssssss"},
				{Name: "ssssss"},
			},
		},
	}
	bootstrapBaseurl := "a/b"
	configmanager.GlobalConfig = &globalconfig.GlobalConfig{
		PersistDir: "/",
	}

	ci := NewCloudinit(clusterAsset, bootstrapBaseurl)

	t.Run("GenerateBootConfig_fail", func(t *testing.T) {
		err := ci.GenerateBootConfig()
		if err != nil {
			t.Log("test fail", err)
			return
		}
		t.Log("success")
	})
	configmanager.GlobalConfig.PersistDir = "./"

	t.Run("GenerateBootConfig_master_fail", func(t *testing.T) {
		clusterAsset.Runtime = constants.Containerd

		clusterAsset.SSHKey = "sasdsddssdsfd"
		err := ci.GenerateBootConfig()
		if err != nil {
			t.Log("test fail", err)
			return
		}
		t.Log("success")
	})
	clusterAsset.SSHKey = "./assets.go"

	t.Run("GenerateBootConfig", func(t *testing.T) {
		clusterAsset.Runtime = constants.Crio
		configmanager.GlobalConfig.PersistDir = "./"
		err := ci.GenerateBootConfig()
		if err != nil {
			t.Log("test fail", err)
			return
		}
		t.Log("success")

		if err := os.RemoveAll(clusterAsset.ClusterID); err != nil {
			t.Logf("Failed to remove cluster folder: %v", err)
		}
	})
}
