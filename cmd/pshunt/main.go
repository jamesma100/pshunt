package main

import (
	"github.com/jamesma100/pshunt/pkg/parser"
	"github.com/jamesma100/pshunt/pkg/runner"
	"github.com/jamesma100/pshunt/pkg/ui"
)

func main() {
	contents := runner.GetPsList()
	psList := parser.ParseList(contents)

	ui.StartUI(psList)

}
