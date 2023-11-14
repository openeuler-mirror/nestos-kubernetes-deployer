package upgrade

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
)

const namespaces = "housekeeper-system"

// todo: 创建CRD 、RBAC、controller资源

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
