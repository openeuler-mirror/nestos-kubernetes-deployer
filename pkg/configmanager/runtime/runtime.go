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

package runtime

import (
	"nestos-kubernetes-deployer/pkg/api"
	"nestos-kubernetes-deployer/pkg/constants"
	"strings"

	"github.com/pkg/errors"
)

var (
	mapRuntime = map[string]api.Runtime{
		constants.Isulad:     &isuladRuntime{},
		constants.Docker:     &dockerRuntime{},
		constants.Crio:       &crioRuntime{},
		constants.Containerd: &containerdRuntime{},
	}
)

func GetRuntime(runtime string) (api.Runtime, error) {
	runtime = strings.ToLower(runtime)
	if runtime == "" {
		return mapRuntime[constants.Isulad], nil
	}

	rt, ok := mapRuntime[runtime]
	if !ok {
		return nil, errors.New("unsupported runtime")
	}

	return rt, nil
}
