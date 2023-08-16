package playbook

import (
	"nestos-kubernetes-deployer/pkg/deployer"
	"nestos-kubernetes-deployer/pkg/deployer/phase"
)

type Playbook struct {
	clusterInfo deployer.ClusterInfo
	phases      []phase.Phase
}

func (p Playbook) Start() {
	for _, i := range p.phases {
		p.Run(i)
	}
}

func (p *Playbook) AddPhase(phase phase.Phase) {
	p.phases = append(p.phases, phase)
}

func (p Playbook) Run(phase phase.Phase) {

}
