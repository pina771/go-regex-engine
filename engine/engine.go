package engine

import (
	"pina771/regex-eng/lexer"
	"pina771/regex-eng/nfa"
	"pina771/regex-eng/parser"
	"slices"
)

func Match(text, regex string) bool {
	parser := parser.New(lexer.New(regex))
	machine := parser.ToNFA()
	endStates := machine.End

	cNfaState := nfa.NfaState{machine.Start: true}
	for _, v := range text {
		vString := string(v)

		// IMPORTANT: GREŠKA, SHADOWING, 100 Go mistakes to avoid
		// cNfaState := nfa.Step(cNfaState, vString)
		//
		// CORRECT:
		cNfaState = nfa.Step(cNfaState, vString)
		if len(cNfaState) == 0 {
			break
		}
	}
	for state := range cNfaState {
		if slices.Contains(endStates, state) {
			return true
		}
	}
	return false
}
