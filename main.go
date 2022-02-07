package main

import (
	"github.com/mitchellh/cli"
	"log"
	"os"
	"todo/server"
)

func main() {
	c := cli.NewCLI("todo", "1.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory {
		"front": server.CommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
