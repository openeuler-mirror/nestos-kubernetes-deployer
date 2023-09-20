# Housekeeper

## 概述

云原生领域主要采用容器技术与容器编排技术实现了业务发布、运维，与底层环境高度解耦，但同时带来运维技术栈的不统一，造成了k8s和底层操作系统分别独立管理，运维平台重复建设等问题。为了应对这些问题，NKD集成了housekeeper模块，实现了业务与NestOS云底座操作系统一致性运维，采用了容器化的方式进行运维管理。housekeeper的主要更新流程是当操作系统或k8s基础组件需要升级维护时，使用镜像构建工具重新构建新版系统镜像，并在查询到新版镜像后，向集群创建housekeeper CR资源。集群中的housekeeper服务按照配置逐次对集群节点进行升级，完成整个集群的升级工作。

## 架构介绍
housekeeper的架构如图
![housekeeper-arch](/docs/figures/housekeeper-arch.jpg)
如图所示housekeeper主要包含三个组件housekeeper-operator-manager、housekeeper-controller-manager、housekeeper-daemon
- housekeeper-operator-manager: 以Deployment形式运行在Master节点上，负责协调所有Machines进行升级（不负责直接更新），并标记准备升级的节点。
- housekeeper-controller-manager：以DaemonSet形式运行在集群中的所有节点上，负责驱逐业务pod，以及转发升级信息到housekeeper-daemon。
- housekeeper-daemon: 接收来自housekeeper-controller-manager的信息，并根据指令执行OS的原子性更新或者kubernetes版本的升级。
