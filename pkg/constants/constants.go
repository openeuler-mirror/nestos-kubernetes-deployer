/*
Copyright 2023 KylinSoft  Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package constants

import "os"

const (
	// 节点类型标识符
	Controlplane = "controlplane"
	Master       = "master"
	Worker       = "worker"

	// 容器运行类型标识符
	Isulad     = "isulad"
	Docker     = "docker"
	Crio       = "crio"
	Containerd = "containerd"

	// 操作系统引导文件存储文件夹
	BootConfigSaveDir = "bootconfig"
	// 操作系统引导阶段可能启用的服务
	InitClusterService       = "init-cluster.service"
	JoinMasterService        = "join-master.service"
	JoinWorkerService        = "join-worker.service"
	KubeletService           = "kubelet.service"
	ReleaseImagePivotService = "release-image-pivot.service"
	SetKernelPara            = "set-kernel-para.service"
	BootConfigSystemdPath    = "bootconfig/systemd"
	// 操作系统引导阶段可能写入的文件
	InitClusterYaml       = "/etc/nkdfiles/init-config.yaml.template"
	IsuladConfig          = "/etc/isulad/daemon.json.template"
	DockerConfig          = "/etc/docker/daemon.json"
	ReleaseImagePivotFile = "/etc/nkdfiles/node-pivot.sh.template"
	SetKernelParaConf     = "/etc/sysctl.d/kubernetes.conf"
	KubeletServiceConf    = "/etc/systemd/system/kubelet.service.d/10-kubeadm.conf.template"
	Hosts                 = "/etc/hosts.template"
	HookFilesPath         = "/etc/nkdfiles/hookfiles/"
	BootConfigFilesPath   = "bootconfig/files"
	// 引导配置文件名称
	ControlplaneIgn      = "controlplane.ign"
	ControlplaneMergeIgn = "controlplane-merge.ign"
	MasterIgn            = "master.ign"
	MasterMergeIgn       = "master-merge.ign"
	WorkerIgn            = "worker.ign"
	WorkerMergeIgn       = "worker-merge.ign"

	CertsFiles = "certs.json"

	SystemdServiceMode os.FileMode = 0644
	StorageFilesMode   os.FileMode = 0755
	SaveFileDirMode    os.FileMode = 0750
	BootConfigFileMode os.FileMode = 0644
)
