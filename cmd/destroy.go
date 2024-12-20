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
package cmd

import (
	"fmt"
	"nestos-kubernetes-deployer/cmd/command"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/infraasset"
	"nestos-kubernetes-deployer/pkg/infra"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewDestroyCommand() *cobra.Command {
	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy a kubernetes cluster",
		RunE:  runDestroyCmd,
	}
	command.SetupDestroyCmdOpts(destroyCmd)

	return destroyCmd
}

func runDestroyCmd(cmd *cobra.Command, args []string) error {
	cleanup := command.SetuploggerHook(opts.Opts.RootOptDir)
	defer cleanup()

	clusterID, err := cmd.Flags().GetString("cluster-id")
	if err != nil {
		logrus.Errorf("Failed to get cluster id: %v", err)
		return err
	}
	if clusterID == "" {
		logrus.Errorf("cluster-id is not provided: %v", err)
		return err
	}

	if err := configmanager.Initial(&opts.Opts); err != nil {
		logrus.Errorf("Failed to initialize configuration parameters: %v", err)
		return err
	}

	clusterConfig, err := configmanager.GetClusterConfig(clusterID)
	if err != nil {
		logrus.Errorf("Failed to get cluster config using the cluster id: %v", err)
		return err
	}

	var infrastructure infra.Infrastructure
	switch strings.ToLower(clusterConfig.Platform) {
	case "libvirt":
		persistDir := configmanager.GetPersistDir()

		infrastructure = &infra.Libvirt{
			PersistDir: persistDir,
			ClusterID:  clusterID,
			Node:       "worker",
			Count:      0,
		}
		if err := infrastructure.Destroy(); err != nil {
			logrus.Errorf("Failed to destroy worker nodes:%v", err)
			return err
		}

		infrastructure = &infra.Libvirt{
			PersistDir: persistDir,
			ClusterID:  clusterID,
			Node:       "master",
			Count:      0,
		}
		if err := infrastructure.Destroy(); err != nil {
			logrus.Errorf("Failed to destroy master nodes:%v", err)
			return err
		}
	case "openstack":
		persistDir := configmanager.GetPersistDir()

		infrastructure = &infra.OpenStack{
			PersistDir: persistDir,
			ClusterID:  clusterID,
			Node:       "worker",
			Count:      0,
		}
		if err := infrastructure.Destroy(); err != nil {
			logrus.Errorf("Failed to destroy worker nodes:%v", err)
			return err
		}

		infrastructure = &infra.OpenStack{
			PersistDir: persistDir,
			ClusterID:  clusterID,
			Node:       "master",
			Count:      0,
		}
		if err := infrastructure.Destroy(); err != nil {
			logrus.Errorf("Failed to destroy master nodes:%v", err)
			return err
		}
	case "pxe":
		logrus.Println("If necessary, manually destroy the config for the PXE server:\n",
			"1. Stop dhcpd service\n",
			fmt.Sprintf("2. Delete http root dir: %s\n", clusterConfig.InfraPlatform.(*infraasset.PXEAsset).HTTPRootDir),
			fmt.Sprintf("3. Delete tftp root dir: %s", clusterConfig.InfraPlatform.(*infraasset.PXEAsset).TFTPRootDir),
		)
	case "ipxe":
		logrus.Println("If necessary, manually destroy the config for the iPXE server:\n",
			"1. Stop dhcpd service\n",
			fmt.Sprintf("2. Delete ipxe config: %s\n", clusterConfig.InfraPlatform.(*infraasset.IPXEAsset).FilePath),
			fmt.Sprintf("3. Delete OS install tree: %s", clusterConfig.InfraPlatform.(*infraasset.IPXEAsset).OSInstallTreePath),
		)
	default:
		logrus.Errorf("unsupported platform")
		return err
	}

	// delete asset files
	if err := configmanager.Delete(clusterID); err != nil {
		logrus.Errorf("Failed to clean the asset files")
		return err
	}

	return nil
}
