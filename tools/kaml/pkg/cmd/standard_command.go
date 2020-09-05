package cmd

import (
	"context"
	"fmt"
	"os"

	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type standardCommandOptions struct {
	Filters []yaml.Filter
}

func runStandardCommand(ctx context.Context, options standardCommandOptions) error {
	io := kio.ByteReadWriter{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}

	nodes, err := io.Read()
	if err != nil {
		return fmt.Errorf("failed to parse yaml: %w", err)
	}

	var filtered []*yaml.RNode
	for _, obj := range nodes {
		_, err = obj.Pipe(options.Filters...)
		if err != nil {
			return err
		}
		filtered = append(filtered, obj)
	}

	if err := io.Write(filtered); err != nil {
		return err
	}
	return nil
}
