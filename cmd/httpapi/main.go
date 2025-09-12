package main

import (
	"fmt"
	"net/url"

	"github.com/inklabs/cqrs/api/httpserver"

	"github.com/inklabs/vote"
)

func main() {
	fmt.Println("Vote - HTTP API")

	app := vote.NewProdApp()
	httpActionDecoder := vote.NewHTTPActionDecoder()

	baseURI := url.URL{
		Scheme: "https",
		Host:   "api.vote.inklabs.dev",
	}

	httpserver.Start(
		app,
		httpActionDecoder,
		vote.ValidationRules,
		vote.DomainBytes,
		baseURI.String(),
		vote.Version,
	)
}
