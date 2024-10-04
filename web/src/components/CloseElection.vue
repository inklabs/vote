<template>
  <v-btn
    class="ma-2"
    @click="closeElection()"
  >
    <v-icon
      icon="mdi-clock-start"
      start
    ></v-icon>
    Close Election
  </v-btn>

  <AsyncCommandProgress
    :asyncCommandID="asyncCommandID"
    successTitle="Success Closing Election"
    errorTitle="Error Closing Election"
  />

</template>

<script>
export default {
  props: ['electionID'],
  data() {
    return {
      loading: true,
      asyncCommandID: "",
    }
  },
  methods: {
    closeElection() {
      const asyncCommandID = this.$uuid.v4()
      fetch('http://localhost:8080/election/CloseElectionByOwner', {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
          ID: asyncCommandID,
          ElectionID: this.electionID,
        })
      })
        .then(response => {
          response.json().then((body) => {
            if (body.data.attributes.Status === "QUEUED") {
              this.asyncCommandID = body.data.attributes.ID
            }
          })
        })
        .catch(error => {
          console.error('Error closing election:', error);
        });
    },
  }
}
</script>
