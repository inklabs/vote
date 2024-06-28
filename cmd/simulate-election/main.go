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

	totalElections := flag.Int("totalElections", 10, "Total # of elections to simulate")
	maxProposals := flag.Int("maxProposals", 7, "Max # of proposals per election")
	totalVoters := flag.Int("totalVoters", 10000, "Total # of voters")
	host := flag.String("host", "127.0.0.1:8081", "Vote gRPC host address")
	flag.Parse()

	fmt.Printf("Simulating election until stopped\n")
	fmt.Printf("totalElections: %d\n", *totalElections)
	fmt.Printf("maxProposals: %d\n", *maxProposals)
	fmt.Printf("totalVoters: %d\n", *totalVoters)

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client := goclient.NewClient(conn)
	simulator := NewSimulator(client, *totalElections, *maxProposals, *totalVoters)
	simulator.Start(ctx)

	fmt.Printf("\nDone\n")
}

func NewSimulator(client *goclient.GoClient, totalElections, maxProposals, totalVoters int) *simulator {
	return &simulator{
		client:          client,
		totalElections:  totalElections,
		maxProposals:    maxProposals,
		totalVoters:     totalVoters,
		electionAdminID: uuid.NewString(),
		elections:       make(map[string][]string),
	}
}

type simulator struct {
	client          *goclient.GoClient
	totalElections  int
	maxProposals    int
	totalVoters     int
	electionAdminID string

	// elections electionID => proposalIDs
	elections map[string][]string
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
}

func (s *simulator) setupElections(ctx context.Context) error {
	for i := 0; i < s.totalElections; i++ {
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

func random(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func shuffle(slice []string) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
