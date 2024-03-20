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

package asset

import (
	"fmt"
	"strings"
)

var (
	mapRuntime = map[string]string{
		"isulad": "/var/run/isulad.sock",
		"docker": "/var/run/dockershim.sock",
		"crio":   "unix:///var/run/crio/crio.sock",
	}
)

func GetRuntimeCriSocket(runtime string) (string, error) {
	if content, ok := mapRuntime[strings.ToLower(runtime)]; ok {
		return content, nil
	}
	return "", fmt.Errorf("runtime %s not found", runtime)
}
