<template>
  <v-container class="pa-0">
    <v-card>
      <v-card-title class="pb-0"> {{ problem_title }} </v-card-title>
      <v-card-text>
        <v-text-field
          v-model="userInput"
          :rules="[(v) => !!v || '请输入答案']"
          :success-messages="this.result === 'true' ? ['回答正确'] : []"
          :error-messages="this.result === 'false' ? ['回答错误'] : []"
          validate-on-blur
          label="输入你的答案"
          hide-details="auto"
        ></v-text-field>
      </v-card-text>
      <v-card-actions>
        <v-btn :color="primary" text :disabled="this.userInput === ''" @click="submit">提交</v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script>
export default {
  name: 'ExerciseInput',
  props: {
    problem_id: {
      type: Number,
      required: true
    },
    problem_title: {
      type: String,
      required: true
    },
    result: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      userInput: ''
    };
  },
  methods: {
    submit() {
      this.$emit('submit', this.userInput);
    }
  }
};
</script>
