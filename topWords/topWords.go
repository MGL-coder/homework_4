package topWords

import (
	"regexp"
	"strings"
)

// TopWords does not return error, so it will return meaningful results at specific edge cases
func TopWords(s string, n int) []string {
	if n <= 0 {
		return []string{}
	}

	exp := regexp.MustCompile("\\b[a-zA-Z]+\\b")
	frequency := make(map[string]int)

	for _, word := range exp.FindAllString(s, -1) {
		word := strings.ToLower(word)
		frequency[word]++
	}

	if len(frequency) < n {
		n = len(frequency)
	}

	words := make([]string, 0, n)
	for i := 0; i < n; i++ {
		max := 0
		topWord := ""
		for word := range frequency {
			if max < frequency[word] {
				max = frequency[word]
				topWord = word
			}
		}
		delete(frequency, topWord)
		words = append(words, topWord)
	}

	return words
}
