# Infra

NKD采用开源的基础设施即代码（Infrastructure as Code，IaC）工具OpenTofu，连接基础设施提供商，为用户提供自动化的基础设施操作。目前，NKD已在Libvirt和OpenStack平台实现对NestOS的支持，包括但不限于NestOS实例的创建、网络配置、存储管理等功能。

## 创建基础设施

用户可使用NKD支持的各种方式，例如命令行、配置文件等，连接基础设施提供商，创建所需配置的NestOS实例。可配参数如下：

| 参数     | 描述                          |
| -------- | ----------------------------- |
| OSImage  | NestOS镜像的下载地址/本地存储路径 |
| Hostname | NestOS实例的主机名称              |
| IP       | NestOS实例的IP地址                |
| CPU      | NestOS实例的CPU规格               |
| RAM      | NestOS实例的内存规格              |
| Disk     | NestOS实例的硬盘规格              |
| Ign_Path | NestOS实例启动所需的Ignition配置文件路径  |

注：需根据不同的基础设施提供商，提供相应的鉴权信息。