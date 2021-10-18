package convert

import (
	"context"
	"fmt"
	"sort"
	"strings"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative/pkg/manifest"
)

type BuildRoleOptions struct {
	Name               string
	Namespace          string
	ServiceAccountName string
	Supervisory        bool

	// CRD is the name of the CRD to generate permissions for.
	CRD string

	// LimitResourceNames specifies that RBAC permissions should restrict to resource names in the manifest.
	LimitResourceNames bool

	// LimitNamespaces specifies that RBAC permissions should restrict to resource names in the manifest.
	LimitNamespaces bool

	// Format specifies the format we should write in (yaml or kubebuilder)
	Format string
}

type ruleSet struct {
	rules []v1.PolicyRule
}

func (r *ruleSet) Add(rules ...v1.PolicyRule) {
	r.rules = append(r.rules, rules...)
}

func BuildRole(manifestStr string, opt BuildRoleOptions) ([]runtime.Object, error) {
	var objects []runtime.Object

	ctx := context.Background()
	objs, err := manifest.ParseObjects(ctx, manifestStr)
	if err != nil {
		return nil, err
	}
	if len(objs.Blobs) != 0 {
		return nil, fmt.Errorf("unable to parse manifest fully")
	}

	ruleMap := make(map[string]*ruleSet)
	ruleMap[""] = &ruleSet{}

	// to deal with duplicates, we keep a map of all the kinds that has been added so far
	kindMap := make(map[string]bool)

	for _, obj := range objs.Items {
		ruleSetKey := ""
		if opt.LimitNamespaces {
			ruleSetKey = obj.Namespace
		}

		rules := ruleMap[ruleSetKey]
		if rules == nil {
			rules = &ruleSet{}
			ruleMap[ruleSetKey] = rules
		}

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
				return nil, err
			}
			rules.Add(newClusterRole.Rules...)
		}

		// needs plural of kind
		resource := ResourceFromKind(obj.Kind)

		if opt.LimitResourceNames {
			rules.Add(v1.PolicyRule{
				APIGroups:     []string{obj.Group},
				Resources:     []string{resource},
				ResourceNames: []string{obj.Name},
				Verbs:         []string{"update", "delete", "patch"},
			})

			rules.Add(v1.PolicyRule{
				APIGroups: []string{obj.Group},
				Resources: []string{resource},
				Verbs:     []string{"create"},
			})

			rules.Add(v1.PolicyRule{
				APIGroups: []string{obj.Group},
				Resources: []string{resource},
				Verbs:     []string{"get", "list", "watch"},
			})

		} else if !kindMap[obj.Group+"::"+obj.Kind] {
			newRule := v1.PolicyRule{
				APIGroups: []string{obj.Group},
				Resources: []string{resource},
				Verbs:     []string{"create", "update", "delete", "get"},
			}
			rules.Add(newRule)
			kindMap[obj.Group+"::"+obj.Kind] = true
		}
	}

	if opt.CRD != "" {
		gr := schema.ParseGroupResource(opt.CRD)

		// TODO: Should we assume namespace scoped?
		rules := ruleMap[""]

		rules.Add(v1.PolicyRule{
			APIGroups: []string{gr.Group},
			Resources: []string{gr.Resource},
			Verbs:     []string{"get", "list", "patch", "update", "watch"},
		})
		rules.Add(v1.PolicyRule{
			APIGroups: []string{gr.Group},
			Resources: []string{gr.Resource + "/status"},
			Verbs:     []string{"get", "patch", "update"},
		})
	}

	for ns, rules := range ruleMap {
		rules.normalize()

		if ns != "" {
			role := &v1.Role{
				ObjectMeta: metav1.ObjectMeta{
					Name:      opt.Name,
					Namespace: ns,
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "Role",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				Rules: rules.rules,
			}
			objects = append(objects, role)

			// if saName is passed in, generate YAML for rolebinding
			if opt.ServiceAccountName != "" {
				roleBinding := v1.RoleBinding{
					TypeMeta: metav1.TypeMeta{
						Kind:       "RoleBinding",
						APIVersion: "rbac.authorization.k8s.io/v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      opt.Name + "-binding",
						Namespace: ns,
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
						Kind:     "Role",
						Name:     opt.Name,
					},
				}

				objects = append(objects, &roleBinding)
			}
		} else {
			role := &v1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name: opt.Name,
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterRole",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				Rules: rules.rules,
			}
			objects = append(objects, role)

			// if saName is passed in, generate YAML for rolebinding
			if opt.ServiceAccountName != "" {
				roleBinding := v1.ClusterRoleBinding{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ClusterRoleBinding",
						APIVersion: "rbac.authorization.k8s.io/v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: opt.Name + "-binding",
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

				objects = append(objects, &roleBinding)
			}
		}
	}

	return objects, err
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

func (r *ruleSet) normalize() {
	r.rules = normalizeRules(r.rules)

	r.rules = combineRBACRules(r.rules)

	sort.Slice(r.rules, func(i, j int) bool { return ruleLT(&r.rules[i], &r.rules[j]) })
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

func combineRBACRules(rules []v1.PolicyRule) []v1.PolicyRule {
	rules = foldResources(rules)
	rules = foldVerbs(rules)
	return rules
}

func foldResources(rules []v1.PolicyRule) []v1.PolicyRule {
	var out []v1.PolicyRule

	ruleMap := make(map[string]v1.PolicyRule)

	for _, rule := range rules {
		if len(rule.NonResourceURLs) != 0 {
			out = append(out, rule)
			continue
		}

		key := "groups=" + strings.Join(rule.APIGroups, ",")
		//key += ";resources=" + strings.Join(rule.Resources, ",")
		key += ";resourceNames=" + strings.Join(rule.ResourceNames, ",")
		key += ";verbs=" + strings.Join(rule.Verbs, ",")

		existing, found := ruleMap[key]
		if !found {
			ruleMap[key] = rule
			continue
		}

		existing.Resources = append(existing.Resources, rule.Resources...)
		ruleMap[key] = existing
	}

	for _, rule := range ruleMap {
		out = append(out, rule)
	}

	return out
}

func foldVerbs(rules []v1.PolicyRule) []v1.PolicyRule {
	var out []v1.PolicyRule

	ruleMap := make(map[string]v1.PolicyRule)

	for _, rule := range rules {
		if len(rule.NonResourceURLs) != 0 {
			out = append(out, rule)
			continue
		}

		key := "groups=" + strings.Join(rule.APIGroups, ",")
		key += ";resources=" + strings.Join(rule.Resources, ",")
		key += ";resourceNames=" + strings.Join(rule.ResourceNames, ",")
		// key += ";verbs=" + strings.Join(rule.Verbs, ",")

		existing, found := ruleMap[key]
		if !found {
			ruleMap[key] = rule
			continue
		}

		existing.Verbs = append(existing.Verbs, rule.Verbs...)
		ruleMap[key] = existing
	}

	for _, rule := range ruleMap {
		out = append(out, rule)
	}

	return out
}

func normalizeRules(rules []v1.PolicyRule) []v1.PolicyRule {
	for i := range rules {
		rule := &rules[i]

		rule.APIGroups = normalizeStringSlice(rule.APIGroups)
		rule.NonResourceURLs = normalizeStringSlice(rule.NonResourceURLs)
		rule.ResourceNames = normalizeStringSlice(rule.ResourceNames)
		rule.Resources = normalizeStringSlice(rule.Resources)
		rule.Verbs = normalizeStringSlice(rule.Verbs)
	}

	return rules
}

func normalizeStringSlice(in []string) []string {
	var out []string

	done := make(map[string]bool)
	for _, s := range in {
		if done[s] {
			continue
		}
		out = append(out, s)
	}

	sort.Strings(out)
	return out
}
