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
package common

import (
	"fmt"
	"os"
	"strings"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ReadWriterClient is Kubernetes API
type ReadWriterClient interface {
	client.Reader
	client.StatusClient
	client.Writer
}

var (
	// controller do not requeue
	NoRequeue = ctrl.Result{}
	// controller requeue
	RequeueNow   = ctrl.Result{Requeue: true}
	RequeueAfter = ctrl.Result{Requeue: true, RequeueAfter: time.Second * 20}
)

func IsFileExist(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	if fileInfo.IsDir() {
		return false
	}
	return true
}

func ExtractImageTag(imageURL string) (string, error) {
	parts := strings.Split(imageURL, "/")
	lastPart := parts[len(parts)-1]
	tagParts := strings.Split(lastPart, ":")
	if len(tagParts) > 1 {
		return tagParts[len(tagParts)-1], nil
	}
	return "", fmt.Errorf("unable to extract the mirror tag from image URL: %s", imageURL)
}
