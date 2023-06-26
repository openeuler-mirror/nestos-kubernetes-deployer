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

package terraform

import (
	"context"
	"path/filepath"

	prov "gitee.com/openeuler/nestos-kubernetes-deployer/pkg/infra/terraform/providers"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/openshift/installer/data"
	"github.com/pkg/errors"
)

func unpack(workingDir string, platform string, target string) (err error) {
	err = data.Unpack(workingDir, filepath.Join(platform, target))
	if err != nil {
		return err
	}

	return nil
}

// terraform init
func TFInit(workingDir string, platform string, target string, terraformBinary string, providers []prov.Provider) (err error) {
	err = unpack(workingDir, platform, target)
	if err != nil {
		return errors.Wrap(err, "failed to unpack terraform modules")
	}

	tf, err := newTFExec(workingDir, terraformBinary)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec.")
	}

	// 如果想导入本地已有插件，需构建相应目录
	// 如openstack插件所需目录为plugins/registry.terraform.io/terraform-provider-openstack/openstack/1.51.1/linux_arm64/terraform-provider-openstack_v1.51.1
	return errors.Wrap(
		tf.Init(context.Background(), tfexec.PluginDir(filepath.Join(terraformBinary, "plugins"))),
		"failed doing terraform init.",
	)
}
