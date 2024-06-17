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

package bootconfig

import "os"

type File struct {
	Node
	FileEmbedded1
}
type Node struct {
	Overwrite *bool  `json:"overwrite,omitempty"`
	Path      string `json:"path"`
}
type FileEmbedded1 struct {
	Contents Resource    `json:"contents,omitempty"`
	Mode     os.FileMode `json:"mode,omitempty"`
}
type Resource struct {
	Source []byte `json:"source,omitempty"`
}

type Systemd struct {
	Units []Unit `json:"units,omitempty"`
}
type Unit struct {
	Contents string `json:"contents,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty"`
	Name     string `json:"name"`
}
