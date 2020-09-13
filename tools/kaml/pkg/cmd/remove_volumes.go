package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type RemoveVolumes struct {
	Volumes []string
}

func RunRemoveVolumes(ctx context.Context, opt RemoveVolumes) error {
	keepPredicate := func(node *yaml.RNode) bool {
		nameField := node.Field("name")
		if nameField.IsNilOrEmpty() {
			return true
		}

		name := yaml.GetValue(nameField.Value)

		for _, volume := range opt.Volumes {
			if volume == name {
				return false // remove
			}
		}
		return true
	}

	var standardOptions standardCommandOptions
	standardOptions.Filters = append(standardOptions.Filters, xform.RemoveVolumes(keepPredicate))

	return runStandardCommand(ctx, standardOptions)
}

func BuildRemoveVolumesCommand(parent *cobra.Command) {
	var opt RemoveVolumes

	cmd := &cobra.Command{
		Use:     "remove-volumes",
		Aliases: []string{"remove-volume"},
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.Volumes = append(opt.Volumes, args...)
			return RunRemoveVolumes(cmd.Context(), opt)
		},
	}

	parent.AddCommand(cmd)
}
