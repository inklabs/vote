<template>
  <v-container v-show="loading">
    <v-progress-linear
      color="light-blue"
      height="10"
      v-model="progress"
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
    <v-alert
      v-for="(error, index) in errors"
      :key="index"
      :text="error"
      border="start"
      type="error"
      variant="tonal"
      ></v-alert>
  </v-container>
</template>

<script>
export default {
  props: [
    "asyncCommandID",
    "successTitle",
    "successMsg",
    "errorTitle",
    "errorMsg",
  ],
  data() {
    return {
      loading: false,
      showSuccess: false,
      showError: false,
      errors: [],
      progress: 0,
    }
  },
  methods: {
    async loadProgress() {
      this.loading = true;
      try {
        const body = await this.$sdk.AsyncCommandStatus({
          AsyncCommandID: this.asyncCommandID,
        });

        this.progress = body.data.attributes.PercentDone
        if (!body.data.attributes.IsFinished) {
          setTimeout(this.loadProgress, 500);
          return
        }

        this.loading = false;
        if (body.data.attributes.IsSuccess) {
          this.showSuccess = true;
        } else {
          this.showLogsWithError()
        }
      } catch (error) {
        console.error('Error fetching async command status:', error);
      }
    },
    async showLogsWithError() {
      try {
        const body = await this.$sdk.AsyncCommandStatus({
          AsyncCommandID: this.asyncCommandID,
          IncludeLogs: true,
        });

        body.included.forEach(asyncCommandLog => {
          this.errors.push(asyncCommandLog.attributes.Message)
        })
        this.showError = true;

      } catch (error) {
        console.error('Error fetching async command status logs:', error);
      }
    }
  },
  watch: {
    asyncCommandID() {
        if (!this.asyncCommandID) {
          console.log('received invalid asyncCommandID')
          return
        }

        this.showSuccess = false;
        this.showError = false;
        this.errors = [];
        this.loadProgress();
      }
  }
}
</script>
