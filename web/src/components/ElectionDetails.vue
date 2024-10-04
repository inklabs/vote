<template>
  <v-container>
    <div class="text-center">
      <h1 class="text-h2 font-weight-bold">{{ election.Name }}</h1>
      <p>{{ election.Description }}</p>
    </div>
  </v-container>
</template>

<script>
export default {
  props: ['electionID'],
  data() {
    return {
      election: {},
    }
  },
  mounted() {
    this.fetchElection(this.electionID);
  },
  methods: {
    fetchElection(electionID) {
      fetch(`http://localhost:8080/election/GetElection?ElectionID=${electionID}`, {
        method: "GET",
      })
        .then(response => {
          response.json().then((body) => {
            this.election = body.data.attributes;
          })
        })
        .catch(error => {
          console.error('Error fetching election:', error);
        });
    },
  }
}
</script>
