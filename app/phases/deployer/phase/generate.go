package phase

import "nestos-kubernetes-deployer/pkg/infra"

type GeneratePhase struct {
	Mode string // master, node
}

func (p GeneratePhase) GenerateAssets() infra.Assets {
	//managers = [
	//"etcd",
	//"manifests",
	//"certs",
	//]
	//assets = dict()
	//for i in managers:
	//	assets.merge(i.generateAssets())
	//return assets
	return infra.Assets{}
}
