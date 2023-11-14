package manager

import (
	"nestos-kubernetes-deployer/pkg/configmanager/asset/cluster"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/global"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func Init(cmd *cobra.Command) error {
	globalAsset := &global.GlobalAsset{}
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
	globalAsset.Initial()

	// Initialize cluster assets.
	clusterAsset.Initial()

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
