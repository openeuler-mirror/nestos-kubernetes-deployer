# 配置文件说明

NestOS镜像下载地址见[官网](https://nestos.openeuler.org/)
``` shell
cluster_id: cluster                                 # 集群名称
architecture: amd64                                 # 部署集群的机器架构,支持amd64或者arm64
platform: libvirt                                   # 部署平台为libvirt
infraplatform
  uri: qemu:///system                                
  osimage: https://nestos.org.cn/nestos20230928/nestos-for-container/x86_64/NestOS-For-Container-22.03-LTS-SP2.20230928.0-qemu.{arch}.qcow2                                             # 指定部署集群机器的操作系统镜像地址，支持架构x86_64或者aarch64
  cidr: 192.168.132.0/24                            # 路由地址
  gateway: 192.168.132.1                            # 网关地址
username: root                                      # 指定 ssh 登录所配置节点的用户名
password: $1$yoursalt$UGhjCXAJKpWWpeN8xsF.c/        # 指定 ssh 登录所配置节点的密码
sshkey: "/root/.ssh/id_rsa.pub"                     # ssh 免密登录的密钥存储文件的路径
master:                                             # 配置master节点的列表
- hostname: k8s-master01                            # 该节点的名称
  hardwareinfo:                                     # 该节点配置的硬件资源信息
    cpu: 4                                          # 该节点CPU的核数
    ram: 8192                                       # 该节点的内存大小
    disk: 50                                        # 该节点的磁盘大小
  ip: "192.168.132.11"                              # 该节点的IP地址
  ign_data:                                         # 该节点的Ignition文件的路径
worker:                                             # 配置worker节点的列表
- hostname: k8s-worker01            
  hardwareinfo:
    cpu: 4
    ram: 8192
    disk: 50
  ip: ""                                            # 如果不设置worker节点IP地址，则由dhcp自动分配，默认为空
  ign_data: "/etc/nkd/cluster/ignition"
kubernetes:                                         # 集群相关配置列表
  kubernetes_version: "v1.23.10"                    # 部署集群的版本
  apiserver_endpoint: "192.168.132.11:6443"         # 对外暴露的APISERVER服务的地址或域名   
  image_registry: "k8s.gcr.io"                      # 下载容器镜像时使用的镜像仓库的mirror站点地址
  pause_image: "pause:3.6"                          # 容器运行时的pause容器的容器镜像名称
  release_image_url: "hub.oepkgs.net/nestos/nestos:22.03-LTS-SP2.20230928.0-{arch}-k8s-v1.23.10"                             # 包含K8S二进制组件的NestOS发布镜像的地址，支持架构x86_64或者aarch64
  token: ""                                         # 启动引导过程中使用的令牌，默认自动生成
  adminkubeconfig: /etc/nkd/cluster/admin.config    # 集群管理员配置文件admin.conf的路径
  certificatekey: ""                                # 添加新的控制面节点时用来解密所下载的Secret中的证书的秘钥
  network:                                          # k8s集群网络配置
    service_subnet: "10.96.0.0/16"                  # k8s创建的service的IP地址网段
    pod_subnet: "10.100.0.0/16"                     # k8s集群网络的IP地址网段
    coredns_image_version: "v1.8.6"                 # coredns镜像版本
housekeeper:                                                                                          # housekeeper相关配置列表
  deployhousekeeper: false                                                                            # 是否部署housekeeper
  operatorimageurl: "hub.oepkgs.net/nestos/housekeeper/{arch}/housekeeper-operator-manager:{tag}"     # housekeeper-operator镜像的地址，支持架构amd64或者arm64
  controllerimageurl: "hub.oepkgs.net/nestos/housekeeper/{arch}/housekeeper-controller-manager:{tag}" # housekeeper-controller镜像的地址，支持架构amd64或者arm64   
certasset:                                          # 配置外部证书文件路径列表，默认自动生成
  rootcacertpath: ""                
  rootcakeypath: ""
  etcdcacertpath: ""
  etcdcakeypath: ""
  frontproxycacertpath: ""
  frontproxycakeypath: ""
  sapub: ""
  sakey: ""
```

设置部署平台为openstack，需要重新设置“infraplatform”字段配置参数
``` shell
platform: openstack                                   # 部署平台为openstack
infraplatform                      
	username:                                           # openstack用户名，需要有创建资源权限                                       
	password:                                           # openstack登录密码，用于登录openstack平台
	tenant_name:                                        # openstack租户名，用户所属的合集，例如：admin
	auth_url:                                           # openstack鉴权地址，例如：http://{ip}:{port}/v3
	region:                                             # openstack地区，用于资源隔离，例如：RegionOne
	internal_network:                                   # openstack内部网络名称，用户自定义内部网络名称
	external_network:                                   # openstack外部网络名称，用户自定义外部网络名称
	glance_name:                                        # 创建openstack实例的qcow2镜像
	availability_zone:                                  # 可用域，默认nova
```