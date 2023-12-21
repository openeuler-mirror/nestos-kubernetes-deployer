# 用户操作手册

## 准备工作 

* 环境要求
  * Linux x86_64/aarch64
  * 安装tofu软件包
  ``` shell
  # 安装amd64版本
  wget https://github.com/opentofu/opentofu/releases/download/v1.6.0-rc1/tofu_1.6.0-rc1_amd64.rpm
  rpm -ivh tofu_1.6.0-rc1_amd64.rpm
  # 安装
  wget https://github.com/opentofu/opentofu/releases/download/v1.6.0-rc1/tofu_1.6.0-rc1_arm.rpm
  rpm -ivh tofu_1.6.0-rc1_arm.rpm
  ``` 
  * 选择openstack平台部署集群，需要提前搭建好openstack环境
  * 选择libvirt平台部署集群，需要提前安装libvirt虚拟化环境

* 安装NKD
  * 选择拷贝编译好的NKD二进制文件直接使用
  * 根据以下编译安装说明编译安装NKD

## 编译安装

* 编译环境：Linux x86_64/aarch64
* 进行编译需要以下软件包：
  * golang >= 1.17
  * git
  ``` shell
  sudo yum install golang git
  ```  
* 使用git获取本项目的源码
  ``` shell
  sudo git clone https://gitee.com/openeuler/nestos-kubernetes-deployer
  ```
* 编译二进制
  ``` shell
  cd nestos-kubernetes-deployer
  sudo go build -mod=vendor -tags release --ldflags="-w -s" -o nkd nkd.go
  ```

## 基本功能

在“部署集群”章节中有部署集群的具体过程，这里列出了NKD的基本执行指令：
  ``` shell
  # 应用配置文件部署集群
  $ nkd deploy -f cluster_config.yaml

  # 销毁指定集群
  $ nkd destroy --cluster-id < your-cluster-id >

  # 扩展指定集群节点数量
  $ nkd extend --cluster-id < your-cluster-id > --num 10

  # 升级指定集群
  # --cluster-id string: 指定要升级的集群的唯一标识符
  # --force: 强制驱逐Pod，这可能导致数据丢失或服务中断，请谨慎使用
  # --imageurl string: 指定用于升级的容器镜像的地址
  # --kube-version string: 选择特定的Kubernetes版本进行升级
  # --kubeconfig string: 指定访问Kubeconfig文件的路径，默认为 "/etc/nkd/pki/kubeconfig/admin.conf"
  # --maxunavailable uint: 同时升级的节点的最大数量
  $ nkd upgrade --cluster-id < your-cluster-id > --imageurl < your-image-url > --kube-version < your-k8s-version > 
  ```
除了应用配置文件部署集群外，支持应用配置项参数部署集群
  ``` shell
  $ nkd deploy --help
    --CertificateKey string         指定要添加到主节点的证书密钥
    --cluster-id string             指定集群的唯一标识符
    --controller-image-url string   指定Housekeeper控制器组件的容器镜像地址
    --deploy-housekeeper            是否部署Housekeeper Operator，默认false
    -f, --file string               指定集群部署配置文件的位置
    --image-registry string         指定用于拉取Kubernetes组件容器镜像的地址
    --image-version string          指定CoreDNS容器镜像的版本
    --kubeversion string            指定要部署的Kubernetes版本
    --master-cpu uint               设置主节点的CPU（单位：核）
    --master-disk uint              设置主节点磁盘大小（单位：GB）
    --master-hostname stringArray   设置主节点主机名，默认为：[k8s-master01, k8s-master02, k8s-master03]
    --master-igns stringArray       设置主节点的Ignition文件路径
    --master-ips stringArray        设置主节点IP地址
    --master-password string        登录主节点的密码
    --master-ram uint               设置主节点的RAM（单位：GB）
    --master-sshkey string          设置主节点的SSH私钥，用于身份验证
    --master-username string        登录主节点的用户名，默认为 "root"
    --operator-image-url string     指定Housekeeper Operator组件的容器镜像地址
    --pause-image string            指定pause容器的镜像
    --platform string               选择用于部署集群的基础设施平台
    --pod-subnet string             指定Kubernetes Pod的子网
    --release-image-url string      指定包含Kubernetes组件的NestOS容器镜像的URL，仅支持qcow2格式
    --service-subnet string         指定Kubernetes服务的子网，默认为 "10.96.0.0/16"
    --token string                  指定用于访问资源的身份验证令牌
    --worker-cpu uint               设置工作节点的CPU（单位：核心）
    --worker-disk uint              设置工作节点磁盘大小（单位：GB）
    --worker-hostname stringArray   设置工作节点主机名，默认为：[k8s-worker01, k8s-worker02, k8s-worker03]
    --worker-igns stringArray       设置工作节点的Ignition文件路径
    --worker-ips stringArray        设置工作节点IP地址
    --worker-password string        登录工作节点的密码
    --worker-ram uint               设置工作节点的RAM（单位：GB）
    --worker-sshkey string          设置工作节点的SSH私钥，用于身份验证
    --worker-username string        登录工作节点的用户名，默认为 "root"
  ```

## 部署过程展示

调整集群部署配置文件
![](./figures/cluster_config.mp4)

应用配置文件部署集群
![](./figures/cluster_deploy.mp4)

## 镜像构建

* NestOS容器镜像支持利用Dockerfile在原来的基础上构建新的容器镜像
* 制作注意事项
    * 请确保已安装docker。
    * 基础镜像需从NestOS官网下载最新版本容器镜像。
    * 制作部署镜像，需提前下载相对应版本的kubeadm、kubelet、crictl二进制文件并复制到/usr/bin目录，以及将calico网络插件的yaml文件复制到/etc/nkd目录。
    * 软件包的安装需要使用rpm-ostree命令。
 * Dockerfiles示例如下
      ``` dockerfile
      FROM nestos_base_image
      COPY kube* /usr/bin/
      COPY calico.yaml /etc/nkd/
      RUN ostree container commit
      ```
备注：部署集群前用户需要自定义构建部署镜像

## 部署集群

1. 准备安装配置文件，具体的cluster_config.yaml集群配置文件见[配置文件说明](/docs/config_file_desc.md)
2. 执行命令部署集群
    ```
    nkd deploy -f cluster_config.yaml
    ```


