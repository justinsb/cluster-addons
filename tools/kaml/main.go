package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform/annotations"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform/labels"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
)

func main() {
	ctx := context.Background()
	if cmd := os.Getenv("EXEC_KRM_FUNCTION"); cmd != "" {
		if err := runKrmFunction(ctx, cmd); err != nil {
			fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
			os.Exit(1)
		}
		return
	}
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	rootCmd := BuildRootCommand()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		return err
	}

	return nil
}

func runKrmFunction(ctx context.Context, commandName string) error {
	resourceList := &framework.ResourceList{}

	switch commandName {
	case "set-annotations":
		resourceList.FunctionConfig = &annotations.SetAnnotations{}

	case "remove-labels":
		resourceList.FunctionConfig = &labels.RemoveLabel{}

	default:
		return fmt.Errorf("unknown KRM command %q", commandName)
	}

	cmd := framework.Command(resourceList, func() error {
		xform := resourceList.FunctionConfig.(xform.Runnable)
		if err := xform.Run(ctx, resourceList); err != nil {
			return err
		}
		return nil
	})
	if err := cmd.Execute(); err != nil {
		return err
	}
	return nil
}

func BuildRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "kaml",
	}

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	labels.AddRemoveLabelsCommand(rootCmd)
	annotations.AddSetAnnotationsCommand(rootCmd)

	return rootCmd
}
