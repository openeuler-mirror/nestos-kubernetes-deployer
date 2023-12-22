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
	"runtime"
)

type InfraAsset interface {
}

func InitInfraAsset(clusterAsset *ClusterAsset, opts *opts.OptionsList) (InfraAsset, error) {
	if err := checkStringValue(&clusterAsset.Platform, opts.Platform, "platform"); err != nil {
		return nil, err
	}
	setStringValue(&clusterAsset.Architecture, opts.Arch, runtime.GOARCH)
	setStringValue(&clusterAsset.Platform, opts.Platform, "libvirt")
	switch clusterAsset.Platform {
	case "openstack", "Openstack", "OpenStack":
		openstackAsset, ok := convertMap(clusterAsset.InfraPlatform, "openstack")
		if !ok {
			return nil, errors.New("failed to get openstack asset")
		}
		infraAsset, err := initOpenStackAssetFromMap(openstackAsset, opts)
		if err != nil {
			return nil, err
		}
		return infraAsset, nil
	case "libvirt", "Libvirt":
		libvirtAsset, ok := convertMap(clusterAsset.InfraPlatform, "libvirt")
		if !ok {
			return nil, errors.New("failed to get libvirt asset")
		}
		infraAsset, err := initLibvirtAssetFromMap(libvirtAsset, opts, clusterAsset.Architecture)
		if err != nil {
			return nil, err
		}
		return infraAsset, nil
	default:
		return nil, errors.New("unsupported platform")
	}
}

func convertMap(inputMap interface{}, platform string) (map[string]interface{}, bool) {
	resultMap := make(map[string]interface{})

	if inputMap == nil {
		// If inputMap is nil, return an empty map corresponding to the platform structure.
		switch platform {
		case "openstack", "Openstack", "OpenStack":
			return map[string]interface{}{
				"username":          "",
				"password":          "",
				"tenant_name":       "",
				"auth_url":          "",
				"region":            "",
				"internal_network":  "",
				"external_network":  "",
				"glance_name":       "",
				"availability_zone": "",
			}, true
		case "libvirt", "Libvirt":
			return map[string]interface{}{
				"uri":          "",
				"osimage_path": "",
				"cidr":         "",
				"gateway":      "",
			}, true
		default:
			return resultMap, false
		}
	}

	// Check if the inputMap is of type map[interface{}]interface{}.
	if inputMap, ok := inputMap.(map[interface{}]interface{}); ok {
		for key, value := range inputMap {
			keyStr, ok := key.(string)
			if !ok {
				return resultMap, false
			}
			resultMap[keyStr] = value
		}
	} else {
		// If not, handle other types as needed.
		return resultMap, false
	}

	return resultMap, true
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

func initOpenStackAssetFromMap(openstackMap map[string]interface{}, opts *opts.OptionsList) (InfraAsset, error) {
	openstackAsset := &OpenStackAsset{}

	updateFieldFromMap("username", &openstackAsset.UserName, openstackMap)
	updateFieldFromMap("password", &openstackAsset.Password, openstackMap)
	updateFieldFromMap("tenant_name", &openstackAsset.Tenant_Name, openstackMap)
	updateFieldFromMap("auth_url", &openstackAsset.Auth_URL, openstackMap)
	updateFieldFromMap("region", &openstackAsset.Region, openstackMap)
	updateFieldFromMap("internal_network", &openstackAsset.Internal_Network, openstackMap)
	updateFieldFromMap("external_network", &openstackAsset.External_Network, openstackMap)
	updateFieldFromMap("glance_name", &openstackAsset.Glance_Name, openstackMap)
	updateFieldFromMap("availability_zone", &openstackAsset.Availability_Zone, openstackMap)

	if err := checkStringValue(&openstackAsset.UserName, opts.InfraPlatform.OpenStack.UserName, "openstack_username"); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.Password, opts.InfraPlatform.OpenStack.Password, "openstack_password"); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.Tenant_Name, opts.InfraPlatform.OpenStack.Tenant_Name, "openstack_tenant_name"); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.Auth_URL, opts.InfraPlatform.OpenStack.Auth_URL, "openstack_auth_url"); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.Region, opts.InfraPlatform.OpenStack.Region, "openstack_region"); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.Internal_Network, opts.InfraPlatform.OpenStack.Internal_Network, "openstack_internal_network"); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.External_Network, opts.InfraPlatform.OpenStack.External_Network, "openstack_external_network"); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.Glance_Name, opts.InfraPlatform.OpenStack.Glance_Name, "openstack_glance_name"); err != nil {
		return nil, err
	}
	if err := checkStringValue(&openstackAsset.Availability_Zone, opts.InfraPlatform.OpenStack.Availability_Zone, "openstack_availability_zone"); err != nil {
		return nil, err
	}

	return openstackAsset, nil
}

type LibvirtAsset struct {
	URI     string
	OSImage string
	CIDR    string
	Gateway string
}

func initLibvirtAssetFromMap(libvirtMap map[string]interface{}, opts *opts.OptionsList, arch string) (InfraAsset, error) {
	libvirtAsset := &LibvirtAsset{}

	updateFieldFromMap("uri", &libvirtAsset.URI, libvirtMap)
	updateFieldFromMap("osimage", &libvirtAsset.OSImage, libvirtMap)
	updateFieldFromMap("cidr", &libvirtAsset.CIDR, libvirtMap)
	updateFieldFromMap("gateway", &libvirtAsset.Gateway, libvirtMap)

	osImage := "https://nestos.org.cn/nestos20230928/nestos-for-container/x86_64/NestOS-For-Container-22.03-LTS-SP2.20230928.0-qemu.x86_64.qcow2"
	if arch == "arm64" || arch == "aarch64" {
		osImage = "https://nestos.org.cn/nestos20230928/nestos-for-container/aarch64/NestOS-For-Container-22.03-LTS-SP2.20230928.0-qemu.aarch64.qcow2"
	}

	setStringValue(&libvirtAsset.URI, opts.InfraPlatform.Libvirt.URI, "qemu:///system")
	setStringValue(&libvirtAsset.OSImage, opts.InfraPlatform.Libvirt.OSImage, osImage)
	setStringValue(&libvirtAsset.CIDR, opts.InfraPlatform.Libvirt.CIDR, "192.168.132.0/24")
	setStringValue(&libvirtAsset.Gateway, opts.InfraPlatform.Libvirt.Gateway, "192.168.132.1")

	return libvirtAsset, nil
}

func updateFieldFromMap(fieldName string, fieldValue *string, inputMap map[string]interface{}) {
	if value, ok := inputMap[fieldName]; ok {
		if strValue, ok := value.(string); ok && *fieldValue == "" {
			*fieldValue = strValue
		}
	}
}
