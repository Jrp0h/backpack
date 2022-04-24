package utils

func JoinWithSeparatorAndLast(values []string, separator string, last string) string {
	lastIndex := len(values) - 1
	result := ""

	for i, text := range values {
		switch {
		case i == 0:
			result += text

		case i == lastIndex:
			result += " " + last + " " + text

		default:
			result += separator + " " + text
		}
	}

	return result

}

func JoinSliceWithSeparatorAndLast[T any](values []T, getString func(T) string, separator string, last string) string {
	s := make([]string, len(values))

	for i, v := range values {
		s[i] = getString(v)
	}

	return JoinWithSeparatorAndLast(s, separator, last)
}

func JoinSliceAsSentanceQuestion[T any](values []T, getString func(T) string) string {
	return JoinSliceWithSeparatorAndLast(values, getString, ",", "or")
}

func JoinSliceAsSentanceStatement[T any](values []T, getString func(T) string) string {
	return JoinSliceWithSeparatorAndLast(values, getString, ",", "and")
}

func JoinAsSentanceQuestion(values []string) string {
	return JoinWithSeparatorAndLast(values, ",", "or")
}

func JoinAsSentanceStatement(values []string) string {
	return JoinWithSeparatorAndLast(values, ",", "and")
}
