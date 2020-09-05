package xform

import "sigs.k8s.io/kustomize/kyaml/yaml"

func ClearLabel(predicate func(key string) bool) FieldClearer {
	return FieldClearer{
		FieldPaths: []FieldPath{
			ParseFieldPath("metadata.labels"),
			ParseFieldPath("spec.selector"),
			ParseFieldPath("spec.selector.matchLabels"),
			ParseFieldPath("spec.template.metadata.labels"),
		},
		Predicate: predicate,
	}
}

func AddLabels(add map[string]string) FieldAdder {
	m := make(map[string]*yaml.Node)
	for k, v := range add {
		m[k] = StringNode(v)
	}

	return FieldAdder{
		FieldPaths: []FieldPath{
			ParseFieldPath("metadata.labels"),
			ParseFieldPath("spec.selector"),
			ParseFieldPath("spec.selector.matchLabels"),
			ParseFieldPath("spec.template.metadata.labels"),
		},
		Add: m,
	}
}

// fieldAdderFilter sets the specified fields if they are not found.
type fieldAdderFilter struct {
	add map[string]*yaml.Node
}

func (c fieldAdderFilter) Filter(rn *yaml.RNode) (*yaml.RNode, error) {
	if err := yaml.ErrorIfInvalid(rn, yaml.MappingNode); err != nil {
		return nil, err
	}

	var newContents []*yaml.Node

	content := rn.Content()

	done := make(map[string]bool)

	for i := 0; i < len(content); i += 2 {
		key := content[i]
		done[key.Value] = true

		newValue, found := c.add[key.Value]
		if found {
			newContents = append(newContents, key)
			newContents = append(newContents, newValue)
		} else {
			newContents = append(newContents, key)
			if len(content) > i+1 {
				newContents = append(newContents, content[i+1])
			}
		}
	}

	for k, v := range c.add {
		if done[k] {
			continue
		}
		newContents = append(newContents, StringNode(k))
		newContents = append(newContents, v)
	}

	rn.YNode().Content = newContents
	return nil, nil
}

func StringNode(s string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Value: s}
}

// FieldAdder sets or adds fields or map entries.
type FieldAdder struct {
	FieldPaths []FieldPath

	// Add contains the fields to be set.
	Add map[string]*yaml.Node
}

func (c FieldAdder) Filter(rn *yaml.RNode) (*yaml.RNode, error) {
	for _, fieldPath := range c.FieldPaths {
		_, err := rn.Pipe(
			yaml.PathGetter{Path: fieldPath},
			fieldAdderFilter{add: c.Add})
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
