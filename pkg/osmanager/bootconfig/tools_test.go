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
package bootconfig

import (
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/globalconfig"
	"nestos-kubernetes-deployer/pkg/constants"
	"os"
	"testing"
)

func TestTool(t *testing.T) {
	err := os.Chdir("../../../data")
	if err != nil {
		t.Fatal(err)
	}
	clusterAsset := &asset.ClusterAsset{
		ClusterID: "cluster",
		Master: []asset.NodeAsset{
			{Hostname: "k8s-master01", IP: ""},
			{Hostname: "k8s-master02", IP: ""},
		},
		Worker: []asset.NodeAsset{
			{Hostname: "k8s-worker01", IP: ""},
		},
		SSHKey: "./assets.go",
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
			BootstrapIgnPort: "1234",
		},
	}

	var tmplData *TmplData
	t.Run("GetTmplData", func(t *testing.T) {
		tmplData, err = GetTmplData(clusterAsset)
		if err != nil {
			t.Log("test fail", err)
			return
		}
		t.Log("success")
	})

	t.Run("AppendStorageFiles Success", func(t *testing.T) {
		var files []File
		err := AppendStorageFiles(&files, "/", constants.BootConfigFilesPath, tmplData, []string{constants.InitClusterService})
		if err != nil {
			t.Log("test fail", err)
			return
		}
		t.Log("success")
	})

	t.Run("AppendStorageFiles Fail", func(t *testing.T) {
		var files []File
		err := AppendStorageFiles(&files, "/", "invalid/path", tmplData, []string{constants.InitClusterService})
		if err == nil {
			t.Log("expected failure, got success")
		}
	})

	t.Run("AppendSystemdUnits Success", func(t *testing.T) {
		var systemd Systemd
		err := AppendSystemdUnits(&systemd, constants.BootConfigSystemdPath, tmplData, []string{constants.InitClusterService})
		if err != nil {
			t.Log("test fail", err)
			return
		}
		t.Log("success")
	})

	t.Run("AppendSystemdUnits Fail", func(t *testing.T) {
		var systemd Systemd
		err := AppendSystemdUnits(&systemd, "invalid/path", tmplData, []string{constants.InitClusterService})
		if err == nil {
			t.Log("expected failure, got success")
		}
	})

	t.Run("GetSavePath Success", func(t *testing.T) {
		GetSavePath("cluster")
		t.Log("success")
	})

	t.Run("SaveYAML Success", func(t *testing.T) {
		err := SaveYAML(clusterAsset, "/tmp", "test.yaml", "header")
		if err != nil {
			t.Log("test fail", err)
		}
		t.Log("success")
	})

	t.Run("SaveYAML Fail", func(t *testing.T) {
		err := SaveYAML(clusterAsset, "/invalid/path", "", "header")
		if err == nil {
			t.Log("expected failure, got success")
		}
	})

	t.Run("SaveJSON Success", func(t *testing.T) {
		err := SaveJSON(clusterAsset, "/tmp", "test.json")
		if err != nil {
			t.Log("test fail", err)
		}
		t.Log("success")
	})

	t.Run("SaveJSON Fail", func(t *testing.T) {
		err := SaveJSON("", "", "")
		if err == nil {
			t.Log("expected failure, got success")
		}
	})

	t.Run("Marshal Success", func(t *testing.T) {
		_, err := Marshal(clusterAsset)
		if err != nil {
			t.Log("test fail", err)
		}
		t.Log("success")
	})

	t.Run("SaveFile Success", func(t *testing.T) {
		tbyte := []byte("This is a test content")
		err := SaveFile(tbyte, "/tmp", "testfile")
		if err != nil {
			t.Log("test fail", err)
		}
		t.Log("success")
	})

	t.Run("SaveFile Fail", func(t *testing.T) {
		err := SaveFile(nil, "/invalid/path", "testfile")
		if err == nil {
			t.Log("expected failure, got success")
		}
	})

	t.Run("CreateSetHostnameUnit Success", func(t *testing.T) {
		CreateSetHostnameUnit()
		t.Log("success")
	})

}
