package nfa

import (
	"slices"
	"testing"
)

func TestSingleCharacter(t *testing.T) {
	input := "a"

	frag := SingleChar(input)
	if len(frag.Start.outs) != 1 {
		t.Fatalf("Test failed. Resulting fragment start state does not have a" +
			"single outgoing state. Should have 'a'->charState.")
	}
	if frag.Start.outs[input] == nil {
		t.Fatalf("Test failed. Resulting fragment start state does not have an" +
			"outgoing state corresponding to the character." +
			"Should have 'a'->charState.")
	}

	if frag.Start.outs[input] != frag.End[0] {
		t.Fatalf("Test failed: Fragment end state incorrect. Should be start.outs['%s']",
			input)
	}
}

func TestConcat(t *testing.T) {
	in1, in2 := "a", "b"
	f1 := SingleChar(in1)
	f2 := SingleChar(in2)
	f := Concat(f1, f2)
	if f.Start != f1.Start {
		t.Fatalf("Test failed: Concat Fragment starting state not equal to f1 starting state.")
	}

	if _, ok := f.Start.outs["a"]; !ok {
		t.Fatalf("Test failed: Init state does not have 'a' out to middle state.")
	}

	f1EndState := f.Start.outs["a"]
	if f1EndState.eOuts[0] != f2.Start {
		t.Fatalf("Test failed: f1 end state does not have e-transition to f2 start state")
	}

	if !slices.Equal(f.End, f2.End) {
		t.Fatalf("Test failed: f end states and f2 end states should match."+
			"got=%+v \t%+v", f.End, f2.End)
	}
}

func TestAlternation(t *testing.T) {
	in1, in2 := "a", "b"
	f1 := SingleChar(in1)
	f2 := SingleChar(in2)
	f := Alternate(f1, f2)

	if !slices.Contains(f.Start.eOuts, f1.Start) {
		t.Fatalf("Fail: Starting state does not have e-out to f1.start")
	}
	if !slices.Contains(f.Start.eOuts, f2.Start) {
		t.Fatalf("Fail: Starting state does not have e-out to f2.start")
	}

	if len(f1.End[0].eOuts) != 1 {
		t.Fatalf("Fail: f1 should have a single eOut.")
	}
	if f1.End[0].eOuts[0] != f.End[0] {
		t.Fatalf("Fail: f1 end state eOut should point to final fragment end state.")
	}
	if len(f2.End[0].eOuts) != 1 {
		t.Fatalf("Fail: f2 should have a single eOut.")
	}
	if f2.End[0].eOuts[0] != f.End[0] {
		t.Fatalf("Fail: f1 end state eOut should point to final fragment end state.")
	}
}

func TestCombination(t *testing.T) {
	in1, in2, in3 := "a", "b", "c"
	f1 := SingleChar(in1)
	f2 := SingleChar(in2)
	f3 := SingleChar(in3)

	fConcat := Concat(f1, f2)

	// ab|c
	f := Alternate(fConcat, f3)
	if !slices.Contains(f.Start.eOuts, f1.Start) {
		t.Fatalf("Fail: f.start should have eOut to f1.start")
	}
	if !slices.Contains(f.Start.eOuts, f3.Start) {
		t.Fatalf("Fail: f.start should have eOut to f3.start")
	}
	if slices.Contains(f.Start.eOuts, f2.Start) {
		t.Fatalf("Fail: f.start should _not_ have eOut to f2.start. This is the 2nd character " +
			"in the concatenate fragment.")
	}

	if len(fConcat.End) != 1 || len(fConcat.End[0].eOuts) != 1 {
		t.Fatalf("Fail: fConcat end state wrong. Should be a single state with single eOut")
	}

	if fConcat.End[0].eOuts[0] != f.End[0] {
		t.Fatalf("Fail: fConcat end state single eOut should be the single end state of the final fragment.")
	}

	if len(f3.End[0].eOuts) != 1 || f3.End[0].eOuts[0] != f.End[0] {
		t.Fatalf("Fail: f3 single end state should have a single eOut to the single end state of the final fragment.")
	}

}

func TestStar(t *testing.T) {
	f1 := SingleChar("a")
	star := Star(f1)

	if len(star.End) != 1 {
		t.Fatalf("Fail: Star Fragment should have only a single end state. got=%d", len(star.End))
	}
	// Star start should have e-outs to f1.start & star.end
	if !slices.Contains(star.Start.eOuts, f1.Start) {
		t.Fatalf("Fail: Star Fragment start should have eOut to inner fragment start")
	}

	if !slices.Contains(star.Start.eOuts, star.End[0]) {
		t.Fatalf("Fail: Star Fragment start should have eOut to inner fragment start")
	}
}

func TestPlus(t *testing.T) {
	f1 := SingleChar("a")
	plus := Plus(f1)

	if len(plus.End) != 1 {
		t.Fatalf("Fail: Plus fragment should have only a single end state. got=%d", len(plus.End))
	}

	// Plus start should have only a single e-out to f1.start
	if !slices.Contains(plus.Start.eOuts, f1.Start) {
		t.Fatalf("Fail: Plus Fragment start should have eOut to inner fragment start")
	}

	// Inner fragment end state should have e-out to plus start and plus end
	if !slices.Contains(f1.End[0].eOuts, plus.Start) {
		t.Fatalf("Fail: Inner fragment must have e-out to plus.Start")
	}
	if !slices.Contains(f1.End[0].eOuts, plus.End[0]) {
		t.Fatalf("Fail: Inner fragment must have e-out to plus.End")
	}
}

func TestEpsEnv(t *testing.T) {
	nfa := FromNotation("/test/eEpsEnv.txt")
	s1 := nfa[1]
	s1EpsEnv := s1.EpsEnv()

	expectedStates := []*State{nfa[1], nfa[2], nfa[3], nfa[5]}
	for _, expState := range expectedStates {
		if s1EpsEnv[expState] != true {
			t.Fatalf("S1 e-Env must contain: %p", expState)
		}
	}
}

func TestStep(t *testing.T) {
	nfa := FromNotation("/test/eEpsEnv.txt")
	s0 := nfa[0]
	startState := NfaState{s0: true}
	afterStep := Step(startState, "a")

	expectedStates := []*State{nfa[1], nfa[2], nfa[3], nfa[5]}
	for _, expState := range expectedStates {
		if afterStep[expState] != true {
			t.Fatalf("NfaState after step must contain state %p, %+v", expState, expState)
		}
	}

	expectedStates = []*State{nfa[4]}
	afterStep = Step(afterStep, "b")
	for _, expState := range expectedStates {
		if afterStep[expState] != true {
			t.Fatalf("NfaState after step must contain state %p, %+v", expState, expState)
		}
	}
}
