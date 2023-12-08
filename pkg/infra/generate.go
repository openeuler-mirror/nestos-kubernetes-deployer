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
	"io"
	"nestos-kubernetes-deployer/data"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

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
		openstack.Username = openstackAsset.UserName
		openstack.Password = openstackAsset.Password
		openstack.Tenant_Name = openstackAsset.Tenant_Name
		openstack.Auth_URL = openstackAsset.Auth_URL
		openstack.Region = openstackAsset.Region
		openstack.Internal_Network = openstackAsset.Internal_Network
		openstack.External_Network = openstackAsset.External_Network
		openstack.Glance_Name = openstackAsset.Glance_Name
		openstack.Availability_Zone = openstackAsset.Availability_Zone
	}
}

type Libvirt struct {
	Username     string
	Remote_IP    string
	OSImage_Path string
}

func (libvirt *Libvirt) SetPlatform(infraAsset asset.InfraAsset) {
	if libvirtAsset, ok := infraAsset.(*asset.LibvirtAsset); ok {
		libvirt.Username = libvirtAsset.UserName
		libvirt.Remote_IP = libvirtAsset.Remote_IP
		libvirt.OSImage_Path = libvirtAsset.OSImage_Path
	}
}

type Infra struct {
	Platform
	Master Node
	Worker Node
}

type Node struct {
	Count    int
	CPU      []string
	RAM      []string
	Disk     []string
	Hostname []string
	IP       []string
	Ign_Path []string
}

func (infra *Infra) Generate(conf *asset.ClusterAsset, node string) (err error) {
	var (
		master_cpu      []uint
		master_ram      []uint
		master_disk     []uint
		master_hostname []string
		master_ip       []string
		master_ignPath  []string

		worker_cpu      []uint
		worker_ram      []uint
		worker_disk     []uint
		worker_hostname []string
		worker_ip       []string
		worker_ignPath  []string
	)

	switch conf.Platform {
	case "openstack", "Openstack", "OpenStack":
		infra.Platform = &OpenStack{}
	case "libvirt", "Libvirt":
		infra.Platform = &Libvirt{}
	default:
		return errors.New("unsupported platform")
	}

	infra.Platform.SetPlatform(conf.InfraPlatform)

	infra.Master.Count = len(conf.Master)
	for _, master := range conf.Master {
		master_cpu = append(master_cpu, master.CPU)
		master_ram = append(master_ram, master.RAM)
		master_disk = append(master_disk, master.Disk)
		master_hostname = append(master_hostname, master.Hostname)
		master_ip = append(master_ip, master.IP)
		master_ignPath = append(master_ignPath, master.Ign_Path)
	}
	infra.Master.CPU, err = convertSliceToStrings(master_cpu)
	if err != nil {
		return err
	}
	infra.Master.RAM, err = convertSliceToStrings(master_ram)
	if err != nil {
		return err
	}
	infra.Master.Disk, err = convertSliceToStrings(master_disk)
	if err != nil {
		return err
	}
	infra.Master.Hostname, err = convertSliceToStrings(master_hostname)
	if err != nil {
		return err
	}
	infra.Master.IP, err = convertSliceToStrings(master_ip)
	if err != nil {
		return err
	}
	infra.Master.Ign_Path, err = convertSliceToStrings(master_ignPath)
	if err != nil {
		return err
	}

	infra.Worker.Count = len(conf.Worker)
	for _, worker := range conf.Worker {
		worker_cpu = append(worker_cpu, worker.CPU)
		worker_ram = append(worker_ram, worker.RAM)
		worker_disk = append(worker_disk, worker.Disk)
		worker_hostname = append(worker_hostname, worker.Hostname)
		worker_ip = append(worker_ip, worker.IP)
		worker_ignPath = append(worker_ignPath, worker.Ign_Path)
	}
	infra.Worker.CPU, err = convertSliceToStrings(worker_cpu)
	if err != nil {
		return err
	}
	infra.Worker.RAM, err = convertSliceToStrings(worker_ram)
	if err != nil {
		return err
	}
	infra.Worker.Disk, err = convertSliceToStrings(worker_disk)
	if err != nil {
		return err
	}
	infra.Worker.Hostname, err = convertSliceToStrings(worker_hostname)
	if err != nil {
		return err
	}
	infra.Worker.IP, err = convertSliceToStrings(worker_ip)
	if err != nil {
		return err
	}
	infra.Worker.Ign_Path, err = convertSliceToStrings(worker_ignPath)
	if err != nil {
		return err
	}

	persistDir := configmanager.GetPersistDir()
	if err := os.MkdirAll(filepath.Join(persistDir, conf.Cluster_ID, node), 0644); err != nil {
		return err
	}

	outputFile, err := os.Create(filepath.Join(persistDir, conf.Cluster_ID, node, fmt.Sprintf("%s.tf", node)))
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

func convertSliceToStrings(slice interface{}) ([]string, error) {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input is not a slice")
	}

	quotedStrs := make([]string, sliceValue.Len())
	for i := 0; i < sliceValue.Len(); i++ {
		element := sliceValue.Index(i)
		switch element.Kind() {
		case reflect.Uint:
			quotedStrs[i] = fmt.Sprintf(`"%d",`, element.Interface())
		case reflect.String:
			quotedStrs[i] = fmt.Sprintf(`"%s",`, element.Interface())
		default:
			return nil, fmt.Errorf("unsupported type in slice")
		}
	}

	if len(quotedStrs) > 0 {
		quotedStrs[len(quotedStrs)-1] = strings.TrimRight(quotedStrs[len(quotedStrs)-1], ",")
	}

	return quotedStrs, nil
}
