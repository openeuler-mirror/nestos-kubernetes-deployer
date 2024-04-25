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
package ignition

import (
	"net/url"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
)

// set-hostname.service用于动态设置节点hostname，
// ${hostname}变量会通过基础设施模块(terraform)为变量赋值(参见https://developer.hashicorp.com/terraform/language/functions/templatefile)
// terraform模板文件修改部分：templatefile(var.instance_ign[count.index], { hostname = var.instance_hostname[count.index] })
func createSetHostnameUnit() string {
	unit := `[Unit]
Description=Set hostname
ConditionPathExists=!/var/log/set-hostname.stamp
[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=/usr/bin/hostnamectl set-hostname ${hostname}
ExecStart=/bin/touch /var/log/set-hostname.stamp
[Install]
WantedBy=multi-user.target`
	return unit
}

func generateMergeIgnition(bootstrapIgnitionHost string, role string) *igntypes.Config {
	setHostnameUnit := createSetHostnameUnit()

	ign := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
			Config: igntypes.IgnitionConfig{
				Merge: []igntypes.Resource{{
					Source: ignutil.StrToPtr(func() *url.URL {
						return &url.URL{
							Scheme: "http",
							Host:   bootstrapIgnitionHost,
							Path:   role,
						}
					}().String()),
				}},
			},
		},
		Systemd: igntypes.Systemd{
			Units: []igntypes.Unit{
				{
					Contents: &setHostnameUnit,
					Name:     "set-hostname.service",
					Enabled:  ignutil.BoolToPtr(true),
				},
			},
		},
	}
	return &ign
}
