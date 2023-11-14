from __future__ import print_function
import sys
sys.path.append("grpc/python")
from google.protobuf.json_format import MessageToJson
from electionpb.election_pb2 import ListOpenElectionsRequest
from electionpb.election_pb2 import CommenceElectionRequest
from electionpb import election_pb2_grpc

import logging
import grpc


def run():
    """
    This function demonstrates how to use the ElectionStub.

    Examples:
    >>> run()  # Assuming gRPC server is running locally on port 8081
    client received: {
      "openElections": [
        {
          "electionId": "E1",
          "organizerUserId": "U1",
          "name": "Election Name 1",
          "description": "Election Description 1",
          "commencedAt": "1699900000"
        },
        {
          "electionId": "E2",
          "organizerUserId": "U1",
          "name": "Election Name 2",
          "description": "Election Description 2",
          "commencedAt": "1699900001"
        }
      ]
    }
    Note: The actual output might vary depending on the gRPC server's response.
    """
    with grpc.insecure_channel("localhost:8081") as channel:
        stub = election_pb2_grpc.ElectionStub(channel)
        stub.CommenceElection(
            CommenceElectionRequest(
                election_id="E1",
                organizer_user_id="U1",
                name="Election Name 1",
                description="Election Description 1"
            )
        )
        stub.CommenceElection(
            CommenceElectionRequest(
                election_id="E2",
                organizer_user_id="U1",
                name="Election Name 2",
                description="Election Description 2"
            )
        )
        response = stub.ListOpenElections(
            ListOpenElectionsRequest(
                page=1,
                items_per_page=10,
                sort_by="Name",
                sort_direction="ascending"
            )
        )
    print("client received: " + MessageToJson(response))


if __name__ == "__main__":
    logging.basicConfig()
    run()
