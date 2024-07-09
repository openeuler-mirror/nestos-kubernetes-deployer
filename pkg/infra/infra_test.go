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

package infra

import (
	"testing"
)

func TestInfra(t *testing.T) {
	p := InfraPlatform{}

	t.Run("TestLibvirt", func(t *testing.T) {
		libvirtMaster := &Libvirt{
			PersistDir: "/tmp",
			ClusterID:  "cluster",
			Node:       "master",
			Count:      1,
		}
		p.SetInfra(libvirtMaster)
		if err := p.Deploy(); err != nil {
			t.Log("test fail", err)
		}
		if err := p.Extend(); err != nil {
			t.Log("test fail", err)
		}
		if err := p.Destroy(); err != nil {
			t.Log("test fail", err)
		}
	})

	t.Run("TestOpenStack", func(t *testing.T) {
		openstackMaster := &OpenStack{
			PersistDir: "tmp",
			ClusterID:  "cluster",
			Node:       "master",
			Count:      1,
		}
		p.SetInfra(openstackMaster)
		if err := p.Deploy(); err != nil {
			t.Log("test fail", err)
		}
		if err := p.Extend(); err != nil {
			t.Log("test fail", err)
		}
		if err := p.Destroy(); err != nil {
			t.Log("test fail", err)
		}
	})
}
