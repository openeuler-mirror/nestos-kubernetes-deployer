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

package infra

import (
	"fmt"
	"html/template"
	"io"
	"nestos-kubernetes-deployer/data"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Infra struct {
	Platform
	Master
	Worker
}

type Platform struct {
	OpenStack
}

type OpenStack struct {
	Username          string
	Password          string
	Tenant_Name       string
	Auth_URL          string
	Region            string
	Internal_Network  string
	External_Network  string
	Glance_Name       string
	Availability_Zone string
}

type Master struct {
	Count int
	CPU   []int
	RAM   []int
	Disk  []int
}

type Worker struct {
	Count int
	CPU   []int
	RAM   []int
	Disk  []int
}

func (infra *Infra) Generate(conf *asset.ClusterAsset, node string) error {
	openstackAsset, ok := conf.InfraPlatform.(*asset.OpenStackAsset)
	if !ok {
		return errors.New("unsupported platform")
	}

	infra.Platform.OpenStack = OpenStack{
		Username:    openstackAsset.UserName,
		Password:    openstackAsset.Password,
		Tenant_Name: openstackAsset.Tenant_Name,
		Auth_URL:    openstackAsset.Auth_URL,
		Region:      openstackAsset.Region,
	}

	infra.Master.Count = conf.Master.Count
	for i := 0; i < infra.Master.Count; i++ {
		infra.Master.CPU = append(infra.Master.CPU, conf.Master.NodeAsset[i].CPU)
		infra.Master.RAM = append(infra.Master.RAM, conf.Master.NodeAsset[i].RAM)
		infra.Master.Disk = append(infra.Master.Disk, conf.Master.NodeAsset[i].Disk)
	}

	infra.Worker.Count = conf.Worker.Count
	for i := 0; i < infra.Worker.Count; i++ {
		infra.Worker.CPU = append(infra.Worker.CPU, conf.Worker.NodeAsset[i].CPU)
		infra.Worker.RAM = append(infra.Worker.RAM, conf.Worker.NodeAsset[i].RAM)
		infra.Worker.Disk = append(infra.Worker.Disk, conf.Worker.NodeAsset[i].Disk)
	}

	outputFile, err := os.Create(filepath.Join(node, fmt.Sprintf("%s.tf", node)))
	if err != nil {
		return errors.Wrap(err, "failed to create terraform config file")
	}
	defer outputFile.Close()

	tfFilePath := filepath.Join("terraform", conf.Platform, fmt.Sprintf("%s.tf.template", node))
	tfFile, err := data.Assets.Open(tfFilePath)
	if err != nil {
		return err
	}
	defer tfFile.Close()

	tfData, err := io.ReadAll(tfFile)
	if err != nil {
		return err
	}
	tmpl, err := template.New("terraform").Parse(string(tfData))
	if err != nil {
		return errors.Wrap(err, "failed to create terraform config template")
	}

	// 将填充后的数据写入文件
	if err = tmpl.Execute(outputFile, infra); err != nil {
		return errors.Wrap(err, "failed to write terraform config")
	}
	return nil
}
