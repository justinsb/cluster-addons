package xform

func RemoveNamespace() FieldClearer {
	return FieldClearer{
		FieldPaths: []FieldPath{
			ParseFieldPath("metadata"),
		},
		Predicate: func(name string) bool {
			return name == "namespace"
		},
	}
}
