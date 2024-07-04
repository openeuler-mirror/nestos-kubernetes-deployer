/*
Copyright 2024 KylinSoft  Co., Ltd.

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
package osmanager

import (
	"fmt"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/osmanager/generalos"
	"nestos-kubernetes-deployer/pkg/osmanager/nestos"
	"strings"
)

const (
	nestosType    = "nestos"
	generalosType = "generalos"
)

type osmanager struct {
	config *asset.ClusterAsset
}

func NewOSManager(config *asset.ClusterAsset) *osmanager {
	return &osmanager{
		config: config,
	}
}

func (o *osmanager) GenerateOSConfig() error {
	if o.IsNestOS() {
		osDep, err := nestos.NewNestOS(o.config)
		if err != nil {
			return fmt.Errorf("error creating NestOS osmanager instance: %v", err)
		}
		if err := osDep.GenerateResourceFiles(); err != nil {
			return fmt.Errorf("error generating NestOS resource files: %v", err)
		}
		return nil
	} else if o.IsGeneralOS() {
		osDep, err := generalos.NewGeneralOS(o.config)
		if err != nil {
			return fmt.Errorf("error creating GeneralOS osmanager instance: %v", err)
		}
		if err := osDep.GenerateResourceFiles(); err != nil {
			return fmt.Errorf("error generating GeneralOS resource files: %v", err)
		}
		return nil
	}
	return fmt.Errorf("unsupported OS type: %s", o.config.OSImage.Type)
}

func (o *osmanager) IsNestOS() bool {
	if strings.ToLower(o.config.OSImage.Type) == nestosType {
		o.config.OSImage.IsNestOS = true
		return true
	}
	return false
}

func (o *osmanager) IsGeneralOS() bool {
	if strings.ToLower(o.config.OSImage.Type) == generalosType {
		o.config.OSImage.IsGeneralOS = true
		return true
	}
	return false
}
