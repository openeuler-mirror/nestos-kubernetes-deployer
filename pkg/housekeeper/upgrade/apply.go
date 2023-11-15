package upgrade

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const namespaces = "housekeeper-system"

func deployOperator(kubeconfig string){
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
	// todo: 部署CRD
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
