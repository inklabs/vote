package vote_test

import (
	"context"
	"fmt"
	"time"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/cqrstest"
	"github.com/protocolbuffers/txtpbfmt/parser"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	"github.com/inklabs/vote"
	"github.com/inklabs/vote/grpc/go/asynccommandpb"
	"github.com/inklabs/vote/grpc/go/electionpb"
	voteserver "github.com/inklabs/vote/grpc/grpcserver"
	"github.com/inklabs/vote/sdk/go/goclient"
)

func ExampleApp_grpcListOpenElections() {
	app := newTestApp()
	defer app.Stop()
	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()
	voteserver.RegisterServers(grpcServer, app)
	conn := startBufferedGRPCServer(grpcServer)
	ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
	defer done()
	client := goclient.NewClient(conn)

	_, _ = client.Election.CommenceElection(ctx, &electionpb.CommenceElectionRequest{
		ElectionId:      "E1",
		OrganizerUserId: "U1",
		Name:            "Election Name 1",
		Description:     "Election Description 1",
	})

	_, _ = client.Election.CommenceElection(ctx, &electionpb.CommenceElectionRequest{
		ElectionId:      "E2",
		OrganizerUserId: "U1",
		Name:            "Election Name 2",
		Description:     "Election Description 2",
	})

	_, _ = client.Election.CommenceElection(ctx, &electionpb.CommenceElectionRequest{
		ElectionId:      "E3",
		OrganizerUserId: "U1",
		Name:            "Election Name 3",
		Description:     "Election Description 3",
	})

	response, _ := client.Election.ListOpenElections(ctx, &electionpb.ListOpenElectionsRequest{
		Page:          cqrs.Int64(1),
		ItemsPerPage:  cqrs.Int64(2),
		SortBy:        cqrs.String("Name"),
		SortDirection: cqrs.SortAscending,
	})

	fmt.Print(MarshalTextString(response))

	// Output:
	//  open_elections: {
	//   election_id: "E1"
	//   organizer_user_id: "U1"
	//   name: "Election Name 1"
	//   description: "Election Description 1"
	//   commenced_at: 1699900000
	// }
	// open_elections: {
	//   election_id: "E2"
	//   organizer_user_id: "U1"
	//   name: "Election Name 2"
	//   description: "Election Description 2"
	//   commenced_at: 1699900001
	// }
	// total_results: 3
}

func ExampleApp_grpcCloseElectionByOwner() {
	recordingEventDispatcher := cqrstest.NewRecordingEventDispatcher()
	app := newTestApp(
		vote.WithEventDispatcher(recordingEventDispatcher),
	)
	defer app.Stop()
	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()
	voteserver.RegisterServers(grpcServer, app)
	conn := startBufferedGRPCServer(grpcServer)
	ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
	defer done()
	client := goclient.NewClient(conn)

	recordingEventDispatcher.Add(4)

	_, _ = client.Election.CommenceElection(ctx, &electionpb.CommenceElectionRequest{
		ElectionId:      "E1",
		OrganizerUserId: "U1",
		Name:            "Election Name",
		Description:     "Election Description",
	})

	_, _ = client.Election.MakeProposal(ctx, &electionpb.MakeProposalRequest{
		ElectionId:  "E1",
		ProposalId:  "P1",
		OwnerUserId: "U2",
		Name:        "Proposal Name",
		Description: "Proposal Description",
	})

	_, _ = client.Election.CastVote(ctx, &electionpb.CastVoteRequest{
		ElectionId:        "E1",
		UserId:            "U3",
		RankedProposalIDs: []string{"P1"},
	})

	asyncCommandResponse, _ := client.Election.CloseElectionByOwner(ctx, &electionpb.CloseElectionByOwnerRequest{
		Id:         "AC1",
		ElectionId: "E1",
	})
	fmt.Printf("asyncCommandResponse:\n%s", MarshalTextString(asyncCommandResponse))

	recordingEventDispatcher.Wait(ctx)

	status, _ := client.AsyncCommand.Status(ctx, &asynccommandpb.StatusRequest{
		CommandId:   "AC1",
		IncludeLogs: true,
	})
	fmt.Printf("\n%s", MarshalTextString(status))

	// Output:
	// asyncCommandResponse:
	// id: "AC1"
	// status: "QUEUED"
	// has_been_queued: true
	//
	// async_command_status: {
	//   created_at: 1699900003
	//   modified_at: 1699900007
	//   started_at_micro: 1699900004000000
	//   finished_at_micro: 1699900007000000
	//   execution_duration: "3s"
	//   total_to_process: 1
	//   total_processed: 1
	//   percent_done: 100
	//   is_success: true
	//   is_finished: true
	//   close_election_by_owner: {
	//     id: "AC1"
	//     election_id: "E1"
	//   }
	// }
	// logs: {
	//   type: "INFO"
	//   created_at_micro: 1699900006000000
	//   message: "Closing election with winner: P1"
	// }
}

func MarshalTextString(response proto.Message) string {
	unstableOutput := prototext.Format(response)
	out, _ := parser.Format([]byte(unstableOutput))
	return string(out)
}
