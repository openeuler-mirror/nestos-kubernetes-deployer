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
	"fmt"
	"io"
	"nestos-kubernetes-deployer/data"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/infraasset"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Infra struct {
	ClusterID   string
	Platform    interface{}
	Master      Node
	Worker      Node
	MachineType string
}

type OpenStack struct {
	Username         string
	Password         string
	TenantName       string
	AuthURL          string
	Region           string
	InternalNetwork  string
	ExternalNetwork  string
	GlanceName       string
	AvailabilityZone string
}

type Libvirt struct {
	URI     string
	OSImage string
	CIDR    string
	Gateway string
}

type Node struct {
	Count      int
	CPU        []string
	RAM        []string
	Disk       []string
	Hostname   []string
	IP         []string
	BootConfig []string
}

func (infra *Infra) Generate(conf *asset.ClusterAsset, node string) (err error) {
	infra.ClusterID = conf.ClusterID

	switch strings.ToLower(conf.Platform) {
	case "libvirt":
		libvirtAsset := conf.InfraPlatform.(*infraasset.LibvirtAsset)
		infra.Platform = &Libvirt{
			URI:     libvirtAsset.URI,
			OSImage: libvirtAsset.OSPath,
			CIDR:    libvirtAsset.CIDR,
			Gateway: libvirtAsset.Gateway,
		}
	case "openstack":
		openstackAsset := conf.InfraPlatform.(*infraasset.OpenStackAsset)
		infra.Platform = &OpenStack{
			Username:         openstackAsset.UserName,
			Password:         openstackAsset.Password,
			TenantName:       openstackAsset.TenantName,
			AuthURL:          openstackAsset.AuthURL,
			Region:           openstackAsset.Region,
			InternalNetwork:  openstackAsset.InternalNetwork,
			ExternalNetwork:  openstackAsset.ExternalNetwork,
			GlanceName:       openstackAsset.GlanceName,
			AvailabilityZone: openstackAsset.AvailabilityZone,
		}
	default:
		logrus.Errorf("unsupported platform")
		return err
	}

	if node == "master" {
		var (
			master_cpu        []uint
			master_ram        []uint
			master_disk       []uint
			master_hostname   []string
			master_ip         []string
			master_bootConfig []string
		)

		infra.Master.Count = len(conf.Master)
		for i, master := range conf.Master {
			master_cpu = append(master_cpu, master.CPU)
			master_ram = append(master_ram, master.RAM)
			master_disk = append(master_disk, master.Disk)
			master_hostname = append(master_hostname, master.Hostname)
			master_ip = append(master_ip, master.IP)
			if i == 0 {
				master_bootConfig = append(master_bootConfig, conf.BootConfig.Controlplane.Path)
				continue
			}
			master_bootConfig = append(master_bootConfig, conf.BootConfig.Master.Path)
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
		infra.Master.BootConfig, err = convertSliceToStrings(master_bootConfig)
		if err != nil {
			return err
		}
	} else if node == "worker" {
		var (
			worker_cpu        []uint
			worker_ram        []uint
			worker_disk       []uint
			worker_hostname   []string
			worker_ip         []string
			worker_bootConfig []string
		)

		infra.Worker.Count = len(conf.Worker)
		for _, worker := range conf.Worker {
			worker_cpu = append(worker_cpu, worker.CPU)
			worker_ram = append(worker_ram, worker.RAM)
			worker_disk = append(worker_disk, worker.Disk)
			if worker.IP == "" {
				worker.IP = "null"
			}
			worker_ip = append(worker_ip, worker.IP)
			worker_hostname = append(worker_hostname, worker.Hostname)
			worker_bootConfig = append(worker_bootConfig, conf.BootConfig.Worker.Path)
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
		infra.Worker.BootConfig, err = convertSliceToStrings(worker_bootConfig)
		if err != nil {
			return err
		}
	}

	var arch string
	switch conf.Architecture {
	case "amd64", "x86_64":
		infra.MachineType = "pc"
		arch = "x86_64"
	case "arm64", "aarch64":
		infra.MachineType = "virt"
		arch = "aarch64"
	default:
		logrus.Errorf("unsupported architecture")
		return err
	}

	persistDir := configmanager.GetPersistDir()
	if err := os.MkdirAll(filepath.Join(persistDir, conf.ClusterID, node), 0644); err != nil {
		return err
	}

	outputFile, err := os.Create(filepath.Join(persistDir, conf.ClusterID, node, fmt.Sprintf("%s.tf", node)))
	if err != nil {
		return errors.Wrap(err, "failed to create terraform config file")
	}
	defer outputFile.Close()

	// Read template.
	tfFilePath := filepath.Join("terraform", arch, strings.ToLower(conf.OSImage.Type), conf.Platform, fmt.Sprintf("%s.tf.template", node))
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
