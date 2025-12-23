package flinch

import (
	"context"
	"io"
)

type Context struct {
	context.Context

	logger  Logger
	screen  Screen
	scripts Scripts
	time    Time
}

func NewContext(ctx context.Context, writer io.Writer) *Context {
	return &Context{
		Context: ctx,
		logger:  NewLogger(writer),
		screen:  NewScreen(),
		scripts: NewScripts(),
		time:    NewTime(),
	}
}

func (ctx *Context) Update() {
	ctx.time.Tick()
	ctx.scripts.Update(ctx.time.Delta())
}

func (ctx *Context) Logger() Logger {
	return ctx.logger
}

func (ctx *Context) Screen() Screen {
	return ctx.screen
}

func (ctx *Context) Scripts() Scripts {
	return ctx.scripts
}

func (ctx *Context) Time() Time {
	return ctx.time
}
