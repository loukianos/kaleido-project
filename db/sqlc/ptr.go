package db

// Ptr returns a pointer to v, for filling nullable columns.
func Ptr[T any](v T) *T {
	return &v
}
