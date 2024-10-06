<template>
  <v-container>
    <h2 class="text-h4 font-weight-bold">Commence Election</h2>

    <v-form v-model="isFormValid">
      <v-text-field
        label="Name"
        v-model="election.name"
        :rules="nameRules"
        required
      ></v-text-field>

      <v-text-field
        label="Description"
        v-model="election.description"
        :rules="descriptionRules"
        required
      ></v-text-field>

      <v-btn :disabled="!isFormValid" @click="startElection">Start Election</v-btn>
    </v-form>
  </v-container>
</template>

<script>
export default {
  data() {
    return {
      isFormValid: false,
      election: {
        name: '',
        description: ''
      },
      nameRules: [v => !!v || 'Name is required'],
      descriptionRules: [v => !!v || 'Description is required']
    }
  },
  methods: {
    async startElection() {
      const electionID = this.$uuid.v4()

      try {
        await this.$sdk.election.CommenceElection({
          ElectionID: electionID,
          Name: this.election.name,
          Description: this.election.description,
        });

        this.$router.push(`/elections/${electionID}`)
      } catch (error) {
        console.error('Error commencing election:', error);
      }
    }
  }
}
</script>
