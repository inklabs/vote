package main

import (
	"fmt"

	"github.com/inklabs/cqrs/api/grpcserver"
	"google.golang.org/grpc"

	"github.com/inklabs/vote"
	"github.com/inklabs/vote/grpc/grpcservergen"
)

func main() {
	app := vote.NewApp()

	fmt.Println("Hello World - gRPC API")
	grpcserver.Start(app, func(grpcServer *grpc.Server) {
		grpcservergen.RegisterServers(grpcServer, app)
	})
}
