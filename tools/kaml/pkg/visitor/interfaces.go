package visitor

import (
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type visitor interface {
	VisitSequence(path Path, node *yaml.Node) error
	VisitMap(path Path, node *yaml.Node) error
	VisitScalar(path Path, node *yaml.Node) error
}

type Visitor struct {
}

func (v *Visitor) VisitSequence(path Path, node *yaml.Node) error {
	return nil
}

func (v *Visitor) VisitScalar(path Path, node *yaml.Node) error {
	return nil
}

func (v *Visitor) VisitMap(path Path, node *yaml.Node) error {
	return nil
}
