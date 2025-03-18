<template>
  <v-container>
    <v-card flat>
      <v-tabs v-model="tabs" centered grow>
        <v-tab>浏览习题</v-tab>
        <v-tab>浏览题组</v-tab>
        <v-tab v-if="editable()"> 管理习题</v-tab>
        <v-tab v-if="editable()"> 管理题组</v-tab>
      </v-tabs>
      <v-tabs-items v-model="tabs">
        <!-- 浏览习题标签页 -->
        <v-tab-item>
          <exercise-list :exercise-list="exerciseList" @clickRow="clickExercise"></exercise-list>
        </v-tab-item>
        <!-- 浏览题组标签页 -->
        <v-tab-item>
          <exercise-set-list :exercise-set-list="exerciseSetList" @clickRow="clickExam"></exercise-set-list>
        </v-tab-item>
        <v-tab-item v-if="editable()">
          <editable-exercise-list
            :exercise-list="exerciseList"
            @deleteExerciseItemConfirm="deleteExerciseItemConfirm"
            @invokeExerciseEdit="invokeExerciseEdit"
          ></editable-exercise-list>
          <exercise-editor
            :show="showExerciseEditor"
            :content="editExercise"
            :exercise_id="exercise_id"
            @handleExerciseEdit="handleExerciseEdit"
            @cancel="cancelExerciseEdit"
            ref="exerciseEditorRef"
          />
        </v-tab-item>
        <v-tab-item v-if="editable()">
          <editable-exam-list
            :exercise-set-list="exerciseSetList"
            @deleteExamItemConfirm="deleteExamItemConfirm"
            @invokeExamEdit="invokeExamEdit"
          ></editable-exam-list>
          <exam-editor
            :show="showExamEditor"
            :content="editExam"
            :exam_id="exam_id"
            @submit="handleExamEdit"
            @cancel="cancelExamEdit"
            ref="examEditorRef"
          />
        </v-tab-item>
      </v-tabs-items>
    </v-card>
  </v-container>
</template>

<script>
import ExerciseList from '@/components/exercise/ExerciseList.vue';
import ExerciseSetList from '@/components/exercise/ExamList.vue';
import EditableExerciseList from '@/components/exercise/EditableExerciseList.vue';
import EditableExamList from '@/components/exercise/EditableExamList.vue';
import ExerciseEditor from '@/components/exercise/ExerciseEditor.vue';
import ExamEditor from '@/components/exercise/ExamEditor.vue';
import {
  createNewExercise,
  createNewExerciseSet,
  deleteExercise,
  deleteExerciseSet,
  editExercise,
  editExerciseSet,
  getAllExercise,
  getAllExerciseSet
  // getExerciseDetail,
  // getExerciseSetDetail
} from '@/api/exercise';

export default {
  components: {
    ExerciseList,
    ExerciseSetList,
    EditableExerciseList,
    EditableExamList,
    ExerciseEditor,
    ExamEditor
  },
  data() {
    return {
      dialogVisible: false,
      dialogSetVisible: false,
      tabs: null,
      exerciseList: [],
      exerciseSetList: [],
      showExamEditor: false,
      showExerciseEditor: false,
      exercise_id: -1,
      editExercise: {
        title: '',
        content: '',
        type: '',
        answer: '',
        tags: '',
        choices: ''
      },
      exam_id: -1,
      editExam: {
        title: '',
        description: '',
        questions: ''
      }
    };
  },
  computed: {},
  methods: {
    clickExercise(item) {
      this.$router.push('/exercise/' + item.id);
    },
    clickExam(item) {
      this.$router.push('/exerciseSet/' + item.id);
    },
    editable() {
      if (this.$store.state.user.role === '管理员' && !this.$vuetify.breakpoint.mobile) {
        return true;
      } else if (this.$store.state.user.role === '老师' && !this.$vuetify.breakpoint.mobile) {
        return true;
      } else return false;
    },
    updateExercises() {
      getAllExercise().then((re) => {
        // // console.log('获得所有习题：', re.data);
        this.exerciseList = re.data;
      });
    },
    updateExams() {
      getAllExerciseSet().then((re) => {
        // // console.log('获得所有题组', re.data);
        this.exerciseSetList = re.data;
      });
    },
    async deleteExerciseItemConfirm(id) {
      await deleteExercise(id).then(() =>
        // this.$store.commit('snackbar/open', {
        //   msg: '删除成功',
        //   color: 'success'
        // })
        this.$store.dispatch('snackbar/success', '删除题目成功')
      );
      this.updateExercises();
    },
    async deleteExamItemConfirm(id) {
      await deleteExerciseSet(id).then(() =>
        // this.$store.commit('snackbar/open', {
        //   msg: '删除成功',
        //   color: 'success'
        // })
        this.$store.dispatch('snackbar/success', '删除题组成功')
      );
      this.updateExams();
    },
    /*Exercise*/
    async invokeExerciseEdit(exercise_id) {
      this.exercise_id = exercise_id;
      // if (exercise_id !== -1) {
      //   await getExerciseDetail(exercise_id).then((re) => {
      //     this.editExercise.title = re.data.title;
      //     this.editExercise.content = re.data.content;
      //     this.editExercise.type = re.data.type;
      //     this.editExercise.answer = re.data.answer;
      //     this.editExercise.tags = re.data.tags;
      //     this.editExercise.choices = re.data.choices;
      //   });
      // } else {
      //   this.editExercise.title = '';
      //   this.editExercise.content = '';
      //   this.editExercise.type = '';
      //   this.editExercise.answer = '';
      //   this.editExercise.tags = '';
      //   this.editExercise.choices = '';
      // }
      await this.$refs.exerciseEditorRef.createDialog(exercise_id !== -1, exercise_id);
      this.showExerciseEditor = true;
    },
    async handleExerciseEdit(id, exercise) {
      if (id === -1) {
        await createNewExercise(exercise).then(() => this.$store.dispatch('snackbar/success', '创建题目成功'));
      } else {
        await editExercise(id, exercise).then(() => this.$store.dispatch('snackbar/success', '修改题目成功'));
      }
      this.showExerciseEditor = false;
      this.updateExercises();
    },
    cancelExerciseEdit() {
      this.showExerciseEditor = false;
    },
    /*Exam*/
    async invokeExamEdit(exam_id) {
      this.exam_id = exam_id;
      // if (exam_id === -1) {
      //   this.editExam = {
      //     title: '',
      //     description: '',
      //     questions: ''
      //   };
      // } else {
      //   await getExerciseSetDetail(exam_id).then((re) => {
      //     this.editExam.title = re.data.title;
      //     this.editExam.description = re.data.description;
      //     this.editExam.questions = re.data.questions;
      //   });
      // }
      await this.$refs.examEditorRef.createDialog(exam_id !== -1, exam_id);
      this.showExamEditor = true;
    },
    async handleExamEdit(exam_id, data) {
      if (exam_id === -1) {
        await createNewExerciseSet(data).then(() => this.$store.dispatch('snackbar/success', '创建题组成功'));
      } else {
        await editExerciseSet(exam_id, data).then(() => this.$store.dispatch('snackbar/success', '修改题组成功'));
      }
      this.showExamEditor = false;
      this.updateExams();
    },
    cancelExamEdit() {
      this.showExamEditor = false;
    }
  },

  created() {
    // // console.log(this.$store.state.user.role);
    getAllExercise().then((re) => {
      // // console.log('获得所有习题：', re.data);
      this.exerciseList = re.data;
    });
    getAllExerciseSet().then((re) => {
      // // console.log('获得所有题组', re.data);
      this.exerciseSetList = re.data;
    });
  }
};
</script>
