/*
Copyright 2024 KylinSoft  Co., Ltd.

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

package osmanager

import (
	"errors"
	"nestos-kubernetes-deployer/pkg/cert"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/infra"
	"nestos-kubernetes-deployer/pkg/osmanager/bootconfig/ignition"

	"github.com/sirupsen/logrus"
)

type NestOS struct {
	conf         *asset.ClusterAsset
	certs        *cert.CertGenerator
	ignitionFile *ignition.Ignition
	infraMaster  *infra.Infra
	infraWorker  *infra.Infra
}

func NewNestOS(conf *asset.ClusterAsset) (*NestOS, error) {
	if conf == nil {
		return nil, errors.New("cluster asset config is nil")
	}
	if len(conf.Master) == 0 {
		return nil, errors.New("master node config is empty")
	}

	certGenerator := cert.NewCertGenerator(conf.Cluster_ID, &conf.Master[0])
	ignitionFile := ignition.NewIgnition(conf, configmanager.GetBootstrapIgnHostPort())
	return &NestOS{
		conf:         conf,
		certs:        certGenerator,
		ignitionFile: ignitionFile,
		infraMaster:  &infra.Infra{},
		infraWorker:  &infra.Infra{},
	}, nil
}

func (n *NestOS) GenerateResourceFiles() error {
	if err := n.certs.GenerateAllFiles(); err != nil {
		logrus.Errorf("Error generating all certs files: %v", err)
		return err
	}
	n.conf.CaCertHash = n.certs.CaCertHash

	if err := n.ignitionFile.GenerateBootConfig(); err != nil {
		logrus.Errorf("failed to generate ignition file: %v", err)
		return err
	}

	if err := n.infraMaster.Generate(n.conf, "master"); err != nil {
		logrus.Errorf("Failed to generate master terraform file")
		return err
	}
	if err := n.infraWorker.Generate(n.conf, "worker"); err != nil {
		logrus.Errorf("Failed to generate worker terraform file")
		return err
	}

	return nil
}
