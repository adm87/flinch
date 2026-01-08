package flinch

type Input interface {
	Update(ctx *Context) error
}

type input struct {
}

func NewInput() Input {
	return &input{}
}

func (i *input) Update(ctx *Context) error {
	// Placeholder for input update logic
	return nil
}
