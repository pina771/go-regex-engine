package engine

import "testing"

func TestMatch(t *testing.T) {

	// Inline structs
	inputs := []struct {
		text  string
		regex string
	}{
		{"c", "ab|c"},
		{"ab", "ab|c"},
		{"there", "hello|there"},
		{"axb", "a.b"},
		{"bLt", "b.t|a"},
		{"tera", "...."},
		{"xtaaab", "xta*b"},
	}
	for _, input := range inputs {
		if matches := Match(input.text, input.regex); !matches {
			t.Fatalf("Input %s should match regex %s. Got=%v",
				input.text, input.regex, matches)
		}

	}
}
