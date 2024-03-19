# Ignition

Ignition is a utility used by NestOS during the initramfs phase to manipulate disks. This includes tasks such as creating users, adding trusted SSH keys, writing files (regular files, systemd services, etc.), and configuring networks. Upon the first boot, Ignition reads its configuration and applies it. Ignition utilizes JSON configuration files to represent the set of changes to be made. The format of this configuration is detailed in the specification [here](https://coreos.github.io/ignition/specs/).

## Provide configuration
The Ignition file (controlplane.ign) required to generate and deploy the cluster is approximately 90KB in size, containing declarations for systemd services, certificates, and other necessities for cluster create. When deploying the cluster on the OpenStack platform, the Nova user data limit is 64KB, making direct infrastructure deployment impossible. To address this issue, NKD has created a smaller Ignition file (*-merge.ign) to serve as a bootstrapping configuration file for creating the infrastructure. The main Ignition file will be loaded into memory accessible via HTTP service, enabling the smaller Ignition file to be automatically loaded during the system boot phase.

When the infrastructure of the cluster nodes is successfully created, the operating system is in the bootstrapping phase. At this point, Ignition retrieves configuration information and writes configuration files (including users, cluster certificates, cluster deployment services, etc.) to the disks of the node machines, and then starts systemd services. Once the "release-image-pivot.service" service is started successfully, the node machines switch the file system to a customized version of NestOS based on the rpm-ostree mechanism. Afterward, the cluster creation service is executed, continuing until the completion of cluster creation, at which point the service is shut down.

### ControlPlane Node
The configuration information of the control plane node's Ignition file is as shown in the image:
![ignition_design_1](/docs/en/figures/ignition_design_1.jpg)

### Master Node
The configuration information of the master node's Ignition file is as shown in the image：
![ignition_design_2](/docs/en/figures/ignition_design_2.jpg)

### Worker Node
The configuration information of the worker node's Ignition file is as shown in the image：
![ignition_design_3](/docs/en/figures/ignition_design_3.jpg)

The generated Ignition file directory structure is as follows：
``` shell
$ tree
.
├── controlplane.ign
├── controlplane-merge.ign
├── master.ign
├── master-merge.ign
├── worker.ign
└── worker-merge.ign
```
```shell
The content of the smaller "*-merge.ign" Ignition file is as follows:
Where:
    - IP: The address where the HTTP service is deployed, defaulting to the host machine's IP address.
    - port: The access port for the HTTP service, defaulting to port 9080.
    - {*.ign}: The primary Ignition files for executing cluster creation.
{
  "ignition": {
    "config": {
      "merge": [
        {
          "source": "http://{IP}:{port}/{*.ign}"
        }
      ]
    },
    "version": "3.2.0"
  }
}
```