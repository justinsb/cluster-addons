package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
)

type RemoveAnnotations struct {
	Annotations []string
}

func RunRemoveAnnotations(ctx context.Context, opt RemoveAnnotations) error {
	predicate := func(key string) bool {
		for _, annotation := range opt.Annotations {
			if annotation == key {
				return true
			}
		}
		return false
	}

	var standardOptions standardCommandOptions
	standardOptions.Filters = append(standardOptions.Filters, xform.ClearAnnotation(predicate))

	return runStandardCommand(ctx, standardOptions)
}

func BuildRemoveAnnotationsCommand(parent *cobra.Command) {
	var opt RemoveAnnotations

	cmd := &cobra.Command{
		Use: "remove-annotations",
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.Annotations = append(opt.Annotations, args...)
			return RunRemoveAnnotations(cmd.Context(), opt)
		},
	}

	parent.AddCommand(cmd)
}
