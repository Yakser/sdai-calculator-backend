package typing

// Ref returns a pointer to the value
func Ref[T any](value T) *T {
	return &value
}

// DerefOrZero returns the dereference value, or zero value if it was nil.
func DerefOrZero[P *V, V any](ptr P) V {
	if ptr == nil {
		return Zero[V]()
	}
	return *ptr
}

// DerefOrDefault returns the dereferenced value, or passed default if it was nil.
func DerefOrDefault[P *V, V any](ptr P, defaultV V) V {
	if ptr == nil {
		return defaultV
	}
	return *ptr
}
