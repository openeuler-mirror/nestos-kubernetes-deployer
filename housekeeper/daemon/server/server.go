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
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	pb "housekeeper.io/pkg/connection/proto"
)

const ostreeImage = "ostree-unverified-image:docker://"

type Server struct {
	pb.UnimplementedUpgradeClusterServer
	mu sync.Mutex
}

// Implements the Upgrade
func (s *Server) Upgrade(_ context.Context, req *pb.UpgradeRequest) (*pb.UpgradeResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(req.OsVersion) > 0 {
		if err := upgradeOSVersion(req); err != nil {
			logrus.Errorf("upgrade os version error: %v", err)
			return &pb.UpgradeResponse{}, err
		}
	}
	return &pb.UpgradeResponse{}, nil
}

func upgradeOSVersion(req *pb.UpgradeRequest) error {
	//upgrade os
	customImageURL := fmt.Sprintf("%s%s", ostreeImage, req.OsImageUrl)
	args := []string{"rebase", "--experimental", customImageURL, "--bypass-driver"}
	if err := runCmd("rpm-ostree", args...); err != nil {
		logrus.Errorf("failed to upgrade os to %s : %w", req.OsVersion, err)
		return err
	}
	// todo：skipping reboot
	rebootArgs := []string{"-c", "systemctl reboot"}
	if err := runCmd("/bin/sh", rebootArgs...); err != nil {
		logrus.Errorf("failed to run reboot: %v", err)
		return err
	}
	return nil
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running %s %s: %s: %w", name, strings.Join(args, " "), string(stderr.Bytes()), err)
	}
	return nil
}
