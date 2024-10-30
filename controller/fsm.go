package controller

// State interface to define required state functions.
type State interface {
	Enter()  // When entering a state
	Update() // Update logic for each state
	Exit()   // When exiting a state
}

// FSM struct to manage states, includes a reference to the player.
type FSM struct {
	Current  State
	UserData any // Reference to the player
}

// SetState transitions to a new state.
func (fsm *FSM) SetState(newState State) {
	if fsm.Current != nil {
		fsm.Current.Exit()
	}
	fsm.Current = newState
	fsm.Current.Enter()
}

// Update the current state.
func (fsm *FSM) Update() {
	if fsm.Current != nil {
		fsm.Current.Update()
	}
}
