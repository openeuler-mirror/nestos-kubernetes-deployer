# 快速使用指导

## 镜像构建

* 部署housekeeper可以使用已发布的容器镜像版本，同时支持用户自行编译构建
* housekeeper-operator-manager和housekeeper-controller-manager容器镜像的0.1.0版本已发布，镜像拉取信息如下：
  ``` shell
  docker pull hub.oepkgs.net/nestos/nkd/amd64/housekeeper-operator-manager:0.1.0
  docker pull hub.oepkgs.net/nestos/nkd/amd64/housekeeper-controller-manager:0.1.0
  docker pull hub.oepkgs.net/nestos/nkd/arm64/housekeeper-operator-manager:0.1.0
  docker pull hub.oepkgs.net/nestos/nkd/arm64/housekeeper-controller-manager:0.1.0
  ```  

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
  cd nestos-kubernetes-deployer/housekeeper
  sudo make
  ```
* 用户需要自定义编写Dockerfile文件构建容器镜像
  * housekeeper-operator-manager以及housekeeper-controller-manager容器镜像的构建，Dockerfiles示例如下：
      ``` dockerfile
      FROM base_image
      COPY ./bin/housekeeper-controller-manager /housekeeper-controller-manager
      ENTRYPOINT ["/housekeeper-controller-manager"]

      FROM base_image
      COPY ./bin/housekeeper-operator-manager /housekeeper-operator-manager
      ENTRYPOINT ["/housekeeper-operator-manager"]
      ```

### 升级镜像制作

* NestOS容器镜像支持利用Dockerfile在原来的基础上构建新的容器镜像
* 制作注意事项
    * 请确保已安装docker
    * 基础镜像需从NestOS官网下载最新版本容器镜像
    * 制作kubernetes版本升级镜像，需提前下载相对应版本的kubeadm、kubelet二进制文件并复制到/usr/bin目录，以及将calico网络插件的yaml文件复制到/etc/nkd目录
    * 软件包的安装需要使用rpm-ostree命令
 * Dockerfiles示例如下
      ``` dockerfile
      FROM nestos_base_image
      COPY kube* /usr/bin/
      COPY calico.yaml /etc/nkd/
      RUN ostree container commit
      ```

## 部署指导

* 环境要求
  * Linux x86_64/aarch64
  * NKD已部署kubernetes集群
* 部署
  * 使用kubernetes的声明式API进行配置，部署CRD（CustomResourceDefinition），housekeeper-operator-manager、housekeeper-controller-manager控制器，以及RBAC(Role-Based Access Control)机制的YAML。
  * YAML的示例模板在/docs/user_guide/housekeeper/config目录下，用户可以根据需要进行简单的修改
  * 如果已经编辑好了YAML，执行部署crd、manager和rbac命令：
    ``` shell
    kubectl apply -f confg/crd
    kubectl apply -f config/rbac 
    kubectl apply -f config/manager
    ```
  * 部署后通过下面命令行查看集群各个组件是否正常启动，如果所有组件的STATUS都是Running，即部署完成
    ``` shell
    kubectl get pods -A
    ```

## 使用指导

### 注意事项

  * housekeeper不支持kubernetes跨大版本升级
  * 高可用集群Master节点为逐一节点进行升级
  
### 参数说明

创建CRD对象的参数字段及说明如下：
  | 参数           |参数类型  | 参数说明                                                  | 使用说明 | 是否必选         |
  | -------------- | ------  | -----------------------------------------------------------| ----- | ---------------- |
  | osVersion      | string  | 用于升级容器镜像的NestOS版本           | 需为NestOS version格式，例如NestOS 23.03.20230511.0 | 是         |
  | osImageURL      | string  | 用于升级容器镜像的地址           | 需要为容器镜像格式 REPOSITORY/NAME[:TAG@DIGEST] | 是         |
  | kubeVersion      | string  | 用于升级kubernetes的版本号           | 如果仅升级OS版本，此项需填空 | 是         |
  | evictPodForce      | bool  | 用于升级时是否强制驱逐pod           | 需为true或者false | 是         |
  | maxUnavailable      | int  | 用于进行升级的最大节点数           | 仅限制worker节点同时升级的节点数量 | 是         |

### 升级指导

* 升级前请先制作完成升级所需的容器镜像，并编辑config/samples/housekeeper.io_v1alpha1_update.yaml文件，示例如下：
    ``` yaml
    apiVersion: housekeeper.io/v1alpha1
    kind: Update
    metadata:
      name: housekeeper-upgrade
    spec:
      osVersion: os.version
      osImageURL: image.url
      kubeVersion: kubernetes.version
      evictPodForce: false
      maxUnavailable: 2
    ```
* 查看未升级节点的OS版本及kubernetes版本
    ``` shell
    kubectl get nodes
    ```
* 执行下面命令，在集群中部署cr实例后，节点会根据配置参数进行升级
    ``` shell
    kubectl apply -f housekeeper.io_v1alpha1_update.yaml
    ```
* 在次查看节点的OS版本及kubernetes版本确认是否升级完成
    ``` shell
    kubectl get nodes
    ```
* 如后续再次升级，需调整housekeeper.io_v1alpha1_update.yaml中osVersion、osImageURL、kubeVersion字段的相应信息，并且制作升级所需的容器镜像。