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

package config

import (
	"fmt"
	"nestos-kubernetes-deployer/app/apis/nkd"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadOrDefaultInitConfiguration(cfgPath string) (interface{}, string, error) {
	if cfgPath != "" {
		cfg, nodetype, err := LoadInitConfigurationFromFile(cfgPath)
		if err != nil {
			return nil, "", err
		}

		return cfg, nodetype, nil
	}

	return DefaultinitConfiguration()
}

func LoadInitConfigurationFromFile(cfg string) (interface{}, string, error) {
	node := new(nkd.Node)
	yamlFile, err := os.ReadFile(cfg)
	if err != nil {
		return nil, "", err
	}
	err = yaml.Unmarshal(yamlFile, node)
	if err != nil {
		return nil, "", err
	}

	switch node.Node {
	case "master":
		master := new(nkd.Master)
		masterInfo, err := Unmarshal(master, yamlFile)
		if err != nil {
			return nil, "", err
		}
		return masterInfo, "master", nil
	case "worker":
		worker := new(nkd.Worker)
		workerInfo, err := Unmarshal(worker, yamlFile)
		if err != nil {
			return nil, "", err
		}
		return workerInfo, "worker", nil
	default:
		return nil, "", fmt.Errorf("unsupported node type: %s", node.Node)
	}
}

func Unmarshal(nodeinfo interface{}, cfg []byte) (interface{}, error) {
	err := yaml.Unmarshal(cfg, nodeinfo)
	if err != nil {
		return nil, err
	}
	return nodeinfo, nil
}

func DefaultinitConfiguration() (*interface{}, string, error) {
	return nil, "", nil
}
