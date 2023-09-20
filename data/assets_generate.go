//go:build tools
// +build tools

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
