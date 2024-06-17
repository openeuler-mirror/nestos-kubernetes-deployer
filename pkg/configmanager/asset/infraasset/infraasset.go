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

package infraasset

import (
	"errors"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"runtime"
	"strings"
)

type InfraAsset interface {
	InitAsset(assetMap map[string]interface{}, opts *opts.OptionsList, args ...interface{}) (InfraAsset, error)
}

func InitInfraAsset(clusterAsset *asset.ClusterAsset, opts *opts.OptionsList) (InfraAsset, error) {
	asset.SetStringValue(&clusterAsset.Architecture, opts.Arch, runtime.GOARCH)
	asset.SetStringValue(&clusterAsset.Platform, opts.Platform, "libvirt")

	switch strings.ToLower(clusterAsset.Platform) {
	case "openstack":
		assetMap, ok := convertMap(clusterAsset.InfraPlatform, "openstack")
		if !ok {
			return nil, errors.New("failed to get openstack asset")
		}

		openstackAsset := &OpenStackAsset{}
		infraAsset, err := openstackAsset.InitAsset(assetMap, opts)
		if err != nil {
			return nil, err
		}
		return infraAsset, nil
	case "libvirt":
		assetMap, ok := convertMap(clusterAsset.InfraPlatform, "libvirt")
		if !ok {
			return nil, errors.New("failed to get libvirt asset")
		}

		libvirtAsset := &LibvirtAsset{}
		infraAsset, err := libvirtAsset.InitAsset(assetMap, opts, clusterAsset.Architecture)
		if err != nil {
			return nil, err
		}
		return infraAsset, nil
	case "pxe":
		assetMap, ok := convertMap(clusterAsset.InfraPlatform, "pxe")
		if !ok {
			return nil, errors.New("failed to get pxe asset")
		}

		pxeAsset := &PXEAsset{}
		infraAsset, err := pxeAsset.InitAsset(assetMap, opts)
		if err != nil {
			return nil, err
		}
		return infraAsset, nil
	case "ipxe":
		assetMap, ok := convertMap(clusterAsset.InfraPlatform, "ipxe")
		if !ok {
			return nil, errors.New("failed to get ipxe asset")
		}

		ipxeAsset := &IPXEAsset{}
		infraAsset, err := ipxeAsset.InitAsset(assetMap, opts)
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
		switch strings.ToLower(platform) {
		case "libvirt":
			return map[string]interface{}{
				"uri":     "",
				"osImage": "",
				"cidr":    "",
				"gateway": "",
			}, true
		case "openstack":
			return map[string]interface{}{
				"username":         "",
				"password":         "",
				"tenantName":       "",
				"authURL":          "",
				"region":           "",
				"internalNetwork":  "",
				"externalNetwork":  "",
				"glanceName":       "",
				"availabilityZone": "",
			}, true
		case "pxe":
			return map[string]interface{}{
				"httpServerPort": "",
				"httpRootDir":    "",
				"tftpServerIP":   "",
				"tftpServerPort": "",
				"tftpRootDir":    "",
			}, true
		case "ipxe":
			return map[string]interface{}{
				"ipxePort":              "",
				"ipxeFilePath":          "",
				"ipxeOSInstallTreePath": "",
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

func updateFieldFromMap(fieldName string, fieldValue *string, inputMap map[string]interface{}) {
	if value, ok := inputMap[fieldName]; ok {
		if strValue, ok := value.(string); ok && *fieldValue == "" {
			*fieldValue = strValue
		}
	}
}
