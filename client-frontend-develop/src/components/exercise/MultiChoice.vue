<template>
  <v-container class="pa-0">
    <v-card>
      <v-card-title class="pb-0"> {{ problem_title }} </v-card-title>
      <v-container class="pb-0">
        <v-checkbox
          v-for="(choice, index) in choices"
          :key="index"
          :value="index"
          :label="choice.toString()"
          :color="result === 'true' ? 'success' : result === 'false' ? 'error' : 'primary'"
          v-model="userChoice"
          height="0"
          class="pa-0"
        ></v-checkbox>
      </v-container>
      <v-card-actions>
        <v-btn color="primary" text @click="submit">提交</v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script>
export default {
  name: 'ExerciseCheckedBox',
  props: {
    problem_id: {
      type: Number,
      required: true
    },
    problem_title: {
      type: String,
      required: true
    },
    choices: {
      type: Array,
      required: true
    },
    result: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      userChoice: []
    };
  },
  methods: {
    submit() {
      if (this.userChoice.length === 0) {
        alert('请至少选择一个选项');
      } else {
        let answer = '';
        this.userChoice.sort();
        this.userChoice.forEach((item) => {
          answer += item + ';';
        });
        this.$emit('submit', answer);
      }
    }
  }
};
</script>
