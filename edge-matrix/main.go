package main

import (
	_ "embed"
	"github.com/emc-protocol/edge-matrix/command/root"
)

func main() {
	root.NewRootCommand().Execute()
}
