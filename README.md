# Example Voting CQRS App

Demo application using the [Go CQRS](https://github.com/inklabs/cqrs) application framework

## Design

- Events
    - ElectionHasCommenced: ElectionID, OrganizerUserID, Name, Description, ts
    - ProposalWasMade: ElectionID, ProposalID, OwnerUserID, Name, Description, ts
    - VoteWasCast: ElectionID, ProposalID, UserID, []RankedProposalIDs{1, 2}, ts
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

## Links

- https://github.com/inklabs/cqrs
