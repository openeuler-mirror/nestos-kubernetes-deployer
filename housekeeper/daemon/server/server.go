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

package server

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"housekeeper.io/pkg/common"
	pb "housekeeper.io/pkg/connection/proto"
	utilVersion "k8s.io/apimachinery/pkg/util/version"
)

const (
	ostreeImage            = "ostree-unverified-image:docker://"
	kubeadmCmd             = "/usr/local/bin/kubeadm"
	upgradeControlPlaneCmd = "/usr/local/bin/kubeadm upgrade apply -y"
	upgradeNodesCmd        = "/usr/local/bin/kubeadm upgrade node"
	kubeletUpdateCmd       = "systemctl daemon-reload && systemctl restart kubelet"
)

type Server struct {
	pb.UnimplementedUpgradeClusterServer
	mu sync.Mutex
}

// Implements the Upgrade
func (s *Server) Upgrade(_ context.Context, req *pb.UpgradeRequest) (*pb.UpgradeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	markFile := fmt.Sprintf("%s%s%s", "/var/housekeeper/", req.KubeVersion, ".stamp")
	// upgrade os
	if len(req.OsVersion) > 0 {
		//Checking for os version
		if err := checkOsVersion(req); err != nil {
			return &pb.UpgradeResponse{}, err
		}
	}
	if common.IsFileExist(markFile) {
		return &pb.UpgradeResponse{}, nil
	}
	// upgrade kubernetes
	if len(req.KubeVersion) > 0 {
		//Checking for kubernetes version
		if err := checkKubeVersion(req); err != nil {
			return &pb.UpgradeResponse{}, err
		}
	}
	// todo: check upgrade successfully
	if err := markNode(markFile); err != nil {
		logrus.Errorf("failed to mark node: %v", err)
		return &pb.UpgradeResponse{}, err
	}
	return &pb.UpgradeResponse{}, nil
}

func checkOsVersion(req *pb.UpgradeRequest) error {
	args := []string{"-c", "cat /etc/os-release | grep 'VERSION=' | head -n 1 | awk -F 'VERSION=' '{print $2}'"}
	osVersion, err := runCmd("/bin/sh", args...)
	if err != nil {
		logrus.Errorf("failed to get os version: %v", err)
		return err
	}
	if cmp, err := utilVersion.MustParseSemantic(string(osVersion)).Compare(req.OsVersion); err != nil {
		logrus.Errorf("failed to parse os version: %v", err)
		return err
	} else if cmp == 0 {
		logrus.Infof("The current os version %s and the desired upgrade version %s are the same",
			string(osVersion), req.OsVersion)
		return nil
	}
	//Compare the current os version with the desired version.
	//If different, the update command is executed
	if err := upgradeOSVersion(req); err != nil {
		logrus.Errorf("upgrade os version error: %v", err)
		return err
	}
	return nil
}

func checkKubeVersion(req *pb.UpgradeRequest) error {
	args := []string{"version", "-o", "short"}
	kubeadmVersion, err := runCmd(kubeadmCmd, args...)
	if err != nil {
		logrus.Errorf("kubeadm get version failed: %v", err)
		return err
	}
	if cmp, err := utilVersion.MustParseSemantic(string(kubeadmVersion)).Compare(req.KubeVersion); err != nil {
		logrus.Errorf("failed to parse kubeadm version: %v", err)
		return err
	} else if cmp == -1 {
		logrus.Infof("The request upgraded version %s is larger than kubeadm's version %s",
			req.KubeVersion, string(kubeadmVersion))
		return nil
	}
	//If the version of kubeadm is not less than the version requested for upgrade,
	//the upgrade command is executed
	if err := upgradeKubeVersion(req); err != nil {
		logrus.Errorf("upgrade kubernetes version error: %v", err)
		return err
	}
	return nil
}

func upgradeOSVersion(req *pb.UpgradeRequest) error {
	//upgrade os
	customImageURL := fmt.Sprintf("%s%s", ostreeImage, req.OsImageUrl)
	args := []string{"rebase", "--experimental", customImageURL, "--bypass-driver"}
	if _, err := runCmd("rpm-ostree", args...); err != nil {
		logrus.Errorf("failed to upgrade os to %s : %w", req.OsVersion, err)
		return err
	}
	// todo：skipping restart system
	if err := exec.Command("/bin/sh", "-c", "systemctl reboot").Run(); err != nil {
		logrus.Errorf("failed to run reboot: %v", err)
		return err
	}
	return nil
}

func upgradeKubeVersion(req *pb.UpgradeRequest) error {
	// todo: pull images before performing the upgrade
	kubeVersion := req.KubeVersion
	if ok, err := isControlPlaneNode(); err != nil {
		return err
	} else if ok {
		if err = upgradeControlPlaneNode(kubeVersion); err != nil {
			logrus.Errorf("failed to upgrade master nodes: %v", err)
			return err
		}
	} else {
		if err = upgradeNodes(); err != nil {
			logrus.Errorf("failed to upgrade worker nodes: %v", err)
			return err
		}
	}
	return nil
}

func upgradeControlPlaneNode(version string) error {
	args := []string{"-c", upgradeControlPlaneCmd, version}
	if err := exec.Command("/bin/sh", args...).Run(); err != nil {
		logrus.Errorf("failed to upgrade nodes: %w", err)
		return err
	}
	if err := exec.Command("/bin/sh", "-c", kubeletUpdateCmd).Run(); err != nil {
		logrus.Errorf("failed to restart kubelet: %w", err)
		return err
	}
	return nil
}

func upgradeNodes() error {
	if err := exec.Command("/bin/sh", "-c", upgradeNodesCmd).Run(); err != nil {
		logrus.Errorf("failed to upgrade nodes: %w", err)
		return err
	}
	if err := exec.Command("/bin/sh", "-c", kubeletUpdateCmd).Run(); err != nil {
		logrus.Errorf("failed to restart kubelet: %w", err)
		return err
	}
	return nil
}

func isControlPlaneNode() (bool, error) {
	if !common.IsFileExist("/etc/kubernetes/admin.conf") {
		return false, nil
	}
	ipArgs := []string{"-c", "ifconfig | grep 'inet' | grep 'broadcast'| awk '{print $2}'"}
	ipAddress, err := runCmd("/bin/sh", ipArgs...)
	if err != nil {
		return false, err
	}
	adminArgs := []string{"-c", "cat /etc/kubernetes/admin.conf | grep 'server' |grep -E -o '([0-9]{1,3}.){3}[0-9]{1,3}'"}
	ipControlPlane, err := runCmd("/bin/sh", adminArgs...)
	if err != nil {
		return false, err
	}
	return string(ipControlPlane) == string(ipAddress), nil
}

func markNode(file string) error {
	if err := os.MkdirAll("/var/housekeeper", 0644); err != nil {
		return err
	}
	args := []string{"-c", "touch", file}
	_, err := runCmd("/bin/sh", args...)
	if err != nil {
		return err
	}
	return nil
}

func runCmd(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		logrus.Errorf("error running  %s: %s: %w", name, strings.Join(args, " "), err)
		return nil, err
	}
	return output, nil
}
