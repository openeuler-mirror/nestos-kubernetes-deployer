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
package phases

import "gitee.com/openeuler/nestos-kubernetes-deployer/app/cmd/phases/workflow"

func NewGenerateIgnCmd() workflow.Phase {
	return workflow.Phase{
		Name:  "ign",
		Short: "run ign to genertate ignition file",
		Run:   runGenerateIgnConfig,
	}
}

func runGenerateIgnConfig(r workflow.RunData) error {

	return nil
}
