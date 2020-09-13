package xform

import (
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// ObjectFilterFunc removes the field or map key.
// Returns a RNode with the removed field or map entry.
type ObjectFilterFunc struct {
	//Kind string

	// KeepPredicate matches against the name of the field or key in the map.
	KeepPredicate func(node *yaml.RNode) bool

	// IfEmpty bool
}

func (c ObjectFilterFunc) Filter(rn *yaml.RNode) (*yaml.RNode, error) {
	if err := yaml.ErrorIfInvalid(rn, yaml.SequenceNode); err != nil {
		return nil, err
	}

	var keep []*yaml.Node

	content := rn.Content()

	for i := 0; i < len(content); i++ {
		if c.KeepPredicate(yaml.NewRNode(content[i])) {
			keep = append(keep, content[i])
		}
	}
	rn.YNode().Content = keep
	return nil, nil
}

// ObjectFilter removes an annotation at metadata.annotations.
// Returns nil if the annotation or field does not exist.
type ObjectFilter struct {
	FieldPaths []FieldPath

	// Predicate matches against the name of the field or key in the map.
	KeepPredicate func(*yaml.RNode) bool
}

func (c ObjectFilter) Filter(rn *yaml.RNode) (*yaml.RNode, error) {
	for _, fieldPath := range c.FieldPaths {
		_, err := rn.Pipe(
			yaml.PathGetter{Path: fieldPath},
			ObjectFilterFunc{KeepPredicate: c.KeepPredicate})
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
