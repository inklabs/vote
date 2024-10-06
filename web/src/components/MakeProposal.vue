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
    async makeProposal() {
      try {
        await this.$sdk.election.MakeProposal({
          ElectionID: this.electionID,
          ProposalID: this.$uuid.v4(),
          Name: this.proposal.name,
          Description: this.proposal.description,
        })
        location.reload();
      } catch (error) {
        console.error('Error making proposal:', error);
      }
    }
  }
}
</script>
