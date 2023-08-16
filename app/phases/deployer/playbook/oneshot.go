package playbook

import "nestos-kubernetes-deployer/pkg/deployer/phase"

func OneshotPlaybook() Playbook {
	return Playbook{
		phase.InitPhase{},
	}
}

//managers = [
//"etcd",
//"manifests",
//"certs",
//]
