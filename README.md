# nestos-kubernetes-deployer
![ignition_design_2](/docs/logo/nkd-logo.png)

## 介绍
NKD（NestOS Kubernetes Deployer）是NestOS团队面向容器云场景开发的集群部署运维工具。涵盖了基础设施和Kubernetes核心组件的部署、更新和配置管理等一系列功能，为用户提供了一站式的解决方案。支持自定义多种容器运行时（包括crio、iSulad、docker和containerd），并且兼容多种平台部署，确保用户能够轻松应对各种复杂的部署需求。此外，NKD具备集群自签名证书创建的能力，并支持部署多种版本的Kubernetes集群，从而覆盖实际使用中可能遇到的各种场景。

## 软件架构
详细内容请见[软件架构说明](docs/zh/overall_design.md)

## 安装部署
详细内容请见[用户操作手册](docs/zh/manual.md)

## 未来规划
NKD的最终愿景是通过长期驻留的服务形式，为运维提供支持，并同时支持多个集群的高效管理。它将提供持久的配置变更记录、证书管理、多种更新升级策略以及镜像源频道等丰富功能。未来，我们将持续精进NKD的功能和性能，并引入更多智能化特性，例如自动化故障处理和资源优化等。我们的目标是将NKD塑造成NestOS生态的核心组成部分，为云原生场景下的运维工作提供全方位支持，从而进一步推动云原生技术的发展和广泛应用。

## 参与贡献
非常欢迎对本项目感兴趣的伙伴加入我们，并参与贡献。

## License
Apache License 2.0