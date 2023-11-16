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

// Set global data
var GlobalConfig *GlobalAsset

// ========== Structure method ==========

type GlobalAsset struct {
	Log_Level string
}

// TODO: Initial inits the global asset.
func (ga *GlobalAsset) Initial(cmd *cobra.Command) error {
	if err := ga.setGlobalAsset(cmd); err != nil {
		return err
	}
	GlobalConfig = ga

	return nil
}

func (ga *GlobalAsset) setGlobalAsset(cmd *cobra.Command) error {
	log_level, _ := cmd.Flags().GetString("log-level")
	if log_level != "" {
		ga.Log_Level = log_level
	} else {
		ga.Log_Level = "default log level"
	}

	return nil
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
