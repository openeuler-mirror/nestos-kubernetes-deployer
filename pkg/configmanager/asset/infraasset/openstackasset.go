/*
Copyright 2024 KylinSoft  Co., Ltd.

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
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
)

type OpenStackAsset struct {
	UserName         string `yaml:"username"`
	Password         string `yaml:"password"`
	TenantName       string `yaml:"tenantName"`
	AuthURL          string `yaml:"authURL"`
	Region           string `yaml:"region"`
	InternalNetwork  string `yaml:"internalNetwork"`
	ExternalNetwork  string `yaml:"externalNetwork"`
	GlanceName       string `yaml:"glanceName"`
	AvailabilityZone string `yaml:"availabilityZone"`
}

func (oa *OpenStackAsset) InitAsset(openstackMap map[string]interface{}, opts *opts.OptionsList, args ...interface{}) (InfraAsset, error) {
	updateFieldFromMap("username", &oa.UserName, openstackMap)
	asset.SetStringValue(&oa.UserName, opts.InfraPlatform.OpenStack.UserName, "admin")

	updateFieldFromMap("password", &oa.Password, openstackMap)
	if err := asset.CheckStringValue(&oa.Password, opts.InfraPlatform.OpenStack.Password, "openstack-password"); err != nil {
		return nil, err
	}

	updateFieldFromMap("tenantName", &oa.TenantName, openstackMap)
	asset.SetStringValue(&oa.TenantName, opts.InfraPlatform.OpenStack.TenantName, "admin")

	updateFieldFromMap("authURL", &oa.AuthURL, openstackMap)
	asset.SetStringValue(&oa.AuthURL, opts.InfraPlatform.OpenStack.AuthURL, "http://controller:5000/v3")

	updateFieldFromMap("region", &oa.Region, openstackMap)
	asset.SetStringValue(&oa.Region, opts.InfraPlatform.OpenStack.Region, "RegionOne")

	updateFieldFromMap("internalNetwork", &oa.InternalNetwork, openstackMap)
	if err := asset.CheckStringValue(&oa.InternalNetwork, opts.InfraPlatform.OpenStack.InternalNetwork, "openstack-internalNetwork"); err != nil {
		return nil, err
	}

	updateFieldFromMap("externalNetwork", &oa.ExternalNetwork, openstackMap)
	if err := asset.CheckStringValue(&oa.ExternalNetwork, opts.InfraPlatform.OpenStack.ExternalNetwork, "openstack-externalNetwork"); err != nil {
		return nil, err
	}

	updateFieldFromMap("glanceName", &oa.GlanceName, openstackMap)
	if err := asset.CheckStringValue(&oa.GlanceName, opts.InfraPlatform.OpenStack.GlanceName, "openstack-glanceName"); err != nil {
		return nil, err
	}

	updateFieldFromMap("availabilityZone", &oa.AvailabilityZone, openstackMap)
	asset.SetStringValue(&oa.AvailabilityZone, opts.InfraPlatform.OpenStack.AvailabilityZone, "nova")

	return oa, nil
}
