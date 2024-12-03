/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * nestos-kubernetes-deployer licensed under the Apache License, Version 2.0.
 * See LICENSE file for more details.
 * Author: liukuo <liukuo@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */

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
