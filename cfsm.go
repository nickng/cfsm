package cfsm

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/nickng/cfsm/petrify"
)

const (
	// FreeStateID is the ID used to identify a State unattached to any CFSM.
	FreeStateID = -1
)

// System is a set of CFSMs.
type System struct {
	sync.Mutex
	CFSMs   []*CFSM // Individual CFSMs in the communicating system.
	Comment string  // Comments on the System.
}

// NewSystem returns a new communicating system
func NewSystem() *System {
	return &System{CFSMs: make([]*CFSM, 0)}
}

// NewMachine creates a new CFSM in the communicating system and returns it.
func (s *System) NewMachine() *CFSM {
	s.Lock()
	defer s.Unlock()
	cfsm := &CFSM{ID: len(s.CFSMs)}
	s.CFSMs = append(s.CFSMs, cfsm)
	return cfsm
}

// RemoveMachine removes a CFSM with the given id from System.
func (s *System) RemoveMachine(id int) {
	s.Lock()
	defer s.Unlock()
	removed := 0
	for i := range s.CFSMs {
		if s.CFSMs[i-removed].ID == id {
			s.CFSMs = append(s.CFSMs[:i-removed], s.CFSMs[i-removed+1:]...)
			removed++
		}
	}
	s.CFSMs = s.CFSMs[:len(s.CFSMs)-removed]
	for i := range s.CFSMs {
		s.CFSMs[i].ID = i
	}
}

func (s *System) String() string {
	var buf bytes.Buffer
	for _, cfsm := range s.CFSMs {
		buf.WriteString(cfsm.String())
	}
	return buf.String()
}

// CFSM is a single Communicating Finite State Machine.
type CFSM struct {
	ID      int    // Unique identifier.
	Start   *State // Starting state of the CFSM.
	Comment string // Comments on the CFSM.

	states []*State // States in a CFSM.
}

// NewState creates a new State for this CFSM.
func (m *CFSM) NewState() *State {
	state := &State{ID: len(m.states), edges: make(map[Transition]*State)}
	m.states = append(m.states, state)
	return state
}

// NewFreeState creates a new free State for this CFSM.
func (m *CFSM) NewFreeState() *State {
	state := &State{ID: FreeStateID, edges: make(map[Transition]*State)}
	return state
}

// AddState adds an unattached State to this CFSM.
func (m *CFSM) AddState(s *State) {
	if s.ID == FreeStateID {
		s.ID = len(m.states)
		m.states = append(m.states, s)
	} else {
		log.Fatal("CFSM AddState failed:", ErrStateAlias)
	}
}

// States return states defined in the machine.
func (m *CFSM) States() []*State {
	return m.states
}

// IsEmpty returns true if there are no transitions in the CFSM.
func (m *CFSM) IsEmpty() bool {
	return len(m.states) == 0 || (len(m.states) == 1 && len(m.states[0].edges) == 0)
}

func (m *CFSM) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("\n-- Machine #%d\n", m.ID))
	buf.WriteString(fmt.Sprintf("-- %s\n", strings.Replace(m.Comment, "\n", "\n--", -1)))
	buf.WriteString(fmt.Sprintf(".outputs\n.state graph\n"))
	for _, st := range m.states {
		for tr, st2 := range st.edges {
			buf.WriteString(fmt.Sprintf("q%d%d %s q%d%d\n",
				m.ID, st.ID, petrify.Encode(tr.Label()), m.ID, st2.ID))
		}
	}
	if m.Start == nil {
		buf.WriteString("-- Start state not set\n")
	} else {
		buf.WriteString(fmt.Sprintf(".marking q%d%d\n", m.ID, m.Start.ID))
	}
	buf.WriteString(".end\n")
	return buf.String()
}

// State is a state.
type State struct {
	ID    int    // Unique identifier.
	Label string // Free form text label.

	edges map[Transition]*State
}

// NewState creates a new State independent from any CFSM.
func NewState() *State {
	return &State{ID: -1, edges: make(map[Transition]*State)}
}

// Name of a State is a unique string to identify the State.
func (s *State) Name() string {
	return fmt.Sprintf("q%d", s.ID)
}

// AddTransition adds a transition to the current State.
func (s *State) AddTransition(t Transition) {
	s.edges[t] = t.State()
}

// Transition is a transition from a State to another State.
type Transition interface {
	Label() string // Label is the marking on the transition.
	State() *State // State after transition.
}

// Send is a send transition (output).
type Send struct {
	to    *CFSM  // Destination CFSM.
	msg   string // Payload message.
	state *State // State after transition.
}

// NewSend returns a new Send transition.
func NewSend(cfsm *CFSM, msg string) *Send {
	return &Send{to: cfsm, msg: msg}
}

// Label for Send is "!"
func (s *Send) Label() string {
	if s.state == nil {
		log.Fatal("Cannot get Label for Send:", ErrStateUndef)
	}
	return fmt.Sprintf("%d ! %s", s.to.ID, s.msg)
}

// State returns the State after transition.
func (s *Send) State() *State {
	return s.state
}

// SetNext sets the next state of the Send transition.
func (s *Send) SetNext(st *State) {
	s.state = st
}

// Recv is a receive transition (input).
type Recv struct {
	from  *CFSM  // Source CFSM.
	msg   string // Payload message expected.
	state *State // State after transition.
}

// NewRecv returns a new Recv transition.
func NewRecv(cfsm *CFSM, msg string) *Recv {
	return &Recv{from: cfsm, msg: msg}
}

// Label for Recv is "?"
func (r *Recv) Label() string {
	if r.state == nil {
		log.Fatal("Cannot get Label for Recv:", ErrStateUndef)
	}
	return fmt.Sprintf("%d ? %s", r.from.ID, r.msg)
}

// State returns the State after transition.
func (r *Recv) State() *State {
	return r.state
}

// SetNext sets the next state of the Recv transition.
func (r *Recv) SetNext(st *State) {
	r.state = st
}
