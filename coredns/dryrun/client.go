package dryrun

import (
	"context"
	"encoding/json"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func newClient(scheme *runtime.Scheme) *dryrunClient {
	return &dryrunClient{
		scheme:  scheme,
		objects: make(map[objectKey]client.Object),
	}
}

var _ client.Client = &dryrunClient{}

type dryrunClient struct {
	scheme  *runtime.Scheme
	objects map[objectKey]client.Object
}

type objectKey struct {
	Group     string
	Kind      string
	Namespace string
	Name      string
}

func (c *dryrunClient) buildKey(obj client.Object) objectKey {
	gvk := obj.GetObjectKind().GroupVersionKind()
	if gvk.Kind == "" {
		foundGVK, err := apiutil.GVKForObject(obj, c.scheme)
		if err != nil {
			klog.Fatalf("cannot get GVK for %T: %v", obj, err)
		}
		gvk = foundGVK
	}
	return objectKey{
		Group:     gvk.Group,
		Kind:      gvk.Kind,
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
}

func (c *dryrunClient) Create(ctx context.Context, obj client.Object, opt ...client.CreateOption) error {
	if len(opt) != 0 {
		return fmt.Errorf("options not implemented")
	}

	k := c.buildKey(obj)
	if c.objects[k] != nil {
		return fmt.Errorf("duplicate object %v", k)
	}
	c.objects[k] = obj
	return nil
}

func (c *dryrunClient) Get(ctx context.Context, name types.NamespacedName, out client.Object) error {
	k := c.buildKey(out)
	k.Namespace = name.Namespace
	k.Name = name.Name

	obj := c.objects[k]
	if obj == nil {
		return apierrors.NewNotFound(schema.GroupResource{Group: k.Group, Resource: k.Kind}, k.Name)
	}

	// TODO(justinsb): is there a better way to do this?
	j, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	decoder := scheme.Codecs.UniversalDecoder()
	_, _, err = decoder.Decode(j, nil, out)
	return err
}

func (c *dryrunClient) Delete(ctx context.Context, obj client.Object, opt ...client.DeleteOption) error {
	return fmt.Errorf("Delete not implemented")
}

func (c *dryrunClient) DeleteAllOf(ctx context.Context, obj client.Object, opt ...client.DeleteAllOfOption) error {
	return fmt.Errorf("DeleteAllOf not implemented")
}

func (c *dryrunClient) List(ctx context.Context, obj client.ObjectList, opt ...client.ListOption) error {
	return fmt.Errorf("List not implemented")
}

func (c *dryrunClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opt ...client.PatchOption) error {
	return fmt.Errorf("Patch not implemented")
}

func (c *dryrunClient) Update(ctx context.Context, obj client.Object, opt ...client.UpdateOption) error {
	return fmt.Errorf("Update not implemented")
}

func (c *dryrunClient) RESTMapper() meta.RESTMapper {
	panic("RESTMapper not implemented")
}

func (c *dryrunClient) Scheme() *runtime.Scheme {
	return c.scheme
}

func (c *dryrunClient) Status() client.StatusWriter {
	panic("Status not implemented")
}
