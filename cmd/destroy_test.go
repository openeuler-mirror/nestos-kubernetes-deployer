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
package cmd

import (
	"github.com/agiledragon/gomonkey/v2"
	"nestos-kubernetes-deployer/cmd/command/opts"
	"nestos-kubernetes-deployer/pkg/configmanager"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/asset/infraasset"
	"os"
	"testing"
)

func TestDestroy(t *testing.T) {
	str := `logLevel: default log level
clusterConfigPath: ""
persistdir: /etc/nkd
bootstrapurl:
  bootstrapIgnHost: 10.44.55.21
  bootstrapIgnPort: "9080"
`
	opts.Opts.RootOptDir = "../data"
	cmd := NewDestroyCommand()

	args := []string{"--cluster-id", "cluster"}
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Logf("Failed to execute command: %v", err)
	}

	if err := writeToFile(opts.Opts.RootOptDir+"/global_config.yaml", str); err != nil {
		t.Logf("writeToFile : %v", err)
	}
	cc := &asset.ClusterAsset{
		ClusterID:    "cluster",
		Architecture: "amd64",
		Platform:     "pxe",
		InfraPlatform: &infraasset.PXEAsset{
			IP:             "",
			HTTPServerPort: "10",
			HTTPRootDir:    "./",
			TFTPServerPort: "20",
			TFTPRootDir:    "./",
		},
		OSImage:  asset.OSImage{Type: "nestos"},
		UserName: "root",
		SSHKey:   "./test.yaml",
		Password: "123",
		Master: []asset.NodeAsset{
			{
				Hostname: "k8s-master01",
				IP:       "192.168.132.11",
				HardwareInfo: asset.HardwareInfo{
					CPU:  2,
					RAM:  2048,
					Disk: 30,
				},
			},
		},
		Worker: []asset.NodeAsset{
			{
				Hostname: "k8s-worker01",
				IP:       "192.168.132.12",
				HardwareInfo: asset.HardwareInfo{
					CPU:  2,
					RAM:  2048,
					Disk: 30,
				},
			},
		},
		Runtime: "crio",
		Kubernetes: asset.Kubernetes{
			KubernetesVersion:    "v1.29.1",
			KubernetesAPIVersion: "v1beta3",
			ApiServerEndpoint:    "192.168.132.11:6443",
			ImageRegistry:        "registry.k8s.io",
			PauseImage:           "pause:3.9",
			Network: asset.Network{
				ServiceSubnet: "10.96.0.0/16",
				PodSubnet:     "10.244.0.0/16",
			},
		},
	}
	configmanager.ClusterAsset = map[string]*asset.ClusterAsset{
		"cluster": cc,
	}
	p := gomonkey.ApplyFunc(configmanager.Initial, func(*opts.OptionsList) error {
		return nil
	})
	defer p.Reset()

	p2 := gomonkey.ApplyFunc(configmanager.GetClusterConfig, func(string) (*asset.ClusterAsset, error) {
		return cc, nil
	})
	defer p2.Reset()

	t.Run("DestroyCmd Fail", func(t *testing.T) {
		if err := runDestroyCmd(cmd, args); err != nil {
			t.Logf("Expected error, got : %v", err)
		}
		if err := os.RemoveAll("logs"); err != nil {
			t.Logf("Failed to remove logs folder: %v", err)
		}
		if err := os.RemoveAll("logs"); err != nil {
			t.Logf("Failed to remove logs folder: %v", err)
		}
		if err := os.RemoveAll(opts.Opts.RootOptDir + "/global_config.yaml"); err != nil {
			t.Logf("Failed to remove logs folder: %v", err)
		}
	})

	t.Run("DestroyCmd Fail", func(t *testing.T) {
		cc.Platform = "ipxe"
		cc.InfraPlatform = &infraasset.IPXEAsset{
			IP:   "",
			Port: "101",
		}
		mData := map[string]interface{}{
			"libvirt": &infraasset.LibvirtAsset{
				URI:     "www.a.com",
				OSPath:  "a.yaml",
				CIDR:    "1.1.1.1",
				Gateway: "1.1.1.1",
			},
			"openstack": &infraasset.OpenStackAsset{
				UserName: "zhangs",
			},
			"ipxe": &infraasset.IPXEAsset{
				IP:   "",
				Port: "101",
			},
		}

		for k, v := range mData {
			cc.Platform = k
			cc.InfraPlatform = v
			if err := runDestroyCmd(cmd, args); err != nil {
				t.Logf("Expected error, got : %v", err)
			}
			if err := os.RemoveAll("logs"); err != nil {
				t.Logf("Failed to remove logs folder: %v", err)
			}
			if err := os.RemoveAll("logs"); err != nil {
				t.Logf("Failed to remove logs folder: %v", err)
			}
			if err := os.RemoveAll(opts.Opts.RootOptDir + "/global_config.yaml"); err != nil {
				t.Logf("Failed to remove logs folder: %v", err)
			}
			configmanager.ClusterAsset = map[string]*asset.ClusterAsset{
				"cluster": cc,
			}
		}

	})

}

func writeToFile(path string, content string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
