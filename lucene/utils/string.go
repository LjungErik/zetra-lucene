package utils

func CommonPrefixLength(s1, s2 string) int {
	n := len(s1)
	if len(s2) < n {
		n = len(s2)
	}

	for i := range n {
		if s1[i] != s2[i] {
			return i
		}
	}

	return n
}
