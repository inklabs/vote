<template>
  <v-container>
    <h2 class="text-h4 font-weight-bold">Ballot Proposals</h2>
    <p>Drag and drop to reorder your preferred choices and click "Cast This Ballot"</p>

    <div class="py-4"/>

    <VueDraggable v-model="proposals" target=".sort-target" :animation="150">
      <v-table theme="dark">
        <thead>
        <tr>
          <th></th>
          <th>Name</th>
          <th>Description</th>
        </tr>
        </thead>
        <tbody class="sort-target">
        <tr
          v-for="proposal in proposals"
          :key="proposal.ProposalID"
          class="cursor-move"
        >
          <td><v-icon icon="mdi-drag"></v-icon></td>
          <td>{{ proposal.Name }}</td>
          <td>{{ proposal.Description }}</td>
        </tr>
        </tbody>
      </v-table>
    </VueDraggable>

    <div class="py-4"/>

    <v-btn
      class="ma-2"
      @click="castBallot()"
    >
      <v-icon
        icon="mdi-vote"
        start
      ></v-icon>
      Cast This Ballot
    </v-btn>

  </v-container>
</template>

<script>
import {VueDraggable} from 'vue-draggable-plus'

export default {
  components: {
    VueDraggable,
  },
  props: ['electionID'],
  data() {
    return {
      headers: [
        {title: '', key: ''},
        {title: '#', key: ''},
        {title: 'Name', key: 'Name'},
        {title: 'Description', key: 'Description'},
        {title: 'Proposed', key: 'ProposedAt'},
      ],
      pagination: {
        itemsPerPage: 10,
        totalResults: 0,
      },
      proposals: [],
      drag: false,
      loading: true,
    }
  },
  mounted() {
    this.fetchProposals({page: 1, itemsPerPage: 10});
  },
  methods: {
    fetchProposals({page, itemsPerPage}) {
      this.loading = true;
      fetch(`http://localhost:8080/election/ListProposals?ElectionID=${this.electionID}&Page=${page}&ItemsPerPage=${itemsPerPage}`, {
        method: "GET",
      })
        .then(response => {
          response.json().then((body) => {
            this.proposals = body.data.attributes.Proposals;

            if (body.data.attributes.TotalResults > itemsPerPage) {
              console.log("more results than were returned");
            }
            this.loading = false;
          })
        })
        .catch(error => {
          console.error('Error fetching proposals:', error);
        });
    },
    castBallot() {
      const rankedProposalsIDs = this.proposals.map(a => a.ProposalID);

      fetch(`http://localhost:8080/election/CastVote`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
          VoteID: this.$uuid.v4(),
          ElectionID: this.electionID,
          UserID: this.$uuid.v4(),
          RankedProposalIDs: rankedProposalsIDs
        })
      })
        .then(response => {
          response.json().then((body) => {
            if (body.data.attributes.Status !== "OK") {
              console.log("unable to cast vote")
              console.log(body)
              return
            }

            this.$router.push(`/elections`);
          })
        })
        .catch(error => {
          console.error('Error casting vote:', error);
        })
    },
  }
}
</script>
