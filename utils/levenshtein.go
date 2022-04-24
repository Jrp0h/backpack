package utils

import (
	"strings"
)

// taken from https://github.com/spf13/cobra/blob/eb3b6397b1b5d1b0a2cd66a9afe0520f480c0a87/cobra.go#L165
// ld compares two strings and returns the levenshtein distance between them.
func ld(s, t string, ignoreCase bool) int {
	if ignoreCase {
		s = strings.ToLower(s)
		t = strings.ToLower(t)
	}
	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(t)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}

	}
	return d[len(s)][len(t)]
}

type LevenshteinResult struct {
	Text     string
	Distance int
}

type LevenshteinResults []LevenshteinResult

func Levenshtein(needle string, haystack []string, ignoreCase bool) LevenshteinResults {
	Log.Debug("utils/levenshtein: Ignoring case %t", ignoreCase)
	values := make([]LevenshteinResult, 0)

	lowest := MaxInt

	for _, text := range haystack {
		distance := ld(needle, text, ignoreCase)
		result := LevenshteinResult{
			Text:     text,
			Distance: distance,
		}

		Log.Debug("utils/levenshtein: %s has a distance of %d to %s", needle, result.Distance, text)

		switch {
		case distance < lowest:
			values = make([]LevenshteinResult, 1)
			values[0] = result
			lowest = distance

		case distance == lowest:
			values = append(values, result)
		}

	}

	return values
}

func (l LevenshteinResults) AsStatement() string {
	return JoinSliceAsSentanceStatement(l, func(r LevenshteinResult) string {
		return r.Text
	})
}

func (l LevenshteinResults) AsQuestion() string {
	return JoinSliceAsSentanceQuestion(l, func(r LevenshteinResult) string {
		return r.Text
	})
}
