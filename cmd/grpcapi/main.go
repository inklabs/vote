package main

import (
	"fmt"

	"github.com/inklabs/cqrs/api/grpcserver"
	"google.golang.org/grpc"

	"github.com/inklabs/vote"
	voteserver "github.com/inklabs/vote/grpc/grpcserver"
)

func main() {
	fmt.Println("Vote - gRPC API")

	app := vote.NewProdApp()

	grpcserver.Start(app, func(grpcServer *grpc.Server) {
		voteserver.RegisterServers(grpcServer, app)
	})
}
