package util

func UpdateMap[K comparable, V any](target map[K]V, update map[K]V) map[K]V {
	for k, v := range update {
		target[k] = v
	}

	return target
}
