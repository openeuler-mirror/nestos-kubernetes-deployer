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

func InitInfraAsset(opts *opts.OptionsList) (InfraAsset, error) {
	switch opts.Platform {
	case "openstack", "Openstack", "OpenStack":
		openstackAsset, err := initOpenStackAsset(opts)
		if err != nil {
			return nil, err
		}
		return openstackAsset, nil
	case "libvirt", "Libvirt":
		libvirtAsset, err := initLibvirtAsset(opts)
		if err != nil {
			return nil, err
		}
		return libvirtAsset, nil
	default:
		return nil, errors.New("unsupported platform")
	}
}

type OpenStackAsset struct {
	UserName    string
	Password    string
	Tenant_Name string
	Auth_URL    string
	Region      string
}

func initOpenStackAsset(opts *opts.OptionsList) (*OpenStackAsset, error) {
	openstackAsset := &OpenStackAsset{}
	return openstackAsset, nil
}

type LibvirtAsset struct {
}

func initLibvirtAsset(opts *opts.OptionsList) (*LibvirtAsset, error) {
	libvirtAsset := &LibvirtAsset{}
	return libvirtAsset, nil
}
