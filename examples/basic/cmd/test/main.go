package main

import (
	"github.com/src-d/code-annotation/examples/basic"
)

var (
	version string
	build   string
	commit  string
)

func main() {
	basic.HelloWorld(version, build, commit)
}
