package utils

func Map[T any](array []T, cb func(index int, element T) T) []T {
	res := make([]T, 0, len(array)+1)
	for i, e := range array {
		res = append(res, cb(i, e))
	}
	return res
}

func ForEach[T any](array []T, cb func(index int, element T)) {
	for i, e := range array {
		cb(i, e)
	}
}
