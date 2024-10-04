package vote_test

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/inklabs/cqrs/cqrstest"

	"github.com/inklabs/vote"
)

func ExampleApp_cliRoot_help() {
	app := newTestApp()
	defer app.Stop()
	cmd := vote.GetCobraRootCommand(app)
	cmd.SetArgs([]string{"-h"})
	_ = cmd.Execute()

	// Output:
	// CLI application
	//
	// Usage:
	//   cli [flags]
	//   cli [command]
	//
	// Available Commands:
	//   async-command-status Async Command Status
	//   completion           Generate the autocompletion script for the specified shell
	//   election             9 actions: [CastVote, CloseElectionByOwner, CommenceElection, GetElection, GetElectionResults, GetProposalDetails, ListOpenElections, ListProposals, MakeProposal]
	//   help                 Help about any command
	//
	// Flags:
	//   -h, --help   help for cli
	//
	// Use "cli [command] --help" for more information about a command.
}

func ExampleApp_cliElection_help() {
	app := newTestApp()
	defer app.Stop()
	cmd := vote.GetCobraRootCommand(app)
	cmd.SetOut(NewTrimmingWriter(os.Stdout))
	cmd.SetArgs([]string{"election", "-h"})
	_ = cmd.Execute()

	// Output:
	// election subdomain actions
	//
	// Usage:
	//   cli election [command]
	//
	// Available Commands:
	//   CastVote
	//   CloseElectionByOwner
	//   CommenceElection
	//   GetElection
	//   GetElectionResults
	//   GetProposalDetails
	//   ListOpenElections
	//   ListProposals
	//   MakeProposal
	//
	// Flags:
	//   -h, --help   help for election
	//
	// Use "cli election [command] --help" for more information about a command.
}

func ExampleApp_cliElectionCastVote_help() {
	app := newTestApp()
	defer app.Stop()
	cmd := vote.GetCobraRootCommand(app)
	cmd.SetOut(NewTrimmingWriter(os.Stdout))
	cmd.SetArgs([]string{"election", "CastVote", "-h"})
	_ = cmd.Execute()

	// Output:
	// Usage:
	//   cli election CastVote [flags]
	//
	// Flags:
	//       --ElectionID string
	//       --RankedProposalIDs strings
	//       --UserID string
	//       --VoteID string
	//   -h, --help                        help for CastVote
}

func ExampleApp_cliElectionListOpenElections_help() {
	app := newTestApp()
	defer app.Stop()
	cmd := vote.GetCobraRootCommand(app)
	cmd.SetOut(NewTrimmingWriter(os.Stdout))
	cmd.SetArgs([]string{"election", "ListOpenElections", "-h"})
	_ = cmd.Execute()

	// Output:
	// Returns:
	// election.ListOpenElectionsResponse {
	//	OpenElections []OpenElection
	//	TotalResults int
	// }
	//
	// Usage:
	//   cli election ListOpenElections [flags]
	//
	// Flags:
	//       --ItemsPerPage int       (optional) 1 - 50
	//       --Page int               (optional) >= 1
	//       --SortBy string          (optional) Name, CommencedAt
	//       --SortDirection string   (optional) ascending, descending
	//   -h, --help                   help for ListOpenElections
}

func ExampleApp_cliElectionCloseElectionByOwner() {
	recordingEventDispatcher := cqrstest.NewRecordingEventDispatcher()
	app := newTestApp(
		vote.WithEventDispatcher(recordingEventDispatcher),
	)
	defer app.Stop()

	recordingEventDispatcher.Add(4)

	cmd := vote.GetCobraRootCommand(app)
	cmd.SetOut(NewTrimmingWriter(os.Stdout))
	cmd.SetArgs([]string{"election", "CommenceElection",
		"--ElectionID", "E1",
		"--Name", "Election Name",
		"--Description", "Election Description",
		"--OrganizerUserID", "U1",
	})
	_ = cmd.Execute()
	cmd.SetArgs([]string{"election", "MakeProposal",
		"--ElectionID", "E1",
		"--ProposalID", "P1",
		"--Name", "Proposal Name",
		"--Description", "Proposal Description",
		"--OwnerUserID", "U2",
	})
	_ = cmd.Execute()
	cmd.SetArgs([]string{"election", "CastVote",
		"--VoteID", "V1",
		"--ElectionID", "E1",
		"--UserID", "U3",
		"--RankedProposalIDs", "P1",
	})
	_ = cmd.Execute()
	cmd.SetArgs([]string{"election", "CloseElectionByOwner", "--ID", "AC1", "--ElectionID", "E1"})
	_ = cmd.Execute()

	ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
	defer done()
	recordingEventDispatcher.Wait(ctx)

	cmd.SetArgs([]string{"async-command-status", "--ID", "AC1", "--logs"})
	_ = cmd.Execute()

	// Output:
	// Response: *cqrs.CommandResponse {
	//   "Status": "OK"
	// }
	// Response: *cqrs.CommandResponse {
	//   "Status": "OK"
	// }
	// Response: *cqrs.CommandResponse {
	//   "Status": "OK"
	// }
	// Response: *cqrs.AsyncCommandResponse {
	//   "ID": "AC1",
	//   "Status": "QUEUED",
	//   "HasBeenQueued": true
	// }
	// AsyncCommandStatus: *cqrs.AsyncCommandStatus {
	//   "Command": {
	//     "ID": "AC1",
	//     "ElectionID": "E1"
	//   },
	//   "CreatedAt": 1699900003,
	//   "ModifiedAt": 1699900007,
	//   "StartedAtMicro": 1699900004000000,
	//   "FinishedAtMicro": 1699900007000000,
	//   "ExecutionDuration": "3s",
	//   "TotalToProcess": 1,
	//   "TotalProcessed": 1,
	//   "PercentDone": 100,
	//   "IsSuccess": true,
	//   "IsFinished": true
	// }
	// AsyncCommandLogs: []cqrs.AsyncCommandLog [
	//   {
	//     "Type": "INFO",
	//     "CreatedAtMicro": 1699900006000000,
	//     "Message": "Closing election with winner: P1"
	//   }
	// ]
}

func ExampleApp_cliListOpenElections() {
	recordingEventDispatcher := cqrstest.NewRecordingEventDispatcher()
	app := newTestApp(
		vote.WithEventDispatcher(recordingEventDispatcher),
	)
	defer app.Stop()

	recordingEventDispatcher.Add(4)

	cmd := vote.GetCobraRootCommand(app)
	cmd.SetOut(io.Discard)
	cmd.SetArgs([]string{"election", "CommenceElection",
		"--ElectionID", "E1",
		"--Name", "Election Name 1",
		"--Description", "Election Description 1",
		"--OrganizerUserID", "U1",
	})
	_ = cmd.Execute()
	cmd.SetArgs([]string{"election", "CommenceElection",
		"--ElectionID", "E2",
		"--Name", "Election Name 2",
		"--Description", "Election Description 2",
		"--OrganizerUserID", "U1",
	})
	_ = cmd.Execute()
	cmd.SetArgs([]string{"election", "CommenceElection",
		"--ElectionID", "E3",
		"--Name", "Election Name 3",
		"--Description", "Election Description 3",
		"--OrganizerUserID", "U1",
	})
	_ = cmd.Execute()

	cmd.SetOut(NewTrimmingWriter(os.Stdout))
	cmd.SetArgs([]string{"election", "ListOpenElections",
		"--ItemsPerPage", "2",
		"--Page", "1",
	})
	_ = cmd.Execute()

	// Output:
	// Response: election.ListOpenElectionsResponse {
	//   "OpenElections": [
	//     {
	//       "ElectionID": "E1",
	//       "OrganizerUserID": "U1",
	//       "Name": "Election Name 1",
	//       "Description": "Election Description 1",
	//       "CommencedAt": 1699900000
	//     },
	//     {
	//       "ElectionID": "E2",
	//       "OrganizerUserID": "U1",
	//       "Name": "Election Name 2",
	//       "Description": "Election Description 2",
	//       "CommencedAt": 1699900001
	//     }
	//   ],
	//   "TotalResults": 3
	// }
}
