package ptrconv

// Ptr returns pointer to value
func Ptr[T any](val T) *T {
	return &val
}
