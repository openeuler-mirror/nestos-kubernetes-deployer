/*
Copyright 2023 KylinSoft  Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package workflow

import (
	"strings"
)

type Runner struct {
	Phases       []Phase
	phaseRunners []*phaseRunner
}

type phaseRunner struct {
	Phase
	parent        *phaseRunner
	level         int
	selfPath      []string
	generatedName string
	use           string
}

func NewRunner() *Runner {
	return &Runner{
		Phases: []Phase{},
	}
}

func (r *Runner) Run() error {
	r.prepareForExcution()

	err := r.VisitAll(func(p *phaseRunner) error {
		if p.Run != nil {
			if err := p.Run(); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (r *Runner) AppendPhase(t Phase) {
	r.Phases = append(r.Phases, t)
}

func (r *Runner) VisitAll(fn func(*phaseRunner) error) error {
	for _, currentRunner := range r.phaseRunners {
		if err := fn(currentRunner); err != nil {
			return err
		}
	}
	return nil
}
func (r *Runner) prepareForExcution() {
	r.phaseRunners = []*phaseRunner{}
	var parentRunner *phaseRunner
	for _, phase := range r.Phases {
		addPhaseRunner(r, parentRunner, phase)
	}
}

func addPhaseRunner(e *Runner, parentRunner *phaseRunner, phase Phase) {
	use := cleanName(phase.Name)
	generatedName := use
	selfPath := []string{generatedName}
	currentRunner := &phaseRunner{
		Phase:         phase,
		parent:        parentRunner,
		level:         len(selfPath) - 1,
		selfPath:      selfPath,
		generatedName: generatedName,
		use:           use,
	}
	e.phaseRunners = append(e.phaseRunners, currentRunner)

	for _, childPhase := range phase.Phases {
		addPhaseRunner(e, currentRunner, childPhase)
	}

}

func cleanName(name string) string {
	ret := strings.ToLower(name)
	return ret
}
