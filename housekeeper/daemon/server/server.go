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
	"os/exec"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"housekeeper.io/pkg/common"
	pb "housekeeper.io/pkg/connection/proto"
)

const (
	ostreeImage      = "ostree-unverified-image:docker://"
	kubeadmCmd       = "/usr/bin/kubeadm"
	upgradeMasterCmd = "/usr/bin/kubeadm upgrade apply -y"
	upgradeWorkerCmd = "/usr/bin/kubeadm upgrade node"
	kubeletUpdateCmd = "systemctl daemon-reload && systemctl restart kubelet"
	adminFile        = "/etc/kubernetes/admin.conf"
)

type Server struct {
	pb.UnimplementedUpgradeClusterServer
	mu sync.Mutex
}

// Implements the Upgrade
func (s *Server) Upgrade(_ context.Context, req *pb.UpgradeRequest) (*pb.UpgradeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// upgrade os
	if len(req.OsVersion) > 0 {
		//Checking for os version
		if err := checkOsVersion(req); err != nil {
			return &pb.UpgradeResponse{}, err
		}
	}
	// upgrade kubernetes
	if len(req.KubeVersion) > 0 {
		markFile := fmt.Sprintf("%s%s%s", "/run/housekeeper-daemon/", req.KubeVersion, ".stamp")
		if common.IsFileExist(markFile) {
			return &pb.UpgradeResponse{}, nil
		}
		//Checking for kubernetes version
		if err := checkKubeVersion(req); err != nil {
			return &pb.UpgradeResponse{}, err
		}
		// todo: check upgrade successfully
		if err := markNode(markFile); err != nil {
			logrus.Errorf("failed to mark node: %v", err)
			return &pb.UpgradeResponse{}, err
		}
	}
	return &pb.UpgradeResponse{}, nil
}

func checkOsVersion(req *pb.UpgradeRequest) error {
	osVersionBytes, err := runCmd("awk", "-F=", "/PRETTY_NAME/ {gsub(/\"/, \"\", $2); print $2}", "/etc/os-release")
	if err != nil {
		logrus.Errorf("failed to get os version: %v", err)
		return err
	}
	osVersion := strings.TrimSpace(string(osVersionBytes))
	reqOsVersion := strings.TrimSpace(req.OsVersion)
	if osVersion == reqOsVersion {
		logrus.Infof("The current OS version %s and the desired upgrade version %s are the same", osVersion, reqOsVersion)
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
	kubeadmVersionBytes, err := runCmd(kubeadmCmd, args...)
	if err != nil {
		logrus.Errorf("kubeadm get version failed: %v", err)
		return err
	}
	kubeadmVersion := strings.TrimSpace(string(kubeadmVersionBytes))
	KubeVersion := strings.TrimSpace(req.KubeVersion)
	if kubeadmVersion == KubeVersion {
		logrus.Infof("The current k8s version %s and the desired upgrade version %s are the same", string(kubeadmVersion), req.KubeVersion)
		return nil
	}
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
	if isMasterNode() {
		if err := upgradeMasterNodes(req.KubeVersion); err != nil {
			logrus.Errorf("failed to upgrade master nodes: %v", err)
			return err
		}
	} else {
		if err := upgradeWorkerNodes(); err != nil {
			logrus.Errorf("failed to upgrade worker nodes: %v", err)
			return err
		}
	}
	return nil
}

func upgradeMasterNodes(version string) error {
	if err := exec.Command("/bin/sh", "-c", kubeletUpdateCmd).Run(); err != nil {
		logrus.Errorf("failed to restart kubelet: %w", err)
		return err
	}
	args := []string{"-c", upgradeMasterCmd, version}
	if err := exec.Command("/bin/sh", args...).Run(); err != nil {
		logrus.Errorf("failed to upgrade nodes: %w", err)
		return err
	}
	return nil
}

func upgradeWorkerNodes() error {
	if err := exec.Command("/bin/sh", "-c", kubeletUpdateCmd).Run(); err != nil {
		logrus.Errorf("failed to restart kubelet: %w", err)
		return err
	}
	if err := exec.Command("/bin/sh", "-c", upgradeWorkerCmd).Run(); err != nil {
		logrus.Errorf("failed to upgrade nodes: %w", err)
		return err
	}
	return nil
}

func isMasterNode() bool {
	return common.IsFileExist(adminFile)
}

func markNode(file string) error {
	_, err := runCmd("touch", file)
	if err != nil {
		logrus.Errorf("failed to create mark file: %v", err)
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
