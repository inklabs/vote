package main

import (
	"fmt"

	"github.com/inklabs/cqrs/api/httpserver"

	"github.com/inklabs/vote"
)

func main() {
	fmt.Println("Vote - HTTP API")

	app := vote.NewApp()
	httpActionDecoder := vote.NewHTTPActionDecoder()

	httpserver.Start(app, httpActionDecoder, vote.DomainBytes)
}
