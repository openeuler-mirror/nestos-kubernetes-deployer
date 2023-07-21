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
	"time"

	"github.com/sirupsen/logrus"
	housekeeperiov1alpha1 "housekeeper.io/operator/api/v1alpha1"
	"housekeeper.io/pkg/common"
	"housekeeper.io/pkg/constants"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
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
	//返回worker节点的数量
	upInstance, nodesNum, err := getUpgradeInstance(ctx, r, req.NamespacedName)
	if err != nil {
		return common.RequeueNow, err
	}
	limit := min(upInstance.Spec.MaxUnavailable, nodesNum)
	if requeueAfter, err := setLabels(ctx, r, req, limit, upInstance); err != nil {
		logrus.Errorf("unable set nodes label: %v", err)
		return common.RequeueNow, err
	} else if requeueAfter {
		return common.RequeueAfter, nil
	}
	return common.RequeueNow, nil
}

func getUpgradeInstance(ctx context.Context, r common.ReadWriterClient, name types.NamespacedName) (
	upInstance housekeeperiov1alpha1.Update, nodeNum int, err error) {
	if err = r.Get(ctx, name, &upInstance); err != nil {
		logrus.Errorf("unable to fetch upgrade instance: %v", err)
		return
	}
	requirement, err := labels.NewRequirement(constants.LabelMaster, selection.DoesNotExist, nil)
	if err != nil {
		logrus.Errorf("unable to create requirement %s: %v"+constants.LabelMaster, err)
		return
	}
	nodesItems, err := getNodes(ctx, r, 0, *requirement)
	if err != nil {
		logrus.Errorf("failed to get nodes list: %v", err)
		return
	}
	nodeNum = len(nodesItems)
	return
}

func setLabels(ctx context.Context, r common.ReadWriterClient, req ctrl.Request, limit int,
	upInstance housekeeperiov1alpha1.Update) (bool, error) {
	reqUpgrade, err := labels.NewRequirement(constants.LabelUpgrading, selection.DoesNotExist, nil)
	if err != nil {
		logrus.Errorf("unable to create upgrade label requirement: %v", err)
		return false, err
	}
	reqMaster, err := labels.NewRequirement(constants.LabelMaster, selection.Exists, nil)
	if err != nil {
		logrus.Errorf("unable to create master label requirement: %v", err)
		return false, err
	}
	reqNoMaster, err := labels.NewRequirement(constants.LabelMaster, selection.DoesNotExist, nil)
	if err != nil {
		logrus.Errorf("unable to create non-master label requirement: %v", err)
		return false, err
	}
	masterNodes, err := getNodes(ctx, r, 1, *reqUpgrade, *reqMaster)
	if err != nil {
		logrus.Errorf("unable to get master node list: %v", err)
		return false, err
	}
	//limit: 限制worker节点每次升级的数量
	noMasterNodes, err := getNodes(ctx, r, limit, *reqUpgrade, *reqNoMaster)
	if err != nil {
		logrus.Errorf("unable to get non-master node list: %v", err)
		return false, err
	}
	needRequeue, err := assignUpdated(ctx, r, masterNodes, upInstance)
	if err != nil {
		logrus.Errorf("unabel to add upgrade label of the master nodes: %v", err)
		return false, err
	} else if needRequeue {
		return true, nil
	}
	if needRequeue, err = assignUpdated(ctx, r, noMasterNodes, upInstance); err != nil {
		logrus.Errorf("unabel to add upgrade label of non-master nodes: %v", err)
		return false, err
	}
	return needRequeue, nil
}

func getNodes(ctx context.Context, r common.ReadWriterClient, limit int, reqs ...labels.Requirement) ([]corev1.Node, error) {
	var nodeList corev1.NodeList
	opts := client.ListOptions{LabelSelector: labels.NewSelector().Add(reqs...), Limit: int64(limit)}
	if err := r.List(ctx, &nodeList, &opts); err != nil {
		logrus.Errorf("unable to list nodes with requirements: %v", err)
		return nil, err
	}
	return nodeList.Items, nil
}

// Add the label to nodes
func assignUpdated(ctx context.Context, r common.ReadWriterClient, nodeList []corev1.Node,
	upInstance housekeeperiov1alpha1.Update) (bool, error) {
	var (
		kubeVersionSpec = upInstance.Spec.KubeVersion
		osVersionSpec   = upInstance.Spec.OSVersion
	)
	if len(osVersionSpec) == 0 {
		logrus.Warning("os version is required")
		return false, nil
	}
	if len(nodeList) == 0 {
		return false, nil
	}
	for _, node := range nodeList {
		if conditionMet(node, kubeVersionSpec, osVersionSpec) {
			node.Labels[constants.LabelUpgrading] = ""
			if err := r.Update(ctx, &node); err != nil {
				logrus.Errorf("unable to add %s label:%v", node.Name, err)
				return false, err
			}
			if err := waitForUpgradeComplete(node, kubeVersionSpec, osVersionSpec); err != nil {
				logrus.Errorf("failed to wait for node upgrade to complete: %v", err)
				return false, err
			}
		} else {
			return false, nil
		}
	}
	return true, nil
}

func conditionMet(node corev1.Node, kubeVersionSpec string, osVersionSpec string) bool {
	var (
		kubeProxyVersion = node.Status.NodeInfo.KubeProxyVersion
		kubeletVersion   = node.Status.NodeInfo.KubeletVersion
		osVersion        = node.Status.NodeInfo.OSImage
	)
	if len(kubeVersionSpec) > 0 {
		if kubeVersionSpec == kubeProxyVersion && kubeVersionSpec == kubeletVersion {
			return false
		}
	} else {
		if osVersionSpec == osVersion {
			return false
		}
	}
	return true
}

func waitForUpgradeComplete(node corev1.Node, kubeVersionSpec string, osVersionSpec string) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.Timeout)
	defer cancel()
	done := make(chan struct{})

	go func() {
		wait.Until(func() {
			if !conditionMet(node, kubeVersionSpec, osVersionSpec) {
				close(done)
			}
		}, 10*time.Second, ctx.Done())
	}()

	select {
	case <-done:
		logrus.Infof("successful upgrade node: %s", node.Name)
	case <-ctx.Done():
		// 上下文超时，跳出循环
		if ctx.Err() == context.DeadlineExceeded {
			logrus.Errorf("failed to upgrade node: %s: %v", node.Name, ctx.Err())
			return ctx.Err()
		}
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// SetupWithManager sets up the controller with the Manager.
func (r *UpdateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&housekeeperiov1alpha1.Update{}).
		Complete(r)
}
