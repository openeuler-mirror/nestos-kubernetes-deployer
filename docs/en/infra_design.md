# Infra

NKD uses OpenTofu, an open source Infrastructure as Code (IaC) tool, to connect infrastructure providers and provide users with automated infrastructure operations. Now, NKD has supported for NestOS on Libvirt and OpenStack platforms, including but not limited to create NestOS instances, config network, manage storage and so on.

## Create infrastructure

Users can use various ways supported by NKD, such as the command line, configuration files, etc., to connect to the infrastructure providers and create NestOS instance. Parameters are as follows:

| Parameter     | Description                          |
| -------- | ----------------------------- |
| OSImage  | Download address / local path of the NestOS image |
| Hostname | Hostname of the NestOS instance              |
| IP       | IP of the NestOS instance               |
| CPU      | CPU of the NestOS instance               |
| RAM      | Memory of the NestOS instance              |
| Disk     | Disk of the NestOS instance              |
| Ign_Path | Ignition configuration file path of the NestOS instance  |

Note: The authentication information depends on the different infrastructure providers.