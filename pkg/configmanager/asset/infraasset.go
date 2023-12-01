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

package asset

import (
	"errors"
	"nestos-kubernetes-deployer/cmd/command/opts"
)

type InfraAsset interface {
}

func InitInfraAsset(clusterAsset *ClusterAsset, opts *opts.OptionsList) (InfraAsset, error) {
	switch opts.Platform {
	case "openstack", "Openstack", "OpenStack":
		openstackAsset := clusterAsset.InfraPlatform.(*OpenStackAsset)
		infraAsset, err := initOpenStackAsset(openstackAsset, opts)
		if err != nil {
			return nil, err
		}
		return infraAsset, nil
	case "libvirt", "Libvirt":
		libvirtAsset := clusterAsset.InfraPlatform.(*LibvirtAsset)
		infraAsset, err := initLibvirtAsset(libvirtAsset, opts)
		if err != nil {
			return nil, err
		}
		return infraAsset, nil
	default:
		return nil, errors.New("unsupported platform")
	}
}

type OpenStackAsset struct {
	UserName          string
	Password          string
	Tenant_Name       string
	Auth_URL          string
	Region            string
	Internal_Network  string
	External_Network  string
	Glance_Name       string
	Availability_Zone string
}

func initOpenStackAsset(openstackAsset *OpenStackAsset, opts *opts.OptionsList) (*OpenStackAsset, error) {
	if err := checkStringValue(&openstackAsset.UserName, opts.InfraPlatform.OpenStack.UserName); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.Password, opts.InfraPlatform.OpenStack.Password); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.Tenant_Name, opts.InfraPlatform.OpenStack.Tenant_Name); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.Auth_URL, opts.InfraPlatform.OpenStack.Auth_URL); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.Region, opts.InfraPlatform.OpenStack.Region); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.Internal_Network, opts.InfraPlatform.OpenStack.Internal_Network); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.External_Network, opts.InfraPlatform.OpenStack.External_Network); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.Glance_Name, opts.InfraPlatform.OpenStack.Glance_Name); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.Availability_Zone, opts.InfraPlatform.OpenStack.Availability_Zone); err != nil {
		return nil, err
	}

	return openstackAsset, nil
}

type LibvirtAsset struct {
}

func initLibvirtAsset(libvirtAsset *LibvirtAsset, opts *opts.OptionsList) (*LibvirtAsset, error) {
	return libvirtAsset, nil
}
