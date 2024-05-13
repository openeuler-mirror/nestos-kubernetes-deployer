package data_test

import (
	"nestos-kubernetes-deployer/data"
	"testing"
)

func TestOpenFile(t *testing.T) {
	file, err := data.Assets.Open("/bootconfig/systemd/init-cluster.service")
	if err != nil {
		t.Errorf("Failed to open file: %v", err)
	}
	defer file.Close()
}
