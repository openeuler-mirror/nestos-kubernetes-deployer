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

package generalos

import (
	"errors"
	"nestos-kubernetes-deployer/pkg/cert"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/osmanager/bootconfig/cloudinit"
	"nestos-kubernetes-deployer/pkg/osmanager/bootconfig/kickstart"
	"nestos-kubernetes-deployer/pkg/terraform"
	"strings"

	"github.com/sirupsen/logrus"
)

type GeneralOS struct {
	conf          *asset.ClusterAsset
	certs         *cert.CertGenerator
	cloudinitFile *cloudinit.Cloudinit
	kickstartFile *kickstart.Kickstart
	infraMaster   *terraform.Infra
	infraWorker   *terraform.Infra
}

func NewGeneralOS(conf *asset.ClusterAsset) (*GeneralOS, error) {
	if conf == nil {
		return nil, errors.New("cluster asset config is nil")
	}
	if len(conf.Master) == 0 {
		return nil, errors.New("master node config is empty")
	}

	certGenerator := cert.NewCertGenerator(conf.ClusterID, &conf.Master[0])
	cloudinitFile := cloudinit.NewCloudinit(conf, configmanager.GetBootstrapIgnHostPort())
	kickstartFile := kickstart.NewKickstart(conf, configmanager.GetBootstrapIgnHostPort())
	return &GeneralOS{
		conf:          conf,
		certs:         certGenerator,
		cloudinitFile: cloudinitFile,
		kickstartFile: kickstartFile,
		infraMaster:   &terraform.Infra{},
		infraWorker:   &terraform.Infra{},
	}, nil
}

func (g *GeneralOS) GenerateResourceFiles() error {
	if err := g.certs.GenerateAllFiles(); err != nil {
		logrus.Errorf("Error generating all certs files: %v", err)
		return err
	}
	g.conf.CaCertHash = g.certs.CaCertHash

	switch strings.ToLower(g.conf.Platform) {
	case "libvirt", "openstack":
		if err := g.cloudinitFile.GenerateBootConfig(); err != nil {
			logrus.Errorf("failed to generate cloudinit file: %v", err)
			return err
		}

		if err := g.infraMaster.Generate(g.conf, "master"); err != nil {
			logrus.Errorf("Failed to generate master terraform file")
			return err
		}
		if err := g.infraWorker.Generate(g.conf, "worker"); err != nil {
			logrus.Errorf("Failed to generate worker terraform file")
			return err
		}
	case "pxe", "ipxe":
		if err := g.kickstartFile.GenerateBootConfig(); err != nil {
			logrus.Errorf("failed to generate kickstart file: %v", err)
			return err
		}
	}

	return nil
}
