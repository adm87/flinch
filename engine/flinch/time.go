package flinch

import (
	stdTime "time"
)

const (
	MaxAccumulatedTime = 0.1
	FixedDelta         = 1.0 / 60.0
)

type Time interface {
	Tick()

	Delta() float64      // Delta time since last frame
	FixedDelta() float64 // Fixed delta time for fixed updates
	FixedSteps() int     // Number of fixed updates to run this frame

	FPS() int      // Frames per second
	FixedFPS() int // Fixed updates per second
}

type time struct {
	lastTime stdTime.Time
	delta    float64

	accumulator float64
	fixedSteps  int

	secondTimer float64
	frameCount  int
	fixedCount  int
	fps         int
	fixedFps    int
}

func NewTime() Time {
	return &time{}
}

func (t *time) Tick() {
	now := stdTime.Now()

	// Initialize on first tick
	if t.lastTime.IsZero() {
		t.lastTime = now
		return
	}

	// Calculate delta time
	t.delta = now.Sub(t.lastTime).Seconds()
	t.lastTime = now

	// Skip if no time has passed
	if t.delta <= 0 {
		return
	}

	// Clamp delta to prevent spiral of death
	if t.delta > MaxAccumulatedTime {
		t.delta = MaxAccumulatedTime
	}

	// Update fixed timestep accumulator
	t.accumulator = min(t.accumulator+t.delta, MaxAccumulatedTime)
	t.fixedSteps = 0
	for t.accumulator >= FixedDelta {
		t.fixedSteps++
		t.accumulator -= FixedDelta
		t.fixedCount++
	}

	// Update FPS counters
	t.frameCount++
	t.secondTimer += t.delta
	if t.secondTimer >= 1.0 {
		t.fps = t.frameCount
		t.fixedFps = t.fixedCount
		t.frameCount = 0
		t.fixedCount = 0
		t.secondTimer -= 1.0
	}
}

func (t *time) Delta() float64 {
	return t.delta
}

func (t *time) FixedDelta() float64 {
	return FixedDelta
}

func (t *time) FixedSteps() int {
	return t.fixedSteps
}

func (t *time) FPS() int {
	return t.fps
}

func (t *time) FixedFPS() int {
	return t.fixedFps
}

// ========================== Timer ==========================

type Timer struct {
	duration  float32
	elapsed   float32
	repeat    bool
	completed bool
}

func NewTimer(duration float32, repeat bool) *Timer {
	return &Timer{
		duration: duration,
		repeat:   repeat,
	}
}

func (tm *Timer) Update(dt float32) {
	if tm.completed {
		return
	}

	tm.elapsed += dt
	if tm.elapsed >= tm.duration {
		if tm.repeat {
			tm.elapsed -= tm.duration
		} else {
			tm.completed = true
		}
	}
}

func (tm *Timer) Completed() bool {
	return tm.completed
}

func (tm *Timer) Reset() {
	tm.elapsed = 0
	tm.completed = false
}
