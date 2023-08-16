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
	"io/ioutil"

	"nestos-kubernetes-deployer/app/apis/nkd"

	"gopkg.in/yaml.v2"
)

func LoadOrDefaultInitConfiguration(cfgPath string, cfg *nkd.Nkd) (*nkd.Nkd, error) {
	if cfgPath != "" {
		cfg, err := LoadInitConfigurationFromFile()
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}
	cfg, err := DefaultinitConfiguration()
	if err != nil {
		return nil, err
	}
	return cfg, nil

}

func LoadInitConfigurationFromFile() (*nkd.Nkd, error) {
	conf := new(nkd.Nkd)
	yamlFile, err := ioutil.ReadFile("test.yaml")
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, conf)
	fmt.Println(conf.Infra.Platform)
	return conf, nil
}

func DefaultinitConfiguration() (*nkd.Nkd, error) {
	return nil, nil
}
