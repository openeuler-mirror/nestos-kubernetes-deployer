//go:build !release
// +build !release

//go:generate go run assets_generate.go

package data

import (
	"net/http"
)

var Assets http.FileSystem

func init() {
	dir := "data"
	Assets = http.Dir(dir)
}
