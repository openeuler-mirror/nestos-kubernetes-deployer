/*
Copyright 2024 KylinSoft  Co., Ltd.

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

package infra

type Infrastructure interface {
	Deploy() error
	Extend() error
	Destroy() error
}

type InfraPlatform struct {
	infra Infrastructure
}

func (p *InfraPlatform) SetInfra(infra Infrastructure) {
	p.infra = infra
}

func (p *InfraPlatform) Deploy() error {
	return p.infra.Deploy()
}

func (p *InfraPlatform) Extend() error {
	return p.infra.Extend()
}

func (p *InfraPlatform) Destroy() error {
	return p.infra.Destroy()
}
