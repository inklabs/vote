package event

type ElectionHasCommenced struct {
	ElectionID      string
	OrganizerUserID string
	Name            string
	Description     string
	OccurredAt      int
}

type ProposalWasMade struct {
	ElectionID  string
	ProposalID  string
	OwnerUserID string
	Name        string
	Description string
	OccurredAt  int
}

type VoteWasCast struct {
	ElectionID        string
	ProposalID        string
	UserID            string
	RankedProposalIDs []string
	OccurredAt        int
}

type ElectionWasClosedByOwner struct {
	ElectionID string
	OccurredAt int
}

type ElectionWinnerWasSelected struct {
	ElectionID        string
	WinningProposalID string
	OccurredAt        int
}
