# nestos-kubernetes-deployer
![ignition_design_2](/docs/logo/nkd-logo.png)

## Introduction

NKD (NestOS Kubernetes Deployer) is a solution designed for deploying and  maintaining Kubernetes clusters on NestOS. It is designed to simplify  the process of deploying and upgrading clusters by providing a range of  services outside the cluster, including deployment, updates, and  configuration management of infrastructure and core components of  Kubernetes. NKD is designed to provide a more convenient cluster  operation experience, allowing users to easily complete complex  management tasks, thereby improving the overall efficiency of deployment and maintenance. 

#### Support Platforms

It supports multiple deployment platforms, and NKD dynamically creates the required IaaS resources by connecting infrastructure providers  according to the needs of the cluster, and currently supports OpenStack  and libvirt platforms. 

## Software architecture

For more information, see [Software Architecture](docs/overall_design.md)

## Install and deploy

Please refer to the [user manual for details](docs/manual.md)

## Planning for the future

NKD's ultimate vision is to support O&M and efficient management of  multiple clusters at the same time in the form of long-term resident  services. It provides a variety of features, such as persistent  configuration change records, certificate management, multiple update  and upgrade policies, and image source channels. In the future, we will  continue to improve the functions and performance of NKD, and introduce  more intelligent features, such as automated fault handling and resource optimization. Our goal is to shape NKD into a core component of the  NestOS ecosystem, providing all-round support for O&M in  cloud-native scenarios, so as to further promote the development and  widespread application of cloud-native technologies. 

## Get involved

Interested partners are very welcome to join us and contribute. 

## License

Apache License 2.0
