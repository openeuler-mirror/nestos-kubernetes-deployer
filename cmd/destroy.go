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
	"nestos-kubernetes-deployer/cmd/command"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager"
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

	platform, err := cmd.Flags().GetString("platform")
	if err != nil {
		logrus.Errorf("Failed to get platform: %v", err)
		return err
	}
	if platform == "" {
		logrus.Errorf("platform is not provided: %v", err)
		return err
	}

	if err := configmanager.Initial(&opts.Opts); err != nil {
		logrus.Errorf("Failed to initialize configuration parameters: %v", err)
		return err
	}

	var infrastructure infra.Infrastructure
	switch platform {
	case strings.ToLower("libvirt"):
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
	case strings.ToLower("openstack"):
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
	case strings.ToLower("pxe"):
		logrus.Println("If necessary, manually delete the configuration file for deploying the PXE server")
		return err
	case strings.ToLower("ipxe"):
		logrus.Println("If necessary, manually delete the configuration file for deploying the iPXE server")
		return err
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
