# Ranked Choice Voting - Go CQRS

This is a demo application that implements
[Ranked Choice Voting](https://fairvote.org/our-reforms/ranked-choice-voting/),
also known as Instant Runoff Voting, where voters rank candidates in order of preference, and
a single winner is selected.

The sole purpose of this project is to demonstrate the capabilities of the
[Go CQRS](https://github.com/inklabs/cqrs) application framework.

## Design

The architecture leverages Command Query Responsibility Segregation
([CQRS](https://learn.microsoft.com/en-us/azure/architecture/patterns/cqrs)).
Business logic can be written in Commands, Queries, or Listeners.

### Actions

Actions are Commands that modify state, or Queries that retrieve data. Commands can
be asynchronously processed, with the status and log messages available for later retrieval.

- Commands
    - [CommenceElection](action/election/commence_election.go)
    - [MakeProposal](action/election/make_proposal.go)
    - [CastVote](action/election/cast_vote.go)
- AsyncCommands
    - [CloseElectionByOwner](action/election/close_election_by_owner.go)
- Queries
    - [ListOpenElections](action/election/list_open_elections.go)
    - [ListProposals](action/election/list_proposals.go)
    - [GetProposalDetails](action/election/get_proposal_details.go)
    - [GetElectionResults](action/election/get_election_results.go)

### Events

Events are used to remove temporal coupling between unrelated modules. Commands raise events,
and Listeners subscribe to them.

- [Events](event/election_events.go)
  - ElectionHasCommenced
  - ProposalWasMade
  - VoteWasCast
  - ElectionWasClosedByOwner
  - ElectionWinnerWasSelected

### Listeners

Listeners subscribe to Events and execute code asynchronously.

  - [ElectionWinnerVoterNotification](listener/election_winner_voter_notification.go)
    - TODO: notify voters via SMS, Slack, or email
  - [ElectionWinnerMediaNotification](listener/election_winner_media_notification.go)
    - TODO: send press release email

## APIs

APIs are exposed through HTTP, gRPC, and CLI tooling.

### Examples:

- [CLI Examples](cli_test.go)
- [HTTP Schema Examples](http_schema_test.go)
- [HTTP API Examples](http_test.go)

## SDK

SDKs are auto-generated leveraging the gRPC API. Currently supported:

- Go
- Python

## Test

```
go generate .
go test ./...
```

## Run

```
go run cmd/httpapi/main.go
go run cmd/grpcapi/main.go
go run cmd/cli-local/main.go --help
```

## Test Python

```
from __future__ import print_function
from google.protobuf.json_format import MessageToJson
from electionpb.election_pb2 import ListOpenElectionsRequest
from electionpb import election_pb2_grpc

import logging
import grpc


def run():
    print("Will try to greet world ...")
    with grpc.insecure_channel("localhost:8081") as channel:
        stub = election_pb2_grpc.ElectionStub(channel)
        response = stub.ListOpenElections(ListOpenElectionsRequest())
    print("client received: " + MessageToJson(response))


if __name__ == "__main__":
    logging.basicConfig()
    run()
```

## Links

- https://github.com/inklabs/cqrs
