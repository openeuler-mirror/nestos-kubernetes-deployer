# 点火设计
当集群基础设施顺利完成创建后，操作系统随即进入启动流程。在此过程中，点火文件会自动获取所需的配置信息，并将配置文件（如用户信息、集群证书、集群部署服务等）写入到节点机器的磁盘中。随后，systemd服务将启动，确保系统能够按照预设的规则和配置稳定运行。当"release-image-pivot.service"服务正常启动后，将会为集群部署搭建所需的环境。在环境配置完成后，K8S集群部署服务将正式启动，这标志着集群部署工作的正式开始。整个过程将持续进行，直到集群部署完成后服务自动关闭。NKD目前支持多种点火文件包括Ignition、cloud-init和Kickstart文件。

## Ignition
Ignition是不可变基础设施操作系统（如NestOS、Fedora CoreOS）在initramfs期间用来操作磁盘的实用程序，这包括创建用户、添加受信的SSH密钥、写入文件（常规文件、systemd服务...）、网络配置等。首次启动时，Ignition读取其配置并应用该配置。Ignition使用JSON配置文件来表示要进行的更改集。此配置的格式在[规范](https://coreos.github.io/ignition/specs/)中有详细说明。

NKD生成部署集群所需的Ignition（controlplane.ign）文件约为90KB，该文件声明了集群部署过程中所需的关键组件，包括必要的systemd服务和集群证书等。然而，当在OpenStack平台上部署集群时，由于Nova用户数据限制为64KB，因此无法直接创建实例。为了解决这个问题，NKD创建了一个精简版的Ignition（*-merge.ign）文件，作为创建实例时的引导配置文件。同时，NKD将主要的Ignition文件（controlplane.ign）加载到了一个可通过HTTP服务访问的内存存储中。

生成的Ignition文件目录结构如下：
``` shell
$ tree
.
├── controlplane.ign
├── controlplane-merge.ign
├── master.ign
├── master-merge.ign
├── worker.ign
└── worker-merge.ign
```
```shell
较小的"*-merge.ign"Ignition文件的内容如下：
其中:
    - IP：部署HTTP服务的地址，默认为宿主机IP地址
    - port：HTTP服务的访问端口，默认端口为9080
    - {*.ign}：执行集群部署的主要Ignition文件
{
  "ignition": {
    "config": {
      "merge": [
        {
          "source": "http://{IP}:{port}/{*.ign}"
        }
      ]
    },
    "version": "3.2.0"
  }
}
```

## cloud-init
cloud-init是专为云计算环境中虚拟机实例初始化而开发的一款开源工具。在使用NKD进行集群部署过程中，选择底层操作系统为通用操作系统（例如：openeuler），在虚拟化平台上进行集群部署时将生成cloud-init文件，该文件中配置了集群部署所需的环境和集群部署服务，例如：

- 配置主机名
- 在实例上安装软件包
- 配置集群环境
- 运行集群安装脚本

## Kickstart
在使用NKD进行集群部署的过程中，当选择底层操作系统为通用操作系统，并基于PXE（预启动执行环境）平台上进行集群部署时将生成kickstart文件，该文件中配置了集群部署所需的环境和集群部署服务，实现了一种无人值守的安装方式。