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

package main

import (
	"flag"
	"os"
	"path/filepath"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	housekeeperiov1alpha1 "housekeeper.io/operator/api/v1alpha1"
	"housekeeper.io/operator/housekeeper-controller/controllers"
	"housekeeper.io/pkg/connection"
	"housekeeper.io/pkg/constants"
	"housekeeper.io/pkg/version"
	//+kubebuilder:scaffold:imports
)

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(housekeeperiov1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var err error
	opts := zap.Options{}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: "0",
	})
	if err != nil {
		logrus.Errorf("unable to start manager: %v", err)
		os.Exit(1)
	}

	reconciler := controllers.NewUpdateReconciler(mgr)
	if reconciler.Connection, err = connection.New("unix://" + filepath.Join(constants.SockDir, constants.SockName)); err != nil {
		logrus.Errorf("unable running housekeeper-controller: %v", err)
	}
	if err = reconciler.SetupWithManager(mgr); err != nil {
		logrus.Error(err, "unable to create controller", "controller", "Update")
		os.Exit(1)
	}

	logrus.Info("starting housekeeper-controller manager version:", version.Version)
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		logrus.Errorf("problem running housekeeper-controller manager: %v", err)
		os.Exit(1)
	}
}
