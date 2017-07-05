package command

type Command interface {
	Name() string
	Spec(prefix string) string
	Summary() string
	Handle(ctx *Context) error
}

type commmand struct {
	name       string
	argSpec    string
	summary    string
	handleFunc func(ctx *Context) error
}

var _ Command = new(commmand)

func NewCommand(name, argSpec, summary string, handleFunc func(ctx *Context) error) Command {
	return &commmand{name: name, argSpec: argSpec, summary: summary, handleFunc: handleFunc}
}

func (c *commmand) Name() string {
	return c.name
}

func (c *commmand) Spec(prefix string) string {
	return prefix + c.name + " " + c.argSpec + "\n"
}

func (c *commmand) Summary() string {
	return c.summary
}

func (c commmand) Handle(ctx *Context) error {
	return c.handleFunc(ctx)
}
