package normalize

import (
	"context"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/visitor"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// NormalizeLabels describes a transform that normalizes labels.
type NormalizeLabels struct {
	visitor.Visitor
}

// Run applies the transform.
func (opt NormalizeLabels) Run(ctx context.Context, resourceList *framework.ResourceList) error {
	return visitor.VisitResourceList(resourceList, &opt)
}

func (opt NormalizeLabels) VisitMap(path visitor.Path, node *yaml.Node) error {
	sortKeys := false

	switch path {
	case ".metadata.labels":
		sortKeys = true
	default:
		sortKeys = false
	}

	if !sortKeys {
		return nil
	}

	n := len(node.Content)
	if n%2 != 0 {
		return fmt.Errorf("unexpected content length in MappingNode %v", path)
	}

	var keys []string
	keyMap := make(map[string]*yaml.Node)
	valueMap := make(map[string]*yaml.Node)
	for i := 0; i < n; i += 2 {
		k := node.Content[i]
		v := node.Content[i+1]
		ks, ok := visitor.AsString(k)
		if !ok {
			return fmt.Errorf("expected string key in %v", path)
		}
		keys = append(keys, ks)
		keyMap[ks] = k
		valueMap[ks] = v
	}

	sort.Strings(keys)

	var newContent []*yaml.Node
	for _, k := range keys {
		newContent = append(newContent, keyMap[k], valueMap[k])
	}
	node.Content = newContent
	return nil
}

// AddNormalizeLabelsCommand creates the cobra.Command.
func AddNormalizeLabelsCommand(parent *cobra.Command) {
	var opt NormalizeLabels

	cmd := &cobra.Command{
		Use: "normalize-labels",
		RunE: func(cmd *cobra.Command, args []string) error {
			return xform.RunXform(cmd.Context(), opt.Run)
		},
	}

	parent.AddCommand(cmd)
}
