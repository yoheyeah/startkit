package startkit

import (
	"startkit/starter"
)

type Context struct {
	*starter.Content
}

type StarterFunc func(c *Context)

func Default() *Context {
	c := Context{}
	content := starter.DefaultBuilder()
	c.Content = content
	return &c
}

func New(file string) *Context {
	c := Context{}
	content := starter.CustomBuilder(file)
	c.Content = content
	return &c
}

func (c *Context) Run(funcs ...StarterFunc) {
	for i := range funcs {
		funcs[i](c)
	}
}
