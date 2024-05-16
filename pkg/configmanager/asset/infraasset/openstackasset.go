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

func (oa *OpenStackAsset) InitAsset(openstackMap map[string]interface{}, opts *opts.OptionsList, args ...interface{}) (InfraAsset, error) {
	updateFieldFromMap("username", &oa.UserName, openstackMap)
	if err := asset.CheckStringValue(&oa.UserName, opts.InfraPlatform.OpenStack.UserName, "openstack-username"); err != nil {
		return nil, err
	}

	updateFieldFromMap("password", &oa.Password, openstackMap)
	if err := asset.CheckStringValue(&oa.Password, opts.InfraPlatform.OpenStack.Password, "openstack-password"); err != nil {
		return nil, err
	}

	updateFieldFromMap("tenant_name", &oa.Tenant_Name, openstackMap)
	if err := asset.CheckStringValue(&oa.Tenant_Name, opts.InfraPlatform.OpenStack.Tenant_Name, "openstack-tenant-name"); err != nil {
		return nil, err
	}

	updateFieldFromMap("auth_url", &oa.Auth_URL, openstackMap)
	if err := asset.CheckStringValue(&oa.Auth_URL, opts.InfraPlatform.OpenStack.Auth_URL, "openstack-auth-url"); err != nil {
		return nil, err
	}

	updateFieldFromMap("region", &oa.Region, openstackMap)
	if err := asset.CheckStringValue(&oa.Region, opts.InfraPlatform.OpenStack.Region, "openstack-region"); err != nil {
		return nil, err
	}

	updateFieldFromMap("internal_network", &oa.Internal_Network, openstackMap)
	if err := asset.CheckStringValue(&oa.Internal_Network, opts.InfraPlatform.OpenStack.Internal_Network, "openstack-internal-network"); err != nil {
		return nil, err
	}

	updateFieldFromMap("external_network", &oa.External_Network, openstackMap)
	if err := asset.CheckStringValue(&oa.External_Network, opts.InfraPlatform.OpenStack.External_Network, "openstack-external-network"); err != nil {
		return nil, err
	}

	updateFieldFromMap("glance_name", &oa.Glance_Name, openstackMap)
	if err := asset.CheckStringValue(&oa.Glance_Name, opts.InfraPlatform.OpenStack.Glance_Name, "openstack-glance-name"); err != nil {
		return nil, err
	}

	updateFieldFromMap("availability_zone", &oa.Availability_Zone, openstackMap)
	if err := asset.CheckStringValue(&oa.Availability_Zone, opts.InfraPlatform.OpenStack.Availability_Zone, "openstack-availability-zone"); err != nil {
		return nil, err
	}

	return oa, nil
}
