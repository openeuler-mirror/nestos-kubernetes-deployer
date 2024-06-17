# 集群配置文件说明

``` shell
clusterID: cluster                                  # 集群名称
architecture: amd64                                 # 部署集群的机器架构,支持amd64或者arm64
platform: libvirt                                   # 部署平台为libvirt、openstack、pxe
infraPlatform                                       # 指定基础设施平台类型
                                                    # 需要根据不同的部署平台设置参数
osImage:
  type:                                             # 指定操作系统类型，例如nestos、generalos
username: root                                      # 指定 ssh 登录所配置节点的用户名
password:                                           # 指定 ssh 登录所配置节点的密码
sshKey: "/root/.ssh/id_rsa.pub"                     # ssh 免密登录的密钥存储文件的路径
master:                                             # 配置master节点的列表
- hostname: k8s-master01                            # 该节点的名称
  hardwareinfo:                                     # 该节点配置的硬件资源信息
    cpu: 4                                          # 该节点CPU的核数
    ram: 8192                                       # 该节点的内存大小
    disk: 50                                        # 该节点的磁盘大小
  ip: "192.168.132.11"                              # 该节点的IP地址
worker:                                             # 配置worker节点的列表
- hostname: k8s-worker01            
  hardwareinfo:
    cpu: 4
    ram: 8192
    disk: 50
  ip: ""                                           # 如果不设置worker节点IP地址，则由dhcp自动分配，默认为空
runtime: isulad                                    # 指定容器运行时类型，目前支持 docker、isulad、containerd和crio
kubernetes:                                        # 集群相关配置列表
  kubernetesVersion: "v1.29.1"                     # 部署集群的版本
  kubernetesApiversion: "v1beta3"                  # 指定kubeadm配置文件格式的版本，目前支持 v1beta3、v1beta2、v1beta1
  apiserverEndpoint: "192.168.132.11:6443"         # 对外暴露的APISERVER服务的地址或域名   
  imageRegistry: "registry.k8s.io"                 # Kubeadm初始化时使用的镜像仓库地址
  registryMirror: ""                               # 下载容器镜像时，使用的镜像仓库的 mirror 站点地址
  pauseImage: "pause:3.9"                          # 容器运行时的pause容器的容器镜像名称
  releaseImageUrl: ""                              # 包含K8S二进制组件的NestOS发布镜像的地址，支持架构x86_64或者aarch64
  token: ""                                        # 启动引导过程中使用的令牌，默认自动生成
  adminKubeconfig: /etc/nkd/cluster/admin.config   # 集群管理员配置文件admin.conf的路径
  certificateKey: ""                               # 添加新的控制面节点时用来解密所下载的Secret中的证书的秘钥
  packageList:                                     # 集群环境中需要安装的RPM软件包名称列表
  rpmPackagePath: ""                               # 集群环境中需要安装的RPM软件包文件路径
  network:                                         # k8s集群网络配置
    serviceSubnet: "10.96.0.0/16"                  # k8s创建的service的IP地址网段
    podSubnet: "10.244.0.0/16"                     # k8s集群网络的IP地址网段
    plugin: ""                                     # 网络插件
housekeeper:                                                                                          # housekeeper相关配置列表
  deployHousekeeper: false                                                                            # 是否部署housekeeper
  operatorImageURL: "hub.oepkgs.net/nestos/housekeeper/{arch}/housekeeper-operator-manager:{tag}"     # housekeeper-operator镜像的地址，支持架构amd64或者arm64
  controllerImageURL: "hub.oepkgs.net/nestos/housekeeper/{arch}/housekeeper-controller-manager:{tag}" # housekeeper-controller镜像的地址，支持架构amd64或者arm64   
certAsset:                                          # 配置外部证书文件路径列表，默认自动生成
  rootCACertPath: ""
  rootCAKeyPath: ""
  etcdCACertPath: ""
  etcdCAKeyPath: ""
  frontProxyCACertPath: ""
  frontProxyCAKeyPath: ""
  saPub: ""
  saKey: ""
```

指定部署平台为libvirt配置参数示例：
``` shell
platform: libvirt                                   # 部署平台为libvirt
infraPlatform
  uri: qemu:///system                                
  osPath:                                           # 指定部署集群机器的操作系统镜像地址，支持架构x86_64或者aarch64
  cidr: 192.168.132.0/24                            # 路由地址
  gateway: 192.168.132.1                            # 网关地址
```

指定部署平台为openstack配置参数示例：
``` shell
platform: openstack                                 # 部署平台为openstack
infraPlatform                      
	username:                                         # openstack用户名，需要有创建资源权限                                       
	password:                                         # openstack登录密码，用于登录openstack平台
	tenantName:                                       # openstack租户名，用户所属的合集，例如：admin
	authURL:                                          # openstack鉴权地址，例如：http://{ip}:{port}/v3
	region:                                           # openstack地区，用于资源隔离，例如：RegionOne
	internalNetwork:                                  # openstack内部网络名称，用户自定义内部网络名称
	externalNetwork:                                  # openstack外部网络名称，用户自定义外部网络名称
	glanceName:                                       # 创建openstack实例的qcow2镜像
	availabilityZone:                                 # 可用域，默认nova
```

指定部署平台为pxe时配置参数示例：
``` shell
platform: pxe                                        # 部署平台为pxe
infraPlatform
  ip:                                                # http服务器的ip地址
  httpServerPort: "9080"                             # http服务器的端口号
  httpRootDir: /var/www/html/                        # 设置 HTTP 服务器的根目录
  tftpServerPort: "69"                               # TFTP服务器端口号
  tftpRootDir: /var/lib/tftpboot/                    # TFTP服务器的根目录
```

## 镜像下载地址

- NestOS镜像下载地址见[官网](https://nestos.openeuler.org/)，需下载NestOS For Container版本
- Openeuler镜像下载地址见[官网](https://www.openeuler.org/)

## 密码密文生成方式：

- 指定集群底层操作系统为nestos时需使用密文密码，其生成方式：
  ``` shell
  openssl passwd -1 -salt yoursalt
  Password: qwer1234!@#$
  $1$yoursalt$UGhjCXAJKpWWpeN8xsF.c/
  ```

- 部署平台为pxe时需使用密文密码，其生成方式：
  ``` shell
  # python3  
  Python 3.7.9 (default, Mar  2 2021, 02:43:11)
  [GCC 7.3.0] on linux
  Type "help", "copyright", "credits" or "license" for more information.  
  >>> import crypt  
  >>> passwd = crypt.crypt("myPasswd")  
  >>> print (passwd)  
  $6$sH1qri2n14V1VCv/$fWnV3rPv95gWHJ3wZu6o0bBGy.SnllSw4a2HuoP45jXfI9fCrwe60AULO/0aXS7dWTSwvwdqqY4yFhwUdJcb.0
  ```

