# 用户操作手册

## 准备工作 

* 环境要求
  * Linux x86_64/aarch64
  * 安装tofu软件包
    ``` shell
    # 安装amd64版本
    $ wget https://github.com/opentofu/opentofu/releases/download/v1.6.2/tofu_1.6.2_amd64.rpm
    $ rpm -ivh tofu_1.6.2_amd64.rpm
    ``` 
    ``` shell
    # 安装arm64版本
    $ wget https://github.com/opentofu/opentofu/releases/download/v1.6.2/tofu_1.6.2_arm64.rpm
    $ rpm -ivh tofu_1.6.2_arm64.rpm
    ``` 

* 安装NKD
  * 选择拷贝编译好的NKD二进制文件直接使用
  * 根据以下编译安装说明编译安装NKD

备注：
为确保NKD部署的顺利运行，其所在的部署环境需能够与集群节点机器网络正常通信。如果存在防火墙，需正确配置以允许NKD与集群之间的通信，如开放特定的http服务端口。若采用域名进行通信，需确保DNS服务器配置正确，并且NKD所在的环境能够访问DNS服务器。

## 支持平台

### libvirt
在libvirt平台部署集群时，需要提前安装libvirt虚拟化环境

### OpenStack
在OpenStack平台部署集群时，需要提前搭建好OpenStack环境

### 裸金属
在裸金属平台部署集群时，需要提前准备物理机

## 编译安装

* 编译环境：Linux x86_64/aarch64
* 进行编译需要以下软件包：
  * golang >= 1.17
  * git
  ``` shell
  $ sudo yum install golang git
  ```  
* 使用git获取本项目的源码
  ``` shell
  sudo git clone https://gitee.com/openeuler/nestos-kubernetes-deployer
  ```
* 编译二进制
  ``` shell
  $ sh hack/build.sh
  ```

## 配置管理

### 全局配置
全局配置文件用于管理整个集群（多集群）的配置，具体的配置项参数和默认配置详见[全局配置文件说明](./globalconfig_file_desc.md)

#### 点火服务配置参数：
NKD部署集群过程中集群节点需要访问NKD提供的点火服务，通过以下全局配置参数对点火服务进行配置：
* bootstrapIgnHost：点火服务地址（域名或ip，一般为NKD运行环境）
* bootstrapIgnPort：点火服务端口（默认9080，需自行开放防火墙端口）

为适配多网卡环境，点火服务真实监听地址为0.0.0.0。
* 简单网络环境下，部署集群节点可直接访问NKD服务，"bootstrap_ign_host"参数项可以为空，此时NKD会探测路由表默认最高优先级的IP地址作为访问点火服务URL的host；
* 复杂网络环境下，部署集群节点无法直接访问NKD的运行环境，"bootstrap_ign_host"参数项需要配置为对外映射ip或域名，用户需自行配置NAT映射或DNS服务，以确保集群节点可访问到NKD点火服务。

"bootstrap_ign_port"参数当前被点火服务监听端口和访问点火服务URL端口复用，简单网络环境下这两个值保持一致，但复杂网络环境下，需保证NKD服务对外映射端口与本地监听端口保持一致。

### 集群配置
集群配置文件用于对每个集群独立配置，具体的配置项参数和默认配置详见[集群配置文件说明](./config_file_desc.md)

## 基本功能

在“部署集群”章节中有部署集群的具体过程，这里列出了NKD的基本执行指令：
  ``` shell
  # 生成默认配置模板
  $ nkd template -f cluster_config.yaml

  # 应用配置文件部署集群
  $ nkd deploy -f cluster_config.yaml 

  # 销毁指定集群
  $ nkd destroy --cluster-id [your-cluster-id]

  # 扩展指定集群节点数量
  $ nkd extend --cluster-id [your-cluster-id] --num 10

  # 升级指定集群
  # --cluster-id string: 指定要升级的集群的唯一标识符
  # --force: 强制驱逐Pod，这可能导致数据丢失或服务中断，请谨慎使用
  # --imageurl string: 指定用于升级的容器镜像的地址
  # --kube-version string: 选择特定的Kubernetes版本进行升级
  # --kubeconfig string: 指定访问Kubeconfig文件的路径，默认为 "/etc/nkd/[your-cluster-id]/admin.config"
  # --maxunavailable uint: 同时升级的节点的最大数量
  $ nkd upgrade --cluster-id [your-cluster-id] --imageurl [your-image-url] --kube-version [your-k8s-version] 
  ```
除了应用配置文件部署集群外，支持应用配置项参数部署集群
  ``` shell
  $ nkd deploy --help
    --arch string                       部署集群的机器架构（例如，amd64或者arm64）
    --bootstrap-ign-host string         指定点火服务地址（域名或者IP地址）
    --bootstrap-ign-port string         指定点火服务端口（默认：9080）
    --certificateKey string             用于在加入新的Master节点后，从 secret 下载的证书进行解密的密钥。
                                        （证书密钥是一个十六进制编码的字符串，是一个大小为 32 字节的 AES 密钥）
    --clusterID string                  指定集群的唯一标识符                 
    --controller-image-url string       指定Housekeeper控制器组件的容器镜像地址
    --deploy-housekeeper                是否部署Housekeeper Operator，默认false
    -f, --file string                   指定集群部署配置文件的位置
    --image-registry string             指定用于拉取Kubernetes组件容器镜像的地址
    --ipxe-filePath string              ipxe配置文件路径
    --ipxe-osInstallTreePath string     ipxe所需操作系统安装树路径 (默认: /var/www/html/)
    --kubernetes-apiversion uint        指定Kubernetes API版本。可接受的参考数值为：
                                        - 1 用于Kubernetes版本 < v1.15.0;
                                        - 2 用于Kubernetes版本 >= v1.15.0 && < v1.22.0;
                                        - 3 用于Kubernetes版本 >= v1.22.0;
    --kubeversion string                指定要部署的Kubernetes版本
    --libvirt-cidr string               用于libvirt平台的CIDR (默认: 192.168.132.0/24)
    --libvirt-gateway string            用于libvirt平台的网关 (默认: 192.168.132.1)
    --libvirt-osPath string             libvirt 平台下的操作系统路径
    --libvirt-uri string                用于libvirt的URI (默认: qemu:///system)
    --master-cpu uint                   设置主节点的CPU（单位：核）
    --master-disk uint                  设置主节点磁盘大小（单位：GB）
    --master-hostname stringArray       设置主节点主机名
    --master-ips stringArray            设置主节点IP地址
    --master-ram uint                   设置主节点的RAM（单位：MB）
    --network-plugin-url                部署网络插件yaml的URL
    --openstack-authURL string          OpenStack的鉴权地址 (默认: http://controller:5000/v3)
    --openstack-availabilityZone string OpenStack的可用域 (默认: nova)
    --openstack-externalNetwork string  OpenStack的外部网络
    --openstack-glanceName string       OpenStack的镜像名称
    --openstack-internalNetwork string  OpenStack的内部网络
    --openstack-password string         OpenStack的密码
    --openstack-region string           OpenStack的地区(默认: RegionOne)
    --openstack-tenantName string       OpenStack的租户名称(默认: admin)
    --openstack-username string         OpenStack的用户名(默认: admin)
    --operator-image-url string         指定Housekeeper Operator组件的容器镜像地址
    --os-type string                    指定集群节点的操作系统类型（例如：nestos、openeuler）
    --password string                   指定 ssh 登录所配置节点的密码
    --pause-image string                指定pause容器的镜像
    --platform string                   选择用于部署集群的基础设施平台（支持libvirt或者openstack平台）
    --pod-subnet string                 指定Kubernetes Pod的子网（默认：10.244.0.0/16）
    --posthook-yaml string              指定一个 YAML 文件或目录，在集群部署后使用 'kubectl apply' 应用
    --prehook-script string             指定一个脚本文件或目录，在集群部署前执行
    --pxe-httpRootDir string            PXE平台下 HTTP 服务器的根目录 (默认: /var/www/html/)
    --pxe-ip string                     PXE本地服务器的IP地址
    --pxe-tftpRootDir string            PXE平台下TFTP服务器的根目录 (默认: /var/lib/tftpboot/)
    --release-image-url string          指定包含Kubernetes组件的NestOS容器镜像的URL，仅支持qcow2格式
    --runtime string                    指定容器运行时类型（docker、isulad 或 crio）
    --service-subnet string             指定Kubernetes服务的子网（默认："10.96.0.0/16"）
    --sshkey string                     ssh 免密登录的密钥存储文件的路径（默认：~/.ssh/id_rsa.pub）
    --token string                      用于验证从控制平面获取的集群信息，非控制平面节点用于加入集群
    --username string                   需要部署 k8s 集群的机器的 ssh 登录用户名
    --worker-cpu uint                   设置工作节点的CPU（单位：核心）
    --worker-disk uint                  设置工作节点磁盘大小（单位：GB）
    --worker-hostname stringArray       设置工作节点主机名  
    --worker-ips stringArray            设置工作节点IP地址
    --worker-ram uint                   设置工作节点的RAM（单位：MB）
  全局参数：
    --dir string         文件生成目录 (默认 "/etc/nkd")
    --log-level string   日志级别 (例如 "debug | info | warn | error") (默认 "info")

  # 应用可选配置项参数部署集群
  $ nkd deploy --platform [platform] --master-ips [master-ip-01] --master-ips [master-ip-02] --master-hostname [master-hostname-01] --master-hostname [master-hostname-02] --master-cpu [master-cpu-cores] --worker-hostname [worker-hostname-01] --worker-disk [worker-disk-size] ...
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
    * 基础镜像需从NestOS官网下载最新版本容器镜像，官网镜像未包含kubernetes相关二进制组件
    * 制作部署镜像，需提前下载相对应版本的kubeadm、kubelet、crictl二进制文件并拷贝到/usr/bin目录。
    * 软件包的安装需要使用rpm-ostree命令。
 * Dockerfiles示例如下
      ``` dockerfile
      FROM nestos_base_image
      COPY kube* /usr/bin/
      COPY crictl /usr/bin/
      RUN ostree container commit
      ```
备注：如果集群底层操作系统选择NestOS，用户在部署集群前需要自定义构建部署镜像。

## 部署集群

 - 添加可选参数项部署集群，命令示例：
    ``` shell
    $ nkd deploy --master-ips 192.168.132.11 --master-ips 192.168.132.12 --master-hostname k8s-master01 --master-hostname k8s-master02 --master-cpu 8 --worker-hostname k8s-worker01 --worker-disk 50 ...
    ```
 - 此外更精细化的配置，可以通过集群配置文件部署集群，详情见配置管理。
    ``` shell
    $ nkd deploy -f cluster_config.yaml
    ```
