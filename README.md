# Example Voting CQRS App

Demo application using the [Go CQRS](https://github.com/inklabs/cqrs) application framework

Ranked Choice Voting - https://fairvote.org/our-reforms/ranked-choice-voting/

## Design

- Events
    - ElectionHasCommenced: ElectionID, OrganizerUserID, Name, Description, ts
    - ProposalWasMade: ElectionID, ProposalID, OwnerUserID, Name, Description, ts
    - VoteWasCast: ElectionID, UserID, []RankedProposalIDs{1, 2}, ts
    - ElectionWasClosedByOwner: ElectionID, OwnerUserID, ts
    - ElectionWinnerWasSelected: ElectionID, WinningProposalID, ts
- Commands
    - CommenceElection -> ElectionHasCommenced
    - MakeProposal -> ProposalWasMade
    - CastVote -> VoteWasCast
- AsyncCommands
    - CloseElectionByOwner -> ElectionWasClosedByOwner, tabulate results -> ElectionWinnerWasSelected
- Queries
    - ListOpenElections:
    - ListProposals: ElectionID
    - GetProposalDetails: ProposalID
    - GetElectionResults: ElectionID
- Listeners
    - ElectionWinnerVoterNotification: ElectionWinnerWasSelected -> notify voters via SMS, Slack, or email
    - ElectionWinnerMediaNotification: ElectionWinnerWasSelected -> send press release email

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
