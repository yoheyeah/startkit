package startkit

import (
	"startkit/starter"

	"golang.org/x/sync/errgroup"
)

type Context struct {
	*starter.Content
}

type StarterFunc func(c *Context) error

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

func (c *Context) Run(funcs ...StarterFunc) (err error) {
	var g errgroup.Group
	for i := range funcs {
		g.Go(func() error {
			return funcs[i](c)
		})
	}
	return g.Wait()
}
