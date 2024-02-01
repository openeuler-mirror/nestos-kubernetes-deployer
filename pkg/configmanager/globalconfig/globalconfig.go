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
)

func InitGlobalConfig(opts *opts.OptionsList) (*GlobalConfig, error) {
	globalAsset := &GlobalConfig{
		Log_Level: "default log level",
	}

	if opts.NKD.Log_Level != "" {
		globalAsset.Log_Level = opts.NKD.Log_Level
	}

	persistDir := opts.RootOptDir
	if err := os.MkdirAll(persistDir, 0644); err != nil {
		return nil, err
	}
	globalAsset.PersistDir = persistDir

	globalAsset.BootstrapIgnHost = opts.NKD.BootstrapIgnHost
	globalAsset.BootstrapIgnPort = opts.NKD.BootstrapIgnPort

	if globalAsset.BootstrapIgnHost == "" {
		if ip, err := utils.GetLocalIP(); err != nil {
			return nil, fmt.Errorf("failed to get local IP: %v", err)
		} else {
			globalAsset.BootstrapIgnHost = ip
		}
	}

	if globalAsset.BootstrapIgnPort == "" {
		// HTTP service default port
		globalAsset.BootstrapIgnPort = "9080"
	}

	if !utils.IsPortOpen(globalAsset.BootstrapIgnPort) {
		return nil, fmt.Errorf("The port %s is occupied.", globalAsset.BootstrapIgnPort)
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

// TODO: Delete deletes the global asset.
func (ga *GlobalConfig) Delete() error {
	return nil
}

// TODO: Persist persists the global asset.
func (ga *GlobalConfig) Persist() error {
	// TODO
	return nil
}
