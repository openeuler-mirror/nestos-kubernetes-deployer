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
	checkStringValue(&clusterAsset.Platform, opts.Platform)

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
		infraAsset, err := initLibvirtAssetFromMap(libvirtAsset, opts)
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
				"username":     "",
				"remote_ip":    "",
				"osimage_path": "",
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
	openstackAsset := &OpenStackAsset{
		UserName:          opts.InfraPlatform.OpenStack.UserName,
		Password:          opts.InfraPlatform.OpenStack.Password,
		Tenant_Name:       opts.InfraPlatform.OpenStack.Tenant_Name,
		Auth_URL:          opts.InfraPlatform.OpenStack.Auth_URL,
		Region:            opts.InfraPlatform.OpenStack.Region,
		Internal_Network:  opts.InfraPlatform.OpenStack.Internal_Network,
		External_Network:  opts.InfraPlatform.OpenStack.External_Network,
		Glance_Name:       opts.InfraPlatform.OpenStack.Glance_Name,
		Availability_Zone: opts.InfraPlatform.OpenStack.Availability_Zone,
	}

	updateFieldFromMap("username", &openstackAsset.UserName, openstackMap)
	updateFieldFromMap("password", &openstackAsset.Password, openstackMap)
	updateFieldFromMap("tenant_name", &openstackAsset.Tenant_Name, openstackMap)
	updateFieldFromMap("auth_url", &openstackAsset.Auth_URL, openstackMap)
	updateFieldFromMap("region", &openstackAsset.Region, openstackMap)
	updateFieldFromMap("internal_network", &openstackAsset.Internal_Network, openstackMap)
	updateFieldFromMap("external_network", &openstackAsset.External_Network, openstackMap)
	updateFieldFromMap("glance_name", &openstackAsset.Glance_Name, openstackMap)
	updateFieldFromMap("availability_zone", &openstackAsset.Availability_Zone, openstackMap)

	return openstackAsset, nil
}

type LibvirtAsset struct {
	UserName     string
	Remote_IP    string
	OSImage_Path string
}

func initLibvirtAssetFromMap(libvirtMap map[string]interface{}, opts *opts.OptionsList) (InfraAsset, error) {
	libvirtAsset := &LibvirtAsset{
		UserName:     opts.InfraPlatform.Libvirt.UserName,
		Remote_IP:    opts.InfraPlatform.Libvirt.RemoteIP,
		OSImage_Path: opts.InfraPlatform.Libvirt.OSImagePath,
	}

	updateFieldFromMap("username", &libvirtAsset.UserName, libvirtMap)
	updateFieldFromMap("remote_ip", &libvirtAsset.Remote_IP, libvirtMap)
	updateFieldFromMap("osimage_path", &libvirtAsset.OSImage_Path, libvirtMap)

	return libvirtAsset, nil
}

func updateFieldFromMap(fieldName string, fieldValue *string, inputMap map[string]interface{}) {
	if value, ok := inputMap[fieldName]; ok {
		if strValue, ok := value.(string); ok && *fieldValue == "" {
			*fieldValue = strValue
		}
	}
}
