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

type IgnAsset struct {
	KubernetesVersion string
}

// TODO: Init inits the ign asset.
func (ia *IgnAsset) Initial() error {
	return nil
}

// TODO: Delete deletes the ign asset.
func (ia *IgnAsset) Delete() error {
	return nil
}

// TODO: Persist persists the ign asset.
func (ia *IgnAsset) Persist() error {
	// TODO: Serialize the cluster asset to json or yaml.
	return nil
}
