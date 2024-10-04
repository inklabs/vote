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
    startElection() {
      fetch('http://localhost:8080/election/CommenceElection', {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
          ElectionID: this.$uuid.v4(),
          Name: this.election.name,
          Description: this.election.description,
        })
      })
        .then(response => {
          response.json().then((body) => {
            if (body.data.attributes.Status === "OK") {
              const electionID = body.meta.request.attributes.ElectionID
              if (electionID === "") {
                console.error("unable to get election id")
                return
              }

              this.$router.push(`/elections/${electionID}`)
            }
            console.log(body)
          })
        })
        .catch(error => {
          console.error('Error commencing election:', error);
        });
    }
  }
}
</script>
