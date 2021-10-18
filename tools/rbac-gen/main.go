package main

import (
	"bytes"
	"context"
	goflags "flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"sigs.k8s.io/cluster-addons/tools/rbac-gen/pkg/convert"
	"sigs.k8s.io/yaml"
)

func main() {
	err := run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	rootCommand := BuildGenerateCommand(ctx)

	fs := goflags.NewFlagSet("", goflags.PanicOnError)
	klog.InitFlags(fs)
	rootCommand.PersistentFlags().AddGoFlagSet(fs)

	rootCommand.SilenceErrors = true
	rootCommand.SilenceUsage = true

	if err := rootCommand.Execute(); err != nil {
		return err
	}
	return nil
}

func BuildGenerateCommand(ctx context.Context) *cobra.Command {
	yamlFile := "manifest.yaml"
	out := ""

	var opt convert.BuildRoleOptions
	opt.Name = "generated-role"
	opt.Namespace = "kube-system"

	cmd := &cobra.Command{
		Use: "generate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunGenerate(ctx, yamlFile, out, opt)
		},
	}

	cmd.Flags().StringVar(&yamlFile, "yaml", yamlFile, "yaml file from which the rbac will be generated.")
	cmd.Flags().StringVar(&opt.Name, "name", opt.Name, "name of role to be generated")
	cmd.Flags().StringVar(&opt.ServiceAccountName, "sa-name", opt.ServiceAccountName, "name of service account the role should be bound to")
	cmd.Flags().StringVar(&opt.Namespace, "ns", opt.Namespace, "namespace of the role to be generated")
	cmd.Flags().StringVar(&out, "out", out, "name of output file")
	cmd.Flags().BoolVar(&opt.Supervisory, "supervisory", opt.Supervisory, "outputs role for operator in supervisory mode")
	cmd.Flags().StringVar(&opt.CRD, "crd", opt.CRD, "CRD to generate")
	cmd.Flags().BoolVar(&opt.LimitResourceNames, "limit-resource-names", opt.LimitResourceNames, "Limit to resource names in the manifest")

	return cmd
}

func RunGenerate(ctx context.Context, yamlFile string, out string, opt convert.BuildRoleOptions) error {
	//	read yaml file passed in from cmd
	in := ""
	if yamlFile == "-" {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		in = string(b)
	} else {
		b, err := ioutil.ReadFile(yamlFile)
		if err != nil {
			return err
		}
		in = string(b)
	}

	// generate Group and Kind

	objects, err := convert.BuildRole(in, opt)
	if err != nil {
		return err
	}

	y, err := toYAML(objects)
	if err != nil {
		return err
	}

	var conv KubebuilderConverter
	if err := conv.VisitObjects(objects); err != nil {
		return err
	}

	if out == "" {
		_, err = os.Stdout.Write(y)
	} else {
		err = ioutil.WriteFile(out, y, 0644)
	}

	return err
}

func toYAML(objects []runtime.Object) ([]byte, error) {
	var buf bytes.Buffer

	for i, obj := range objects {
		if i != 0 {
			buf.WriteString("\n---\n\n")
		}
		b, err := yaml.Marshal(obj)
		if err != nil {
			return nil, err
		}

		buf.Write(b)
	}

	return buf.Bytes(), nil
}

type KubebuilderConverter struct {
}

func (c *KubebuilderConverter) VisitObjects(objects []runtime.Object) error {
	for _, obj := range objects {
		if clusterRole, ok := obj.(*v1.ClusterRole); ok {
			if err := c.visitClusterRole(clusterRole); err != nil {
				return err
			}
			continue
		}
		return fmt.Errorf("unhandled type %T", obj)
	}
	return nil
}

func (c *KubebuilderConverter) visitClusterRole(obj *v1.ClusterRole) error {
	for _, rule := range obj.Rules {
		//+kubebuilder:rbac:groups=addons.kope.io,resources=networkings,verbs=get;list;watch;update;patch
		//+kubebuilder:rbac:groups=addons.kope.io,resources=networkings/status,verbs=get;update;patch

		def := "//+kubebuilder:rbac:"
		if len(rule.APIGroups) != 0 {
			def += "groups=" + strings.Join(rule.APIGroups, ";")
		}
		if len(rule.Resources) != 0 {
			def += ",resources=" + strings.Join(rule.Resources, ";")
		}
		if len(rule.ResourceNames) != 0 {
			def += ",resourceNames=" + strings.Join(rule.ResourceNames, ";")
		}
		if len(rule.Verbs) != 0 {
			def += ",verbs=" + strings.Join(rule.Verbs, ";")
		}
		if len(rule.NonResourceURLs) != 0 {
			def += ",urls=" + strings.Join(rule.NonResourceURLs, ";")
		}

		klog.Infof("rule: %s", def)
	}

	return nil
}
