/*
Copyright 2020 The Kubernetes Authors.

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
	"context"
	"flag"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	addonsv1alpha1 "sigs.k8s.io/cluster-addons/coredns/api/v1alpha1"
	"sigs.k8s.io/cluster-addons/coredns/controllers"
	"sigs.k8s.io/cluster-addons/coredns/dryrun"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = addonsv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	rbacMode := "reconcile"
	flag.StringVar(&rbacMode, "rbac-mode", rbacMode, "The mode to use for RBAC reconciliation.")

	dryRun := false
	flag.BoolVar(&dryRun, "dry-run", dryRun, "If set, will take an object on stdin and generate a manifest on stdout")

	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	addon.Init()

	if dryRun {
		reconcilerFactory := func(ctx context.Context, mgr manager.Manager, obj declarative.DeclarativeObject) (*declarative.Reconciler, error) {
			kind := obj.GetObjectKind().GroupVersionKind().Kind
			switch kind {
			case "CoreDNS":
				r := &controllers.CoreDNSReconciler{
					Client: mgr.GetClient(),
				}
				if err := r.SetupReconciler(mgr); err != nil {
					return nil, fmt.Errorf("error building reconciler: %w", err)
				}
				return &r.Reconciler, nil
			}

			return nil, fmt.Errorf("unknown kind %v", kind)
		}
		if err := dryrun.DoDryRun(context.Background(), scheme, reconcilerFactory); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		return
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "f4f34b31.x-k8s.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.CoreDNSReconciler{
		Client:   mgr.GetClient(),
		Log:      ctrl.Log.WithName("controllers").WithName("CoreDNS"),
		Scheme:   mgr.GetScheme(),
		RBACMode: rbacMode,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CoreDNS")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
