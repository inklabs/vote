package vote_test

import (
	"os"

	"github.com/inklabs/vote"
)

func ExampleApp_cliRoot() {
	app := vote.NewProdApp()
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
	app := vote.NewProdApp()
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
	app := vote.NewProdApp()
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
	app := vote.NewProdApp()
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
