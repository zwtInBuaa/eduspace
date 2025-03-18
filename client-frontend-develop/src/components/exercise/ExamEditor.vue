<!--添加题组组件-->
<template>
  <v-dialog v-model="dialogSetVisible" persistent width="1024">
    <v-app-bar color="cyan">
      <v-app-bar-title>编辑题组</v-app-bar-title>
      <v-spacer></v-spacer>
      <v-btn
        icon
        dark
        @click="
          () => {
            this.reset();
            this.cancel();
          }
        "
      >
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-app-bar>
    <v-card>
      <!--      <v-btn-->
      <!--        v-show="exam_id !== -1"-->
      <!--        @click="-->
      <!--          () => {-->
      <!--            this.createDialog();-->
      <!--          }-->
      <!--        "-->
      <!--      >-->
      <!--        获取原题组信息-->
      <!--      </v-btn>-->
      <v-app>
        <v-sheet class="pa-5">
          <v-form ref="form">
            <v-text-field label="题组名" v-model="exam_title" :counter="20" :rules="nameRules"></v-text-field>
            <v-card-subtitle>题组内容</v-card-subtitle>
            <MarkdownInput ref="myEditor"></MarkdownInput>
            <v-card-subtitle>选择加入题组的题目</v-card-subtitle>
            <select-exercise-list
              :exercise-list="exerciseList"
              :row-on-click="rowOnClick3"
              ref="selectExerciseList"
            ></select-exercise-list>
            <div class="d-flex flex-column">
              <v-btn color="success" class="mt-4" block @click="validate" :disabled="submitDisabled()">提交</v-btn>
              <v-btn color="error" class="mt-4" block @click="reset"> 重置表单</v-btn>
            </div>
          </v-form>
        </v-sheet>
      </v-app>
    </v-card>
  </v-dialog>
</template>

<script>
import MarkdownInput from '../MarkdownInput.vue';
import SelectExerciseList from './SelectExerciseList.vue';
import { getExerciseSetDetail, getAllExercise } from '@/api/exercise';

export default {
  name: 'NewExam',
  components: {
    MarkdownInput,
    SelectExerciseList
  },
  props: {
    show: {
      type: Boolean,
      require: false,
      default: false
    },
    exam_id: {
      type: Number,
      require: false,
      default: -1
    },
    content: {}
  },
  computed: {
    dialogSetVisible() {
      return this.show;
    }
  },
  data() {
    return {
      exam_title: '',
      exam_content: '555',
      exercises: [],
      exerciseList: [
        {
          id: 1,
          title: '题目1',
          tags: ['tag1', 'tag2']
        }
      ],
      nameRules: [(v) => !!v || '必须填写题目名称', (v) => (v && v.length <= 20) || '题目名不能超过20个字符']
    };
  },
  watch: {
    // show(val) {
    //   if (val) {
    //     if (this.exam_id !== -1) {
    //       this.exam_title = this.content.title;
    //       this.$refs.myEditor.content = this.content.description;
    //       this.$refs.selectExerciseList.selected = this.content.questions;
    //     }
    //   }
    // }
  },
  methods: {
    rowOnClick3() {
      // this.$router.push('/exercise/' + item.id);
    },
    reset() {
      this.$refs.myEditor.content = '';
      this.exam_title = '';
      this.$refs.selectExerciseList.selected = [];
    },
    async validate() {
      const valid = await this.$refs.form.validate();
      if (valid) {
        let IdList = new Set();
        this.$refs.selectExerciseList.selected.forEach((e) => {
          IdList.add(e.id);
        });
        let data = {
          name: this.exam_title,
          description: this.$refs.myEditor.content,
          questions: Array.from(IdList)
        };
        this.$emit('submit', this.exam_id, data);
      }
    },
    createDialog(isEdit, examID) {
      if (isEdit) {
        getExerciseSetDetail(examID).then((re) => {
          this.exam_title = re.data.title;
          this.$refs.myEditor.content = re.data.description;
          this.$refs.selectExerciseList.selected = re.data.questions;
        });
      }
    },
    cancel() {
      this.$emit('cancel');
    },
    submitDisabled() {
      if (!(this.exam_title && this.exam_title.length > 0 && this.exam_title.length < 20)) {
        return true;
      } else {
        let IdList = new Set();
        this.$refs.selectExerciseList.selected.forEach((e) => {
          IdList.add(e.id);
        });
        if (IdList.size < 1) {
          return true;
        }
      }
      return false;
    }
  },

  created() {
    getAllExercise().then((re) => {
      this.exerciseList = re.data;
    });
  }
};
</script>
