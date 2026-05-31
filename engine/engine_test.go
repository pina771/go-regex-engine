package engine

import "testing"

func TestMatch(t *testing.T) {

	// Inline structs
	inputs := []struct {
		text        string
		regex       string
		shouldMatch bool
	}{
		{"c", "ab|c", true},
		{"ab", "ab|c", true},
		{"there", "hello|there", true},
		{"axb", "a.b", true},
		{"bLt", "b.t|a", true},
		{"tera", "....", true},
		{"xtaaab", "xta*b", true},
		{"aaaabbcc", "a*b*c*d*", true},
		{"bc", "a+bc", false},
		{"aaabc", "a+bc", true},
		{"abc", "a+bc", true},
		{"aabb", "a+b+", true},
		{"aabbc", "a+b+c", true},
		{"aabb", "a+b+c", false},
	}
	for i := range inputs {
		if matches := Match(inputs[i].text, inputs[i].regex); matches != inputs[i].shouldMatch {
			t.Fatalf("fail[%d]. String:%s\tRegex:%s",
				i,
				inputs[i].text,
				inputs[i].regex,
			)
		}

	}
}
