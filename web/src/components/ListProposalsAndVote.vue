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
          <td>
            <v-icon icon="mdi-drag"></v-icon>
          </td>
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

    <CommandStatus
      :status="castBallotStatus"
      :loading="castBallotLoading"
      successTitle="Success Casting Ballot"
      errorTitle="Error Casting Ballot"
    />

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
      castBallotStatus: "",
      castBallotLoading: false,
      drag: false,
      loading: true,
    }
  },
  mounted() {
    this.fetchProposals({page: 1, itemsPerPage: 10});
  },
  methods: {
    async fetchProposals({page, itemsPerPage}) {
      this.loading = true;
      try {
        const body = await this.$sdk.election.ListProposals({
          ElectionID: this.electionID,
          Page: page,
          ItemsPerPage: itemsPerPage,
        });

        this.proposals = body.data.attributes.Proposals;
      } catch (error) {
        console.error('Error fetching proposals:', error);
        this.$showSnackbar("Error fetching proposals");
      }

      this.loading = false;
    },
    async castBallot() {
      this.castBallotLoading = true;
      const rankedProposalsIDs = this.proposals.map(a => a.ProposalID);

      try {
        await this.$sdk.election.CastVote({
          VoteID: this.$uuid.v4(),
          ElectionID: this.electionID,
          UserID: this.$uuid.v4(),
          RankedProposalIDs: rankedProposalsIDs
        });
        this.castBallotStatus = 'success';
      } catch (error) {
        this.castBallotStatus = 'error';
        console.error('Error casting vote:', error);
        this.$showSnackbar("Error casting vote");
      }

      this.castBallotLoading = false;
    },
  }
}
</script>
