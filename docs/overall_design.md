# 方案设计

## 整体架构
![arch](./figures/overall_arch.jpg)

模块说明：
- http server：提供用户交互的HTTP接口以及友好的前端管理界面，使用户能够轻松进行操作和配置
- infra-manager：负责创建和删除基础设施
- config-manager：管理集群配置信息，包括创建、更新、删除等操作
- cert-manager：负责创建和更新集群及节点证书，维护系统安全，确保证书的有效性和合规性
- healthz-worker：实时监测系统的健康状况，及时发现并报告问题，确保系统稳定可靠的运行
- installer：执行系统点火阶段的任务，负责部署和创建K8S集群
- HKO (Housekeeper Operator)：部署在集群中，负责集群级操作的组件
- HKD (Housekeeper Daemon)：集成在NestOS镜像中，属于HKO的组成部分
- 镜像构建工具链：用于构建NestOS镜像的工具链，支持系统的自定义镜像生成
- 配置仓库：存储和管理配置信息的数据库
- 容器镜像仓库：用于存储和管理私有化部署的容器镜像，保障应用程序的可靠性和安全性

备注：http server、healthz-worker暂未支持

## 详细设计
NKD模块交互关系图
![detailed_design](/docs/figures/detailed_design.jpg)

### housekeeper模块设计
在集群部署阶段，用户可以选择是否部署housekeeper
详细内容见[设计文档](/docs/housekeeper_design.md)

