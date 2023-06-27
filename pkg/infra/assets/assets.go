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

package assets

// path : contents
type Assets map[string][]byte

func (a Assets) ToDir(dirname string) error {
	return nil
}

func (a *Assets) Merge(b Assets) *Assets {
	return a
}

type AssetsGenerator interface {
	GenerateAssets() Assets
}

type File struct {
	Filename string
	Data     []byte
}
