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

package utils

import (
	"fmt"
	"net"
	"os/user"
	"path/filepath"
	"strings"
)

func GetDefaultPubKeyPath() string {
	return filepath.Join(getSysHome(), ".ssh", "id_rsa.pub")
}

func GetApiServerEndpoint(ip string) string {
	return fmt.Sprintf("%s:%s", ip, "6443")
}

// GetLocalIP retrieves the local IP address
func GetLocalIP() (string, error) {
	// Retrieve route information
	routeOutput, err := RunCommand("ip -o route get 255.0 2>/dev/null")
	if err != nil {
		return "", err
	}

	// Use sed to extract the source IP address
	cmd := "sed -e 's/.*src \\([^ ]*\\).*/\\1/'"
	ipOutput, err := RunCommand(fmt.Sprintf("echo '%s' | %s", routeOutput, cmd))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(ipOutput), nil
}

func getSysHome() string {
	if user, err := user.Current(); err == nil {
		return user.HomeDir
	}
	return "/root"
}

func IsPortOpen(port string) bool {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}
