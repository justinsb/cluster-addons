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
	fmt.Fprintf(os.Stderr, "Hello %q\n", os.Getenv("EXEC_KRM_FUNCTION"))
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

	m := make(map[string]interface{})
	resourceList.FunctionConfig = m

	cmd := framework.Command(resourceList, func() error {
		fmt.Fprintf(os.Stderr, "resourceList.FunctionConfig %+v\n", m)
		// j, err := json.Marshal(m["data"])
		// if err != nil {
		// 	return fmt.Errorf("error marshalling json: %w", err)
		// }

		var runnable xform.Runnable
		switch commandName {
		case "set-annotations":
			a := &annotations.SetAnnotations{}
			a.Annotations = make(map[string]string)
			for k, v := range m["data"].(map[string]interface{}) {
				a.Annotations[k] = v.(string)
			}
			runnable = a
		case "remove-labels":
			runnable = &labels.RemoveLabel{}

		default:
			return fmt.Errorf("unknown KRM command %q", commandName)
		}

		// if err := json.Unmarshal(j, runnable); err != nil {
		// 	return fmt.Errorf("error unmarshalling json: %w", err)
		// }
		fmt.Fprintf(os.Stderr, "runnable %+v\n", runnable)
		if err := runnable.Run(ctx, resourceList); err != nil {
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
