package deployer

import (
	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg"
	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/deployer/phase"
	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra"
)
// TODO 断点续传机制

func InitCluster(filename string) error {
	// TODO cluster和vm的状态机制, init, pending,
	// cluster state:  pending, running, stopped
	// infra state: none, pending, created, starting, running, stopping, stopped
	// k8s state: TODO


	clusterInfo := ParseFromFile(filename)
	infraDeployer := infra.GetInfraDeployer(clusterInfo.InfraDriver)
	bootConfigAssembler := infra.GetBootConfigAssembler(clusterInfo.InfraDriver)

	infraSpec := parseSpecFromClusterInfo(clusterInfo)
	assets := generateInitAssets()
	bootConfig := bootConfigAssembler.Assemble(assets)
	_ = infraDeployer.Create(infraSpec, bootConfig)
	apiMonitor := pkg.ApiMonitor{clusterInfo.ApiEndpoint}
	apiMonitor.WaitForMastersReady(0)
	return nil
}



func JoinToCluster(filename string) {
	clusterInfo := ParseFromFile(filename)
	infraDeployer := infra.GetInfraDeployer(clusterInfo.InfraDriver)
	apiMonitor := pkg.ApiMonitor{clusterInfo.ApiEndpoint}

	generator := infra.GetBootConfigAssembler(clusterInfo.InfraDriver)
	initConfig := generator.Assemble(generateJoinAssets())
	spec := parseSpecFromClusterInfo(clusterInfo)
	_ = infraDeployer.Create(spec, initConfig)
	apiMonitor.WaitForWorkersReady(0)
}

func OneShot(filename string){
	clusterInfo := ParseFromFile(filename)
	infraDeployer := infra.GetInfraDeployer(clusterInfo.InfraDriver)
	apiMonitor := pkg.ApiMonitor{clusterInfo.ApiEndpoint}

	generatePhase := phase.GeneratePhase{"master"}
	assets := generatePhase.GenerateAssets()

	generator := infra.GetBootConfigAssembler(clusterInfo.OsType)
	masterInitConfig := generator.Assemble(assets)
	masterSpec := parseSpecFromClusterInfo(clusterInfo)
	_ = infraDeployer.Create(masterSpec, masterInitConfig)
	apiMonitor.WaitForMastersReady(0)


	initConfig := generator.Assemble(generateJoinAssets())
	spec := parseSpecFromClusterInfo(clusterInfo)
	_ = infraDeployer.Create(spec, initConfig)
	apiMonitor.WaitForWorkersReady(0)

}

func parseSpecFromClusterInfo(info ClusterInfo) infra.InfraSpec {
	return infra.InfraSpec{}
}



func generateJoinAssets() infra.Assets{
	return infra.Assets{}
}
