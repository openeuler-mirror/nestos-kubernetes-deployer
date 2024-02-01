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

package globalconfig

import (
	"fmt"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/utils"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const GlobalConfigFile = "global_config.yaml"

func InitGlobalConfig(opts *opts.OptionsList) (*GlobalConfig, error) {
	globalAsset := &GlobalConfig{
		Log_Level:          "default log level",
		ClusterConfig_Path: "",
		PersistDir:         opts.RootOptDir, // default persist directory
		BootstrapUrl: BootstrapUrl{
			BootstrapIgnPort: "9080", // default port
		},
	}
	configFile := filepath.Join(globalAsset.PersistDir, GlobalConfigFile)
	if _, err := os.Stat(configFile); err == nil {
		configData, err := os.ReadFile(configFile)
		if err != nil {
			logrus.Errorf("Failed to read config file: %s\n", err)
			return nil, err
		}
		err = yaml.Unmarshal(configData, globalAsset)
		if err != nil {
			logrus.Errorf("Failed to unmarshal config data: %s\n", err)
			return nil, err
		}
	}

	if opts.NKD.Log_Level != "" {
		globalAsset.Log_Level = opts.NKD.Log_Level
	}
	if opts.NKD.BootstrapIgnHost != "" {
		globalAsset.BootstrapIgnHost = opts.NKD.BootstrapIgnHost
	}

	if opts.NKD.BootstrapIgnPort != "" {
		globalAsset.BootstrapIgnPort = opts.NKD.BootstrapIgnPort
	}
	if !utils.IsPortOpen(globalAsset.BootstrapIgnPort) {
		return nil, fmt.Errorf("The port %s is occupied.", globalAsset.BootstrapIgnPort)
	}

	if globalAsset.BootstrapIgnHost == "" {
		if ip, err := utils.GetLocalIP(); err != nil {
			logrus.Errorf("failed to get local IP: %v", err)
			return nil, err
		} else {
			globalAsset.BootstrapIgnHost = ip
		}
	}

	if err := os.MkdirAll(globalAsset.PersistDir, 0644); err != nil {
		return nil, err
	}

	if err := globalAsset.Persist(); err != nil {
		return nil, err
	}

	return globalAsset, nil
}

// ========== Structure method ==========

type GlobalConfig struct {
	Log_Level          string
	ClusterConfig_Path string
	PersistDir         string // default: /etc/nkd
	BootstrapUrl
}

type BootstrapUrl struct {
	BootstrapIgnHost string `yaml:"bootstrap_ign_host"`
	BootstrapIgnPort string `yaml:"bootstrap_ign_port"`
}

// Delete deletes the global asset.
func (ga *GlobalConfig) Delete(persistFilePath string) error {
	if _, err := os.Stat(persistFilePath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(persistFilePath); err != nil {
		logrus.Errorf("failed to delete global config file: %v", err)
		return err
	}

	return nil
}

func (ga *GlobalConfig) Persist() error {
	globalConfigData, err := yaml.Marshal(ga)
	if err != nil {
		logrus.Errorf("failed to marshal global config: %v", err)
		return err
	}

	if err := os.WriteFile(ga.PersistDir, globalConfigData, 0644); err != nil {
		logrus.Errorf("failed to write global config file: %v", err)
		return err
	}
	return nil
}
