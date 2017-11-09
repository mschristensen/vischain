package core

func ConcatBytes(a []byte, b []byte) []byte {
	var result []byte
	for _, i := range a {
		result = append(result, i)
	}
	for _, i := range b {
		result = append(result, i)
	}
	return result
}
