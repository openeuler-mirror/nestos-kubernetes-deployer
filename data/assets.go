//go:build !release
// +build !release

//go:generate go run assets_generate.go

/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * nestos-kubernetes-deployer licensed under the Apache License, Version 2.0.
 * See LICENSE file for more details.
 * Author: liukuo <liukuo@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */

package data

import (
	"net/http"
)

var Assets http.FileSystem

func init() {
	dir := "data"
	Assets = http.Dir(dir)
}
