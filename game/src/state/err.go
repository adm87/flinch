package state

import "fmt"

type Error struct {
	Current StateID
	Next    StateID
	Exit    StateExitCondition
	Msg     string
}

func NewError(current StateID, exit StateExitCondition, next StateID, msg string) *Error {
	return &Error{
		Current: current,
		Next:    next,
		Exit:    exit,
		Msg:     msg,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s (current: %d, exit: %d, next: %d)", e.Msg, e.Current, e.Exit, e.Next)
}
