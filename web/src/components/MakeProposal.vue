<template>
  <v-container>
    <h2 class="text-h4 font-weight-bold">Make A New Proposal</h2>

    <v-form v-model="isFormValid">
      <v-text-field
        label="Name"
        v-model="proposal.name"
        :rules="nameRules"
        required
      ></v-text-field>

      <v-text-field
        label="Description"
        v-model="proposal.description"
        :rules="descriptionRules"
        required
      ></v-text-field>

      <v-btn :disabled="!isFormValid" @click="makeProposal">Make Proposal</v-btn>
    </v-form>
  </v-container>
</template>

<script>
export default {
  props: ['electionID'],
  data() {
    return {
      isFormValid: false,
      proposal: {
        name: '',
        description: ''
      },
      nameRules: [v => !!v || 'Name is required'],
      descriptionRules: [v => !!v || 'Description is required']
    }
  },
  methods: {
    makeProposal() {
      fetch('http://localhost:8080/election/MakeProposal', {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
          ElectionID: this.electionID,
          ProposalID: this.$uuid.v4(),
          Name: this.proposal.name,
          Description: this.proposal.description,
        })
      })
        .then(response => {
          response.json().then((body) => {
            if (body.data.attributes.Status === "OK") {
              console.log("Make Proposal successfully created")
              location.reload();
            }
            console.log(body)
          })
        })
        .catch(error => {
          console.error('Error making proposal:', error);
        });
    }
  }
}
</script>
