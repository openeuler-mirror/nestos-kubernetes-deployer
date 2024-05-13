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

package api

type OperatingSystem interface {
	GenerateResourceFiles() error
}

type BootConfigAPI interface {
	GenerateBootConfig() error
}

// Runtime 接口定义了获取运行时相关信息的方法。
type Runtime interface {
	// GetRuntimeCriSocket 返回与运行时相关的 CRI (Container Runtime Interface) 套接字地址
	GetRuntimeCriSocket() string
}
