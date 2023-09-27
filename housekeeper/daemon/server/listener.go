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
	"fmt"
	"net"
	"os"
	"path/filepath"
	"syscall"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	pb "housekeeper.io/pkg/connection/proto"
	"housekeeper.io/pkg/constants"
)

func NewListener(dir, name string) (l net.Listener, err error) {
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, err
	}

	addr := filepath.Join(dir, name)
	gid := os.Getgid()
	if err = syscall.Unlink(addr); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	const socketPermission = 0640
	mask := syscall.Umask(^socketPermission & int(os.ModePerm))
	defer syscall.Umask(mask)

	l, err = net.Listen("unix", addr)
	if err != nil {
		return nil, err
	}

	if err := os.Chown(addr, 0, gid); err != nil {
		if err := l.Close(); err != nil {
			return nil, fmt.Errorf("close listener error %w", err)
		}
		return nil, err
	}
	return l, nil
}

func Run() error {
	lis, err := NewListener(constants.SockDir, constants.SockName)
	if err != nil {
		logrus.Errorf("listen error: %v", err)
		return err
	}
	//get grpc server
	s := grpc.NewServer()
	pb.RegisterUpgradeClusterServer(s, &Server{})
	logrus.Info("housekeeper-daemon start serving")
	if err := s.Serve(lis); err != nil {
		logrus.Errorf("housekeeper-daemon server error: %v", err)
		return err
	}
	return nil
}
