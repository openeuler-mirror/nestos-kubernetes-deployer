# Ignition

Ignition是NestOS在initramfs期间用来操作磁盘的实用程序。其中包括创建用户、添加受信的SSH密钥、写入文件（常规文件、systemd服务...）、网络配置等。首次启动时，Ignition读取其配置并应用该配置。Ignition使用JSON配置文件来表示要进行的更改集。此配置的格式在[规范](https://coreos.github.io/ignition/specs/)中有详细说明。

## 提供配置
当集群节点的基础设施正常创建后，操作系统处于引导阶段，Ignition将获取配置信息并将配置文件（用户、集群证书、集群部署服务等）写入到节点机器的磁盘中，然后启动systemd服务。当"release-image-pivot.service"服务正常启动后，节点机器通过rpm-ostree机制将文件系统切换为基于NestOS的kubernetes定制版本，之后启动K8S集群部署服务，开始正式的集群部署任务，直到集群部署完成后服务关闭。

### ControlPlane Node
controlplane节点的Ignition文件配置信息如图:
![ignition_design_1](/docs/figures/ignition_design_1.jpg)

### Master Node
Master节点的Ignition文件配置信息如图：
![ignition_design_2](/docs/figures/ignition_design_2.jpg)

### Worker Node
Worker节点的Ignition文件配置信息如图：
![ignition_design_3](/docs/figures/ignition_design_3.jpg)
