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
package terraform

import (
	"fmt"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/infraasset"
	"nestos-kubernetes-deployer/pkg/configmanager/globalconfig"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	//err := os.Chdir("../../data")
	//if err != nil {
	//	t.Fatal(err)
	//}

	conf := &asset.ClusterAsset{
		Platform: "libvirt",
		Master: []asset.NodeAsset{
			{Hostname: "master"},
			{Hostname: "master02"},
		},
		Worker: []asset.NodeAsset{
			{Hostname: "worker"},
		},
		ClusterID: "k8s-001",
	}
	configmanager.GlobalConfig = &globalconfig.GlobalConfig{
		PersistDir: "./",
	}

	inf := &Infra{}

	node := "master"

	t.Run("Generate_libvirt", func(t *testing.T) {

		libvirt := &infraasset.LibvirtAsset{
			URI: "sss",
		}
		conf.InfraPlatform = libvirt
		err := inf.Generate(conf, node)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("success")
	})
	node = "worker"
	conf.Architecture = "x86_64"
	t.Run("Generate_openstack", func(t *testing.T) {
		openstack := &infraasset.OpenStackAsset{
			UserName: "ssss",
		}
		conf.InfraPlatform = openstack

		conf.Platform = "openstack"
		err := inf.Generate(conf, node)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("success")
	})
	rmDir := configmanager.GlobalConfig.PersistDir + conf.ClusterID
	err := os.RemoveAll(rmDir)
	if err != nil {
		fmt.Printf("Error removing directory: %s\n", err)
	} else {
		fmt.Println("Directory removed successfully")
	}
}
