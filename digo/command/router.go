package command

import (
	"errors"
)

type Router struct {
	name     string
	summary  string
	children []Command
}

var _ Command = &Router{}

func NewRouter(name, summary string) *Router {
	return &Router{name: name, summary: summary}
}

func (r *Router) Append(cmd Command) {
	r.children = append(r.children, cmd)
}

func (r *Router) Name() string {
	return r.name
}

func (r *Router) Spec(prefix string) (result string) {
	if r.name != "" {
		prefix = prefix + r.name + " "
	}
	for i, cmd := range r.children {
		_, isRouter := cmd.(*Router)
		if isRouter && i != 0 {
			result += "\n"
		}
		result += cmd.Spec(prefix)
		if isRouter && i != len(r.children)-1 {
			result += "\n"
		}
	}
	return result
}

func (r *Router) Summary() string {
	return r.summary
}

func (r *Router) Handle(ctx *Context) error {
	err := ctx.S.Scan()
	if err != nil {
		return err
	}
	token := ctx.S.Token()
	for _, cmd := range r.children {
		if cmd.Name() == token {
			return cmd.Handle(ctx)
		}
	}
	// TODO show docs or something?
	return errors.New("Unknown command.")
}
