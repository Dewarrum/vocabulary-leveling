package utils

func Filter[T any](slc []T, f func(T) bool) (res []T) {
	for _, s := range slc {
		if f(s) {
			res = append(res, s)
		}
	}
	return
}
