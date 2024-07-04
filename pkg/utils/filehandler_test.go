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
package utils

import (
	"nestos-kubernetes-deployer/data"
	"nestos-kubernetes-deployer/pkg/constants"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestOsmanager(t *testing.T) {
	// Change to the data directory where the test files are located
	err := os.Chdir("../../data")
	if err != nil {
		t.Fatal(err)
	}

	url := filepath.Join("bootconfig/systemd", constants.InitClusterService)
	tmplData := map[string]interface{}{
		"Key": "Value", // Adjust according to the template requirements
	}
	t.Run("FetchAndUnmarshalUrl Success", func(t *testing.T) {
		_, err := FetchAndUnmarshalUrl(url, tmplData)
		if err != nil {
			t.Error("test fail", err)
			return
		}

	})

	t.Run("FetchAndUnmarshalUrl Fail", func(t *testing.T) {
		_, err := FetchAndUnmarshalUrl("", "")
		if err == nil {
			t.Error("expected failure, got success")
		}
	})

	t.Run("GetCompleteFile Success", func(t *testing.T) {
		file, err := data.Assets.Open(url)
		if err != nil {
			logrus.Errorf("Error opening file %s: %v\n", url, err)
			return
		}
		defer file.Close()
		GetCompleteFile(url, file, tmplData)
	})

	t.Run("GetCompleteFile Fail", func(t *testing.T) {
		file, err := data.Assets.Open(url)
		if err != nil {
			logrus.Errorf("Error opening file %s: %v\n", url, err)
			return
		}
		defer file.Close()
		GetCompleteFile("", file, tmplData)
	})
}
