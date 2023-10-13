package main

import (
	"fmt"

	"github.com/inklabs/cqrs/api/httpcmd"

	"github.com/inklabs/vote"
)

func main() {
	app := vote.NewApp()

	httpActionDecoder := vote.NewHTTPActionDecoder()

	fmt.Println("Hello World - HTTP API")
	httpcmd.Start(app, httpActionDecoder, vote.DomainBytes)
}
