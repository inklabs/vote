package vote_test

import (
	"os"
	"time"

	"github.com/inklabs/cqrs"
	"github.com/inklabs/cqrs/asynccommandstore"
	"github.com/inklabs/cqrs/pkg/clock/provider/incrementingclock"

	"github.com/inklabs/vote"
)

func ExampleApp_cliRoot() {
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
	//   election             8 actions: [CastVote, CloseElectionByOwner, CommenceElection, GetElectionResults, GetProposalDetails, ListOpenElections, ListProposals, MakeProposal]
	//   help                 Help about any command
	//
	// Flags:
	//   -h, --help   help for cli
	//
	// Use "cli [command] --help" for more information about a command.
}

func ExampleApp_cliElection() {
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

func ExampleApp_cliElectionCastVote() {
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
	//   -h, --help                        help for CastVote
}

func ExampleApp_cliElectionListOpenElections() {
	app := newTestApp()
	defer app.Stop()
	cmd := vote.GetCobraRootCommand(app)
	cmd.SetOut(NewTrimmingWriter(os.Stdout))
	cmd.SetArgs([]string{"election", "ListOpenElections", "-h"})
	_ = cmd.Execute()

	// Output:
	// Returns:
	// election.ListOpenElectionsResponse {
	// 	OpenElections []OpenElection
	// }
	//
	// Usage:
	//   cli election ListOpenElections [flags]
	//
	// Flags:
	//       --ItemsPerPage int       (optional) 1 - 10
	//       --Page int               (optional) >= 1
	//       --SortBy string          (optional) Name, CommencedAt
	//       --SortDirection string   (optional) ascending, descending
	//   -h, --help                   help for ListOpenElections
}

func ExampleApp_cliElectionCloseElectionByOwner() {
	app := newTestApp()
	defer app.Stop()
	cmd := vote.GetCobraRootCommand(app)
	cmd.SetOut(NewTrimmingWriter(os.Stdout))
	cmd.SetArgs([]string{"election", "CommenceElection",
		"--ElectionID", "E1",
		"--Name", "Election Name",
		"--Description", "Election Description",
		"--OrganizerUserID", "U2",
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
		"--ElectionID", "E1",
		"--UserID", "U3",
		"--RankedProposalIDs", "P1",
	})
	_ = cmd.Execute()
	cmd.SetArgs([]string{"election", "CloseElectionByOwner", "--ID", "AC1", "--ElectionID", "E1"})
	_ = cmd.Execute()
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
	//   "TotalToProcess": 0,
	//   "TotalProcessed": 0,
	//   "PercentDone": 0,
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

func newTestApp() cqrs.App {
	startTime := time.Unix(1699900000, 0)
	seededClock := incrementingclock.New(startTime)

	return vote.NewApp(
		vote.WithAsyncCommandStore(asynccommandstore.NewInMemory()),
		vote.WithSyncLocalAsyncCommandBus(),
		vote.WithClock(seededClock),
	)
}
