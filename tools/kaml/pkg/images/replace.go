package images

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/visitor"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// ReplaceImage describes a transform that replaces labels.
type ReplaceImage struct {
	visitor.Visitor

	Replacements []ImageReplacement `json:"image"`
}

type ImageReplacement struct {
	Name    string `json:"name,omitempty"`
	NewName string `json:"newName,omitempty"`
	NewTag  string `json:"newTag,omitempty"`
}

// Run applies the transform.
func (opt ReplaceImage) Run(ctx context.Context, resourceList *framework.ResourceList) error {
	return visitor.VisitResourceList(resourceList, &opt)
}

func (opt ReplaceImage) VisitScalar(ctx *visitor.Context, path visitor.Path, node *yaml.Node) error {
	if path == ".spec.template.spec.containers[].image" {
		s, ok := visitor.AsString(node)
		if !ok {
			klog.Warningf("non-string value for image: %#v", node)
			return nil
		}
		imageTokens := strings.Split(s, ":")
		if len(imageTokens) != 2 {
			klog.Warningf("cannot parse image value %q", s)
			return nil
		}
		for _, replacement := range opt.Replacements {
			if imageTokens[0] == replacement.Name {
				node.Value = replacement.NewName + ":" + replacement.NewTag
				return nil
			}
		}
	}
	klog.V(4).Infof("path %v", path)
	return nil
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
