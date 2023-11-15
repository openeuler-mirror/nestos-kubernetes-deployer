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

var (
	GlobalConfig *GlobalAsset
	IsInitial    = false
)

// ========== 包方法 ==========

func GetGlobalConfig() (*GlobalAsset, error) {
	return GlobalConfig, nil
}

// ========== 模块方法 ==========

type GlobalAsset struct {
	NKD_Version string
}

// TODO: Init inits the global asset.
func (ga *GlobalAsset) Initial() error {
	ga.NKD_Version = "0.2.0"

	GlobalConfig = ga
	IsInitial = true
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
