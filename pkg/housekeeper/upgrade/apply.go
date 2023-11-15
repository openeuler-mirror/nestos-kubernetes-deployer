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

package upgrade

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"nestos-kubernetes-deployer/data"
)

const namespaces = "housekeeper-system"

type yamlTmlpData struct{
	operatorImageUrl string
	controllerImageUrl string
}

//todo:获取yamlTmlpData数据

func DeployOperator(kubeconfig string){
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logrus.Errorf("Error building kubeconfig: %v", err)
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Errorf("Error creating kubernetes client: %v", err)
		return err
	}
	//todo:获取yaml内容，部署资源
	
}

func applyYAML(clientset *kubernetes.Clientset, yamlContent) error {
	_, err := clientset.RESTClient().
		Post().
		Resource("").
		Body([]byte(yamlContent)).
		Do(context.TODO()).
		Get()

	if err != nil {
		logrus.ErrorF("Error applying YAML: %v\n", err)
		return err
	}
	return nil
}

func getYamlContent(uri string, yamlTmlpData interface{}) ([]byte, error){
	file, err := data.Assets.Open(uri)
	if err != nil {
		logrus.Errorf("Error opening file %s: %v\n", uri, err)
		return nil,err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		logrus.Errorf("Error getting file info for %s: %v\n", uri, err)
		return nil,err
	}
	_, data, err := readFile(info.Name(), file, yamlTmlpData)
	if err != nil {
		logrus.Errorf("Error reading file %s: %v\n", uri, err)
		return nil, err
	}
	return data,nil
}

// Read data from the file
func readFile(name string, file io.Reader, tmplData interface{}) (realName string, data []byte, err error) {
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