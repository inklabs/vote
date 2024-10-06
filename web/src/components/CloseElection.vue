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
    async closeElection() {
      const asyncCommandID = this.$uuid.v4()

      try {
        await this.$sdk.election.CloseElectionByOwner({
          ID: asyncCommandID,
          ElectionID: this.electionID,
        })

        this.asyncCommandID = asyncCommandID
      } catch (error) {
        console.error('Error closing election:', error);
      }
    },
  }
}
</script>
