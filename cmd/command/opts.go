package command

var (
	RootOptDir string
)

// 部署集群可选配置参数集合
var ClusterOpts struct {
	ClusterId string
	GatherDeployOpts
}

type GatherDeployOpts struct {
	SSHKey   string
	Platform string
	//
}
