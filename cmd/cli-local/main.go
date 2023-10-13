package main

import (
	"fmt"
	"log"
	"os"

	"github.com/inklabs/vote"
)

func main() {
	fmt.Println("Vote - Local CLI")

	app := vote.NewApp()
	defer app.Stop()

	command := vote.GetCobraRootCommand(app)
	command.SetOut(os.Stdout)
	err := command.Execute()
	if err != nil {
		log.Fatalf("unable to execute application: %v", err)
	}
}
