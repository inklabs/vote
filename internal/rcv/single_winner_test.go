package rcv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inklabs/vote/internal/rcv"
)

const (
	A = "A"
	B = "B"
	C = "C"
	D = "D"
)

func TestSingleWinner(t *testing.T) {
	tests := []struct {
		name    string
		ballots rcv.Ballots
		winner  string
	}{
		{
			name: "round 1: 1 ballot, 1 proposal",
			ballots: rcv.Ballots{
				{A},
			},
			winner: A,
		},
		{
			name: "round 1: 1 ballot, 3 proposals",
			ballots: rcv.Ballots{
				{A, B, C},
			},
			winner: A,
		},
		{
			name: "round 1: 3 ballots",
			ballots: rcv.Ballots{
				{A, B, C},
				{A, B, C},
				{A, B, C},
			},
			winner: A,
		},
		{
			name: "round 2: 5 ballots",
			ballots: rcv.Ballots{
				{A, B, C},
				{B, A, C},
				{C, B, A},
				{A, B, C},
				{B, A, C},
			},
			winner: B,
		},
		{
			//https://github.com/BrightSpots/rcv/blob/develop/src/test/resources/network/brightspots/rcv/test_data/minimum_threshold_test/minimum_threshold_test_expected_summary.csv
			// TODO: Fix this. Current implementation is non-deterministic.
			name: "round 3: minimum threshold",
			ballots: rcv.Ballots{
				{A, B, C},
				{A, B, C},
				{A, B, C},
				{A, B, C},
				{B, C},
				{B, D, A},
				{D, B, A},
				{C, B, A},
				{C, B, A},
				{C, B, A},
			},
			winner: A,
		},
		{
			name: "3 rounds",
			ballots: rcv.Ballots{
				{A},
				{A},
				{A},
				{A},
				{B},
				{B},
				{B},
				{C, A},
				{C, A},
				{D, B},
			},
			winner: A,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			tabulator := rcv.NewSingleWinner(tc.ballots)

			// When
			winningProposalID, err := tabulator.GetWinningProposal()

			// Then
			require.NoError(t, err)
			assert.Equal(t, tc.winner, winningProposalID)
		})
	}
}

func TestSingleWinner_WinnerNotFound(t *testing.T) {
	tests := []struct {
		name    string
		ballots rcv.Ballots
	}{
		{
			name: "equal votes for all proposals with no majority",
			ballots: rcv.Ballots{
				{A},
				{B},
				{C},
			},
		},
		{
			name: "2nd round tie",
			ballots: rcv.Ballots{
				{A},
				{A},
				{A},
				{A},
				{A},
				{B},
				{B},
				{B, A},
				{B, A},
				{B, A},
				{B, A},
				{C, A},
				{C, A},
				{C, A},
				{C, A},
				{C, A},
				{C, A},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			tabulator := rcv.NewSingleWinner(tc.ballots)

			// When
			winningProposalID, err := tabulator.GetWinningProposal()

			// Then
			assert.Equal(t, rcv.ErrWinnerNotFound, err)
			assert.Equal(t, "", winningProposalID)
		})
	}
}
