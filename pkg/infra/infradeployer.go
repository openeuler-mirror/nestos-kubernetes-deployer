package infra

type InfraSpec struct{
	diskSize string
	memorySize string
	image string
	clusterCIDR string
	initFileContent string
}


type InfraDeployer interface {
	Create(spec InfraSpec, config InitConfig) error
}



