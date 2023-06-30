package stages

import (
	"fmt"

	"gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/terraform/providers"
)

type StageOption func(*SplitStage)

func NewStage(platform, name string, providers []providers.Provider, opts ...StageOption) SplitStage {
	s := SplitStage{
		platform:  platform,
		name:      name,
		providers: providers,
	}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

type SplitStage struct {
	platform  string
	name      string
	providers []providers.Provider
}

func (s SplitStage) Name() string {
	return s.name
}

func (s SplitStage) Providers() []providers.Provider {
	return s.providers
}

func (s SplitStage) StateFilename() string {
	return fmt.Sprintf("terraform.%s.tfstate", s.name)
}

func (s SplitStage) OutputsFilename() string {
	return fmt.Sprintf("%s.tfvars.json", s.name)
}
