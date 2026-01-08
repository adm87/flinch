package state

import (
	"github.com/adm87/flinch/engine/encoding"
)

type StateID uint64

func (id StateID) IsNil() bool {
	return id == NilStateID
}

const NilStateID StateID = 0

func NewStateID[T any]() StateID {
	return StateID(encoding.HashType[T]())
}

type StateExitCondition uint8

func (sec StateExitCondition) IsNil() bool {
	return sec == NilExitCondition
}

const NilExitCondition StateExitCondition = 0

type State[Context any] interface {
	Enter(ctx *Context) error
	Exit(ctx *Context) error
	Process(ctx *Context) (StateExitCondition, error)
}

type StateFactory[Context any] func() State[Context]

type States[Context any] map[StateID]StateFactory[Context]

type Transitions map[StateID]map[StateExitCondition]StateID

type FSM[Context any] struct {
	current     StateID
	next        StateID
	transitions Transitions

	states States[Context]
	state  State[Context]

	isTransitioning bool
}

func (fsm *FSM[Context]) IsTransitioning() bool {
	return fsm.isTransitioning
}

func (fsm *FSM[Context]) SetNext(stateID StateID) {
	fsm.next = stateID
}

func (fsm *FSM[Context]) SetTransitions(transitions Transitions) {
	fsm.transitions = transitions
}

func (fsm *FSM[Context]) CurrentID() StateID {
	return fsm.current
}

func (fsm *FSM[Context]) State() State[Context] {
	return fsm.state
}

func (fsm *FSM[Context]) Process(ctx *Context) error {
	if !fsm.next.IsNil() && fsm.next != fsm.current {
		return changeState(fsm, ctx)
	}

	fsm.isTransitioning = false

	exitCondition, err := runState(fsm, ctx)
	if err != nil {
		return err
	}

	if exitCondition.IsNil() {
		return nil
	}

	nextStateID, err := getNextStateID(fsm, exitCondition)
	if err != nil {
		return err
	}
	fsm.next = nextStateID

	return nil
}

func (fsm *FSM[Context]) AddState(stateID StateID, factory StateFactory[Context]) *FSM[Context] {
	fsm.states[stateID] = factory
	return fsm
}

func NewFSM[Context any]() *FSM[Context] {
	return &FSM[Context]{
		states: make(States[Context]),
	}
}

func Register[T, Context any](fsm *FSM[Context], factory StateFactory[Context]) StateID {
	id := NewStateID[T]()
	fsm.AddState(id, factory)
	return id
}

func changeState[Context any](fsm *FSM[Context], ctx *Context) error {
	if fsm.state != nil {
		if err := fsm.state.Exit(ctx); err != nil {
			return NewError(fsm.current, NilExitCondition, fsm.next, err.Error())
		}
	}

	factory, exists := fsm.states[fsm.next]
	if !exists {
		return NewError(fsm.current, NilExitCondition, fsm.next, "next state ID does not exist in states map")
	}

	fsm.state = factory()
	if err := fsm.state.Enter(ctx); err != nil {
		return NewError(fsm.current, NilExitCondition, fsm.next, err.Error())
	}

	fsm.current = fsm.next
	fsm.next = NilStateID
	fsm.isTransitioning = true

	return nil
}

func runState[Context any](fsm *FSM[Context], ctx *Context) (StateExitCondition, error) {
	if fsm.state == nil {
		return NilExitCondition, NewError(fsm.current, NilExitCondition, NilStateID, "current state is nil")
	}

	exitCondition, err := fsm.state.Process(ctx)
	if err != nil {
		return NilExitCondition, NewError(fsm.current, exitCondition, NilStateID, err.Error())
	}

	return exitCondition, nil
}

func getNextStateID[Context any](fsm *FSM[Context], exitCondition StateExitCondition) (StateID, error) {
	if fsm.transitions == nil {
		return NilStateID, nil
	}

	transitions, exists := fsm.transitions[fsm.current]
	if !exists {
		return NilStateID, NewError(fsm.current, NilExitCondition, NilStateID, "no transitions defined for current state")
	}

	nextStateID, exists := transitions[exitCondition]
	if !exists {
		return NilStateID, NewError(fsm.current, exitCondition, NilStateID, "no transition defined for given exit condition")
	}

	return nextStateID, nil
}
