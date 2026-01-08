package flinch

import (
	"context"
	"io"
)

type Context struct {
	context.Context

	input  Input
	logger Logger
	screen Screen
	script Script
	time   Time
}

func NewContext(ctx context.Context, writer io.Writer) *Context {
	return &Context{
		Context: ctx,
		input:   NewInput(),
		logger:  NewLogger(writer),
		screen:  NewScreen(),
		script:  NewScript(),
		time:    NewTime(),
	}
}

func (ctx *Context) Update() error {
	// Update time before any other systems
	ctx.time.Tick()

	ctx.input.Update(ctx)
	ctx.script.Update(ctx)

	return nil
}

func (ctx *Context) Logger() Logger {
	return ctx.logger
}

func (ctx *Context) Screen() Screen {
	return ctx.screen
}

func (ctx *Context) Script() Script {
	return ctx.script
}

func (ctx *Context) Time() Time {
	return ctx.time
}
