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

package nkd

import (
	"time"
)

type System struct {
	Count    int
	HostName string
	Ips      []string
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSHKey   string
}

type Openstack struct {
	User_name        string
	Password         string
	Tenant_name      string
	Auth_url         string
	Region           string
	Internal_network string
	External_network string
	Glance           string
	Flavor           string
}

type Libvirt struct {
}

type Size struct {
	Vcpus int
	Ram   int
	Disk  int
}

type Infra struct {
	Platform  string
	Openstack Openstack
	Vmsize    Size
}

type Cluster struct {
	Name string `yaml:"cluster-name"`
}

type Repo struct {
	Secret   []map[string]string `yaml:"secret"`
	Registry string              `yaml:"registry"`
}

type TypeMeta struct {
	// Kind       string
	ApiVersion string
}

type BootstrapTokenString struct {
	ID     string `json:"-"`
	Secret string `json:"-" datapolicy:"token"`
}

type Duration struct {
	time.Duration `protobuf:"varint,1,opt,name=duration,casttype=time.Duration"`
}

type APIEndpoint struct {
	AdvertiseAddress string
	BindPort         int32
}

type BootstrapToken struct {
	Token  string
	Groups []string
	TTL    string
	Usages []string
}
type TaintEffect string

type PullPolicy string

type Time struct {
	time.Time `protobuf:"-"`
}
type Taint struct {
	// Required. The taint key to be applied to a node.
	Key string `json:"key" protobuf:"bytes,1,opt,name=key"`
	// The taint value corresponding to the taint key.
	// +optional
	Value string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
	// Required. The effect of the taint on pods
	// that do not tolerate the taint.
	// Valid effects are NoSchedule, PreferNoSchedule and NoExecute.
	Effect TaintEffect `json:"effect" protobuf:"bytes,3,opt,name=effect,casttype=TaintEffect"`
	// TimeAdded represents the time at which the taint was added.
	// It is only written for NoExecute taints.
	// +optional
	TimeAdded *Time `json:"timeAdded,omitempty" protobuf:"bytes,4,opt,name=timeAdded"`
}
type NodeRegistrationOptions struct {
	Name                  string
	CRISocket             string
	Taints                []Taint
	KubeletExtraArgs      map[string]string
	IgnorePreflightErrors []string
	ImagePullPolicy       PullPolicy `json:"imagePullPolicy,omitempty"`
}

type ComponentConfigMap map[string]ComponentConfig

type ComponentConfig interface {
	// DeepCopy should create a new deep copy of the component config in place
	DeepCopy() ComponentConfig

	// Marshal is marshalling the config into a YAML document returned as a byte slice
	Marshal() ([]byte, error)

	// Unmarshal loads the config from a document map. No config in the document map is no error.

	// Default patches the component config with kubeadm preferred defaults
	Default(cfg *ClusterConfiguration, localAPIEndpoint *APIEndpoint, nodeRegOpts *NodeRegistrationOptions)

	// IsUserSupplied indicates if the component config was supplied or modified by a user or was kubeadm generated
	IsUserSupplied() bool

	// SetUserSupplied sets the state of the component config "user supplied" flag to, either true, or false.
	SetUserSupplied(userSupplied bool)

	// Mutate allows applying pre-defined modifications to the config before it's marshaled.
	Mutate() error

	// Set can be used to set the internal configuration in the ComponentConfig
	Set(interface{})

	// Get can be used to get the internal configuration in the ComponentConfig
	Get() interface{}
}

type Etcd struct {

	// Local provides configuration knobs for configuring the local etcd instance
	// Local and External are mutually exclusive
	Local *LocalEtcd

	// External describes how to connect to an external etcd cluster
	// Local and External are mutually exclusive
	External *ExternalEtcd
}

type ImageMeta struct {
	// ImageRepository sets the container registry to pull images from.
	// if not set, the ImageRepository defined in ClusterConfiguration will be used instead.
	ImageRepository string

	// ImageTag allows to specify a tag for the image.
	// In case this value is set, kubeadm does not change automatically the version of the above components during upgrades.
	ImageTag string

	//TODO: evaluate if we need also a ImageName based on user feedbacks
}

// LocalEtcd describes that kubeadm should run an etcd cluster locally
type LocalEtcd struct {
	// ImageMeta allows to customize the container used for etcd
	ImageMeta `json:",inline"`

	// DataDir is the directory etcd will place its data.
	// Defaults to "/var/lib/etcd".
	DataDir string

	// ExtraArgs are extra arguments provided to the etcd binary
	// when run inside a static pod.
	// A key in this map is the flag name as it appears on the
	// command line except without leading dash(es).
	ExtraArgs map[string]string

	// ServerCertSANs sets extra Subject Alternative Names for the etcd server signing cert.
	ServerCertSANs []string
	// PeerCertSANs sets extra Subject Alternative Names for the etcd peer signing cert.
	PeerCertSANs []string
}

// ExternalEtcd describes an external etcd cluster
type ExternalEtcd struct {

	// Endpoints of etcd members. Useful for using external etcd.
	// If not provided, kubeadm will run etcd in a static pod.
	Endpoints []string
	// CAFile is an SSL Certificate Authority file used to secure etcd communication.
	CAFile string
	// CertFile is an SSL certification file used to secure etcd communication.
	CertFile string
	// KeyFile is an SSL key file used to secure etcd communication.
	KeyFile string
}

type Networking struct {
	// ServiceSubnet is the subnet used by k8s services. Defaults to "10.96.0.0/12".
	ServiceSubnet string
	// PodSubnet is the subnet used by pods.
	PodSubnet string
	// DNSDomain is the dns domain used by k8s services. Defaults to "cluster.local".
	DNSDomain string
}
type HostPathType string
type HostPathMount struct {
	// Name of the volume inside the pod template.
	Name string
	// HostPath is the path in the host that will be mounted inside
	// the pod.
	HostPath string
	// MountPath is the path inside the pod where hostPath will be mounted.
	MountPath string
	// ReadOnly controls write access to the volume
	ReadOnly bool
	// PathType is the type of the HostPath.
	PathType HostPathType
}
type ControlPlaneComponent struct {
	// ExtraArgs is an extra set of flags to pass to the control plane component.
	// A key in this map is the flag name as it appears on the
	// command line except without leading dash(es).
	// TODO: This is temporary and ideally we would like to switch all components to
	// use ComponentConfig + ConfigMaps.
	ExtraArgs map[string]string

	// ExtraVolumes is an extra set of host volumes, mounted to the control plane component.
	ExtraVolumes []HostPathMount
}
type APIServer struct {
	ControlPlaneComponent

	// CertSANs sets extra Subject Alternative Names for the API Server signing cert.
	CertSANs []string

	// TimeoutForControlPlane controls the timeout that we use for API server to appear
	TimeoutForControlPlane string
}
type DNS struct {
	// ImageMeta allows to customize the image used for the DNS component
	ImageMeta `json:",inline"`
}
type ClusterConfiguration struct {
	TypeMeta

	// ComponentConfigs holds component configs known to kubeadm, should long-term only exist in the internal kubeadm API
	// +k8s:conversion-gen=false
	ComponentConfigs ComponentConfigMap

	// Etcd holds configuration for etcd.
	Etcd Etcd

	// Networking holds configuration for the networking topology of the cluster.
	Networking Networking

	// KubernetesVersion is the target version of the control plane.
	KubernetesVersion string

	// CIKubernetesVersion is the target CI version of the control plane.
	// Useful for running kubeadm with CI Kubernetes version.
	// +k8s:conversion-gen=false
	CIKubernetesVersion string

	// ControlPlaneEndpoint sets a stable IP address or DNS name for the control plane; it
	// can be a valid IP address or a RFC-1123 DNS subdomain, both with optional TCP port.
	// In case the ControlPlaneEndpoint is not specified, the AdvertiseAddress + BindPort
	// are used; in case the ControlPlaneEndpoint is specified but without a TCP port,
	// the BindPort is used.
	// Possible usages are:
	// e.g. In a cluster with more than one control plane instances, this field should be
	// assigned the address of the external load balancer in front of the
	// control plane instances.
	// e.g.  in environments with enforced node recycling, the ControlPlaneEndpoint
	// could be used for assigning a stable DNS to the control plane.
	ControlPlaneEndpoint string

	// APIServer contains extra settings for the API server control plane component
	APIServer APIServer

	// ControllerManager contains extra settings for the controller manager control plane component
	ControllerManager ControlPlaneComponent

	// Scheduler contains extra settings for the scheduler control plane component
	Scheduler ControlPlaneComponent

	// DNS defines the options for the DNS add-on installed in the cluster.
	DNS DNS

	// CertificatesDir specifies where to store or look for all required certificates.
	CertificatesDir string

	// ImageRepository sets the container registry to pull images from.
	// If empty, `registry.k8s.io` will be used by default; in case of kubernetes version is a CI build (kubernetes version starts with `ci/`)
	// `gcr.io/k8s-staging-ci-images` will be used as a default for control plane components and for kube-proxy, while `registry.k8s.io`
	// will be used for all the other images.
	ImageRepository string

	// CIImageRepository is the container registry for core images generated by CI.
	// Useful for running kubeadm with images from CI builds.
	// +k8s:conversion-gen=false
	CIImageRepository string

	// FeatureGates enabled by the user.
	FeatureGates map[string]bool

	// The cluster name
	ClusterName string
}

type Patches struct {
	// Directory is a path to a directory that contains files named "target[suffix][+patchtype].extension".
	// For example, "kube-apiserver0+merge.yaml" or just "etcd.json". "target" can be one of
	// "kube-apiserver", "kube-controller-manager", "kube-scheduler", "etcd", "kubeletconfiguration".
	// "patchtype" can be one of "strategic" "merge" or "json" and they match the patch formats supported by kubectl.
	// The default "patchtype" is "strategic". "extension" must be either "json" or "yaml".
	// "suffix" is an optional string that can be used to determine which patches are applied
	// first alpha-numerically.
	Directory string
}

type Kubeadm struct {
	ClusterConfiguration
	TypeMeta
	BootstrapTokens  []BootstrapToken
	LocalAPIEndpoint APIEndpoint
	NodeRegistration NodeRegistrationOptions
	CertificateKey   string
	Patches          *Patches
}

type Addon struct {
	Addons []map[string]string
}

// type Nkd struct {
// 	Cluster Cluster
// 	System  []System
// 	Repo    Repo
// 	Kubeadm Kubeadm
// 	Addon   Addon
// 	Infra   Infra
// }

type Master struct {
	Node    string
	Cluster Cluster
	System  System
	Repo    Repo
	Kubeadm Kubeadm
	Addon   Addon
	Infra   Infra
}

type WorkerK8s struct {
	Discovery        Discovery
	CaCertPath       string
	NodeRegistration NodeRegistrationOptions
}

type BootstrapTokenDiscovery struct {
	Token                    string
	APIServerEndpoint        string
	UnsafeSkipCAVerification bool
}

type Discovery struct {
	BootstrapToken    *BootstrapTokenDiscovery
	Timeout           string
	TlsBootstrapToken string
}

type Worker struct {
	Node   string
	Repo   Repo
	System System
	Infra  Infra
	Addon  Addon
	Worker WorkerK8s
}

type Node struct {
	Node string
}
