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
package utils

import (
	"bytes"
	"io"
	"nestos-kubernetes-deployer/data"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

// FetchAndUnmarshalURL fetches content from a specified URL, unmarshals it into the provided structure,
func FetchAndUnmarshalUrl(url string, tmplData interface{}) ([]byte, error) {
	file, err := data.Assets.Open(url)
	if err != nil {
		logrus.Errorf("Error opening file %s: %v\n", url, err)
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		logrus.Errorf("Error getting file info for %s: %v\n", url, err)
		return nil, err
	}
	_, data, err := GetCompleteFile(info.Name(), file, tmplData)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetCompleteFile(name string, file io.Reader, tmplData interface{}) (realName string, data []byte, err error) {
	data, err = io.ReadAll(file)
	if err != nil {
		logrus.Errorf("Error reading file %s: %v\n", name, err)
		return "", nil, err
	}
	if filepath.Ext(name) == ".template" {
		name = strings.TrimSuffix(name, ".template")
		tmpl := template.New(name)
		tmpl, err := tmpl.Parse(string(data))
		if err != nil {
			logrus.Errorf("Error parsing template for file %s: %v\n", name, err)
			return "", nil, err
		}
		stringData := applyTmplData(tmpl, tmplData)
		data = []byte(stringData)
	}

	return name, data, nil
}

func applyTmplData(tmpl *template.Template, data interface{}) string {
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, data); err != nil {
		logrus.Errorf("Error applying template: %v\n", err)
		panic(err)
	}
	return buf.String()
}
