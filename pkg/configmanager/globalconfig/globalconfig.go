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

package globalconfig

import (
	"github.com/spf13/cobra"
)

// ========== Structure method ==========

type GlobalAsset struct {
	Log_Level string
}

func InitGlobalConfig(cmd *cobra.Command) (*GlobalAsset, error) {
	globalAsset := &GlobalAsset{}

	log_level, _ := cmd.Flags().GetString("log-level")
	if log_level != "" {
		globalAsset.Log_Level = log_level
	} else {
		globalAsset.Log_Level = "default log level"
	}

	return globalAsset, nil
}

// TODO: Delete deletes the global asset.
func (ga *GlobalAsset) Delete() error {
	return nil
}

// TODO: Persist persists the global asset.
func (ga *GlobalAsset) Persist() error {
	// TODO
	return nil
}
