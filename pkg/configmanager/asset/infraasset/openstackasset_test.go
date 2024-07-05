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
package infraasset

import (
	"nestos-kubernetes-deployer/cmd/command/opts"
	"testing"
)

func TestOpenstack(t *testing.T) {
	opts := &opts.OptionsList{
		RootOptDir: "/tmp",
		InfraPlatform: opts.InfraPlatform{
			OpenStack: opts.OpenStack{
				UserName:         "admin",
				Password:         "",
				TenantName:       "admin",
				AuthURL:          "http://controller:5000/v3",
				Region:           "RegionOne",
				InternalNetwork:  "internal-net",
				ExternalNetwork:  "provider-flat-net",
				GlanceName:       "nestos.qcow2",
				AvailabilityZone: "Phytium",
			},
		},
	}

	oa := OpenStackAsset{}

	t.Run("InitAsset Fail", func(t *testing.T) {
		_, err := oa.InitAsset(nil, opts, nil)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	opts.InfraPlatform.OpenStack.Password = "123456"
	t.Run("InitAsset Success", func(t *testing.T) {
		_, err := oa.InitAsset(nil, opts, nil)
		if err != nil {
			t.Errorf("InitAsset failed: %v", err)
		}
	})
}
