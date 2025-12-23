package flinch

// ========================== Sequence ==========================

// Action defines a function type for sequence actions.
type Action func(dt float64) bool

// Sequence manages a sequence of timed actions.
type Sequence struct {
	actions  []Action
	started  bool
	complete bool
	current  int
}

func NewSequence(actions ...Action) *Sequence {
	return &Sequence{actions: actions}
}

func (tl *Sequence) Update(dt float64) {
	if !tl.started || tl.complete {
		return
	}

	if len(tl.actions) == 0 || tl.current >= len(tl.actions) {
		tl.complete = true
		return
	}

	if tl.actions[tl.current](dt) {
		tl.current++
	}
}

func (tl *Sequence) Start() *Sequence {
	tl.started = true
	return tl
}

func (tl *Sequence) Started() bool {
	return tl.started
}

func (tl *Sequence) Pause() {
	tl.started = false
}

func (tl *Sequence) Stop() {
	tl.started = false
	tl.complete = true
}

func (tl *Sequence) Complete() bool {
	return tl.complete
}

func (tl *Sequence) Reset() {
	tl.current = 0
	tl.started = false
	tl.complete = false
}

// ========================== Script ==========================

type Scripts interface {
	Update(dt float64)
	AddSequence(sequence *Sequence)
	RemoveSequence(sequence *Sequence)
}

type scripts struct {
	sequences []*Sequence
}

func NewScripts() Scripts {
	return &scripts{}
}

func (s *scripts) Update(dt float64) {
	if len(s.sequences) == 0 {
		return
	}

	for _, seq := range s.sequences {
		seq.Update(dt)
	}

	write := 0
	for _, seq := range s.sequences {
		if !seq.Complete() {
			s.sequences[write] = seq
			write++
		}
	}

	s.sequences = s.sequences[:write]
}

func (s *scripts) AddSequence(sequence *Sequence) {
	s.sequences = append(s.sequences, sequence)
}

func (s *scripts) RemoveSequence(sequence *Sequence) {
	for i, seq := range s.sequences {
		if seq == sequence {
			s.sequences = append(s.sequences[:i], s.sequences[i+1:]...)
			return
		}
	}
}
