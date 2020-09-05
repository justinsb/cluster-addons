package xform

import "sigs.k8s.io/kustomize/kyaml/yaml"

// FieldClearerFunc removes the field or map key.
// Returns a RNode with the removed field or map entry.
type FieldClearerFunc struct {
	//Kind string

	// Predicate matches against the name of the field or key in the map.
	Predicate func(string) bool

	// IfEmpty bool
}

func (c FieldClearerFunc) Filter(rn *yaml.RNode) (*yaml.RNode, error) {
	if err := yaml.ErrorIfInvalid(rn, yaml.MappingNode); err != nil {
		return nil, err
	}

	var keep []*yaml.Node

	content := rn.Content()

	for i := 0; i < len(content); i += 2 {
		// if name matches, remove these 2 elements from the list because
		// they are treated as a fieldName/fieldValue pair.
		if !c.Predicate(content[i].Value) {
			// if c.IfEmpty {
			// 	if len(rn.Content()[i+1].Content) > 0 {
			// 		continue
			// 	}
			// }

			keep = append(keep, content[i])
			if len(content) > i+1 {
				keep = append(keep, content[i+1])
			}
			/*
				keep = append(keep, rn.Content()[i].Value)

				// save the item we are about to remove
				//removed := yaml.NewRNode(rn.Content()[i+1])
				if len(rn.YNode().Content) > i+2 {
					l := len(rn.YNode().Content)
					// remove from the middle of the list
					rn.YNode().Content = rn.Content()[:i]
					rn.YNode().Content = append(
						rn.YNode().Content,
						rn.Content()[i+2:l]...)
				} else {
					// remove from the end of the list
					rn.YNode().Content = rn.Content()[:i]
				}

				// return the removed field name and value
				//return removed, nil
			*/
		}
	}
	rn.YNode().Content = keep
	// nothing removed
	return nil, nil
}

// FieldClearer removes an annotation at metadata.annotations.
// Returns nil if the annotation or field does not exist.
type FieldClearer struct {
	FieldPaths []FieldPath

	// Predicate matches against the name of the field or key in the map.
	Predicate func(string) bool
}

func (c FieldClearer) Filter(rn *yaml.RNode) (*yaml.RNode, error) {
	for _, fieldPath := range c.FieldPaths {
		_, err := rn.Pipe(
			yaml.PathGetter{Path: fieldPath},
			FieldClearerFunc{Predicate: c.Predicate})
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
