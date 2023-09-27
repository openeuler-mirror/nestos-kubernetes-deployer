# nestos-kubernetes-deployer

## 介绍
nestos-kubernetes-deployer简称NKD，是基于NestOS部署kubernetes集群运维而准备的解决方案。其目标是在集群外提供对集群基础设施（包括操作系统和kubernetes基础组件）的部署、更新和配置管理等服务，从而简化了集群部署和升级的流程。

## 支持平台
NKD根据集群需求，连接基础设施提供商动态创建所需的IaaS资源，支持裸金属和虚拟化场景，目前优先实现openstack场景。

## 软件架构
整体架构如图：
![arch](/docs/figures/arch.jpg)

NKD的整体架构由多个组件构成，主要包括NKDS（NestOS-kubernetes-deployer-service）作为主体、部署到集群中的HKO（housekeeper operator）以及集成在NestOS镜像中的installer。此外，还可以配合NestOS镜像构建工具链、配置管理仓库（如git）和私有化部署的容器镜像仓库，共同完成集群运维任务。目前NKDS以命令行工具提供，暂不提供对外http接口和前端配置页面，证书管理及健康检测模块，但主体功能所需的基础设施管理、配置管理、系统镜像管理等模块已初步形成。HKO主要包括面向集群的HKO组件和集成在NestOS镜像中的HKD（housekeeper daemon）组件。目前installer组件负责在系统点火阶段部署创建K8S集群，未来计划将其功能融合到HKD组件中，使整体方案更加精简，更易于用户根据个性化需求管理所需的K8S基础组件。

## 安装部署
NKD目前提供工具形式，仅支持通过命令行参数或应用配置文件部署k8s集群。后续会提供用户友好的前端配置界面，便于轻松生成所需配置，并提供配置变更版本管理功能。在部署NestOS系统时，需要通过ignition点火机制传入系统部署后所需的动态配置。NKD会将用户提供的kubernetes集群部署所需的配置自动合并到ign文件中，使得节点在部署完成操作系统后引导自动开始创建k8s集群，无需手动干预。详细内容请见[quick-start.md](docs/user_guide/nkds/quick-start.md).

## 升级维护
NKD提供了操作系统或k8s基础组件升级维护的功能。用户可以选择是否部署housekeeper自定义资源，用于后续的维护升级。housekeeper的架构介绍请见 [architecture.md](docs/user_guide/housekeeper/architecture.md)。用户部署指南请见[quick-start.md](docs/user_guide/housekeeper/quick-start.md).

## 未来规划
NKD的最终目标是以长期驻留服务形式提供运维服务，同时支持多个集群的管理。它将提供持久化配置变更记录、证书管理、多种更新升级策略和镜像源频道等功能。未来，我们将持续优化NKD的功能和性能，并引入更多智能化特性，如自动化故障处理和资源优化等。我们的目标是将NKD打造成NestOS生态中的核心组件，为云原生场景下的运维工作提供全方位支持，进一步推动云原生技术的发展和应用。

## 参与贡献
非常欢迎对本项目感兴趣的伙伴加入我们，并参与贡献。

## License
Apache License 2.0