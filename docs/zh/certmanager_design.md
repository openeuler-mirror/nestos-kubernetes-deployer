

## Cert Manager

NKD使用Kubeadm安装Kubernetes集群，默认情况下Kubernetes所需要的PKI证书由Kubeadm生成，这一做法会带来很多安全风险。Kubeadm提供外部CA模式来支持用户自定义PKI证书，为了保障通信安全，NKD通过证书管理模块，遵循Kubernetes最佳实践的[PKI证书和要求](https://kubernetes.io/zh-cn/docs/setup/best-practices/certificates/)来生成所有证书。



### CA证书介绍

| 路径                   | 默认 CN                   | 描述                   |
| ---------------------- | ------------------------- | ---------------------- |
| ca.crt,key             | kubernetes-ca             | Kubernetes 通用 CA     |
| etcd/ca.crt,key        | etcd-ca                   | 与 etcd 相关的所有功能 |
| front-proxy-ca.crt,key | kubernetes-front-proxy-ca | 用于前端代理           |

NKD支持在配置文件【certasset】字段自定义CA证书，如果用户未提供，NKD将自动生成所需CA证书。

### 所有证书如下

| 默认 CN                       | 建议的密钥路径               | 建议的证书路径               |
| ----------------------------- | ---------------------------- | ---------------------------- |
| etcd-ca                       | etcd/ca.key                  | etcd/ca.crt                  |
| kube-apiserver-etcd-client    | apiserver-etcd-client.key    | apiserver-etcd-client.crt    |
| kubernetes-ca                 | ca.key                       | ca.crt                       |
| kubernetes-ca                 | ca.key                       | ca.crt                       |
| kube-apiserver                | apiserver.key                | apiserver.crt                |
| kube-apiserver-kubelet-client | apiserver-kubelet-client.key | apiserver-kubelet-client.crt |
| front-proxy-ca                | front-proxy-ca.key           | front-proxy-ca.crt           |
| front-proxy-ca                | front-proxy-ca.key           | front-proxy-ca.crt           |
| front-proxy-client            | front-proxy-client.key       | front-proxy-client.crt       |
| etcd-ca                       | etcd/ca.key                  | etcd/ca.crt                  |
| kube-etcd                     | etcd/server.key              | etcd/server.crt              |
| kube-etcd-peer                | etcd/peer.key                | etcd/peer.crt                |
| etcd-ca                       |                              | etcd/ca.crt                  |
| kube-etcd-healthcheck-client  | etcd/healthcheck-client.key  | etcd/healthcheck-client.crt  |

获取用于服务账号管理的密钥对：

| 私钥路径 | 公钥路径 |
| -------- | -------- |
| sa.key   |          |
|          | sa.pub   |

下面提供了自行生成所有密钥和证书时所需要提供的文件路径。

```console
/etc/kubernetes/pki/etcd/ca.key
/etc/kubernetes/pki/etcd/ca.crt
/etc/kubernetes/pki/apiserver-etcd-client.key
/etc/kubernetes/pki/apiserver-etcd-client.crt
/etc/kubernetes/pki/ca.key
/etc/kubernetes/pki/ca.crt
/etc/kubernetes/pki/apiserver.key
/etc/kubernetes/pki/apiserver.crt
/etc/kubernetes/pki/apiserver-kubelet-client.key
/etc/kubernetes/pki/apiserver-kubelet-client.crt
/etc/kubernetes/pki/front-proxy-ca.key
/etc/kubernetes/pki/front-proxy-ca.crt
/etc/kubernetes/pki/front-proxy-client.key
/etc/kubernetes/pki/front-proxy-client.crt
/etc/kubernetes/pki/etcd/server.key
/etc/kubernetes/pki/etcd/server.crt
/etc/kubernetes/pki/etcd/peer.key
/etc/kubernetes/pki/etcd/peer.crt
/etc/kubernetes/pki/etcd/healthcheck-client.key
/etc/kubernetes/pki/etcd/healthcheck-client.crt
/etc/kubernetes/pki/sa.key
/etc/kubernetes/pki/sa.pub
```

### KubeConfig文件

| 文件名                  | 命令                    | 说明                                                       |
| ----------------------- | ----------------------- | ---------------------------------------------------------- |
| admin.conf              | kubectl                 | 配置集群的管理员                                           |
| kubelet.conf            | kubelet                 | 集群中的每个节点都需要一份                                 |
| controller-manager.conf | kube-controller-manager | 必须添加到 `manifests/kube-controller-manager.yaml` 清单中 |
| scheduler.conf          | kube-scheduler          | 必须添加到 `manifests/kube-scheduler.yaml` 清单中          |

下面是前表中所列文件的完整路径：

```console
/etc/kubernetes/admin.conf
/etc/kubernetes/kubelet.conf
/etc/kubernetes/controller-manager.conf
/etc/kubernetes/scheduler.conf
```

