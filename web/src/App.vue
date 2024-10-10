<template>
  <v-app>
    <v-main>
      <router-view />
    </v-main>

    <v-snackbar
      v-model="snackbar"
      color="error"
      timeout="2500"
    >
      {{ snackbarMessage }}

      <template v-slot:actions>
        <v-btn
          variant="text"
          @click="snackbar = false"
        >
          Close
        </v-btn>
      </template>
    </v-snackbar>
  </v-app>
</template>

<script>
  export default {
    name: 'App',
    data: () => ({
      snackbar: false,
      snackbarMessage: '',
    }),
    mounted() {
      this.$onSnackbar(message => {
        this.snackbarMessage = message;
        this.snackbar = true;
      })
    }
  }
</script>
