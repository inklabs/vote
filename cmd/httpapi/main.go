package main

import (
	"fmt"

	"github.com/inklabs/cqrs/api/httpserver"

	"github.com/inklabs/vote"
)

func main() {
	app := vote.NewApp()

	httpActionDecoder := vote.NewHTTPActionDecoder()

	fmt.Println("Hello World - HTTP API")
	httpserver.Start(app, httpActionDecoder, vote.DomainBytes)
}
