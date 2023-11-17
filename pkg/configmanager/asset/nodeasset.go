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

import "github.com/spf13/cobra"

type NodeAsset struct {
	HardwareInfo
	OSImage string
}

type HardwareInfo struct {
	CPU  int
	RAM  int
	Disk int
}

// Initializes the node properties.
func InitNodeAsset(cmd *cobra.Command, nodeType string) NodeAsset {
	node := NodeAsset{}

	cpu, _ := cmd.Flags().GetInt(nodeType + "-cpu")
	ram, _ := cmd.Flags().GetInt(nodeType + "-ram")
	disk, _ := cmd.Flags().GetInt(nodeType + "-disk")

	setIntValue(&node.HardwareInfo.CPU, cpu, 4)
	setIntValue(&node.HardwareInfo.RAM, ram, 8)
	setIntValue(&node.HardwareInfo.Disk, disk, 50)

	return node
}
