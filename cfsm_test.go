package cfsm

import "testing"

// Tests creating new machine.
func TestNewMachine(t *testing.T) {
	sys := NewSystem()
	m0 := sys.NewMachine()
	m1 := sys.NewMachine()
	m2 := sys.NewMachine()
	if m0.ID != 0 {
		t.Errorf("Expected ID of machine 0 to be 0 (got %d).", m0.ID)
	}
	if m1.ID != 1 {
		t.Errorf("Expected ID of machine 1 to be 1 (got %d).", m1.ID)
	}
	if m2.ID != 2 {
		t.Errorf("Expected ID of machine 2 to be 2 (got %d).", m2.ID)
	}
	if len(sys.CFSMs) != 3 {
		t.Errorf("Wrong number of CFSMs in System (expect 4, got %d).", len(sys.CFSMs))
	}
	ids := make(map[int]int)
	for _, m := range sys.CFSMs {
		ids[m.ID]++
	}
	for _, count := range ids {
		if count > 1 {
			t.Error("Machine in System has duplicated ID.")
		}
	}
}

// Tests creating new state.
func TestNewState(t *testing.T) {
	sys := NewSystem()
	m := sys.NewMachine()
	st4 := m.NewFreeState()
	st0 := m.NewState()
	st1 := m.NewState()
	st3 := m.NewFreeState()
	st2 := m.NewState()
	if st0.ID != 0 {
		t.Errorf("Expected ID of state 0 to be 0 (got %d).", st0.ID)
	}
	if st1.ID != 1 {
		t.Errorf("Expected ID of state 1 to be 1 (got %d).", st1.ID)
	}
	if st2.ID != 2 {
		t.Errorf("Expected ID of state 2 to be 2 (got %d).", st2.ID)
	}
	if len(sys.CFSMs[0].states) != 3 {
		t.Errorf("Wrong number of CFSMs in System (expect 3, got %d).", len(sys.CFSMs))
	}
	if st3.ID != FreeStateID {
		t.Error("Expected ID of a free state is", FreeStateID)
	}
	if st4.ID != FreeStateID {
		t.Error("Expected ID of a free state is", FreeStateID)
	}
	m.AddState(st3)
	if len(sys.CFSMs[0].states) != 4 {
		t.Errorf("Wrong number of CFSMs in System (expect 4, got %d).", len(sys.CFSMs))
	}
	m.AddState(st4)
	if len(sys.CFSMs[0].states) != 5 {
		t.Errorf("Wrong number of CFSMs in System (expect 5, got %d).", len(sys.CFSMs))
	}
	if st3.ID != 3 {
		t.Errorf("Expected ID of state 3 to be 3 (got %d).", st3.ID)
	}
	if st4.ID != 4 {
		t.Errorf("Expected ID of state 4 to be 4 (got %d).", st4.ID)
	}
	ids := make(map[int]int)
	for _, m := range sys.CFSMs {
		ids[m.ID]++
	}
	for _, count := range ids {
		if count > 1 {
			t.Error("Machine in System has duplicated ID.")
		}
	}
}

// Tests transition.
func TestTransitions(t *testing.T) {
	sys := NewSystem()
	m0 := sys.NewMachine()
	m1 := sys.NewMachine()
	m0st0 := m0.NewState()
	m0st1 := m0.NewState()
	m1st0 := m1.NewState()
	m1st1 := m1.NewState()

	tr0 := NewSend(m1, "ZeroToOne")
	tr0.SetNext(m0st1)
	m0st0.AddTransition(tr0)
	tr1 := NewRecv(m0, "ZeroToOne")
	tr1.SetNext(m1st1)
	m1st0.AddTransition(tr1)

	m0.Start = m0st0
	m1.Start = m1st0
	if len(m0.states) != 2 {
		t.Errorf("Machine 0 in System has wrong number of states (expected %d.)", len(m0.states))
	}
	if len(m0st0.edges) != 1 {
		t.Errorf("Machine 0 State 0 has %d transitions (expected 1).", len(m0st0.edges))
	}
	if len(m0st1.edges) != 0 {
		t.Errorf("Machine 0 State 1 has %d transitions (expected 0).", len(m0st1.edges))
	}
	if len(m1.states) != 2 {
		t.Errorf("Machine 1 in System has wrong number of states (expected %d.)", len(m1.states))
	}
	if len(m1st0.edges) != 1 {
		t.Errorf("Machine 1 State 0 has %d transitions (expected 1).", len(m1st0.edges))
	}
	if len(m1st1.edges) != 0 {
		t.Errorf("Machine 1 State 1 has %d transitions (expected 0).", len(m1st1.edges))
	}
}

// Test removing last machine
func TestRemoveMachine(t *testing.T) {
	sys := NewSystem()
	m0 := sys.NewMachine()
	if m0.ID != 0 {
		t.Errorf("Machine 0 should have ID 0")
	}
	m1 := sys.NewMachine()
	if m1.ID != 1 {
		t.Errorf("Machine 1 should have ID 1")
	}
	m0st0 := m0.NewState()
	m0st1 := m0.NewState()
	tr0 := NewSend(m1, "Msg")
	tr0.SetNext(m0st1)
	m0st0.AddTransition(tr0)

	m1st0 := m1.NewState()
	m1st1 := m1.NewState()
	tr1 := NewRecv(m0, "Msg")
	tr1.SetNext(m1st1)
	m1st0.AddTransition(tr1)

	m0.Start = m0st0
	m1.Start = m1st0

	m2 := sys.NewMachine()
	if m2.ID != 2 {
		t.Errorf("Machine 2 should have ID 2")
	}
	if len(sys.CFSMs) != 3 {
		t.Errorf("Expects 3 CFSMs but got %d", len(sys.CFSMs))
	}
	sys.RemoveMachine(m2.ID)
	if len(sys.CFSMs) != 2 {
		t.Errorf("Expects 2 CFSMs but got %d", len(sys.CFSMs))
	}
}

// Test removing machine in the middle
func TestRemoveMachine2(t *testing.T) {
	sys := NewSystem()
	m0 := sys.NewMachine()
	if m0.ID != 0 {
		t.Errorf("Machine 0 should have ID 0")
	}
	m1 := sys.NewMachine()
	if m1.ID != 1 {
		t.Errorf("Machine 1 should have ID 1")
	}
	m2 := sys.NewMachine()
	if m2.ID != 2 {
		t.Errorf("Machine 2 should have ID 2")
	}
	m0st0 := m0.NewState()
	m0st1 := m0.NewState()
	tr0 := NewSend(m1, "Msg")
	tr0.SetNext(m0st1)
	m0st0.AddTransition(tr0)

	m1st0 := m1.NewState()
	m1st1 := m1.NewState()
	tr1 := NewRecv(m0, "Msg")
	tr1.SetNext(m1st1)
	m1st0.AddTransition(tr1)

	m0.Start = m0st0
	m1.Start = m1st0

	if len(sys.CFSMs) != 3 {
		t.Errorf("Expects 3 CFSMs but got %d", len(sys.CFSMs))
	}
	sys.RemoveMachine(m2.ID)
	if len(sys.CFSMs) != 2 {
		t.Errorf("Expects 2 CFSMs but got %d", len(sys.CFSMs))
	}
}

// Tests transition with loops.
func TestTransitions2(t *testing.T) {
	sys := NewSystem()
	m0 := sys.NewMachine()
	m1 := sys.NewMachine()
	m0st0 := m0.NewState()
	m0st1 := m0.NewState()
	m1st0 := m1.NewState()
	m1st1 := m1.NewState()

	tr0 := NewSend(m1, "ZeroToOne")
	tr0.SetNext(m0st1)
	m0st0.AddTransition(tr0)
	tr1 := NewSend(m1, "OneToZero")
	tr1.SetNext(m0st0)
	m0st1.AddTransition(tr1)

	tr2 := NewRecv(m0, "ZeroToOne")
	tr2.SetNext(m1st1)
	m1st0.AddTransition(tr2)
	tr3 := NewRecv(m0, "OneToZero")
	tr3.SetNext(m1st0)
	m1st1.AddTransition(tr3)

	m0.Start = m0st0
	m1.Start = m1st0

	if len(m0.states) != 2 {
		t.Errorf("Machine 0 in System has %d states (expected 2).", len(m0.states))
	}
	if len(m0st0.edges) != 1 {
		t.Errorf("Machine 0 State 0 has %d transitions (expected 1).", len(m0st0.edges))
	}
	if len(m0st1.edges) != 1 {
		t.Errorf("Machine 0 State 1 has %d transitions (expected 1).", len(m0st1.edges))
	}
	if len(m1.states) != 2 {
		t.Errorf("Machine 1 in System has %d states (expected 2).", len(m1.states))
	}
	if len(m1st0.edges) != 1 {
		t.Errorf("Machine 1 State 0 has %d transitions (expected 1).", len(m1st0.edges))
	}
	if len(m1st1.edges) != 1 {
		t.Errorf("Machine 1 State 1 has %d transitions (expected 1).", len(m1st1.edges))
	}
}

// To use the CFSM library, first create a system of CFSMs.
// From the system, create machines (i.e. CFSMs) for the system.
// Add states to the machines, and attach transitions to the states.
// Finally set initial state of each machine.
func ExampleSystem() {
	// Create a new system of CFSMs.
	sys := NewSystem()
	alice := sys.NewMachine() // CFSM Alice
	alice.Comment = "Alice"

	bob := sys.NewMachine() // CFSM Bob
	bob.Comment = "Bob"

	a0 := alice.NewState()
	a1 := alice.NewState()
	a01 := NewSend(bob, "int")
	a01.SetNext(a1)
	a0.AddTransition(a01) // Add a transition from a0 --> a1.

	b0 := bob.NewState()
	b1 := bob.NewState()
	b01 := NewRecv(alice, "int")
	b01.SetNext(b1)
	b0.AddTransition(a01) // Add a transition from b0 --> b1.

	// Set initial states of alice and bob.
	alice.Start = a0
	bob.Start = b0
	// Output:
	//
}
