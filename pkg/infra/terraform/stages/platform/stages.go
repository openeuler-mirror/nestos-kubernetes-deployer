package platform

import (
	"errors"
	"fmt"

	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/terraform"
	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/terraform/stages/openstack"
)

func StagesForPlatform(platform string, stage string) ([]terraform.Stage, error) {
	switch platform {
	case "openstack":
		openstack.AddPlatformStage(stage)
		return openstack.PlatformStages, nil
	default:
		return nil, errors.New(fmt.Sprintf("unsupported platform %q", platform))
	}
}
