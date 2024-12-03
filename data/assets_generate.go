//go:build tools
// +build tools

/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * nestos-kubernetes-deployer licensed under the Apache License, Version 2.0.
 * See LICENSE file for more details.
 * Author: liukuo <liukuo@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */

package main

import (
	"nestos-kubernetes-deployer/data"

	"github.com/shurcooL/vfsgen"
	"github.com/sirupsen/logrus"
)

func main() {
	packageName := "data"  // 包名
	buildTags := "release" // 构建标签
	fsName := "Assets"

	err := vfsgen.Generate(data.Assets, vfsgen.Options{
		PackageName:  packageName,
		BuildTags:    buildTags,
		VariableName: fsName,
	})
	if err != nil {
		logrus.Fatalln(err)
	}
}
