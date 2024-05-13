package osmanager

import (
	"fmt"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/osmanager/generalos"
	"nestos-kubernetes-deployer/pkg/osmanager/nestos"
	"strings"
)

const (
	nestosType = "nestos"
	eulerType  = "openeuler"
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
	} else if o.IsOpenEuler() {
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
	return strings.ToLower(o.config.OSImage.Type) == nestosType
}

func (o *osmanager) IsOpenEuler() bool {
	return strings.ToLower(o.config.OSImage.Type) == eulerType
}
