<template>
  <v-container class="pa-0">
    <v-card>
      <v-card-title class="pb-0"> {{ problem_title }} </v-card-title>
      <v-container>
        <v-form class="mx-auto">
          <v-radio-group v-model="userChoice" required dense>
            <v-radio
              v-for="(choice, index) in choices"
              :key="index"
              :value="index"
              :label="choice.toString()"
              :color="result === 'true' ? 'success' : result === 'false' ? 'error' : 'primary'"
            ></v-radio>
          </v-radio-group>
          <v-btn color="primary" text @click="submit">提交</v-btn>
        </v-form>
      </v-container>
    </v-card>
  </v-container>
</template>

<script>
export default {
  name: 'ChoiceExercise',
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
      userChoice: -1
    };
  },
  methods: {
    submit() {
      if (this.userChoice === -1) {
        alert('请至少选择一个选项');
      } else {
        this.$emit('submit', this.userChoice + '');
      }
    }
  }
};
</script>
