package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
)

type RemoveLabels struct {
	Labels []string
}

func RunRemoveLabels(ctx context.Context, opt RemoveLabels) error {
	predicate := func(key string) bool {
		for _, label := range opt.Labels {
			if label == key {
				return true
			}
		}
		return false
	}

	var standardOptions standardCommandOptions
	standardOptions.Filters = append(standardOptions.Filters, xform.ClearLabel(predicate))

	return runStandardCommand(ctx, standardOptions)
}

func BuildRemoveLabelsCommand(parent *cobra.Command) {
	var opt RemoveLabels

	cmd := &cobra.Command{
		Use:     "remove-labels",
		Aliases: []string{"remove-label"},
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.Labels = append(opt.Labels, args...)
			return RunRemoveLabels(cmd.Context(), opt)
		},
	}

	parent.AddCommand(cmd)
}
