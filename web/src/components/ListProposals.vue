<template>
  <v-container>
    <h2 class="text-h4 font-weight-bold">Proposals</h2>

    <div class="py-4"/>

    <v-data-table-server
      v-model:items-per-page="pagination.itemsPerPage"
      :items-per-page-options="[
        {value:5, title: '5'},
        {value:10, title: '10'}
      ]"
      :headers="headers"
      :items="proposals"
      :items-length="pagination.totalResults"
      :loading="loading"
      item-value="ElectionID"
      @update:options="fetchProposals"
      class="elevation-1"
    >
      <template v-slot:item.ProposedAt="{ value }">
        {{ new Date(value * 1000).toLocaleString() }}
      </template>

    </v-data-table-server>
  </v-container>
</template>

<script>
export default {
  props: ['electionID'],
  data() {
    return {
      headers: [
        {title: 'Name', key: 'Name', sortable: false},
        {title: 'Description', key: 'Description', sortable: false},
        {title: 'Proposed', key: 'ProposedAt', sortable: false},
      ],
      pagination: {
        itemsPerPage: 10,
        totalResults: 0,
      },
      proposals: [],
      loading: true,
    }
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
            this.pagination.totalResults = body.data.attributes.TotalResults;
            this.loading = false;
          })
        })
        .catch(error => {
          console.error('Error fetching proposals:', error);
        });
    },
  }
}
</script>
