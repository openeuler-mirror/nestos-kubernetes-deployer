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
package cloudinit

import "os"

type CloudinitConfig struct {
	SSHPasswordAuth   bool           `yaml:"ssh_pwauth"`
	SSHAuthorizedKeys []string       `yaml:"ssh_authorized_keys,omitempty"`
	Chpasswd          ChpasswdConfig `yaml:"chpasswd,omitempty"`
	WriteFiles        []WriteFile    `yaml:"write_files,omitempty"`
	RunCmds           []interface{}  `yaml:"runcmd,omitempty"`
}

type ChpasswdConfig struct {
	List   string `yaml:"list,omitempty"`
	Expire bool   `yaml:"expire"`
}

type WriteFile struct {
	EnCoding    string      `yaml:"encoding,omitempty"`
	Content     string      `yaml:"content,omitempty"`
	Path        string      `yaml:"path,omitempty"`
	Permissions os.FileMode `yaml:"permissions,omitempty"`
}
