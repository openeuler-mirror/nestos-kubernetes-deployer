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

package openstack

// 定义 JSON 数据的结构
type TerraformData struct {
	Openstack struct {
		User_name   string `json:"user_name"`
		Password    string `json:"password"`
		Tenant_name string `json:"tenant_name"`
		Auth_url    string `json:"auth_url"`
		Region      string `json:"region"`
	} `json:"openstack"`
	Flavor struct {
		Name      string `json:"name"`
		Ram       string `json:"ram"`
		Vcpus     string `json:"vcpus"`
		Disk      string `json:"disk"`
		Is_public string `json:"is_public"`
	} `json:"flavor"`
	Instance struct {
		Count      string `json:"count"`
		Name       string `json:"name"`
		Image_name string `json:"image_name"`
		Key_pair   string `json:"key_pair"`
	} `json:"instance"`
	Floatip struct {
		Pool string `json:"pool"`
	} `json:"floatip"`
}
