# 快速使用指导

## 编译及部署

nkd项目在openEuler社区提供rpm软件包以供安装使用，同时也可以自行编译构建
  
### 编译指导

* 编译环境：Linux x86_64/aarch64
* 进行编译需要以下软件包：
  * golang >= 1.17
  * make
  * git
  ``` shell
  sudo yum install golang make git
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

### 部署指导

* 环境要求
  * Linux x86_64/aarch64
  * 已完成搭建OpenStack环境
* 部署
  * 创建工作目录
    ```
    mkdir nkd && cd nkd
    ```
  * 获取当前系统架构的[terraform二进制文件](https://developer.hashicorp.com/terraform/downloads)，保存在工作目录下
  * 获取${terraform-provider-openstack}插件，保存在工作目录下
    ```
    # 本项目可在有网络环境中自动部署${terraform-provider-openstack}插件（该功能由terraform提供），若在离线环境中使用，请在工作目录自行构建如下目录：
    # x86_64
    ./providers/registry.terraform.io/terraform-provider-openstack/openstack/${version}/linux_amd64/${terraform-provider-openstack_${version}}
    # aarch64
    ./providers/registry.terraform.io/terraform-provider-openstack/openstack/${version}/linux_arm64/${terraform-provider-openstack_${version}}
    ```
  * 分别打印master节点与worker节点的配置文件，根据用户配置需求自定义修改
    ``` shell
    nkd config print master
    nkd config print worker
    ```
  * 根据用户自定义修改的配置文件，生成ignition文件
    ``` shell
    nkd config init -c master/master.yaml
    nkd config init -c worker/worker.yaml   
    ```
  * 部署master节点
    ``` shell
    nkd deploy master
    ```
  * 部署worker节点
    ``` shell
    nkd deploy worker
    ```
  * 扩展worker节点。参数为扩展后worker节点的数量，例如扩展后worker节点数量为10
    ``` shell
    nkd extend -n 10
    ```
  * 销毁集群
    ``` shell
    nkd destroy worker
    nkd destroy master
    ```

## 使用指导

### 注意事项
  * 部署集群前用户需要自定义构建部署镜像

### 自定义镜像制作
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

### 参数说明

* 创建基础设施配置参数字段说明如下：
  | 参数           |参数类型  | 参数说明                                                  | 使用说明 | 是否必选         |
  | -------------- | ------  | -----------------------------------------------------------| ----- | ---------------- |
  | platform      | string  | 部署平台名称           | 默认openstack    | 是         |
  | user_name      | string  | openstack用户名           | 需要有创建资源权限    | 是         |
  | password      | string  | openstack登录密码           | 用于登录openstack平台    | 是         |
  | tenant_name      | string  | openstack租户名           | 用户所属的合集，例如：admin    | 是         |
  | auth_url      | string  |  openstack鉴权地址          | 例如：http://{ip}:{port}/v3    | 是         |
  | region      | string  | openstack地区           | 用于资源隔离，例如：RegionOne     | 是         |
  | internal_network      | string  | openstack内部网络名称           | 用户自定义内部网络名称    | 是         |
  | external_network      | string  | openstack外部网络名称           | 用户自定义外部网络名称    | 是         |
  | glance      | string  | 创建openstack实例的qcow2镜像           | 使用openstack版本的nestos镜像    | 是         |
  | zone      | string  | 可用域           | 默认nova    | 是         |
  | vcpus      | int  | 节点机器的cpu数量           | 至少2个cpu    | 是         |
  | ram      | int  |  节点机器的内存大小          | 至少2GB以上的内存    | 是         |
  | disk      | int  | 节点机器的磁盘空间          | 至少20GB以上的磁盘空间    | 是         |

* 创建集群配置参数字段说明如下：
  | 参数           |参数类型  | 参数说明                                                  | 使用说明 | 是否必选         |
  | -------------- | ------  | -----------------------------------------------------------| ----- | ---------------- |
  | node      | string  | 节点类型           | 默认master或者worker    | 是         |
  | pauseimagetag      | string  | pause镜像的版本号          | 根据集群环境确定pause镜像版本，例如：3.6    | 是         |  
  | corednsimagetag      | string  | coredns镜像的版本号           | 需选择kubernetes版本兼容的coredns版本    | 是         |  
  | releaseimageurl      | string  | nestos部署镜像的地址           | 用户需根据构建要求，自行构建带kubernetes组件的部署容器镜像    | 是         |  
  | certificatekey      | string  | 添加新的控制面节点时用来解密所下载的Secret中的证书的秘钥           | 字段被添加到InitConfiguration和JoinConfiguration中    | 是         |    
  | token/bootstraptoken      | string  | 用来在节点与控制面之间建立双向的信任关系，在向集群中添加节点时使用         | 需用户自定义生成，形式为：abcdef.0123456789abcdef       | 是         |
  | tlsbootstraptoken      | string  | 是 TLS 启动引导过程中使用的令牌           |  需用户自定义生成，形式为：abcdef.0123456789abcdef      | 是         |
  | count      | int  | 部署节点的数量           | 部署master节点的数量需要是单数，例如：3       | 是         |  
  | masterhostname      | string  | master节点的hostname           | 用户自定义，例如：k8s-master       | 是         |
  | workerhostname      | string  | wokrer节点的hostname           | 用户自定义，例如：k8s-worker               | 是         |
  | masterips      | string  | master节点的ip地址           | 如部署平台是openstack，此项是内部ip地址    | 是         |
  | username      | string  | 节点的用户名           | 默认root    | 是         |  
  | password      | string  | 节点登录密码           | 密码需为带盐值的hash值，例如："$1$yoursalt$UGhjCXAJKpWWpeN8xsF.c/"    | 否         |
  | sshkey       | string  | ssh公钥           | 部署机器的公钥密码     | 是         |
  | registry       | string  | 镜像仓库地址           | kubernetes组件镜像的拉取地址    | 是         |
  | servicesubnet      | string  |  kubernetes 服务所使用的子网           | 默认值为 10.96.0.0/16    | 是         |
  | podsubnet      | string  |  为pod所使用的子网          | 默认值为10.100.0.0/16    | 是         |
  | kubernetesversion      | string  | kubernetes的版本号           | 例如：v1.23.10     | 是         |
  | apiserverendpoint      | string  | 为API服务器的IP地址       | 第一个master节点的ip地址加端口，例如：{ip}:6443    | 是         |
















