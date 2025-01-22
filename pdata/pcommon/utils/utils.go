package utils

func Ref[T any](v T) *T {
	return &v
}

func GetEmptyPointer[T ~[]E, E comparable](_ T) T {
	return make(T, 0)
}
