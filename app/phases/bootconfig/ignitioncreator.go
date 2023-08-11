package bootconfig

import "gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra"

type IgnitionAssembler struct {
}

func (i IgnitionAssembler) Assemble(assets infra.Assets) infra.InitConfig {
	return infra.InitConfig{}
}
