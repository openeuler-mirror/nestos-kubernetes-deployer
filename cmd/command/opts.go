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

import "github.com/spf13/cobra"

var (
	RootOptDir string
)

var clusterOpts struct {
	clusterId string
	gatherDeployOpts
}

type gatherDeployOpts struct {
	sshkey       string
	platform     string
	deployConfig string
	//
}

var HousekeeperOpts struct {
	OperatorImageUrl   string
	ControllerImageUrl string
}

func SetupDeployCmdOpts(deployCmd *cobra.Command) {
	flags := deployCmd.Flags()
	flags.StringVarP(&clusterOpts.deployConfig, "file", "f", "./deploy-config.yaml", "location of cluster deploy config file, default ./deploy-config.yaml")
	flags.StringVarP(&clusterOpts.clusterId, "cluster-id", "", "", "ClusterID of kubernetes cluster")
	flags.StringVarP(&clusterOpts.sshkey, "sshkey", "", "", "Path to SSH private keys that should be used for authentication.")
	flags.StringVarP(&clusterOpts.platform, "platform", "", "", "Select the infrastructure platform to deploy the cluster")
}
