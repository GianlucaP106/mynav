package utils

func SplitStringByLength(input string, chunkSize int) []string {
	var chunks []string
	for len(input) > 0 {
		if len(input) >= chunkSize {
			chunks = append(chunks, input[:chunkSize])
			input = input[chunkSize:]
		} else {
			chunks = append(chunks, input)
			break
		}
	}
	return chunks
}
