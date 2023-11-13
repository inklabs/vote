package rcv

import (
	"fmt"
)

// Ballots A 2D slice representing the ranked choices of each voter.
// Each inner slice contains the ordered preferences of a voter,
// where the first element is the highest-ranked choice.
type Ballots [][]string

type singleWinner struct {
	totalVotes    int
	threshold     int
	proposalCount map[string]int
	ballots       Ballots
}

// NewSingleWinner is a ranked choice vote tabulator based on the provided
// Ballots. It uses an iterative process to determine the winning proposal,
// considering both majority support and eliminating proposals with the least votes.
// For more information check out [Wikipedia](https://en.wikipedia.org/wiki/Instant-runoff_voting)
// or [FairVote](https://fairvote.org/our-reforms/ranked-choice-voting).
func NewSingleWinner(ballots Ballots) *singleWinner {
	return &singleWinner{
		totalVotes:    len(ballots),
		threshold:     len(ballots) / 2,
		ballots:       ballots,
		proposalCount: make(map[string]int),
	}
}

// GetWinningProposal returns the winning proposal.
// ErrWinnerNotFound is returned if no winner is found.
func (t *singleWinner) GetWinningProposal() (string, error) {
	t.initProposals()
	t.tallyVotes()

	winningProposalID, isFound := t.getWinner()
	if isFound {
		return winningProposalID, nil
	}

	return t.getWinnerFromRemainingProposalIDs()
}

func (t *singleWinner) getWinner() (string, bool) {
	for proposalID, count := range t.proposalCount {
		if count > t.threshold {
			return proposalID, true
		}
	}

	return "", false
}

func (t *singleWinner) initProposals() {
	for _, proposalIDs := range t.ballots {
		for _, proposalID := range proposalIDs {
			if _, ok := t.proposalCount[proposalID]; !ok {
				t.proposalCount[proposalID] = 0
			}
		}
	}
}

// getWinnerFromRemainingProposalIDs eliminates the proposalID with the least votes
// and repeats until a majority winner is found. ErrWinnerNotFound is returned if
// no winner is found.
func (t *singleWinner) getWinnerFromRemainingProposalIDs() (string, error) {
	for len(t.proposalCount) > 1 {
		t.removeMinProposal()
		t.resetProposalCounts()
		t.tallyVotes()

		winningProposalID, isFound := t.getWinner()
		if isFound {
			return winningProposalID, nil
		}
	}

	return "", ErrWinnerNotFound
}

// removeMinProposal removes the lowest ranked proposal.
func (t *singleWinner) removeMinProposal() {
	var minProposalID string
	var minVotes = t.totalVotes

	for proposalID, count := range t.proposalCount {
		if count < minVotes {
			minVotes = count
			minProposalID = proposalID
		}
	}

	delete(t.proposalCount, minProposalID)
}

// tallyVotes increments the count for the next highest-ranked proposal
// still in the running
func (t *singleWinner) tallyVotes() {
	for _, rankedProposalIDS := range t.ballots {
		for _, proposalID := range rankedProposalIDS {
			if _, ok := t.proposalCount[proposalID]; ok {
				t.proposalCount[proposalID]++
				break
			}
		}
	}
}

// resetProposalCounts resets all vote counts to zero.
func (t *singleWinner) resetProposalCounts() {
	for proposalID := range t.proposalCount {
		t.proposalCount[proposalID] = 0
	}
}

var ErrWinnerNotFound = fmt.Errorf("winner not found")
