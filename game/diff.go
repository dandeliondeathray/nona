package game

// Diff returns all letters that are in b but not in a.
func Diff(word, puzzle string) (string, string) {
	if word == puzzle {
		return "", ""
	}
	return diff(word, puzzle), diff(puzzle, word)
}

func diff(a, b string) string {
	lettersInA := make(map[rune]int)
	for _, r := range []rune(a) {
		lettersInA[r] = lettersInA[r] + 1
	}

	for _, r := range []rune(b) {
		occurrences := lettersInA[r]
		if occurrences > 0 {
			lettersInA[r] = occurrences - 1
		}
	}

	result := []rune{}
	for r, occurrences := range lettersInA {
		for i := 0; i < occurrences; i++ {
			result = append(result, r)
		}
	}

	return sortLetters(string(result))
}
