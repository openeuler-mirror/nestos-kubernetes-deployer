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

package command

var (
	RootOptDir string
)

// 部署集群可选配置参数集合
var ClusterOpts struct {
	ClusterId string
	GatherDeployOpts
}

type GatherDeployOpts struct {
	SSHKey   string
	Platform string
	//
}

var HousekeeperOpts struct {
	//todo: 可选配置参数项
}
