package util

func UpdateMap[K comparable, V any](target map[K]V, updates ...map[K]V) map[K]V {
	for _, update := range updates {
		for k, v := range update {
			target[k] = v
		}
	}
	return target
}
