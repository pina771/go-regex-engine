package nfa

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type State struct {
	outs  map[string]*State
	eOuts []*State

	name string
}

func newState() *State {
	return &State{
		outs:  make(map[string]*State),
		eOuts: make([]*State, 0),
	}
}

type Fragment struct {
	Start *State   // Starting fragment state
	End   []*State // Finish states of the fragment
}

// Create a e-NFA for a single character, where i=initial_state and f=final_state
//
// i --a-> f
func SingleChar(char string) *Fragment {
	init, charState := newState(), newState()
	init.outs[char] = charState
	return &Fragment{
		Start: init,
		End:   []*State{charState},
	}
}

// Creates a new fragment by connecting f1 and f2. Adds an e-transition to
// f2.start for every state in f1.end. Return fragment has f1.start as beginning state
// and f2.end as finish state.
func Concat(f1 *Fragment, f2 *Fragment) *Fragment {
	for _, state := range f1.End {
		if !slices.Contains(state.eOuts, f2.Start) {
			state.eOuts = append(state.eOuts, f2.Start)
		}
	}

	return &Fragment{
		Start: f1.Start,
		End:   f2.End,
	}
}

func Alternate(f1 *Fragment, f2 *Fragment) *Fragment {
	newStart := State{
		outs: map[string]*State{},
		eOuts: []*State{
			f1.Start, f2.Start,
		},
	}
	newFinish := State{
		outs:  map[string]*State{},
		eOuts: []*State{},
	}

	for _, state := range f1.End {
		state.eOuts = append(state.eOuts, &newFinish)
	}
	for _, state := range f2.End {
		state.eOuts = append(state.eOuts, &newFinish)
	}
	return &Fragment{
		Start: &newStart,
		End:   []*State{&newFinish},
	}
}

func Star(frag *Fragment) *Fragment {
	newEnd := newState()
	newStart := State{
		outs: map[string]*State{},
		eOuts: []*State{
			frag.Start,
			newEnd,
		},
	}
	for _, endState := range frag.End {
		endState.eOuts = append(endState.eOuts, &newStart, newEnd)
	}

	return &Fragment{
		Start: &newStart,
		End:   []*State{newEnd},
	}
}

func Plus(frag *Fragment) *Fragment {
	newEnd := newState()
	// Difference between Plus and Star is in the starting point. In Plus, the new start does not
	// have a e-transition to the end state, effectively forcing entry into the inner fragment.
	newStart := State{
		outs: map[string]*State{},
		eOuts: []*State{
			frag.Start,
		},
	}
	for _, endState := range frag.End {
		endState.eOuts = append(endState.eOuts, &newStart, newEnd)
	}

	return &Fragment{
		Start: &newStart,
		End:   []*State{newEnd},
	}
}

type NfaState map[*State]bool

// Determines the next NFA State based on current state and an input.
func Step(cState NfaState, in string) NfaState {
	eEnvInitial := make(NfaState, 0)
	for k := range cState {
		eEnv := k.EpsEnv()
		maps.Copy(eEnvInitial, eEnv)
	}

	eEnvMid := make(NfaState, 0)
	for state := range eEnvInitial {
		// dot transition - special case
		if dotState, exists := state.outs["."]; exists {
			eEnvMid[dotState] = true
		}

		if state, exists := state.outs[in]; exists {
			eEnvMid[state] = true
		}
	}

	// NOTE: Is this necessary ?
	eEnvFinal := make(NfaState, 0)
	for state := range eEnvMid {
		eEnv := state.EpsEnv()
		maps.Copy(eEnvFinal, eEnv)
	}
	return eEnvFinal
}

// Calculates the e-environment of the current state. Returns a list of states that
// are possible to access using only e-transitions.
func (s *State) EpsEnv() map[*State]bool {
	eEnvMap := make(map[*State]bool)
	s.expand(eEnvMap)
	return eEnvMap
}

func (s *State) expand(curr map[*State]bool) {
	if curr[s] {
		return
	}

	curr[s] = true
	for _, eout := range s.eOuts {
		if !curr[eout] {
			eout.expand(curr)
		}
	}
}

// Constructs an e-NFA from a text file in the correct notation.
// Each line represents a single state which consists of the following, separated by '|':
// 1. regular/character transitino
// 2. epsilon transitions
// E.g. {"a":1}|[] is a single state with "a" out to the element with index 0.
func FromNotation(path string) []*State {
	cwd, _ := os.Getwd()
	fullPath := filepath.Join(cwd, path)
	fmt.Println("Constructing NFA array from:", fullPath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	states := make([]*State, len(lines))
	for idx := range len(lines) {
		states[idx] = newState()
		states[idx].name = fmt.Sprintf("s%d", idx)
	}

	for cStateIdx, line := range lines {
		parts := strings.Split(line, "|")

		var standardOuts map[string]int
		err := json.Unmarshal([]byte(parts[0]), &standardOuts)
		if err != nil {
			fmt.Println("failed unmarshalling of standardOuts", err)
			os.Exit(1)
		}
		for k, v := range standardOuts {
			// fmt.Printf("state %p, setting out {%s:  %p}\n",
			// 	states[cStateIdx], k, states[v])
			states[cStateIdx].outs[k] = states[v]
		}

		var eOuts []int
		err = json.Unmarshal([]byte(parts[1]), &eOuts)
		if err != nil {
			fmt.Println("failed unmarshalling of eOuts", err)
			os.Exit(1)
		}
		for _, val := range eOuts {
			// fmt.Printf("state %p, adding eOut: %p\n", states[cStateIdx], states[val])
			states[cStateIdx].eOuts = append(states[cStateIdx].eOuts, states[val])
		}
	}

	for idx := range len(lines) {
		fmt.Printf("s%d: %p: %+v\n", idx, states[idx], states[idx])
	}
	return states
}
