package infra

import "gitee.com/openeuler/nestos-kubernetes-deployer/pkg/bootconfig"

func GetInfraDeployer(driverType string) InfraDeployer {
	if driverType == "openstack"{
		return OpenstackDeployer{}
	}
	return nil
}


func GetBootConfigAssembler(osType string) BootConfigAssembler {
	if osType == "NestOS"{
		return bootconfig.IgnitionAssembler{}
	}
	return nil
}