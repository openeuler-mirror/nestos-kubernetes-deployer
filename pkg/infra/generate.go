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

type Platform interface {
	SetPlatform(asset.InfraAsset)
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

func (openstack *OpenStack) SetPlatform(infraAsset asset.InfraAsset) {
	if openstackAsset, ok := infraAsset.(*asset.OpenStackAsset); ok {
		openstack = &OpenStack{
			Username:          openstackAsset.UserName,
			Password:          openstackAsset.Password,
			Tenant_Name:       openstackAsset.Tenant_Name,
			Auth_URL:          openstackAsset.Auth_URL,
			Region:            openstackAsset.Region,
			Internal_Network:  openstackAsset.Internal_Network,
			External_Network:  openstackAsset.External_Network,
			Glance_Name:       openstackAsset.Glance_Name,
			Availability_Zone: openstackAsset.Availability_Zone,
		}
	}
}

type Libvirt struct {
}

func (libvirt *Libvirt) SetPlatform(infraAsset asset.InfraAsset) {
}

type Infra struct {
	Platform
	Master
	Worker
}

type Master struct {
	Count    int
	CPU      []int
	RAM      []int
	Disk     []int
	Hostname []string
	IP       []string
	Ign_Data []string
}

type Worker struct {
	Count    int
	CPU      []int
	RAM      []int
	Disk     []int
	Hostname []string
	IP       []string
	Ign_Data []string
}

func (infra *Infra) Generate(conf *asset.ClusterAsset, node string) error {
	switch conf.Platform {
	case "openstack", "Openstack", "OpenStack":
		infra.Platform = &OpenStack{}
	case "libvirt", "Libvirt":
		infra.Platform = &Libvirt{}
	default:
		return errors.New("unsupported platform")
	}

	infra.Platform.SetPlatform(conf.InfraPlatform)

	infra.Master.Count = conf.Master.Count
	for _, nodeAsset := range conf.Master.NodeAsset {
		infra.Master.CPU = append(infra.Master.CPU, nodeAsset.CPU)
		infra.Master.RAM = append(infra.Master.RAM, nodeAsset.RAM)
		infra.Master.Disk = append(infra.Master.Disk, nodeAsset.Disk)
		infra.Master.Hostname = append(infra.Master.Hostname, nodeAsset.Hostname)
		infra.Master.IP = append(infra.Master.IP, nodeAsset.IP)
		infra.Master.Ign_Data = append(infra.Master.Ign_Data, string(nodeAsset.Ign_Data))
	}

	infra.Worker.Count = conf.Worker.Count
	for _, nodeAsset := range conf.Worker.NodeAsset {
		infra.Worker.CPU = append(infra.Worker.CPU, nodeAsset.CPU)
		infra.Worker.RAM = append(infra.Worker.RAM, nodeAsset.RAM)
		infra.Worker.Disk = append(infra.Worker.Disk, nodeAsset.Disk)
		infra.Worker.Hostname = append(infra.Worker.Hostname, nodeAsset.Hostname)
		infra.Worker.IP = append(infra.Worker.IP, nodeAsset.IP)
		infra.Worker.Ign_Data = append(infra.Worker.Ign_Data, string(nodeAsset.Ign_Data))
	}

	outputFile, err := os.Create(filepath.Join(node, fmt.Sprintf("%s.tf", node)))
	if err != nil {
		return errors.Wrap(err, "failed to create terraform config file")
	}
	defer outputFile.Close()

	// Read template.
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

	// Writes the data to the file.
	if err = tmpl.Execute(outputFile, infra); err != nil {
		return errors.Wrap(err, "failed to write terraform config")
	}
	return nil
}
