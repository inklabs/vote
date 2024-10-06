class VoteSDK {
  constructor(baseURL) {
    this.baseURL = baseURL;
    this.election = new ElectionCommandSDK(baseURL);
  }

  async AsyncCommandStatus({AsyncCommandID, IncludeLogs}) {
    let params = ""
    if (IncludeLogs) {
      params = "?"+new URLSearchParams({"include_logs": "true"}).toString()
    }
    const response = await fetch(`${this.baseURL}/async-command-status/${AsyncCommandID}${params}`, {
      method: "GET",
    });

    if (!response.ok) {
      throw new Error(`unable to get async command status for ${AsyncCommandID}`);
    }

    return await response.json();
  }
}

class ElectionCommandSDK {
  constructor(baseURL) {
    this.baseURL = `${baseURL}/election`;
  }

  async CastVote({VoteID, ElectionID, UserID, RankedProposalIDs}) {
    return this._executeCommand("CastVote", {
      VoteID,
      ElectionID,
      UserID,
      RankedProposalIDs,
    });
  }

  async CommenceElection({ElectionID, Name, Description}) {
    return this._executeCommand("CommenceElection", {
      ElectionID,
      Name,
      Description,
    });
  }

  async MakeProposal({ElectionID, ProposalID, Name, Description}) {
    return this._executeCommand("MakeProposal", {
      ElectionID,
      ProposalID,
      Name,
      Description,
    });
  }

  async CloseElectionByOwner({ID, ElectionID}) {
    return this._executeAsyncCommand("CloseElectionByOwner", {
      ID,
      ElectionID,
    });
  }

  async GetElection({ElectionID}) {
    return this._executeQuery("GetElection", {
      ElectionID,
    });
  }

  async ListOpenElections({Page, ItemsPerPage, SortBy, SortDirection}) {
    return this._executeQuery("ListOpenElections", {
      Page,
      ItemsPerPage,
      SortBy,
      SortDirection,
    });
  }

  async ListProposals({ElectionID, Page, ItemsPerPage}) {
    return this._executeQuery("ListProposals", {
      ElectionID,
      Page,
      ItemsPerPage
    });
  }

  async _executeQuery(queryName, queryParams) {
    const params = new URLSearchParams(queryParams).toString()
    const response = await fetch(`${this.baseURL}/${queryName}?${params}`, {
      method: "GET",
    });

    if (!response.ok) {
      throw new Error(`unable to execute query ${queryName}`);
    }

    const responseBody = await response.json();
    if (responseBody.meta.status !== "OK") {
      throw new Error(`query response was not successful for ${queryName}`);
    }
    return responseBody;
  }

  async _executeAsyncCommand(commandName, payload) {
    const response = await fetch(`${this.baseURL}/${commandName}`, {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify(payload)
    });

    if (!response.ok) {
      throw new Error(`unable to execute async command ${commandName}`);
    }

    const responseBody = await response.json();
    if (responseBody.data.attributes.Status !== "QUEUED" && responseBody.data.attributes.Status !== "OK") {
      throw new Error(`async command response was not successful for ${commandName}`);
    }
    return responseBody;
  }

  async _executeCommand(commandName, payload) {
    const response = await fetch(`${this.baseURL}/${commandName}`, {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify(payload)
    });

    if (!response.ok) {
      throw new Error(`unable to execute command ${commandName}`);
    }

    const responseBody = await response.json();
    if (responseBody.data.attributes.Status !== "OK") {
      throw new Error(`command response was not successful for ${commandName}`);
    }
    return responseBody;
  }
}

export {VoteSDK};
