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
		{"aaaabbcc", "a*b*c*d*"},
	}
	for i := range inputs {
		if matches := Match(inputs[i].text, inputs[i].regex); !matches {
			t.Fatalf("fail[%d]. String:%s\tRegex:%s",
				i,
				inputs[i].text,
				inputs[i].regex,
			)
		}

	}
}
