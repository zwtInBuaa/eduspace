<template>
  <v-container>
    <v-card flat class="pt-2">
      <v-btn icon @click="$router.go(-1)">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <!-- Markdown Display -->
      <v-container class="pa-2">
        <v-card-title>{{ this.title }}</v-card-title>
        <v-card>
          <MarkdownDisplay :content="this.desc" :minHeight="10" />
        </v-card>
      </v-container>
      <!--题目列表-->
      <exercise-list :exercise-list="exerciseList" @clickRow="rowOnClick"></exercise-list>
    </v-card>
  </v-container>
</template>

<script>
import ExerciseList from '@/components/exercise/ExerciseList.vue';
import MarkdownDisplay from '@/components/MarkdownDisplay.vue';
import { getExerciseSetDetail } from '@/api/exercise';

export default {
  components: {
    ExerciseList,
    MarkdownDisplay
  },
  data() {
    return {
      id: 999,
      title: 'avj',
      desc: 'avvkvbdskvbkvbdskvbsdkvjs',
      exerciseList: []
    };
  },
  computed: {},
  methods: {
    rowOnClick(item) {
      // // console.log('rowOnClick' + item.id);
      this.$router.push('/exercise/' + item.id);
    }
  },
  created() {
    this.id = this.$route.params.id;
    getExerciseSetDetail(this.id).then((response) => {
      // // console.log(response);
      this.id = response.id;
      this.title = response.data.name;
      this.desc = response.data.description;
      this.exerciseList = response.data.questions;
    });
  }
};
</script>
