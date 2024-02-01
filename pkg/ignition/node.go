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
package ignition

import (
	"os"
	"path/filepath"

	"github.com/clarketm/json"
	"github.com/sirupsen/logrus"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/vincent-petithory/dataurl"
)

func Marshal(input interface{}) ([]byte, error) {
	return json.Marshal(input)
}

/*
FileWithContents creates an ignition file with the given contents.
Parameters:
  - path (string): The file path.
  - mode (int): The file permissions.
  - contents ([]byte): The file content as a byte slice.

Returns:
  - igntypes.File: Ignition file configuration.
*/
func FileWithContents(path string, mode int, contents []byte) igntypes.File {
	return igntypes.File{
		Node: igntypes.Node{
			Path:      path,
			Overwrite: ignutil.BoolToPtr(true),
		},
		FileEmbedded1: igntypes.FileEmbedded1{
			Mode: &mode,
			Contents: igntypes.Resource{
				Source: ignutil.StrToPtr(dataurl.EncodeBytes(contents)),
			},
		},
	}
}

func AppendFiles(files []igntypes.File, file igntypes.File) []igntypes.File {
	for i, f := range files {
		if f.Node.Path == file.Node.Path {
			files[i] = file
			return files
		}
	}
	files = append(files, file)
	return files
}

/*
Save the ignition config
Parameters:
config - the ignition config to be saved
filePath - the path to save the file
fileName - the name to save the file
*/
func SaveFile(config *igntypes.Config, filePath string, fileName string) error {
	data, err := Marshal(config)
	if err != nil {
		logrus.Errorf("failed to Marshal ignition config: %v", err)
		return err
	}
	fullPath := filepath.Join(filePath, fileName)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0750); err != nil {
		logrus.Errorf("failed to Mkdir: %v", err)
		return err
	}
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		logrus.Errorf("failed to save ignition file: %v", err)
		return err
	}
	return nil
}
