package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/concat"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/images"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/normalize"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform/labels"
)

func main() {
	ctx := context.Background()
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

func BuildRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "kaml",
	}

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	labels.AddRemoveLabelsCommand(rootCmd)

	concat.AddConcatCommand(rootCmd)

	images.AddReplaceImageCommand(rootCmd)

	normalize.AddNormalizeLabelsCommand(rootCmd)

	return rootCmd
}
