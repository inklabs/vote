<template>
  <v-container v-show="isLoading">
    <v-progress-linear
      color="light-blue"
      height="10"
      indeterminate
      striped
    ></v-progress-linear>
  </v-container>

  <v-container v-show="showSuccess">
    <v-alert
      :text="successMsg"
      :title="successTitle"
      type="success"
    />
  </v-container>

  <v-container v-show="showError">
    <v-alert
      :text="errorMsg"
      :title="errorTitle"
      type="error"
    />
  </v-container>
</template>

<script>
export default {
  props: [
    "status",
    "loading",
    "successTitle",
    "successMsg",
    "errorTitle",
    "errorMsg",
  ],
  data() {
    return {
      isLoading: false,
      showSuccess: false,
      showError: false,
    }
  },
  watch: {
    status() {
      if (!this.status) {
        console.log('received invalid status')
        return
      }

      this.showSuccess = false;
      this.showError = false;

      if (this.status === 'success') {
        this.showSuccess = true;
      } else if (this.status === 'error') {
        this.showError = true;
      }
    },
    loading() {
      this.isLoading = this.loading;
    }
  }
}
</script>
