package convert

import (
	"bytes"
	"context"
	"sort"
	"strings"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative/pkg/manifest"
	"sigs.k8s.io/yaml"
)

type BuildRoleOptions struct {
	Name               string
	Namespace          string
	ServiceAccountName string
	Supervisory        bool

	// CRD is the name of the CRD to generate permissions for.
	CRD string
}

func ParseYAMLtoRole(manifestStr string, opt BuildRoleOptions) (string, error) {
	ctx := context.Background()
	objs, err := manifest.ParseObjects(ctx, manifestStr)
	if err != nil {
		return "", err
	}

	clusterRole := v1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      opt.Name,
			Namespace: opt.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
	}
	// to deal with duplicates, we keep a map of all the kinds that has been added so far
	kindMap := make(map[string]bool)

	for _, obj := range objs.Items {
		// The generated role needs the rules from any role or clusterrole
		if obj.Kind == "Role" || obj.Kind == "ClusterRole" {
			if opt.Supervisory {
				continue
			}
			unstruct := obj.UnstructuredObject()
			newClusterRole := v1.ClusterRole{}

			// Converting from unstructured to v1.ClusterRole
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstruct.Object, &newClusterRole)
			if err != nil {
				return "", err
			}
			clusterRole.Rules = append(clusterRole.Rules, newClusterRole.Rules...)
		}

		if !kindMap[obj.Group+"::"+obj.Kind] {
			newRule := v1.PolicyRule{
				APIGroups: []string{obj.Group},
				// needs plural of kind
				Resources: []string{ResourceFromKind(obj.Kind)},
				Verbs:     []string{"create", "update", "delete", "get"},
			}
			clusterRole.Rules = append(clusterRole.Rules, newRule)
			kindMap[obj.Group+"::"+obj.Kind] = true
		}
	}

	if opt.CRD != "" {
		gr := schema.ParseGroupResource(opt.CRD)

		clusterRole.Rules = append(clusterRole.Rules, v1.PolicyRule{
			APIGroups: []string{gr.Group},
			Resources: []string{gr.Resource},
		})
	}

	sort.Slice(clusterRole.Rules, func(i, j int) bool { return ruleLT(&clusterRole.Rules[i], &clusterRole.Rules[j]) })

	output, err := yaml.Marshal(&clusterRole)
	buf := bytes.NewBuffer(output)

	// if saName is passed in, generate YAML for rolebinding
	if opt.ServiceAccountName != "" {
		clusterRoleBinding := v1.ClusterRoleBinding{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ClusterRoleBinding",
				APIVersion: "rbac.authorization.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: opt.Namespace,
				Name:      opt.Name + "-binding",
			},
			Subjects: []v1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      opt.ServiceAccountName,
					Namespace: opt.Namespace,
				},
			},
			RoleRef: v1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     opt.Name,
			},
		}

		outputBinding, err := yaml.Marshal(&clusterRoleBinding)
		if err != nil {
			return "", err
		}
		buf.WriteString("\n---\n\n")
		buf.Write(outputBinding)
	}
	return buf.String(), err
}

func ResourceFromKind(kind string) string {
	if string(kind[len(kind)-1]) == "s" {
		return strings.ToLower(kind) + "es"
	}
	if string(kind[len(kind)-1]) == "y" {
		return strings.ToLower(kind)[:len(kind)-1] + "ies"
	}
	return strings.ToLower(kind) + "s"
}

func ruleLT(l, r *v1.PolicyRule) bool {
	lGroup := firstOrEmpty(l.APIGroups)
	rGroup := firstOrEmpty(r.APIGroups)
	if lGroup != rGroup {
		return lGroup < rGroup
	}
	lResource := firstOrEmpty(l.Resources)
	rResource := firstOrEmpty(r.Resources)
	if lResource != rResource {
		return lResource < rResource
	}
	lVerb := firstOrEmpty(l.Verbs)
	rVerb := firstOrEmpty(r.Verbs)
	if lVerb != rVerb {
		return lVerb < rVerb
	}
	return false
}

func firstOrEmpty(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return s[0]
}
