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
	"os"

	"nestos-kubernetes-deployer/app/apis/nkd"

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
	cfg, err := DefaultinitConfiguration()
	if err != nil {
		return nil, "", err
	}
	return cfg, "", nil

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

	nodetype := node.Node

	if nodetype == "master" {
		master := new(nkd.Master)
		masterinfo, err := Unmarshal(master, yamlFile)
		if err != nil {
			return nil, "", err
		}
		return masterinfo, nodetype, nil

	} else if nodetype == "worker" {
		worker := new(nkd.Worker)
		workerinfo, err := Unmarshal(worker, yamlFile)
		if err != nil {
			return nil, "", err
		}
		return workerinfo, nodetype, nil

	} else {
		return nil, "", err
	}

}

func Unmarshal(nodeinfo interface{}, cfg []byte) (interface{}, error) {
	err := yaml.Unmarshal(cfg, nodeinfo)
	if err != nil {
		return nil, err
	}
	return nodeinfo, nil
}

func DefaultinitConfiguration() (*interface{}, error) {
	return nil, nil
}
