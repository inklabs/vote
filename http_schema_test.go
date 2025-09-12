package vote_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/inklabs/cqrs/generator"
	"github.com/inklabs/cqrs/jsonapi/schemaapi"

	"github.com/inklabs/vote"
)

var domain, _ = generator.LoadDomainFromBytes(vote.DomainBytes)

func ExampleApp_httpSchemaRoot() {
	api, _ := schemaapi.New(domain, vote.ValidationRules, "http://example.com", "1.0.0")

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	api.ServeHTTP(response, request)

	PrettyPrint(response.Body)

	// Output:
	// {
	//   "data": {
	//     "relationships": {
	//       "subdomain": {
	//         "data": [
	//           {
	//             "attributes": {
	//               "name": "election"
	//             },
	//             "links": "http://example.com/election",
	//             "meta": {
	//               "actions": [
	//                 "CastVote",
	//                 "CloseElectionByOwner",
	//                 "CommenceElection",
	//                 "GetElection",
	//                 "GetElectionResults"
	//               ],
	//               "totalActions": 9
	//             },
	//             "type": "Subdomain"
	//           }
	//         ]
	//       }
	//     },
	//     "type": "Domain"
	//   },
	//   "links": {
	//     "self": "http://example.com"
	//   },
	//   "meta": {
	//     "version": "1.0.0"
	//   }
	// }
}

func ExampleApp_httpSchemaElection() {
	api, _ := schemaapi.New(domain, vote.ValidationRules, "http://example.com", "1.0.0")

	request := httptest.NewRequest(http.MethodGet, "/election", nil)
	response := httptest.NewRecorder()
	api.ServeHTTP(response, request)

	PrettyPrint(response.Body)

	// Output:
	// {
	//   "data": {
	//     "attributes": {
	//       "name": "election"
	//     },
	//     "type": "Subdomain"
	//   },
	//   "links": {
	//     "parent": "http://example.com",
	//     "self": "http://example.com/election"
	//   },
	//   "relationships": {
	//     "command": {
	//       "data": [
	//         {
	//           "attributes": {
	//             "name": "CastVote"
	//           },
	//           "links": {
	//             "self": "http://example.com/election/CastVote"
	//           },
	//           "type": "command"
	//         },
	//         {
	//           "attributes": {
	//             "name": "CloseElectionByOwner"
	//           },
	//           "isAsyncCommand": true,
	//           "links": {
	//             "self": "http://example.com/election/CloseElectionByOwner"
	//           },
	//           "type": "command"
	//         },
	//         {
	//           "attributes": {
	//             "name": "CommenceElection"
	//           },
	//           "links": {
	//             "self": "http://example.com/election/CommenceElection"
	//           },
	//           "type": "command"
	//         },
	//         {
	//           "attributes": {
	//             "name": "MakeProposal"
	//           },
	//           "links": {
	//             "self": "http://example.com/election/MakeProposal"
	//           },
	//           "type": "command"
	//         }
	//       ]
	//     },
	//     "query": {
	//       "data": [
	//         {
	//           "attributes": {
	//             "name": "GetElection"
	//           },
	//           "links": {
	//             "self": "http://example.com/election/GetElection"
	//           },
	//           "type": "query"
	//         },
	//         {
	//           "attributes": {
	//             "name": "GetElectionResults"
	//           },
	//           "links": {
	//             "self": "http://example.com/election/GetElectionResults"
	//           },
	//           "type": "query"
	//         },
	//         {
	//           "attributes": {
	//             "name": "GetProposalDetails"
	//           },
	//           "links": {
	//             "self": "http://example.com/election/GetProposalDetails"
	//           },
	//           "type": "query"
	//         },
	//         {
	//           "attributes": {
	//             "name": "ListOpenElections"
	//           },
	//           "links": {
	//             "self": "http://example.com/election/ListOpenElections"
	//           },
	//           "type": "query"
	//         },
	//         {
	//           "attributes": {
	//             "name": "ListProposals"
	//           },
	//           "links": {
	//             "self": "http://example.com/election/ListProposals"
	//           },
	//           "type": "query"
	//         }
	//       ]
	//     }
	//   }
	// }
}

func ExampleApp_httpSchemaElectionCastVote() {
	api, _ := schemaapi.New(domain, vote.ValidationRules, "http://example.com", "1.0.0")

	request := httptest.NewRequest(http.MethodGet, "/election/CastVote", nil)
	response := httptest.NewRecorder()
	api.ServeHTTP(response, request)

	PrettyPrint(response.Body)

	// Output:
	// {
	//   "data": {
	//     "attributes": {
	//       "documentation": "CastVote casts a ballot for a given ElectionID. RankedProposalIDs contains the\nranked candidates in order of preference: first, second, third and so forth. If your\nfirst choice doesn’t have a chance to win, your ballot counts for your next choice.",
	//       "fields": [
	//         {
	//           "isRequired": true,
	//           "name": "VoteID",
	//           "type": "string"
	//         },
	//         {
	//           "isRequired": true,
	//           "name": "ElectionID",
	//           "type": "string"
	//         },
	//         {
	//           "isRequired": true,
	//           "name": "UserID",
	//           "type": "string"
	//         },
	//         {
	//           "isRequired": false,
	//           "name": "RankedProposalIDs",
	//           "type": "[]string"
	//         }
	//       ],
	//       "name": "CastVote",
	//       "subdomainName": "election"
	//     },
	//     "type": "command"
	//   },
	//   "links": {
	//     "parent": "http://example.com/election",
	//     "self": "http://example.com/election/CastVote"
	//   }
	// }
}

func ExampleApp_httpSchemaElectionListOpenElections() {
	api, _ := schemaapi.New(domain, vote.ValidationRules, "http://example.com", "1.0.0")

	request := httptest.NewRequest(http.MethodGet, "/election/ListOpenElections", nil)
	response := httptest.NewRecorder()
	api.ServeHTTP(response, request)

	PrettyPrint(response.Body)

	// Output:
	// {
	//   "data": {
	//     "attributes": {
	//       "documentation": "ListOpenElections returns a paginated result of elections that are still open.",
	//       "fields": [
	//         {
	//           "isRequired": false,
	//           "name": "Page",
	//           "type": "int",
	//           "validationRule": "(optional) >= 1"
	//         },
	//         {
	//           "isRequired": false,
	//           "name": "ItemsPerPage",
	//           "type": "int",
	//           "validationRule": "(optional) 1 - 50"
	//         },
	//         {
	//           "isRequired": false,
	//           "name": "SortBy",
	//           "type": "string",
	//           "validationRule": "(optional) Name, CommencedAt"
	//         },
	//         {
	//           "isRequired": false,
	//           "name": "SortDirection",
	//           "type": "string",
	//           "validationRule": "(optional) ascending, descending"
	//         }
	//       ],
	//       "name": "ListOpenElections",
	//       "subdomainName": "election"
	//     },
	//     "type": "query"
	//   },
	//   "links": {
	//     "parent": "http://example.com/election",
	//     "self": "http://example.com/election/ListOpenElections"
	//   },
	//   "returnType": {
	//     "fields": [
	//       {
	//         "name": "OpenElections",
	//         "type": "[]OpenElection"
	//       },
	//       {
	//         "name": "TotalResults",
	//         "type": "int"
	//       }
	//     ],
	//     "type": "ListOpenElectionsResponse"
	//   }
	// }
}
