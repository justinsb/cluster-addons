/*
Copyright 2021 The Kubernetes Authors.

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
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"github.com/go-logr/logr"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	addonsv1alpha1 "sigs.k8s.io/cluster-addons/coredns/api/v1alpha1"
	"sigs.k8s.io/cluster-addons/coredns/controllers"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(addonsv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	rbacMode := "reconcile"
	flag.StringVar(&rbacMode, "rbac-mode", rbacMode, "The mode to use for RBAC reconciliation.")

	debugAddr := ""
	flag.StringVar(&debugAddr, "debug-bind-address", debugAddr, "The address the debug endpoint binds to.")

	klog.InitFlags(nil)
	flag.Parse()

	ctrl.SetLogger(klogr.New())

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                     scheme,
		MetricsBindAddress:         metricsAddr,
		Port:                       9443,
		HealthProbeBindAddress:     probeAddr,
		LeaderElection:             enableLeaderElection,
		LeaderElectionResourceLock: "leases",
		LeaderElectionID:           "coredns.addons.x-k8s.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if debugAddr != "" {
		if err := mgr.Add(&debugListener{Log: mgr.GetLogger(), Address: debugAddr}); err != nil {
			setupLog.Error(err, "unable to add debug listener")
			os.Exit(1)
		}
	}

	if err = (&controllers.CoreDNSReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		RBACMode: rbacMode,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CoreDNS")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

type debugListener struct {
	Address string
	Log     logr.Logger
}

var _ manager.Runnable = &debugListener{}

func (l *debugListener) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))

	server := &http.Server{
		Addr:         l.Address,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	l.Log.Info("starting debug listener", "address", l.Address)

	// TODO: Close when context closes
	return server.ListenAndServe()
}
