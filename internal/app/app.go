package app

import (
	"fmt"

	"lem-in/internal/core"
	"lem-in/internal/dataops"
	"lem-in/pkg/fsops"
	"lem-in/pkg/serror"
)

func Run() error {

	// recieve os args and interpreter os args

	// print file content in stdout
	if err := fsops.ReadAndPrintFileContent(); err != nil {
		serror.ErrorHandler(err)
	}

	// read file
	data, err := fsops.ReadFile()
	if err != nil {
		serror.ErrorHandler(err)
	}

	// create graph
	graph, err := dataops.CreateGraph(data)
	if err != nil {
		serror.ErrorHandler(err)
	}

	// find all path
	if err := core.FindAllPaths(graph); err != nil {
		serror.ErrorHandler(err)
	}

	// move the ants to the end room
	paths := dataops.MoveAntsToEnd(graph)

	// print result output
	fmt.Printf("\n%v\n", core.ResultOutput(paths, graph))

	return nil
}
