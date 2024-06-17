# nestos-kubernetes-deployer
![ignition_design_2](/docs/logo/nkd-logo.png)

## Introduction

NKD (NestOS Kubernetes Deployer) is a cluster deployment and operation tool developed by the NestOS team for container cloud scenarios. It covers a series of functions such as infrastructure and Kubernetes core component deployment, update, and configuration management, providing users with a one-stop solution. It supports customization of multiple container runtimes including crio, iSulad, docker, and containerd, and is compatible with multiple platform deployments, ensuring that users can easily handle various complex deployment requirements. In addition, NKD has the ability to create cluster self-signed certificates and supports deploying multiple versions of Kubernetes clusters, covering various scenarios that may be encountered in actual use.

## Software architecture

For more information, see [Software Architecture](docs/en/overall_design.md)

## Install and deploy

Please refer to the [user manual for details](docs/en/manual.md)

## Planning for the future

NKD's ultimate vision is to support O&M and efficient management of  multiple clusters at the same time in the form of long-term resident  services. It provides a variety of features, such as persistent  configuration change records, certificate management, multiple update  and upgrade policies, and image source channels. In the future, we will  continue to improve the functions and performance of NKD, and introduce  more intelligent features, such as automated fault handling and resource optimization. Our goal is to shape NKD into a core component of the  NestOS ecosystem, providing all-round support for O&M in  cloud-native scenarios, so as to further promote the development and  widespread application of cloud-native technologies. 

## Get involved

Interested partners are very welcome to join us and contribute. 

## License

Apache License 2.0
