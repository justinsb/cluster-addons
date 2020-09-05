package xform

import "strings"

type FieldPath []string

func ParseFieldPath(k string) FieldPath {
	return FieldPath(strings.Split(k, "."))
}
