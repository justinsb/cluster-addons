package xform

func ClearAnnotation(predicate func(key string) bool) FieldClearer {
	return FieldClearer{
		FieldPaths: []FieldPath{
			ParseFieldPath("metadata.annotations"),
		},
		Predicate: predicate,
	}
}

// // ClearEmptyAnnotations clears the keys, annotations
// // and metadata if they are empty/null
// func ClearEmptyAnnotations(rn *yaml.RNode) error {
// 	_, err := rn.Pipe(yaml.Lookup("metadata"), yaml.FieldClearer{
// 		Name: "annotations", IfEmpty: true})
// 	if err != nil {
// 		return errors.Wrap(err)
// 	}
// 	_, err = rn.Pipe(yaml.FieldClearer{Name: "metadata", IfEmpty: true})
// 	if err != nil {
// 		return errors.Wrap(err)
// 	}
// 	return nil
// }
