package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
)

type RemoveNamespace struct {
}

func RunRemoveNamespace(ctx context.Context, opt RemoveNamespace) error {
	var standardOptions standardCommandOptions
	filter := xform.RemoveNamespace()
	standardOptions.Filters = append(standardOptions.Filters, filter)

	return runStandardCommand(ctx, standardOptions)
}

func BuildRemoveNamespaceCommand(parent *cobra.Command) {
	var opt RemoveNamespace

	cmd := &cobra.Command{
		Use: "remove-namespace",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("unexpected args: %v", args)
			}
			return RunRemoveNamespace(cmd.Context(), opt)
		},
	}

	parent.AddCommand(cmd)
}
