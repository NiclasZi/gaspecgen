package util

func PtrOf[T any](val T) *T {
	return &val
}
