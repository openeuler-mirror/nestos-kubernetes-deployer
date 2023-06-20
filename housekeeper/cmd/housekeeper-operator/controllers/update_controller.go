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

	"github.com/golang/glog"
	housekeeperiov1alpha1 "housekeeper.io/api/v1alpha1"
	"housekeeper.io/pkg/common"
	"housekeeper.io/pkg/constants"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// UpdateReconciler reconciles a Update object
type UpdateReconciler struct {
	client.Client
	Scheme *runtime.Scheme
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
func (r *UpdateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	if r.Client == nil {
		return common.NoRequeue, nil
	}
	ctx = context.Background()
	err := setLabels(ctx, r, req)
	if err != nil {
		return common.RequeueNow, fmt.Errorf("unable set nodes label: %v", err)
	}
	return common.RequeueAfter, nil
}

func setLabels(ctx context.Context, r common.ReadWriterClient, req ctrl.Request) error {
	reqUpgrade, err := labels.NewRequirement(constants.LabelUpgrading, selection.DoesNotExist, nil)
	if err != nil {
		return fmt.Errorf("unable to create upgrade label requirement: %v", err)
	}
	reqMaster, err := labels.NewRequirement(constants.LabelMaster, selection.Exists, nil)
	if err != nil {
		return fmt.Errorf("unable to create master label requirement: %v", err)
	}
	reqNoMaster, err := labels.NewRequirement(constants.LabelMaster, selection.DoesNotExist, nil)
	if err != nil {
		return fmt.Errorf("unable to create non-master label requirement: %v", err)
	}
	masterNodes, err := getNodes(ctx, r, *reqUpgrade, *reqMaster)
	if err != nil {
		return fmt.Errorf("unable to get master node list: %v", err)
	}
	noMasterNodes, err := getNodes(ctx, r, *reqUpgrade, *reqNoMaster)
	if err != nil {
		return fmt.Errorf("unable to get non-master node list: %v", err)
	}
	upgradeCompleted, err := assignUpdated(ctx, r, masterNodes, req.NamespacedName)
	if err != nil {
		return fmt.Errorf("unabel to add the label of the master nodes: %v", err)
	}
	//Make sure the master upgrade is complete before start upgrading non-master nodes
	if upgradeCompleted {
		_, err := assignUpdated(ctx, r, noMasterNodes, req.NamespacedName)
		if err != nil {
			return fmt.Errorf("unabel to add the label of non-master nodes: %v", err)
		}
	}
	return nil
}

func getNodes(ctx context.Context, r common.ReadWriterClient, reqs ...labels.Requirement) ([]corev1.Node, error) {
	var nodeList corev1.NodeList
	opts := client.ListOptions{LabelSelector: labels.NewSelector().Add(reqs...)}
	if err := r.List(ctx, &nodeList, &opts); err != nil {
		return nil, fmt.Errorf("unable to list nodes with requirements: %v", err)
	}
	return nodeList.Items, nil
}

// Add the label to nodes
func assignUpdated(ctx context.Context, r common.ReadWriterClient, nodeList []corev1.Node, name types.NamespacedName) (bool, error) {
	var upInstance housekeeperiov1alpha1.Update
	upgradeNum := -1
	if err := r.Get(ctx, name, &upInstance); err != nil {
		return false, fmt.Errorf("unable to get update Instance %v", err)
	}
	var (
		kubeVersionSpec = upInstance.Spec.KubeVersion
		osVersionSpec   = upInstance.Spec.OSVersion
	)
	// Add the label after kubernetes version upgrade
	labelKubeVersion := fmt.Sprintf("%s%s", constants.LabelKubeVersionPrefix, kubeVersionSpec)
	if len(osVersionSpec) == 0 && len(kubeVersionSpec) == 0 {
		glog.Info("the os version and kube version cannot be all empty")
		return false, nil
	}
	for _, node := range nodeList {
		if len(kubeVersionSpec) > 0 {
			if _, ok := node.Labels[labelKubeVersion]; ok {
				glog.Infof("successfully upgraded the node: %s", node.Name)
				upgradeNum++
				continue
			}
		} else {
			if osVersionSpec == node.Status.NodeInfo.OSImage {
				continue
			}
		}
		node.Labels[constants.LabelUpgrading] = ""
		if err := r.Update(ctx, &node); err != nil {
			glog.Errorf("unable to add %s label:%v", node.Name, err)
		}
	}
	if len(kubeVersionSpec) == 0 {
		return true, nil
	}
	return upgradeNum == len(nodeList), nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UpdateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&housekeeperiov1alpha1.Update{}).
		Complete(r)
}
