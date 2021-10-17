package images

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// ReplaceImage describes a transform that replaces labels.
type ReplaceImage struct {
	Replacements []ImageReplacement `json:"image"`
}

type ImageReplacement struct {
	Name    string `json:"name,omitempty"`
	NewName string `json:"newName,omitempty"`
	NewTag  string `json:"newTag,omitempty"`
}

// Run applies the transform.
func (opt ReplaceImage) Run(ctx context.Context, resourceList *framework.ResourceList) error {
	for _, item := range resourceList.Items {
		if err := opt.visit(item.YNode(), ""); err != nil {
			return err
		}
	}
	return nil
}

func (opt ReplaceImage) visit(item *yaml.Node, path string) error {
	switch item.Kind {
	case yaml.ScalarNode:
		if path == ".spec.template.spec.containers[].image" {
			s, ok := asString(item)
			if !ok {
				klog.Warningf("non-string value for image: %#v", item)
				return nil
			}
			imageTokens := strings.Split(s, ":")
			if len(imageTokens) != 2 {
				klog.Warningf("cannot parse image value %q", s)
				return nil
			}
			for _, replacement := range opt.Replacements {
				if imageTokens[0] == replacement.Name {
					item.Value = replacement.NewName + ":" + replacement.NewTag
					return nil
				}
			}
		}
		klog.V(4).Infof("path %v", path)
		return nil

	case yaml.SequenceNode:
		n := len(item.Content)
		for i := 0; i < n; i += 2 {
			v := item.Content[i]
			childPath := path + "[]"
			if err := opt.visit(v, childPath); err != nil {
				return err
			}
		}
		return nil

	case yaml.MappingNode:
		n := len(item.Content)
		if n%2 != 0 {
			return fmt.Errorf("unexpected content length in MappingNode %v", path)
		}
		for i := 0; i < n; i += 2 {
			k := item.Content[i]
			ks, ok := asString(k)
			if !ok {
				klog.Warningf("ignorning non-string MappingNode key at %v %v", path, k)
				continue
			}
			childPath := path + "." + ks
			v := item.Content[i+1]
			if err := opt.visit(v, childPath); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unhandled yaml node kind %v", item.Kind)
	}
}

func asString(n *yaml.Node) (string, bool) {
	if n.Kind != yaml.ScalarNode {
		return "", false
	}
	if n.Tag == "!!str" || n.Tag == "" {
		return n.Value, true
	}
	klog.Infof("Tag: %v", n.Tag)
	klog.Infof("Tag: %#v", n)
	return "", false
}

// AddReplaceImageCommand creates the cobra.Command.
func AddReplaceImageCommand(parent *cobra.Command) {
	var opt ReplaceImage

	cmd := &cobra.Command{
		Use: "replace-image",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				// We aim to support the kustomize image replacement syntax, but for now we restrict to a very small subset
				var replacement ImageReplacement
				tokens := strings.Split(arg, "=")
				if len(tokens) != 2 {
					return fmt.Errorf("unhandled image replacement specifier: %v", arg)
				}
				replacement.Name = tokens[0]
				newTokens := strings.Split(tokens[1], ":")
				if len(newTokens) != 2 {
					return fmt.Errorf("unhandled image replacement specifier: %v", arg)
				}
				replacement.NewName = newTokens[0]
				replacement.NewTag = newTokens[1]
				opt.Replacements = append(opt.Replacements, replacement)
			}
			return xform.RunXform(cmd.Context(), opt.Run)
		},
	}

	parent.AddCommand(cmd)
}
