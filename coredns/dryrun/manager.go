package dryrun

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative/pkg/manifest"
	"sigs.k8s.io/kustomize/api/filesys"
)

type ReconcilerFactory func(context.Context, manager.Manager, declarative.DeclarativeObject) (*declarative.Reconciler, error)

func DoDryRun(ctx context.Context, scheme *runtime.Scheme, buildReconciler ReconcilerFactory) error {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("error reading from stdin: %w", err)
	}
	input, err := manifest.ParseObjects(ctx, string(b))
	if err != nil {
		return fmt.Errorf("error parsing objects: %w", err)
	}
	if len(input.Items) == 0 {
		return fmt.Errorf("expected at least object, got %d", len(input.Items))
	}
	serializer := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme, scheme, json.SerializerOptions{Yaml: false, Pretty: false, Strict: false})

	crJSON, err := input.Items[0].JSON()
	if err != nil {
		return fmt.Errorf("error building json for object: %w", err)
	}
	cr, _, err := serializer.Decode(crJSON, nil, nil)
	if err != nil {
		return fmt.Errorf("error parsing json for object: %w", err)
	}

	instance, ok := cr.(declarative.DeclarativeObject)
	if !ok {
		return fmt.Errorf("object was not expected type: %T", cr)
	}

	mgr := newManager(newClient(clientgoscheme.Scheme))
	mgr.Scheme = scheme

	// Build our mock object environment
	for _, obj := range input.Items {
		if err := mgr.GetClient().Create(ctx, obj.UnstructuredObject()); err != nil {
			return fmt.Errorf("error creating object: %w", err)
		}
	}

	r, err := buildReconciler(ctx, mgr, instance)
	if err != nil {
		return err
	}

	var fs filesys.FileSystem
	if r.IsKustomizeOptionUsed() {
		fs = filesys.MakeFsInMemory()
	}

	name := types.NamespacedName{
		Namespace: instance.GetNamespace(),
		Name:      instance.GetName(),
	}

	objects, err := r.BuildDeploymentObjectsWithFs(ctx, name, instance, fs)
	if err != nil {
		return fmt.Errorf("error building deployment objects: %w", err)
	}

	y, err := objects.YAMLManifest()
	if err != nil {
		return fmt.Errorf("error building yaml manifest: %w", err)
	}
	if _, err := os.Stdout.Write([]byte(y)); err != nil {
		return fmt.Errorf("error writing to stdout: %w", err)
	}
	return nil
}

// dryrunManager implements manager.Manager for dry-run operations
type dryrunManager struct {
	client        client.Client
	cache         cache.Cache
	config        rest.Config
	Scheme        *runtime.Scheme
	eventRecorder record.EventRecorder
	mapper        meta.RESTMapper
}

var _ manager.Manager = &dryrunManager{}

func newManager(c *dryrunClient) *dryrunManager {
	return &dryrunManager{
		client: c,
		//	cache:  FakeCache{},
	}
}

func (m *dryrunManager) Add(manager.Runnable) error {
	panic("not supported by dry-run manager")
}

func (m *dryrunManager) SetFields(interface{}) error {
	panic("not supported by dry-run manager")
}

func (m *dryrunManager) Start(context.Context) error {
	panic("not supported by dry-run manager")
}

func (m *dryrunManager) GetConfig() *rest.Config {
	return &m.config
}

func (m *dryrunManager) GetScheme() *runtime.Scheme {
	return m.Scheme
}

func (m *dryrunManager) GetClient() client.Client {
	return m.client
}

func (m *dryrunManager) GetFieldIndexer() client.FieldIndexer {
	panic("not supported by dry-run manager")
}

func (m *dryrunManager) GetCache() cache.Cache {
	// if m.cache == nil {
	// 	m.cache = FakeCache{}
	// }
	// return m.cache
	panic("not supported by dry-run manager")
}

func (m *dryrunManager) GetRecorder(name string) record.EventRecorder {
	panic("not supported by dry-run manager")
}

func (m *dryrunManager) GetRESTMapper() meta.RESTMapper {
	return m.mapper
}

func (m *dryrunManager) GetAPIReader() client.Reader {
	panic("not supported by dry-run manager")
}

func (m *dryrunManager) GetEventRecorderFor(name string) record.EventRecorder {
	return m.eventRecorder
}

func (m *dryrunManager) GetWebhookServer() *webhook.Server {
	panic("not supported by dry-run manager")
}

func (m *dryrunManager) AddHealthzCheck(name string, check healthz.Checker) error {
	panic("not supported by dry-run manager")
}

func (m *dryrunManager) AddReadyzCheck(name string, check healthz.Checker) error {
	panic("not supported by dry-run manager")
}

func (m *dryrunManager) AddMetricsExtraHandler(path string, handler http.Handler) error {
	panic("not supported by dry-run manager")
}

func (m *dryrunManager) Elected() <-chan struct{} {
	panic("not supported by dry-run manager")
}

func (m *dryrunManager) GetLogger() logr.Logger {
	panic("not supported by dry-run manager")
}
