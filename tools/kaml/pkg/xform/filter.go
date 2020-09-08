package xform

import "sigs.k8s.io/kustomize/kyaml/yaml"

type FilterOptions struct {
	Invert bool
	Kinds  []string
}

type Matcher func(*yaml.RNode) (bool, error)

func Filter(options FilterOptions) Matcher {
	return func(node *yaml.RNode) (bool, error) {
		match := true
		meta, err := node.GetMeta()
		if err != nil {
			return false, err
		}
		if match {
			if len(options.Kinds) != 0 {
				found := false
				for _, kind := range options.Kinds {
					if kind == meta.Kind {
						found = true
					}
				}
				if !found {
					match = false
				}
			}
		}
		if options.Invert {
			match = !match
		}
		return match, nil
	}
}
