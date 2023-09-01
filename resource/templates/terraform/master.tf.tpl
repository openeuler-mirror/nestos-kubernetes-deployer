terraform {
  required_providers {
    openstack = {
      source = "terraform-provider-openstack/openstack"
      version = "1.52.1"
    }
  }
}

provider "openstack" {
  user_name   = "{{.Infra.Openstack.User_name}}"
  password    = "{{.Infra.Openstack.Password}}"
  tenant_name = "{{.Infra.Openstack.Tenant_name}}"
  auth_url    = "{{.Infra.Openstack.Auth_url}}"
  region      = "{{.Infra.Openstack.Region}}"
}

variable "instance_count" {
  default = "3"
}

variable "instance_name" {
  default = "{{.System.HostName}}"
}

resource "openstack_compute_flavor_v2" "flavor" {
  name      = var.instance_names[count.index]
  ram       = "{{.Infra.Vmsize.Ram}}"
  vcpus     = "{{.Infra.Vmsize.Vcpus}}"
  disk      = "{{.Infra.Vmsize.Disk}}"
  is_public = "true"
}

resource "openstack_compute_secgroup_v2" "secgroup" {
  name        = "k8s_master_secgroup"
  description = "secgroup for k8s master"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = -1
    to_port     = -1
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }
}

resource "openstack_compute_instance_v2" "instance" {
  count           = var.instance_count
  name            = format("${var.instance_name}%02d", count.index + 1)
  image_name      = "{{.Infra.Openstack.Glance}}"
  flavor_name     = openstack_compute_flavor_v2.flavor.name
  security_groups = [openstack_compute_secgroup_v2.secgroup.name]
  user_data       = file("/etc/nkd/instance.ign")

  network {
    name        = "{{.Infra.Openstack.Internal_network}}"
    fixed_ip_v4 = element({{.System.Ips}}, count.index)
  }
}

resource "openstack_networking_floatingip_v2" "floatip" {
  count = length(openstack_compute_instance_v2.instance)
  pool  = "{{.Infra.Openstack.External_network}}"
}

resource "openstack_compute_floatingip_associate_v2" "fip_associate" {
  count       = length(openstack_compute_instance_v2.instance)
  floating_ip = openstack_networking_floatingip_v2.floatip.*.address[count.index]
  instance_id = openstack_compute_instance_v2.instance.*.id[count.index]
}

output "instance_info" {
  value = {
    instance_status = openstack_compute_instance_v2.instance.*.power_state
    floating_ip     = openstack_networking_floatingip_v2.floatip.*.address
  }
}