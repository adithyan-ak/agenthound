package common

// StripJSONComments removes single-line (//) and block (/* */) comments
// from JSON data while preserving slashes inside quoted strings.
func StripJSONComments(data []byte) []byte {
	var result []byte
	inString := false
	escaped := false
	i := 0

	for i < len(data) {
		if escaped {
			result = append(result, data[i])
			escaped = false
			i++
			continue
		}

		if inString {
			if data[i] == '\\' {
				escaped = true
				result = append(result, data[i])
				i++
				continue
			}
			if data[i] == '"' {
				inString = false
			}
			result = append(result, data[i])
			i++
			continue
		}

		if data[i] == '"' {
			inString = true
			result = append(result, data[i])
			i++
			continue
		}

		if data[i] == '/' && i+1 < len(data) {
			if data[i+1] == '/' {
				for i < len(data) && data[i] != '\n' {
					i++
				}
				continue
			}
			if data[i+1] == '*' {
				i += 2
				for i+1 < len(data) {
					if data[i] == '*' && data[i+1] == '/' {
						i += 2
						break
					}
					i++
				}
				continue
			}
		}

		result = append(result, data[i])
		i++
	}

	return result
}
