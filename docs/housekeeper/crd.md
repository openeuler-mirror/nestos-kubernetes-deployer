# housekeeper资源定义

```
apiVersion: housekeeper.io/v1alpha1
kind: UpgradeOSConfig
metadata:
labels:
app.kubernetes.io/name: upgradeosconfig
app.kubernetes.io/instance: upgradeosconfig-sample
app.kubernetes.io/part-of: housekeeper-operator
app.kubernetes.io/managed-by: kustomize
app.kubernetes.io/created-by: housekeeper-operator
name: upgradeosconfig-sample
spec:
    os:待升级的os信息
        image: os镜像，string,required
        version: os版本，string,required
   os_maintain_strategy:OS升级的运维策略
        max_unavailable: 运维策略中同时最多执行os升级的节点数量,int
        drainer_num:每次驱离数量,int
        drainer_interval:驱离时间间隔,秒为单位;int64
        drainer_type:0,串行；1，并行;int
   pod_drain:pod 的驱逐策略
        drainer_force: 是否立即驱逐，bool
        delete_empty_dir_data:是否删除卷上数据，bool
        drace_period_seconds：Pods 在被强制终止前可以运行的最长时间，秒为单位；int64
        drainer_timeout:整个驱离超时时间,秒为单位;int64
   kernel_params:内核参数
        net_ipv4_tcp_syncookies
        net_ipv4_ip_forward
        net_bridge_bridge_nf_call_iptables
        net_bridge_bridge_nf_call_ip6tables
        vm_swappiness
        vm_overcommit_memory
        fs_file_max
        net_ipv4_conf_all_rp_filter
        net_ipv4_conf_default_rp_filter
        net_ipv4_conf_all_accept_redirects
        net_ipv4_conf_default_accept_redirects
        kernel_pid_max
        vm_max_map_count
```