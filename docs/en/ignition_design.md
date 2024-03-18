# Ignition

Ignition是NestOS在initramfs期间用来操作磁盘的实用程序。其中包括创建用户、添加受信的SSH密钥、写入文件（常规文件、systemd服务...）、网络配置等。首次启动时，Ignition读取其配置并应用该配置。Ignition使用JSON配置文件来表示要进行的更改集。此配置的格式在[规范](https://coreos.github.io/ignition/specs/)中有详细说明。

## 提供配置
生成部署集群所需的Ignition（controlplane.ign）文件约为90KB，文件中声明了集群部署所需的systemd服务、证书等。在OpenStack平台上部署集群时，由于Nova用户数据限制为64KB，因此无法直接部署基础设施。为了解决这个问题，NKD创建了一个较小的Ignition（*-merge.ign）文件，作为创建基础设施时的引导配置文件。而主要的Ignition文件将会加载到可以通过 HTTP 服务访问的内存中，使较小的Ignition文件在系统引导阶自动加载主要的Ignition文件。

当集群节点的基础设施正常创建后，操作系统处于引导阶段，Ignition将获取配置信息并将配置文件（用户、集群证书、集群部署服务等）写入到节点机器的磁盘中，然后启动systemd服务。当"release-image-pivot.service"服务正常启动后，节点机器通过rpm-ostree机制将文件系统切换为基于NestOS的kubernetes定制版本，之后启动K8S集群部署服务，开始正式的集群部署任务，直到集群部署完成后服务关闭。

### ControlPlane Node
controlplane节点的Ignition文件配置信息如图:
![ignition_design_1](/docs/en/figures/ignition_design_1.jpg)

### Master Node
Master节点的Ignition文件配置信息如图：
![ignition_design_2](/docs/en/figures/ignition_design_2.jpg)

### Worker Node
Worker节点的Ignition文件配置信息如图：
![ignition_design_3](/docs/en/figures/ignition_design_3.jpg)

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