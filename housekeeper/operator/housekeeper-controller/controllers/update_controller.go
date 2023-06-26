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

package controllers

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	housekeeperiov1alpha1 "housekeeper.io/operator/api/v1alpha1"
	"housekeeper.io/pkg/common"
	"housekeeper.io/pkg/connection"
	"housekeeper.io/pkg/constants"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// UpdateReconciler reconciles a Update object
type UpdateReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	KubeClientSet kubernetes.Interface
	Connection    *connection.Client
	HostName      string
}

//+kubebuilder:rbac:groups=housekeeper.io,resources=updates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=housekeeper.io,resources=updates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=housekeeper.io,resources=updates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Update object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func NewUpdateReconciler(mgr manager.Manager) *UpdateReconciler {
	kubeClientSet, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		logrus.Errorf("failed to build the kubernetes clientset: %v", err)
	}
	reconciler := &UpdateReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		KubeClientSet: kubeClientSet,
		HostName:      os.Getenv("NODE_NAME"),
	}
	return reconciler
}

func (r *UpdateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	ctx = context.Background()
	upInstance, nodeInstance := reqInstance(ctx, r, req.NamespacedName, r.HostName)
	upgradeCluster := checkUpgrade(&nodeInstance, upInstance.Spec.OSVersion, upInstance.Spec.KubeVersion)
	if upgradeCluster {
		if err := r.upgradeNodes(ctx, &upInstance, &nodeInstance); err != nil {
			return common.RequeueNow, err
		}
	} else {
		r.refreshNodes(ctx, &upInstance, &nodeInstance)
	}
	return common.RequeueAfter, nil
}

func (r *UpdateReconciler) upgradeNodes(ctx context.Context, upInstance *housekeeperiov1alpha1.Update,
	node *corev1.Node) error {
	controlPlane := false
	if _, ok := node.Labels[constants.LabelMaster]; ok {
		controlPlane = true
	}
	if _, ok := node.Labels[constants.LabelUpgrading]; ok {
		//todo drain
		pushInfo := &connection.PushInfo{
			KubeVersion:  upInstance.Spec.KubeVersion,
			OSImageURL:   upInstance.Spec.OSImageURL,
			OSVersion:    upInstance.Spec.OSVersion,
			ControlPlane: controlPlane,
		}
		if err := r.Connection.UpgradeKubeSpec(pushInfo); err != nil {
			return err
		}
	}

	return nil
}

func (r *UpdateReconciler) refreshNodes(ctx context.Context, upInstance *housekeeperiov1alpha1.Update, node *corev1.Node) error {
	deleteLabel(ctx, r, node)
	if node.Spec.Unschedulable {
		//todo drain
	}
	return nil
}

func deleteLabel(ctx context.Context, r common.ReadWriterClient, node *corev1.Node) error {
	if _, ok := node.Labels[constants.LabelUpgrading]; ok {
		delete(node.Labels, constants.LabelUpgrading)
		if err := r.Update(ctx, node); err != nil {
			logrus.Errorf("unable to delete %s node label: %w", node.Name, err)
			return err
		}
	}
	return nil
}

func reqInstance(ctx context.Context, r common.ReadWriterClient, name types.NamespacedName,
	HostName string) (upInstance housekeeperiov1alpha1.Update, nodeInstance corev1.Node) {
	if err := r.Get(ctx, name, &upInstance); err != nil {
		logrus.Errorf("unable to fetch update instance: %v", err)
		return
	}
	if err := r.Get(ctx, client.ObjectKey{Name: HostName}, &nodeInstance); err != nil {
		logrus.Errorf("unable to fetch node instance: %v", err)
		return
	}
	return
}

func checkUpgrade(node *corev1.Node, osVersionSpec string, kubeVersionSpec string) bool {
	if len(kubeVersionSpec) > 0 {
		labelKubeVersion := fmt.Sprintf("%s%s", constants.LabelKubeVersionPrefix, kubeVersionSpec)
		if _, ok := node.Labels[labelKubeVersion]; ok {
			return false
		}
	} else {
		return node.Status.NodeInfo.OSImage != osVersionSpec
	}
	return true
}

// SetupWithManager sets up the controller with the Manager.
func (r *UpdateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&housekeeperiov1alpha1.Update{}).
		Complete(r)
}
