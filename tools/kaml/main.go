package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/cmd"
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

	cmd.BuildAddLabelsCommand(rootCmd)
	cmd.BuildFilterCommand(rootCmd)
	cmd.BuildRemoveAnnotationsCommand(rootCmd)
	cmd.BuildRemoveLabelsCommand(rootCmd)
	cmd.BuildRemoveNamespaceCommand(rootCmd)
	cmd.BuildRemoveVolumesCommand(rootCmd)

	return rootCmd
}
