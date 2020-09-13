package xform

import "sigs.k8s.io/kustomize/kyaml/yaml"

func RemoveVolumes(keepPredicate func(*yaml.RNode) bool) ObjectFilter {
	return ObjectFilter{
		FieldPaths: []FieldPath{
			ParseFieldPath("spec.template.spec.containers.[=].volumeMounts"),
			ParseFieldPath("spec.template.spec.volumes"),
		},
		KeepPredicate: keepPredicate,
	}
}
