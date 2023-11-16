/*
Copyright 2023 KylinSoft  Co., Ltd.

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

package global

import (
	"github.com/spf13/cobra"
)

var GlobalConfig *GlobalAsset

// ========== Package method ==========

func GetGlobalConfig() (*GlobalAsset, error) {
	return GlobalConfig, nil
}

// ========== Structure method ==========

type GlobalAsset struct {
	// Having previously set IsInitial as a global variable, it cannot be fixed to true after
	// GlobalAsset executes Initial for the first time, so it is placed in the struct first.
	IsInitial   bool
	NKD_Version string
	Log_Level   string
}

// TODO: Initial inits the global asset.
func (ga *GlobalAsset) Initial(cmd *cobra.Command) error {
	if ga.IsInitial {
		return nil
	}

	nkd_version, _ := cmd.Flags().GetString("version")
	if nkd_version != "" {
		ga.NKD_Version = nkd_version
	} else {
		ga.NKD_Version = "default nkd version"
	}

	log_level, _ := cmd.Flags().GetString("log-level")
	if log_level != "" {
		ga.Log_Level = log_level
	} else {
		ga.Log_Level = "default log level"
	}

	GlobalConfig = ga

	ga.IsInitial = true
	return nil
}

// TODO: Delete deletes the global asset.
func (ga *GlobalAsset) Delete() error {
	return nil
}

// TODO: Persist persists the global asset.
func (ga *GlobalAsset) Persist() error {
	// TODO: Serialize the global asset to json or yaml.
	return nil
}
