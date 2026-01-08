package flinch

// ScriptedAction defines a function type for scriptable actions.
type ScriptedAction func(ctx *Context) error

// ScriptedSequence represents a sequence of scriptable actions.
type ScriptedSequence struct {
	current   int
	started   bool
	completed bool
	actions   []ScriptedAction
}

// NewScriptSequence creates a new ScriptSequence with the given actions.
func NewScriptSequence(actions ...ScriptedAction) *ScriptedSequence {
	return &ScriptedSequence{
		actions: actions,
	}
}

func (ss *ScriptedSequence) Start() *ScriptedSequence {
	ss.started = true
	ss.current = 0
	ss.completed = false
	return ss
}

func (ss *ScriptedSequence) IsCompleted() bool {
	return ss.completed
}

func (ss *ScriptedSequence) Update(ctx *Context) error {
	if !ss.started {
		return nil
	}

	if ss.current >= len(ss.actions) {
		ss.completed = true
		return nil
	}

	defer func() { ss.current++ }()

	if err := ss.actions[ss.current](ctx); err != nil {
		return err
	}

	return nil
}

// Script interface defines methods for script management.
//
// This is a simple system that allows for defining and executing sequences of actions.
type Script interface {
	Update(ctx *Context) error
}

type script struct {
	sequences []*ScriptedSequence
}

func NewScript() Script {
	return &script{}
}

func (s *script) Update(ctx *Context) error {
	for _, seq := range s.sequences {
		if err := seq.Update(ctx); err != nil {
			return err
		}
	}

	// Remove completed sequences
	active := s.sequences[:0]
	for _, seq := range s.sequences {
		if !seq.IsCompleted() {
			active = append(active, seq)
		}
	}
	s.sequences = active

	return nil
}

func (s *script) AddSequence(seq *ScriptedSequence) {
	s.sequences = append(s.sequences, seq)
}

func (s *script) RemoveSequence(seq *ScriptedSequence) {
	for i, existingSeq := range s.sequences {
		if existingSeq == seq {
			s.sequences = append(s.sequences[:i], s.sequences[i+1:]...)
			return
		}
	}
}

func (s *script) ClearSequences() {
	s.sequences = nil
}
