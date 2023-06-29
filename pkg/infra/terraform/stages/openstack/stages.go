package openstack

import (
	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/terraform"
	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/terraform/providers"
	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/terraform/stages"
)

var PlatformStages = []terraform.Stage{}

func AddPlatformStage(stage string) {
	newStage := stages.NewStage(
		"openstack",
		stage,
		[]providers.Provider{providers.OpenStack},
	)

	PlatformStages = append(PlatformStages, newStage)
}
