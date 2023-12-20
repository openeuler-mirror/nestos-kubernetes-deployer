# 配置文件说明

``` shell
cluster_id: cluster                                 // 集群名称
platform: libvirt                                   // 部署平台，可选libvirt或者openstack
infraplatform:
  uri: ""
  osimage_path: "/etc/nkd/NestOS-For-Container-22.03-LTS-SP2.20230928.0-qemu.x86_64.qcow2"
  cidr: "192.168.132.0/24"
  gateway: "192.168.132.1"
master:                                             // 配置master节点的列表
- hostname: k8s-master01                            // 该节点的名称
  hardwareinfo:                                     // 该节点配置的硬件资源信息
    cpu: 4                                          // 该节点CPU的核数
    ram: 8096                                       // 该节点的内存大小
    disk: 50                                        // 该节点的磁盘大小
  username: root                                    // 该节点的ssh登录用户名
  password: "$1$yoursalt$UGhjCXAJKpWWpeN8xsF.c/"    // 该节点的ssh登录密码
  sshkey: "ssh-rsa AAAA***root.localdomain"         // 该节点ssh免密登录的密钥
  ip: "192.168.132.11"                              // 该节点的IP地址
  ign_data:                                         // 该节点的Ignition文件的路径，不配置则生成该文件
- hostname: k8s-master02
  hardwareinfo:
    cpu: 4
    ram: 8096
    disk: 50
  username: root
  password: "$1$yoursalt$UGhjCXAJKpWWpeN8xsF.c/"
  sshkey: "ssh-rsa AAAA***root.localdomain"
  ip: "192.168.132.12"
  ign_data: ""
- hostname: k8s-master03
  hardwareinfo:
    cpu: 4
    ram: 8096
    disk: 50
  username: root
  password: "$1$yoursalt$UGhjCXAJKpWWpeN8xsF.c/"
  sshkey: "ssh-rsa AAAA***root.localdomain"
  ip: "192.168.132.13"
  ign_data: "/etc/nkd/cluster/ignition"
worker:                             // 配置worker节点的列表
- hostname: k8s-worker01            
  hardwareinfo:
    cpu: 4
    ram: 8096
    disk: 50
  username: root
  password: "$1$yoursalt$UGhjCXAJKpWWpeN8xsF.c/"
  sshkey: "ssh-rsa AAAA***root.localdomain"
  ip: "192.168.132.14"
  ign_data: "/etc/nkd/cluster/ignition"
- hostname: k8s-worker02
  hardwareinfo:
    cpu: 4
    ram: 8096
    disk: 50
  username: root
  password: "$1$yoursalt$UGhjCXAJKpWWpeN8xsF.c/"
  sshkey: "ssh-rsa AAAA***root.localdomain"
  ip: "192.168.132.15"
  ign_data: "/etc/nkd/cluster/ignition"
- hostname: k8s-worker03
  hardwareinfo:
    cpu: 4
    ram: 8096
    disk: 50
  username: root
  password: "$1$yoursalt$UGhjCXAJKpWWpeN8xsF.c/"
  sshkey: "ssh-rsa AAAA***root.localdomain"
  ip: "192.168.132.16"
  ign_data: "/etc/nkd/cluster/ignition"
kubernetes:                                         // 集群相关配置列表
  kubernetes_version: "v1.23.10"                    // 部署集群的版本
  apiserver_endpoint: "192.168.132.11:6443"         // 对外暴露的APISERVER服务的地址或域名   
  image_registry: "k8s.gcr.io"                      // 下载容器镜像时使用的镜像仓库的mirror站点地址
  pause_image: "pause:3.6"                          // 容器运行时的pause容器的容器镜像名称
  release_image_url: ""                             // 包含K8S二进制组件的NestOS发布镜像的地址
  token: o0tztj.kjsvjzha417yk5oo                    // 启动引导过程中使用的令牌
  adminkubeconfig: /etc/nkd/cluster/admin.config    // 集群管理员配置文件admin.conf的路径
  certificatekey: ""                                // 添加新的控制面节点时用来解密所下载的Secret中的证书的秘钥
  network:                                          // k8s集群网络配置
    service_subnet: 10.96.0.0/16                    // k8s创建的service的IP地址网段
    pod_subnet: "10.100.0.0/16"                     // k8s集群网络的IP地址网段
    coredns_image_version: "v1.8.6"                 // coredns镜像版本
housekeeper:                                                                                  // housekeeper相关配置列表
  deployhousekeeper: false                                                                    // 是否部署housekeeper
  operatorimageurl: "hub.oepkgs.net/nestos/nkd/{arch}/housekeeper-operator-manager:{tag}"     // housekeeper-operator镜像的地址
  controllerimageurl: "hub.oepkgs.net/nestos/nkd/{arch}/housekeeper-controller-manager:{tag}" // housekeeper-controller镜像的地址   
  kubeversion: ""                                   // 升级的K8S版本
  evictpodforce: false                              // 用于升级时是否强制驱逐pod
  maxunavailable: 2                                 // 用于进行升级的最大节点数
  osimageurl: ""                                    // 用于升级容器镜像的地址, 需要为容器镜像格式 REPOSITORY/NAME[:TAG@DIGEST]
certasset:                                          // 配置外部证书文件路径列表
  rootcacertpath: ""                
  rootcakeypath: ""
  etcdcacertpath: ""
  etcdcakeypath: ""
  frontproxycacertpath: ""
  frontproxycakeypath: ""
  sapub: ""
  sakey: ""
```