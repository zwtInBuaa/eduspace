<template>
  <v-container>
    <v-card flat class="pt-2">
      <v-btn icon @click="$router.go(-1)">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <!-- Markdown Display -->
      <v-container>
        <v-card class="mb-4">
          <MarkdownDisplay :content="this.problem.problem_info" :minHeight="10" />
        </v-card>
        <!-- Commit Box -->
        <choice
          v-if="problem.problem_type === 'single_choice'"
          :problem_title="problem.problem_title"
          :problem_id="problem.problem_id"
          :choices="problem.choices"
          :result="this.result"
          @submit="submit"
          ref="singleChoice"
        ></choice>
        <exercise-input
          v-else-if="problem.problem_type === 'fill_blank'"
          :problem_id="problem.problem_id"
          :problem_title="problem.problem_title"
          :result="this.result"
          @submit="submit"
          ref="fillBlank"
        ></exercise-input>
        <exercise-checked-box
          v-else
          :problem_id="problem.problem_id"
          :problem_title="problem.problem_title"
          :choices="problem.choices"
          :result="this.result"
          @submit="submit"
          ref="multiChoice"
        ></exercise-checked-box>
      </v-container>
    </v-card>
  </v-container>
</template>

<script>
import MarkdownDisplay from '@/components/MarkdownDisplay.vue';
import Choice from '@/components/exercise/SingleChoice.vue';
import ExerciseInput from '@/components/exercise/FillBlank.vue';
import ExerciseCheckedBox from '@/components/exercise/MultiChoice.vue';
import { getExerciseDetail, submitexerciseAnswer } from '@/api/exercise';

export default {
  components: {
    Choice,
    MarkdownDisplay,
    ExerciseInput,
    ExerciseCheckedBox
  },
  data() {
    return {
      tabs: null,
      id: 999,
      problem: {
        problem_type: 'single_choice',
        problem_id: 1,
        problem_title: '题目1',
        problem_info: '吉大三架四级弹我和电话维护队还丢汇合I度海外海',
        tags: ['tag1', 'tag2'],
        choices: [
          {
            index: 1,
            label: '选项1'
          },
          {
            index: 2,
            label: '选项2'
          },
          {
            index: 3,
            label: '选项3'
          }
        ]
      },
      result: ''
    };
  },
  created() {
    this.id = this.$route.params.id;
    getExerciseDetail(this.id).then((re) => {
      // // console.log('题目详情：', re);
      this.problem.problem_id = re.data.id;
      this.problem.problem_title = re.data.title;
      this.problem.problem_info = re.data.content;
      this.problem.problem_type = re.data.type;
      this.problem.choices = re.data.data.choices;
    });
    // // console.log(this.problem);
  },
  computed: {},
  methods: {
    submit(answer) {
      submitexerciseAnswer(this.id, { answer }).then((re) => {
        if (re.data.result) {
          this.$store.dispatch('snackbar/success', '回答正确');
          this.result = 'true';
        } else {
          this.$store.dispatch('snackbar/error', '回答错误');
          this.result = 'false';
        }
      });
    }
  }
};
</script>
