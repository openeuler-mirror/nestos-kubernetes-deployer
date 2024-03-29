# Cluster config file description

For the NestOS image download address, see[website](https://nestos.openeuler.org/)
``` shell
cluster_id: cluster                                 # cluster name
architecture: amd64                                 # deploy cluster architecture, support amd64 or arm64
platform: libvirt                                   # deployment platform is libvirt
infraplatform
  uri: qemu:///system                                
  osimage: https://nestos.org.cn/nestos20230928/nestos-for-container/x86_64/NestOS-For-Container-22.03-LTS-SP2.20230928.0-qemu.{arch}.qcow2                                             # image URL，support amd64 or arm64
  cidr: 192.168.132.0/24
  gateway: 192.168.132.1
username: root                                      # Specify the username for ssh login
password: $1$yoursalt$UGhjCXAJKpWWpeN8xsF.c/        # Specify the password for ssh login
sshkey: "/root/.ssh/id_rsa.pub"                     # The storage path of the ssh-key file
master:                                             # master config
- hostname: k8s-master01
  hardwareinfo:                                     
    cpu: 4
    ram: 8192                                       
    disk: 50                                        
  ip: "192.168.132.11"                              
worker:                                             # worker config
- hostname: k8s-worker01            
  hardwareinfo:
    cpu: 4
    ram: 8192
    disk: 50
  ip: ""                                            # If the worker node IP address is not set, it will be automatically assigned by dhcp and will be empty by default.
runtime: isulad                                     # support docker、isulad、crio
kubernetes:                                         
  kubernetes-version: "v1.23.10"                   
  kubernetes-apiversion: "v1beta3"                  # support v1beta3、v1beta2、v1beta1
  apiserver-endpoint: "192.168.132.11:6443"          
  image-registry: "k8s.gcr.io"                     
  pause-image: "pause:3.6"                         
  release-image-url: "hub.oepkgs.net/nestos/nestos:22.03-LTS-SP2.20230928.0-{arch}-k8s-v1.23.10"                         
  token: ""                                         # automatically generated by default
  adminkubeconfig: /etc/nkd/cluster/admin.config    # path of admin.conf
  certificatekey: ""                                # The key used to decrypt the certificate in the downloaded Secret when adding a new control plane node
  network:                                          
    service-subnet: "10.96.0.0/16"                  
    pod-subnet: "10.244.0.0/16"                     
    plugin: https://projectcalico.docs.tigera.io/archive/v3.22/manifests/calico.yaml # network plugin
housekeeper:                                                                                          # housekeeper
  deployhousekeeper: false                                                                           
  operatorimageurl: "hub.oepkgs.net/nestos/housekeeper/{arch}/housekeeper-operator-manager:{tag}"     # housekeeper-operator image URL
  controllerimageurl: "hub.oepkgs.net/nestos/housekeeper/{arch}/housekeeper-controller-manager:{tag}" # housekeeper-controller image URL  
certasset:                                          # Configure user-defined certificate file path list, automatically generated by default
  rootcacertpath: ""                
  rootcakeypath: ""
  etcdcacertpath: ""
  etcdcakeypath: ""
  frontproxycacertpath: ""
  frontproxycakeypath: ""
  sapub: ""
  sakey: ""
```

To set the deployment platform to openstack, you need to reset the "infraplatform" field configuration parameters.
``` shell
platform: openstack                                   
infraplatform                      
	username:                                           # openstack username, requires permission to create resources                                      
	password:                                           # openstack login password, used to log in to the openstack platform
	tenant_name:                                        # openstack tenant name, the collection the user belongs to, for example: admin
	auth_url:                                           # openstack auth_url，example：http://{ip}:{port}/v3
	region:                                             # Used for resource isolation, for example: RegionOne
	internal_network:                                   
	external_network:                                  
	glance_name:                                        # qcow2 image
	availability_zone:                                  # default nova
```