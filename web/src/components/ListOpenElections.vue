<template>
  <v-container>
    <h2 class="text-h4 font-weight-bold">Open Elections</h2>

    <v-data-table-server
      v-model:items-per-page="pagination.itemsPerPage"
      :items-per-page-options="[
        {value:5, title: '5'},
        {value:10, title: '10'},
        {value:25, title: '25'},
        {value:50, title: '50'}
      ]"
      :headers="headers"
      :items="elections"
      :items-length="pagination.totalResults"
      :loading="loading"
      item-value="ElectionID"
      @update:options="fetchElections"
      class="elevation-1"
    >
      <template v-slot:item.CommencedAt="{ value }">
        {{ new Date(value * 1000).toLocaleString() }}
      </template>
      <template v-slot:item.action="{ item }">
        <v-btn @click="viewProposals(item.ElectionID)" small>View</v-btn>
      </template>
    </v-data-table-server>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      headers: [
        {title: 'Name', key: 'Name'},
        {title: 'Description', key: 'Description', sortable: false},
        {title: 'Commenced', key: 'CommencedAt'},
        {title: '', value: 'action', sortable: false},
      ],
      pagination: {
        itemsPerPage: 10,
        totalResults: 0,
      },
      elections: [],
      loading: true,
    }
  },
  methods: {
    fetchElections({page, itemsPerPage, sortBy}) {
      let sortByVal = "CommencedAt"
      let sortDirection = "descending"
      if (sortBy.length) {
        sortByVal = sortBy[0].key
        sortDirection = sortBy[0].order === "desc" ? "descending" : "ascending";
      }
      this.loading = true;
      fetch(`http://localhost:8080/election/ListOpenElections?SortBy=${sortByVal}&SortDirection=${sortDirection}&Page=${page}&ItemsPerPage=${itemsPerPage}`, {
        method: "GET",
      })
        .then(response => {
          response.json().then((body) => {
            this.elections = body.data.attributes.OpenElections;
            this.pagination.totalResults = body.data.attributes.TotalResults;
            this.loading = false;
          })
        })
        .catch(error => {
          console.error('Error fetching elections:', error);
        });
    },
    viewProposals(electionId) {
      this.$router.push(`/elections/${electionId}`);
    }
  }
}
</script>
