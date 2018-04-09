package main

import (
	_ "net/http/pprof"
)

// Command is a runnable command
type Command interface {
	Name() string
	ShortDesc() string
	LongDesc() string
	Execute(args []string) error
}

type commandDesc struct {
	name      string
	shortDesc string
	longDesc  string
}

func (c *commandDesc) Name() string {
	return c.name
}

func (c *commandDesc) ShortDesc() string {
	return c.shortDesc
}

func (c *commandDesc) LongDesc() string {
	return c.longDesc
}
