<template>
  <v-container>
    <v-responsive
      class="align-centerfill-height mx-auto"
      max-width="900"
    >
      <v-row>
        <v-col cols="12">
          <ElectionDetails :electionID="electionID"/>
        </v-col>
      </v-row>

      <div class="py-4"/>

      <v-row>
        <v-col cols="12">
          <CloseElection :electionID="electionID"/>
        </v-col>
      </v-row>

      <div class="py-4"/>

      <v-row>
        <v-col cols="12">
          <ListProposalsAndVote :electionID="electionID"/>
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="12">
          <MakeProposal :electionID="electionID"/>
        </v-col>
      </v-row>

    </v-responsive>
  </v-container>
  <AppFooter/>
</template>

<script setup>
import {useRoute} from 'vue-router'

const electionID = useRoute().params.id;
</script>
