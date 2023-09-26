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
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	housekeeperiov1alpha1 "housekeeper.io/operator/api/v1alpha1"
	"housekeeper.io/pkg/common"
	"housekeeper.io/pkg/constants"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
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
	var crMutex sync.Mutex
	crMutex.Lock()
	defer crMutex.Unlock()
	ctx = context.Background()
	return reconcile(ctx, r, req)
}

// SetupWithManager sets up the controller with the Manager.
func (r *UpdateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&housekeeperiov1alpha1.Update{}).
		Complete(r)
}

func reconcile(ctx context.Context, r common.ReadWriterClient, req ctrl.Request) (ctrl.Result, error) {
	var update housekeeperiov1alpha1.Update
	if err := r.Get(ctx, req.NamespacedName, &update); err != nil {
		logrus.Errorf("unable to fetch update instance: %v", err)
		return common.NoRequeue, err
	}
	if len(update.Spec.OSImageURL) == 0 {
		logrus.Warning("os upgrade image url is required")
		return common.RequeueAfter, nil
	}

	allNodes, err := getAllNodes(ctx, r)
	if err != nil {
		return common.RequeueNow, err
	}

	allNodesUpgraded := true
	for _, node := range allNodes {
		if _, exists := node.Labels[constants.LabelUpgradeCompleted]; !exists {
			allNodesUpgraded = false
			break
		}
	}
	if allNodesUpgraded {
		for _, node := range allNodes {
			delete(node.Labels, constants.LabelUpgradeCompleted)
			if err := r.Update(ctx, &node); err != nil {
				return common.RequeueNow, err
			}
		}
		return common.NoRequeue, nil // 不重新触发 CR
	}
	masterNodesItems, err := getMasterNodesItems(ctx, r)
	if err != nil {
		return common.RequeueNow, err
	}
	workerNodesItems, err := getWorkerNodesItems(ctx, r)
	if err != nil {
		return common.RequeueNow, err
	}
	if assignUpdated(ctx, r, masterNodesItems, 1, update); err != nil {
		return common.RequeueNow, err
	}
	maxUnavailable := min(update.Spec.MaxUnavailable, len(workerNodesItems))
	if assignUpdated(ctx, r, workerNodesItems, maxUnavailable, update); err != nil {
		return common.RequeueNow, err
	}

	return common.RequeueNow, nil
}

func getMasterNodesItems(ctx context.Context, r common.ReadWriterClient) (
	nodesItems []corev1.Node, err error) {
	reqUpgrade, err := labels.NewRequirement(constants.LabelUpgrading, selection.DoesNotExist, nil)
	if err != nil {
		logrus.Errorf("unable to create requirement %s: %v", reqUpgrade, err)
		return
	}
	reqUpgradeCompleted, err := labels.NewRequirement(constants.LabelUpgradeCompleted, selection.DoesNotExist, nil)
	if err != nil {
		logrus.Errorf("unable to create requirement %s: %v", reqUpgradeCompleted, err)
		return
	}
	reqMaster, err := labels.NewRequirement(constants.LabelMaster, selection.Exists, nil)
	if err != nil {
		logrus.Errorf("unable to create requirement %s: %v", constants.LabelMaster, err)
		return
	}
	nodesItems, err = getNodes(ctx, r, *reqUpgrade, *reqUpgradeCompleted, *reqMaster)
	if err != nil {
		logrus.Errorf("failed to get master nodes list: %v", err)
		return
	}
	return
}

func getWorkerNodesItems(ctx context.Context, r common.ReadWriterClient) (
	nodesItems []corev1.Node, err error) {
	reqUpgrade, err := labels.NewRequirement(constants.LabelUpgrading, selection.DoesNotExist, nil)
	if err != nil {
		logrus.Errorf("unable to create requirement %s: %v", reqUpgrade, err)
		return
	}
	reqUpgradeCompleted, err := labels.NewRequirement(constants.LabelUpgradeCompleted, selection.DoesNotExist, nil)
	if err != nil {
		logrus.Errorf("unable to create requirement %s: %v", reqUpgradeCompleted, err)
		return
	}
	reqWorker, err := labels.NewRequirement(constants.LabelMaster, selection.DoesNotExist, nil)
	if err != nil {
		logrus.Errorf("unable to create requirement %s: %v", constants.LabelMaster, err)
		return
	}
	nodesItems, err = getNodes(ctx, r, *reqUpgrade, *reqUpgradeCompleted, *reqWorker)
	if err != nil {
		logrus.Errorf("failed to get worker nodes list: %v", err)
		return
	}
	return
}

func getNodes(ctx context.Context, r common.ReadWriterClient, reqs ...labels.Requirement) ([]corev1.Node, error) {
	var nodeList corev1.NodeList
	opts := client.ListOptions{LabelSelector: labels.NewSelector().Add(reqs...)}
	if err := r.List(ctx, &nodeList, &opts); err != nil {
		logrus.Errorf("unable to list nodes with requirements: %v", err)
		return nil, err
	}
	return nodeList.Items, nil
}

func getAllNodes(ctx context.Context, r common.ReadWriterClient) ([]corev1.Node, error) {
	var nodeList corev1.NodeList
	if err := r.List(ctx, &nodeList); err != nil {
		logrus.Errorf("unable to list nodes: %v", err)
		return nil, err
	}
	return nodeList.Items, nil
}

// Add the label to nodes
func assignUpdated(ctx context.Context, r common.ReadWriterClient, nodeList []corev1.Node,
	max int, upInstance housekeeperiov1alpha1.Update) error {

	timeoutCtx, cancel := context.WithTimeout(ctx, constants.NodeTimeout)
	defer cancel()
	for _, node := range nodeList {
		if hasUpgradeCompletedLabel(node) {
			continue
		}
		if max <= 0 {
			if err := waitForUpgradeComplete(timeoutCtx, node); err != nil {
				return err
			}
			max = upInstance.Spec.MaxUnavailable
		}
		node.Labels[constants.LabelUpgrading] = ""
		if err := r.Update(ctx, &node); err != nil {
			return err
		}
		max--
	}
	return nil
}

func waitForUpgradeComplete(ctx context.Context, node corev1.Node) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, constants.NodeTimeout)
	defer cancel()

	done := make(chan struct{})
	success := false

	waitFunc := func() {
		if hasUpgradeCompletedLabel(node) {
			success = true
			close(done)
		}
	}

	go wait.Until(waitFunc, 10*time.Second, timeoutCtx.Done())

	select {
	case <-done:
		if success {
			logrus.Infof("successful upgrade node: %s", node.Name)
		} else {
			logrus.Infof("upgrade conditions not met for node: %s", node.Name)
		}
		return nil
	case <-timeoutCtx.Done():
		if timeoutCtx.Err() == context.DeadlineExceeded {
			logrus.Errorf("timeout to upgrade node: %s: %v", node.Name, timeoutCtx.Err())
		}
		return timeoutCtx.Err()
	}
}

func hasUpgradeCompletedLabel(node corev1.Node) bool {
	_, exists := node.Labels[constants.LabelUpgradeCompleted]
	return exists
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
