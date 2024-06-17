# User Operation Manual

## Preparation

* Environment Requirements
  * Linux x86_64/aarch64
  * Installation of the tofu software package
    ``` shell
    # Install amd64 version
    $ wget https://github.com/opentofu/opentofu/releases/download/v1.6.2/tofu_1.6.2_amd64.rpm
    $ rpm -ivh tofu_1.6.2_amd64.rpm
    ``` 
    ``` shell
    # Install arm64 version
    $ wget https://github.com/opentofu/opentofu/releases/download/v1.6.2/tofu_1.6.2_arm64.rpm
    $ rpm -ivh tofu_1.6.2_arm64.rpm
    ``` 

* Install NKD
  * Choose to directly use precompiled NKD binary files.
  * Compile and install NKD according to the following compilation and installation instructions

Note:
To ensure the smooth operation of NKD deployment, the deployment environment where it resides must be able to communicate properly with the cluster node machines. If there is a firewall, it needs to be configured correctly to allow communication between NKD and the cluster, such as opening specific http service ports. If domain names are used for communication, ensure that DNS server configuration is correct and that the environment where NKD resides can access the DNS server.

## Supported Platforms

### libvirt
Deploying clusters on the libvirt platform requires pre-installation of the libvirt virtualization environment.

### openstack
Deploying clusters on the OpenStack platform requires pre-setup of the OpenStack environment.

### PXE

## Compilation and Installation

* Compilation Environment: Linux x86_64/aarch64
* The following software packages are required for compilation:
  * golang >= 1.21
  * git
  ``` shell
  $ sudo yum install golang git
  ```  
* Use git to obtain the source code of this project
  ``` shell
  sudo git clone https://gitee.com/openeuler/nestos-kubernetes-deployer
  ```
* Compile binaries
  ``` shell
  $ sh hack/build.sh
  ```

## Configuration Management

### Global Configuration
The global configuration file is used to manage the configuration of the entire cluster (or multiple clusters). For specific configuration parameters and default configurations, refer to the [Global Configuration File Description](./globalconfig_file_desc.md)

#### Ignition Service Configuration Parameters
During the NKD cluster deployment process, cluster nodes need to access the ignition service provided by NKD. The ignition service is configured through the following global configuration parameters:
* bootstrap_ign_host：Ignition service address (domain name or IP, usually NKD operating environment)
* bootstrap_ign_port：Ignition service port (default 9080, you need to open the firewall port yourself)

To adapt to multi-NIC environments, the actual listening address of the ignition service is 0.0.0.0.
* In a simple network environment, cluster nodes can directly access the NKD service. The "bootstrap_ign_host" parameter can be left empty. In this case, NKD will detect the IP address with the highest priority in the routing table as the host for accessing the ignition service URL.
* In complex network environments where cluster nodes cannot directly access the NKD runtime environment, the "bootstrap_ign_host" parameter needs to be configured as an externally mapped IP or domain name. Users need to configure NAT mapping or DNS services themselves to ensure that cluster nodes can access the NKD ignition service.

The "bootstrap_ign_port" parameter is currently shared by the ignition service listening port and the ignition service URL access port. In a simple network environment, these two values are consistent. However, in complex network environments, it is necessary to ensure that the externally mapped port for the NKD service is consistent with the locally listened port.

### Cluster Configuration
The cluster configuration file is used to configure each cluster independently. For specific configuration parameters and default configurations, please refer to the [Cluster Configuration File Description](./config_file_desc.md)

## Basic Functions

The specific process is outlined in the "Create Cluster" section. Here are the basic execution commands for NKD:
  ``` shell
  # Generate default configuration template
  $ nkd template -f cluster_config.yaml

  # Deploy the cluster using the configuration file
  $ nkd deploy -f cluster_config.yaml

  # Destroy a specific cluster
  $ nkd destroy --cluster-id [your-cluster-id]

  # Scale the number of nodes in a specific cluster
  $ nkd extend --cluster-id [your-cluster-id] --num 10

  # Upgrade a specific cluster
  # --cluster-id string: Unique identifier for the cluster
  # --force: Force eviction of pods even if unsafe. This may result in data loss or service disruption, use with caution (default: false)
  # --imageurl string: The address of the container image to use for upgrading
  # --kube-version string: Choose a specific kubernetes version for upgrading
  # --kubeconfig string: Specify the access path to the Kubeconfig file，default "/etc/nkd/[your-cluster-id]/admin.config"
  # --maxunavailable uint: Number of nodes that are upgraded at the same time (default: 2)
  $ nkd upgrade --cluster-id [your-cluster-id] --imageurl [your-image-url] --kube-version [your-k8s-version] 
  ```
Supports deploying the cluster using application configuration parameters, in addition to deploying it with application configuration files
  ``` shell
  $ nkd deploy --help
      --arch string                   Architecture for Kubernetes cluster deployment (e.g., amd64 or arm64)
      --bootstrap-ign-host string     Ignition service address (domain name or IP)
      --bootstrap-ign-port string     Ignition service port (default: 9080)
      --certificateKey string         The key that is used for decryption of certificates after they are downloaded from the secret upon joining a new master node. (the certificate key is a hex encoded string that is an AES key of size 32 bytes)
      --cluster-id string             Unique identifier for the cluster
      --controller-image-url string   URL of the container image for the housekeeper controller component
      --deploy-housekeeper            Deploy the Housekeeper Operator. (default: false)
  -f, --file string                   Location of the cluster deploy config file
  -h, --help                          help for deploy
      --image-registry string         Registry address for Kubernetes component container images
      --ipxe-filePath string          Path of config file for iPXE
      --ipxe-ip string                IP address of local machine for iPXE
      --ipxe-osInstallTreePath string Path of OS install tree for iPXE. (default: /var/www/html/)
      --kubernetes-apiversion uint    Sets the Kubernetes API version. Acceptable reference values:
                                        - 1 for Kubernetes versions < v1.15.0,
                                        - 2 for Kubernetes versions >= v1.15.0 && < v1.22.0,
                                        - 3 for Kubernetes versions >= v1.22.0
      --kubeversion string            Version of Kubernetes to deploy
      --libvirt-cidr string           CIDR for libvirt (default: 192.168.132.0/24)
      --libvirt-gateway string        Gateway for libvirt (default: 192.168.132.1)
      --libvirt-osPath string         OS path for libvirt
      --libvirt-uri string            URI for libvirt (default: qemu:///system)
      --master-cpu uint               CPU allocation for master nodes (units: cores)
      --master-disk uint              Disk size allocation for master nodes (units: GB)
      --master-hostname stringArray   Hostnames of master nodes (e.g., --master-hostname [master-01] --master-hostname [master-02] ...)
      --master-ips stringArray        IP addresses of master nodes (e.g., --master-ips [master-ip-01] --master-ips [master-ip-02] ...)
      --master-ram uint               RAM allocation for master nodes (units: MB)
      --network-plugin-url string     The deployment yaml URL of the network plugin
      --openstack-authURL string            AuthURL for openstack (default: http://controller:5000/v3)
      --openstack-availabilityZone string   AvailabilityZone for openstack (default: nova)
      --openstack-externalNetwork string    ExternalNetwork for openstack
      --openstack-glanceName string         GlanceName for openstack
      --openstack-internalNetwork string    InternalNetwork for openstack
      --openstack-password string           Password for openstack
      --openstack-region string       Region for openstack (default: RegionOne)
      --openstack-tenantName string   TenantName for openstack (default: admin)
      --openstack-username string     UserName for openstack (default: admin)
      --operator-image-url string     URL of the container image for the housekeeper operator component
      --os-type string                Operating system type for Kubernetes cluster deployment (e.g., nestos or generalos)
      --password string               Password for node login
      --pause-image string            Image for the pause container (e.g., pause:TAG)
      --platform string               Infrastructure platform for deploying the cluster (supports 'libvirt' or 'openstack')
      --pod-subnet string             Subnet used for Kubernetes Pods. (default: 10.244.0.0/16)
      --posthook-yaml string          Specify a YAML file or directory to apply after cluster deployment using 'kubectl apply'
      --prehook-script string         Specify a script file or directory to execute before cluster deployment as hooks
      --pxe-httpRootDir string        Root directory of HTTP server for PXE (default: /var/www/html/)
      --pxe-ip string                 IP address of local machine for PXE
      --pxe-tftpRootDir string        Root directory of TFTP server for PXE (default: /var/lib/tftpboot/)
      --release-image-url string      URL of the NestOS container image containing Kubernetes component
      --runtime string                Container runtime type (docker, isulad or crio)
      --service-subnet string         Subnet used by Kubernetes services. (default: 10.96.0.0/16)
      --sshkey string                 SSH key file path used for node authentication (default: ~/.ssh/id_rsa.pub)
      --token string                  Used to validate the cluster information obtained from the control plane, with non-control plane nodes used for joining the cluster
      --username string               User name for node login
      --worker-cpu uint               CPU allocation for worker nodes (units: cores)
      --worker-disk uint              Disk size allocation for worker nodes (units: GB)
      --worker-hostname stringArray   Hostnames of worker nodes (e.g., --worker-hostname [worker-01] --worker-hostname [worker-02] ...)
      --worker-ips stringArray        IP addresses of worker nodes (e.g., --worker-ips [worker-ip-01] --worker-ips [worker-ip-02] ...)
      --worker-ram uint               RAM allocation for worker nodes (units: MB)
  # Deploying the cluster with optional application configuration parameters
  $ nkd deploy --platform [platform] --master-ips [master-ip-01] --master-ips [master-ip-02] --master-hostname [master-hostname-01] --master-hostname [master-hostname-02] --master-cpu [master-cpu-cores] --worker-hostname [worker-hostname-01] --worker-disk [worker-disk-size] ...
  ```

## Deployment Process Demonstration

Adjusting Cluster Deployment Configuration Files
![](./figures/cluster_config.mp4)

Deploying the Cluster with Application Configuration Files
![](./figures/cluster_deploy.mp4)

## Image Building

* NestOS container images support building new container images based on the existing Dockerfile.
* Considerations for Making
    * Ensure Docker is installed
    * Download the latest version of the base image from the NestOS official website
    * For making deployment images, download the corresponding versions of kubeadm, kubelet, crictl binary files in advance and copy them to the /usr/bin directory. 
    * Installation of packages requires the use of the rpm-ostree command
 * Example Dockerfiles:
      ``` dockerfile
      FROM nestos_base_image
      COPY kube* /usr/bin/
      COPY crictl /usr/bin/
      RUN ostree container commit
      ```
Note: Users need to customize building deployment images before deploying the cluster.

## Create Cluster

 - Deploy the cluster with optional parameters. Example command:
    ``` shell
    $ nkd deploy --master-ips 192.168.132.11 --master-ips 192.168.132.12 --master-hostname k8s-master01 --master-hostname k8s-master02 --master-cpu 8 --worker-hostname k8s-worker01 --worker-disk 50 ...
    ```
 - Additionally, for more fine-grained configurations, you can deploy the cluster using a cluster configuration file. See configuration management for details.
    ``` shell
    $ nkd deploy -f cluster_config.yaml
    ```

## troubleshooting

The logs of NKD are stored in the directory /etc/nkd/logs by default, which facilitates effective troubleshooting in case of any issues encountered during the infrastructure creation process.