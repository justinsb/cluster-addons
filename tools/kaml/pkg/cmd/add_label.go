package cmd

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
)

type AddLabels struct {
	Labels map[string]string
}

func RunAddLabels(ctx context.Context, opt AddLabels) error {
	var standardOptions standardCommandOptions
	filter := xform.AddLabels(opt.Labels)
	standardOptions.Filters = append(standardOptions.Filters, filter)

	return runStandardCommand(ctx, standardOptions)
}

func BuildAddLabelsCommand(parent *cobra.Command) {
	var opt AddLabels

	cmd := &cobra.Command{
		Use: "add-labels",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := make(map[string]string)
			for _, arg := range args {
				tokens := strings.SplitN(arg, "=", 2)
				m[tokens[0]] = tokens[1]
			}
			opt.Labels = m
			return RunAddLabels(cmd.Context(), opt)
		},
	}

	parent.AddCommand(cmd)
}
