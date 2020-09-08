package cmd

import (
	"context"
	"fmt"
	"os"

	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"

	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
)

type standardCommandOptions struct {
	Filters  []yaml.Filter
	Matchers []xform.Matcher
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

	var out []*yaml.RNode
	for _, obj := range nodes {
		_, err = obj.Pipe(options.Filters...)
		if err != nil {
			return err
		}
		out = append(out, obj)
	}

	for _, matcher := range options.Matchers {
		var matching []*yaml.RNode

		for _, obj := range out {
			match, err := matcher(obj)
			if err != nil {
				return err
			}
			if match {
				matching = append(matching, obj)
			}
		}
		out = matching
	}

	if err := io.Write(out); err != nil {
		return err
	}
	return nil
}
