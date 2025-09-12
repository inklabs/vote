package vote_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/inklabs/cqrs/cqrstest"
	"github.com/inklabs/cqrs/jsonapi"

	"github.com/inklabs/vote"
	"github.com/inklabs/vote/action/election"
)

func ExampleApp_httpCloseElectionByOwner() {
	const (
		baseUri       = "http://example.com"
		schemaBaseUri = "http://example.com/schema"
		version       = "1.0.0"
	)
	recordingEventDispatcher := cqrstest.NewRecordingEventDispatcher()
	app := newTestApp(
		vote.WithEventDispatcher(recordingEventDispatcher),
	)
	defer app.Stop()
	api, _ := jsonapi.New(app, vote.NewHTTPActionDecoder(), baseUri, schemaBaseUri, version)

	recordingEventDispatcher.Add(4)

	body, _ := json.Marshal(election.CommenceElection{
		ElectionID:      "E1",
		OrganizerUserID: "U1",
		Name:            "Election Name",
		Description:     "Election Description",
	})
	request := httptest.NewRequest(http.MethodPost, "/election/CommenceElection", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	api.ServeHTTP(response, request)
	PrettyPrint(response.Body)

	body, _ = json.Marshal(election.MakeProposal{
		ElectionID:  "E1",
		ProposalID:  "P1",
		OwnerUserID: "U2",
		Name:        "Proposal Name",
		Description: "Proposal Description",
	})
	request = httptest.NewRequest(http.MethodPost, "/election/MakeProposal", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	api.ServeHTTP(response, request)
	PrettyPrint(response.Body)

	body, _ = json.Marshal(election.CastVote{
		VoteID:            "V1",
		ElectionID:        "E1",
		UserID:            "U3",
		RankedProposalIDs: []string{"P1"},
	})
	request = httptest.NewRequest(http.MethodPost, "/election/CastVote", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	api.ServeHTTP(response, request)
	PrettyPrint(response.Body)

	body, _ = json.Marshal(election.CloseElectionByOwner{
		ID:         "AC1",
		ElectionID: "E1",
	})
	request = httptest.NewRequest(http.MethodPost, "/election/CloseElectionByOwner", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	api.ServeHTTP(response, request)
	PrettyPrint(response.Body)

	ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
	defer done()
	recordingEventDispatcher.Wait(ctx)

	request = httptest.NewRequest(http.MethodGet, "/async-command-status/AC1?include_logs=true", nil)
	response = httptest.NewRecorder()
	api.ServeHTTP(response, request)
	PrettyPrint(response.Body)

	// Output:
	// {
	//   "data": {
	//     "attributes": {
	//       "Status": "OK"
	//     },
	//     "type": "CommandResponse"
	//   },
	//   "links": {
	//     "docs": "http://example.com/schema/election/CommenceElection",
	//     "self": "http://example.com/election/CommenceElection"
	//   },
	//   "meta": {
	//     "request": {
	//       "attributes": {
	//         "ElectionID": "E1",
	//         "OrganizerUserID": "U1",
	//         "Name": "Election Name",
	//         "Description": "Election Description"
	//       },
	//       "type": "election.CommenceElection"
	//     }
	//   }
	// }
	// {
	//   "data": {
	//     "attributes": {
	//       "Status": "OK"
	//     },
	//     "type": "CommandResponse"
	//   },
	//   "links": {
	//     "docs": "http://example.com/schema/election/MakeProposal",
	//     "self": "http://example.com/election/MakeProposal"
	//   },
	//   "meta": {
	//     "request": {
	//       "attributes": {
	//         "ElectionID": "E1",
	//         "ProposalID": "P1",
	//         "OwnerUserID": "U2",
	//         "Name": "Proposal Name",
	//         "Description": "Proposal Description"
	//       },
	//       "type": "election.MakeProposal"
	//     }
	//   }
	// }
	// {
	//   "data": {
	//     "attributes": {
	//       "Status": "OK"
	//     },
	//     "type": "CommandResponse"
	//   },
	//   "links": {
	//     "docs": "http://example.com/schema/election/CastVote",
	//     "self": "http://example.com/election/CastVote"
	//   },
	//   "meta": {
	//     "request": {
	//       "attributes": {
	//         "VoteID": "V1",
	//         "ElectionID": "E1",
	//         "UserID": "U3",
	//         "RankedProposalIDs": [
	//           "P1"
	//         ]
	//       },
	//       "type": "election.CastVote"
	//     }
	//   }
	// }
	// {
	//   "data": {
	//     "attributes": {
	//       "ID": "AC1",
	//       "Status": "QUEUED",
	//       "HasBeenQueued": true
	//     },
	//     "type": "AsyncCommandResponse"
	//   },
	//   "links": {
	//     "docs": "http://example.com/schema/election/CloseElectionByOwner",
	//     "self": "http://example.com/election/CloseElectionByOwner",
	//     "status": "http://example.com/async-command-status/AC1"
	//   },
	//   "meta": {
	//     "request": {
	//       "attributes": {
	//         "ID": "AC1",
	//         "ElectionID": "E1"
	//       },
	//       "type": "election.CloseElectionByOwner"
	//     }
	//   }
	// }
	// {
	//   "data": {
	//     "attributes": {
	//       "Command": {
	//         "ID": "AC1",
	//         "ElectionID": "E1"
	//       },
	//       "CreatedAt": 1699900003,
	//       "ModifiedAt": 1699900007,
	//       "StartedAtMicro": 1699900004000000,
	//       "FinishedAtMicro": 1699900007000000,
	//       "ExecutionDuration": "3s",
	//       "TotalToProcess": 1,
	//       "TotalProcessed": 1,
	//       "PercentDone": 100,
	//       "IsSuccess": true,
	//       "IsFinished": true
	//     },
	//     "type": "AsyncCommandStatus"
	//   },
	//   "included": [
	//     {
	//       "attributes": {
	//         "Type": "INFO",
	//         "CreatedAtMicro": 1699900006000000,
	//         "Message": "Closing election with winner: P1"
	//       },
	//       "type": "AsyncCommandLog"
	//     }
	//   ],
	//   "links": {
	//     "self": "http://example.com/async-command-status/AC1",
	//     "self-include-logs": "http://example.com/async-command-status/AC1?include_logs=true"
	//   }
	// }
}

func ExampleApp_httpListOpenElections() {
	const (
		baseUri       = "http://example.com"
		schemaBaseUri = "http://example.com/schema"
		version       = "1.0.0"
	)
	recordingEventDispatcher := cqrstest.NewRecordingEventDispatcher()
	app := newTestApp(
		vote.WithEventDispatcher(recordingEventDispatcher),
	)
	defer app.Stop()
	api, _ := jsonapi.New(app, vote.NewHTTPActionDecoder(), baseUri, schemaBaseUri, version)

	commands := []election.CommenceElection{
		{
			ElectionID:      "E1",
			OrganizerUserID: "U1",
			Name:            "Election Name 1",
			Description:     "Election Description 1",
		},
		{
			ElectionID:      "E2",
			OrganizerUserID: "U1",
			Name:            "Election Name 2",
			Description:     "Election Description 2",
		},
		{
			ElectionID:      "E3",
			OrganizerUserID: "U1",
			Name:            "Election Name 3",
			Description:     "Election Description 3",
		},
	}

	for _, command := range commands {
		body, _ := json.Marshal(command)
		request := httptest.NewRequest(http.MethodPost, "/election/CommenceElection", bytes.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		api.ServeHTTP(response, request)
	}

	request := httptest.NewRequest(http.MethodGet, "/election/ListOpenElections?SortBy=Name&SortDirection=ascending&ItemsPerPage=2&Page=1", nil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	api.ServeHTTP(response, request)
	PrettyPrint(response.Body)

	// Output:
	// {
	//   "data": {
	//     "attributes": {
	//       "OpenElections": [
	//         {
	//           "ElectionID": "E1",
	//           "OrganizerUserID": "U1",
	//           "Name": "Election Name 1",
	//           "Description": "Election Description 1",
	//           "CommencedAt": 1699900000
	//         },
	//         {
	//           "ElectionID": "E2",
	//           "OrganizerUserID": "U1",
	//           "Name": "Election Name 2",
	//           "Description": "Election Description 2",
	//           "CommencedAt": 1699900001
	//         }
	//       ],
	//       "TotalResults": 3
	//     },
	//     "type": "election.ListOpenElectionsResponse"
	//   },
	//   "links": {
	//     "docs": "http://example.com/schema/election/ListOpenElections",
	//     "self": "http://example.com/election/ListOpenElections?SortBy=Name&SortDirection=ascending&ItemsPerPage=2&Page=1"
	//   },
	//   "meta": {
	//     "request": {
	//       "attributes": {
	//         "Page": 1,
	//         "ItemsPerPage": 2,
	//         "SortBy": "Name",
	//         "SortDirection": "ascending"
	//       },
	//       "type": "election.ListOpenElections"
	//     },
	//     "status": "OK"
	//   }
	// }
}
