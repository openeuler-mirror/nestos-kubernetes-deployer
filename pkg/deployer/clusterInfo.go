package deployer


type ClusterInfo struct {
	Name string `yaml:"name"`
	ApiEndpoint string `yaml:"api_endpoint"`
	Kubernetes Kubernetes `yaml:"kubernetes"`
	ClusterCIDR string `yaml:"cluster_cidr"`
	PodCIDR string `yaml:"pod_cidr"`
	ServiceCIDR string `yaml:"service_cidr"`
	InfraDriver string `yaml:"infra_driver"`
	OsType string `yaml:"os_type"`
	OsImage string `yaml:"os_image"`
	OsConfigList []OsConfig `yaml:"os_config_list"`
	EtcdType string `yaml:"etcd_type"`
}

type Kubernetes struct {
	KubernetesVersion string `yaml:"kubernetes_version"`
}

type OsConfig struct {
	Roles []string `yaml:"roles"`
	DiskSize string `yaml:"disk_size"`
	Image string `yaml:"image"`
	MemorySize string `yaml:"memory_size"`
	Cpus int `yaml:"cpus"`
}




func ParseFromFile(filename string) ClusterInfo{
   return ClusterInfo{}
}


func (* ClusterInfo)GenerateToFile(filename string) error{
	return nil
}