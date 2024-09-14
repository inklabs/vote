package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-faker/faker/v4"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/inklabs/vote/grpc/go/electionpb"
	"github.com/inklabs/vote/sdk/go/goclient"
)

func main() {
	fmt.Println("Election Simulator")

	totalElections := flag.Int("totalElections", 100, "Total # of elections to simulate")
	maxProposals := flag.Int("maxProposals", 7, "Max # of proposals per election")
	totalVoters := flag.Int("totalVoters", 10, "Total # of voters")
	delay := flag.Int("delay", 10, "Delay between calls in milliseconds")
	host := flag.String("host", "127.0.0.1:8081", "Vote gRPC host address")
	flag.Parse()

	fmt.Printf("Simulating election until stopped\n")
	fmt.Printf("totalElections: %d\n", *totalElections)
	fmt.Printf("maxProposals: %d\n", *maxProposals)
	fmt.Printf("totalVoters: %d\n", *totalVoters)
	fmt.Printf("delay: %d\n", *delay)

	dialCtx, connectDone := context.WithTimeout(context.Background(), time.Second*2)
	conn, err := grpc.DialContext(
		dialCtx,
		*host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("unable to dial (%s): %v", *host, err)
	}
	defer func() {
		_ = conn.Close()
		connectDone()
	}()

	client := goclient.NewClient(conn)
	s := NewSimulator(client, *totalElections, *maxProposals, *totalVoters, *delay)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	simulateUntilStopped(ctx, s)

	fmt.Printf("\nDone\n")
}

func simulateUntilStopped(ctx context.Context, simulator *simulator) {
	for {
		fmt.Printf("Starting new Simulation\n")
		simulator.Start(ctx)
		simulator.Errors(ctx)

		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}
	}
}

func NewSimulator(client *goclient.GoClient, totalElections, maxProposals, totalVoters, delay int) *simulator {
	return &simulator{
		client:          client,
		totalElections:  totalElections,
		maxProposals:    maxProposals,
		totalVoters:     totalVoters,
		delay:           time.Duration(delay) * time.Millisecond,
		electionAdminID: uuid.NewString(),
		elections:       make(map[string][]string),
	}
}

type simulator struct {
	client         *goclient.GoClient
	totalElections int
	maxProposals   int
	totalVoters    int
	delay          time.Duration

	electionAdminID string

	// elections electionID => proposalIDs
	elections map[string][]string
}

func (s *simulator) Errors(ctx context.Context) {
	const (
		unknownElectionID = "unknown-election-id"
		unknownUserID     = "unknown-user-id"
	)

	_, _ = s.client.Election.CastVote(ctx, &electionpb.CastVoteRequest{
		ElectionId:        unknownElectionID,
		UserId:            unknownUserID,
		RankedProposalIDs: []string{"unknown-proposal-id"},
	})

	_, _ = s.client.Election.CloseElectionByOwner(ctx, &electionpb.CloseElectionByOwnerRequest{
		Id:         uuid.NewString(),
		ElectionId: unknownElectionID,
	})
}

func (s *simulator) Start(ctx context.Context) {
	err := s.setupElections(ctx)
	if err != nil {
		return
	}

	err = s.castVotes(ctx)
	if err != nil {
		return
	}

	err = s.closeElections(ctx)
	if err != nil {
		return
	}
}

func (s *simulator) setupElections(ctx context.Context) error {
	s.elections = make(map[string][]string)

	for i := 0; i < s.totalElections; i++ {
		time.Sleep(s.delay)
		electionID := uuid.NewString()

		_, err := s.client.Election.CommenceElection(ctx, &electionpb.CommenceElectionRequest{
			ElectionId:      electionID,
			OrganizerUserId: s.electionAdminID,
			Name:            faker.Name(),
			Description:     faker.Sentence(),
		})
		if err != nil {
			return fmt.Errorf("unable to create election: %w", err)
		}

		for j := 0; j < random(1, s.maxProposals); j++ {
			proposalID := uuid.NewString()
			s.elections[electionID] = append(s.elections[electionID], proposalID)

			_, err := s.client.Election.MakeProposal(ctx, &electionpb.MakeProposalRequest{
				ElectionId:  electionID,
				ProposalId:  proposalID,
				OwnerUserId: s.electionAdminID,
				Name:        faker.Name(),
				Description: faker.Sentence(),
			})
			if err != nil {
				return fmt.Errorf("unable to create proposal: %w", err)
			}
		}
	}

	return nil
}

func (s *simulator) castVotes(ctx context.Context) error {
	for i := 0; i < s.totalVoters; i++ {
		time.Sleep(s.delay)
		userID := uuid.NewString()
		electionsToVoteIn := random(1, s.totalElections)

		for electionID, proposals := range s.elections {
			if electionsToVoteIn <= 0 {
				break
			}

			shuffle(proposals)
			_, err := s.client.Election.CastVote(ctx, &electionpb.CastVoteRequest{
				ElectionId:        electionID,
				UserId:            userID,
				RankedProposalIDs: proposals,
			})
			if err != nil {
				return fmt.Errorf("unable to cast vote: %w", err)
			}

			electionsToVoteIn--
		}
	}

	return nil
}

func (s *simulator) closeElections(ctx context.Context) error {
	for electionID := range s.elections {
		time.Sleep(s.delay)
		response, err := s.client.Election.CloseElectionByOwner(ctx, &electionpb.CloseElectionByOwnerRequest{
			Id:         uuid.NewString(),
			ElectionId: electionID,
		})
		if err != nil {
			return fmt.Errorf("unable to close election: %w", err)
		}

		if !response.HasBeenQueued {
			return fmt.Errorf("election %v has not been_queued", electionID)
		}
	}

	return nil
}

func random(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func shuffle(slice []string) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
