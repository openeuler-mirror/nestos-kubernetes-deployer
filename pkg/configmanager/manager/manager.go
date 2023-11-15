package manager

import (
	"nestos-kubernetes-deployer/pkg/configmanager/asset/cluster"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/global"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func Init(cmd *cobra.Command) error {
	clusterAsset := &cluster.ClusterAsset{}

	// Parse command line arguments.
	file, _ := cmd.Flags().GetString("file")

	if file != "" {
		// Parse configuration file.
		yamlFile, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(yamlFile, clusterAsset)
		if err != nil {
			return err
		}
	}

	// TODO: 设置优先级 & 参数解析

	// Initialize global assets.
	if !global.IsInitial {
		globalAsset := &global.GlobalAsset{}
		err := globalAsset.Initial()
		if err != nil {
			return nil
		}
	}

	// Initialize cluster assets.
	err := clusterAsset.Initial()
	if err != nil {
		return nil
	}

	return nil
}

func Search() error {
	return nil
}

func Persist() error {
	return nil
}

func Delete() error {
	return nil
}
