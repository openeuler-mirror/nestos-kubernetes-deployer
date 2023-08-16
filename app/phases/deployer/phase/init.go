package phase

import "nestos-kubernetes-deployer/pkg/infra"

type InitPhase struct {
}

func (p InitPhase) Do() {
	generator := infra.GetBootConfigAssembler(clusterInfo.OsType)
	masterInitConfig := generator.Assemble(generateInitAssets())
	masterSpec := parseSpecFromClusterInfo(clusterInfo)
	_ = infraDeployer.Create(masterSpec, masterInitConfig)
	apiMonitor.WaitForMastersReady(0)

}
