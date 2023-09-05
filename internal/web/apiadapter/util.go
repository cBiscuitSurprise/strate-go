package apiadapter

func MapConvert[I interface{}, O interface{}](input []*I, convert func(*I) *O) []*O {
	output := make([]*O, len(input))
	for i, original := range input {
		output[i] = convert(original)
	}

	return output
}
