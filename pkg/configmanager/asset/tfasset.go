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

package asset

type TFAsset struct {
	OpenStack_Auth_Url string
	Count              int
	User_Data          string
}

// TODO: Init inits the tf asset.
func (tf *TFAsset) Initial() error {
	return nil
}

// TODO: Delete deletes the tf asset.
func (tf *TFAsset) Delete() error {
	return nil
}

// TODO: Persist persists the tf asset.
func (tf *TFAsset) Persist() error {
	// TODO: Serialize the cluster asset to json or yaml.
	return nil
}
