package j

// Property represents JSON property
type Property struct {
	name string
	val  interface{}
}

// Prop creates JSON property
func Prop(name string, val interface{}) *Property {
	return &Property{
		name: name,
		val:  val,
	}
}

// Array creates slice
func Array(values ...interface{}) []interface{} {
	return values
}

// Obj creates map
func Obj(props ...*Property) map[string]interface{} {
	result := make(map[string]interface{}, len(props))

	for _, p := range props {
		result[p.name] = p.val
	}

	return result
}
