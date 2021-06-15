package topWords

import (
	"reflect"
	"testing"
)

func TestTopWords(t *testing.T) {
	testTable := []struct {
		s        string
		n        int
		expected []string
	}{
		{s: "One two two three three three", n: -1, expected: []string{}},
		{s: "One two two three three three", n: 0, expected: []string{}},
		{s: "One two two three three three", n: 1, expected: []string{"three"}},
		{s: "One two two three three three", n: 2, expected: []string{"three", "two"}},
		{s: "One two two three three three", n: 3, expected: []string{"three", "two", "one"}},
		{s: "Leon was going... and going... till... he reached the destination.", n: 1, expected: []string{"going"}},
		{s: "", n: 1, expected: []string{}},
		{s: "    ", n: 3, expected: []string{}},
		{s: "321 213 903 2111", n: 3, expected: []string{}},
		{s: "...Going, going. Till, Till, till, Still.", n: 3, expected: []string{"till", "going", "still"}},
		{s: "Going, going. Till, Till, till, Still, a a a a a a.", n: 4, expected: []string{"a", "till", "going", "still"}},
		{s: "Going, going. Till, Till, till, Still, a a a a a a.", n: 10, expected: []string{"a", "till", "going", "still"}},
	}

	for _, testCase := range testTable {
		result := TopWords(testCase.s, testCase.n)

		t.Logf("Calling TopWords(%s, %d), result %s", testCase.s, testCase.n, result)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Incorrect result. Expect %s, got %s", testCase.expected, result)
		}
	}
}
