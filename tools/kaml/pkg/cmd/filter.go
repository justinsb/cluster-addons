package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
)

type FilterOptions struct {
	xform.FilterOptions
}

func RunFilter(ctx context.Context, opt FilterOptions) error {
	var standardOptions standardCommandOptions
	matcher := xform.Filter(opt.FilterOptions)
	standardOptions.Matchers = append(standardOptions.Matchers, matcher)

	return runStandardCommand(ctx, standardOptions)
}

func BuildFilterCommand(parent *cobra.Command) {
	var opt FilterOptions

	cmd := &cobra.Command{
		Use: "filter",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("unexpected arguments: %v", args)
			}
			return RunFilter(cmd.Context(), opt)
		},
	}

	cmd.Flags().BoolVarP(&opt.Invert, "invert-match", "v", opt.Invert, "select non-matching objects")
	cmd.Flags().StringSliceVar(&opt.Kinds, "kind", opt.Kinds, "select objects of the specified kind")

	parent.AddCommand(cmd)
}
