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
package globalconfig

import (
	"nestos-kubernetes-deployer/cmd/command/opts"
	"testing"
)

func TestGlobalconfig(t *testing.T) {
	opts := &opts.OptionsList{
		RootOptDir: "/tmp",
	}
	file := "/tmp/" + GlobalConfigFile
	gc, err := InitGlobalConfig(opts)
	if err != nil {
		t.Fatalf("InitGlobalConfig returned an error: %v", err)
	}

	t.Run("InitGlobalConfig Success", func(t *testing.T) {
		_, err := InitGlobalConfig(opts)
		if gc == nil {
			t.Fatalf("InitGlobalConfig returned an error: %v", err)
		}
	})

	t.Run("Delete Success", func(t *testing.T) {
		err := gc.Delete(file)
		if err != nil {
			t.Errorf("Delete returned an error: %v", err)
		}
	})

	t.Run("Persist Success", func(t *testing.T) {
		err := gc.Persist()
		if err != nil {
			t.Errorf("Persist returned an error: %v", err)
		}
	})

	t.Run("InitGlobalConfig Fail", func(t *testing.T) {
		opts.RootOptDir = ""
		_, err := InitGlobalConfig(opts)
		if err == nil {
			t.Log("Expected error, got nil")
		}
	})

	t.Run("Delete Fail", func(t *testing.T) {
		err := gc.Delete("")
		if err != nil {
			t.Errorf("Delete returned an error: %v", err)
		}
	})
}
