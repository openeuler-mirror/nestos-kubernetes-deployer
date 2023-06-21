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
package constants

const (
	// LabelUpgrading is the key of the upgrading label for nodes
	LabelUpgrading = "upgrade.housekeeper.io/upgrading"
	// LabelKubeVersionPrefix defines the label associated with kubernetes version
	LabelKubeVersionPrefix = "upgrade.kubernetes.version.io/"
	// LabelMaster defines the label associated with master node.
	LabelMaster = "node-role.kubernetes.io/master"
)

// socket file
const (
	SockDir  = "/run/housekeeper-daemon"
	SockName = "housekeeper-daemon.sock"
)

