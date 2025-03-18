<!--添加题目组件-->
<template>
  <v-dialog v-model="dialogVisible" persistent width="1024">
    <v-app-bar color="cyan">
      <v-app-bar-title>编辑题目</v-app-bar-title>
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
      <v-app>
        <v-sheet class="pa-5">
          <v-form ref="form">
            <v-text-field
              v-model="exercise_title"
              :counter="20"
              :rules="nameRules"
              label="习题标题"
              required
            ></v-text-field>
            <v-card-subtitle>题目内容</v-card-subtitle>
            <!--<p>{{ exercise_content }}</p>-->
            <MarkdownInput ref="myEditor" />

            <v-radio-group label="选择题目类型" v-model="exercise_type" required dense>
              <v-radio
                v-for="choice in typeChoices"
                :key="choice.index"
                :value="choice.name"
                :label="choice.label"
              ></v-radio>
            </v-radio-group>

            <!--            <v-container id="choices"></v-container>-->

            <div v-show="exercise_type === 'fill_blank'">
              <v-text-field
                class="t1"
                label="请输入正确的答案"
                v-model="exercise_answer"
                :counter="100"
                :rules="textRules"
              ></v-text-field>
            </div>

            <div v-show="exercise_type === 'single_choice'">
              <v-card-subtitle class="t2">请输入正确选项后提交</v-card-subtitle>
              <v-radio-group v-model="singleSelectAns">
                <div v-for="choice in radio_choices" :key="choice.index" style="display: inline-flex">
                  <v-radio :value="choice.index"></v-radio>
                  <v-text-field
                    class="t2"
                    label="输入选项内容"
                    v-model="choice.content"
                    :counter="100"
                    :rules="textRules"
                  ></v-text-field>
                </div>
              </v-radio-group>
              <div style="display: inline-flex">
                <v-btn color="green" outlined block @click="addRadio" :disabled="sIndex > 5">增加一个单选项</v-btn>
                <v-btn color="red" outlined block @click="delRadio" :disabled="sIndex < 1">删除一个单选项</v-btn>
              </div>
            </div>

            <div v-show="exercise_type === 'multi_choice'">
              <v-card-subtitle class="t3">请输入正确选项后提交</v-card-subtitle>
              <div class="d-flex pa-4 t3" :key="choice.index" v-for="choice in checkbox_choices">
                <input type="checkbox" v-model="choice.choose" class="pe-2 t3" />
                <v-text-field
                  class="t3"
                  label="输入选项内容"
                  v-model="choice.content"
                  :counter="100"
                  :rules="textRules"
                ></v-text-field>
              </div>
              <div style="display: inline-flex">
                <v-btn color="green" outlined block @click="addCheckBox" :disabled="mIndex > 5">增加一个多选项</v-btn>
                <v-btn color="red" outlined block @click="delCheckBox" :disabled="mIndex < 1">删除一个多选项</v-btn>
              </div>
            </div>
            <v-container fluid>
              <v-card-subtitle>题目标签</v-card-subtitle>
              <v-row>
                <v-col cols="12" sm="4" md="4" v-for="(choice, index) in tagslist" :key="index">
                  <v-checkbox :value="choice" :label="choice" v-model="exercise_tags"></v-checkbox>
                </v-col>
              </v-row>
            </v-container>
            <!--表单处理按钮-->
            <div class="d-flex flex-column">
              <v-btn color="success" class="mt-4" block @click="validate" :disabled="submitDisabled()"> 提交</v-btn>
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
import { getExerciseDetail } from '@/api/exercise';

export default {
  name: 'NewExercise',
  props: {
    content: {},
    exercise_id: {
      type: Number,
      required: false,
      default: -1
    },
    show: {
      type: Boolean,
      required: false,
      default: false
    }
  },
  computed: {
    dialogVisible() {
      return this.show;
    }
  },
  components: {
    MarkdownInput
  },
  data() {
    return {
      singleSelectAns: -1,
      valid: true,
      includeFiles: true,
      exercise_title: '',
      exercise_type: '',
      exercise_answer: '',
      exercise_tags: [],
      exercise_choices: [],
      checkbox_choices: [],
      radio_choices: [],
      exercise_data: {
        choices: {}
      },
      radioAnswer: '',
      tagslist: ['排序', '图', '树', '分治', '基础知识', '循环和分支', '模拟'],
      typeChoices: [
        {
          index: 1,
          label: '填空',
          name: 'fill_blank'
        },
        {
          index: 2,
          label: '单选',
          name: 'single_choice'
        },
        {
          index: 3,
          label: '多选',
          name: 'multi_choice'
        }
      ],
      nameRules: [(v) => !!v || '必须填写题目名称', (v) => (v && v.length <= 20) || '题目名不能超过20个字符'],
      textRules: [(v) => !!v || '建议此字段不要为空', (v) => (v && v.length <= 100) || '建议此字段长度不要超过50'],
      sIndex: 0,
      mIndex: 0
    };
  },
  methods: {
    cancel() {
      this.$emit('cancel');
    },
    choicesToArrStr() {
      let newChoices = [];
      for (let i = 0; i < this.exercise_choices.length; i++) {
        newChoices.push(this.exercise_choices[i].content);
      }
      return newChoices;
    },
    arrStrToChoices(choicesArr, type) {
      if (type === 'single_choice') {
        this.radio_choices = [];
        for (let i = 0; i < choicesArr.length; i++) {
          this.radio_choices.push({
            index: i,
            content: choicesArr[i],
            choose: false
          });
        }
        this.sIndex = choicesArr.length;
      } else if (type === 'multi_choice') {
        this.checkbox_choices = [];
        for (let i = 0; i < choicesArr.length; i++) {
          this.checkbox_choices.push({
            index: i,
            content: choicesArr[i],
            choose: false
          });
        }
        this.mIndex = choicesArr.length;
      }
    },
    transAnswer2Str() {
      if (this.exercise_type === 'single_choice') {
        // console.log(this.singleSelectAns);
        if (this.singleSelectAns !== null && this.singleSelectAns !== -1) {
          try {
            return this.singleSelectAns.toString();
          } catch (e) {
            return '';
          }
        }
        return '';
      } else if (this.exercise_type === 'multi_choice') {
        let answer = '';
        for (let i = 0; i < this.checkbox_choices.length; i++) {
          if (this.checkbox_choices[i].choose) {
            answer += i.toString() + ';';
          }
        }
        return answer;
      } else {
        if (this.exercise_answer) return this.exercise_answer;
        else return '';
      }
    },
    transAnswerStr2Choice(answerStr, type) {
      if (type === 'single_choice') {
        this.singleSelectAns = Number(answerStr);
        this.radio_choices[Number(answerStr)].choose = true;
      } else if (type === 'multi_choice') {
        answerStr.split(';').forEach((value) => {
          if (Number(value)) {
            this.checkbox_choices[Number(value)].choose = true;
          }
        });
      } else {
        this.exercise_answer = answerStr;
      }
    },
    async validate() {
      const valid = true; //await this.$refs.form.validate();
      if (this.exercise_type === 'single_choice') {
        this.exercise_choices = this.radio_choices;
      } else if (this.exercise_type === 'multi_choice') {
        this.exercise_choices = this.checkbox_choices;
      } else {
        this.exercise_choices = [];
      }
      if (valid) {
        this.exercise_answer = this.transAnswer2Str();
        this.exercise_data.choices = this.exercise_choices;
        let newChoices = this.choicesToArrStr();
        let data = {
          title: this.exercise_title,
          content: this.$refs.myEditor.content,
          answer: this.exercise_answer,
          tags: this.exercise_tags,
          course_id: this.$store.state.user.curCourseId,
          type: this.exercise_type,
          data: { choices: newChoices } //this.exercise_data
        };
        this.$emit('handleExerciseEdit', this.exercise_id, data);
      }
    },
    submitDisabled() {
      if (!(this.exercise_title && this.exercise_title.length > 0 && this.exercise_title.length < 20)) {
        return true;
      } else if (!(this.$refs.myEditor && this.$refs.myEditor.content.length > 0)) {
        return true;
      } else if (!(this.transAnswer2Str().length > 0)) {
        return true;
      }
      return false;
    },
    reset() {
      this.$refs.form.reset();
      this.$refs.myEditor.content = '';
      // this.exercise_content = '';
    },
    resetValidation() {
      this.$refs.form.resetValidation();
    },
    addRadio() {
      this.radio_choices.push({
        index: this.sIndex++,
        content: '',
        choose: false
      });
    },
    delRadio() {
      this.radio_choices.pop();
      this.sIndex = this.sIndex - 1 > 0 ? this.sIndex - 1 : 0;
    },
    delCheckBox() {
      this.checkbox_choices.pop();
      this.mIndex = this.mIndex - 1 > 0 ? this.mIndex - 1 : 0;
    },
    addCheckBox() {
      this.checkbox_choices.push({
        index: this.mIndex++,
        content: '',
        choose: false
      });
    },
    createDialog(isEdit, exerciseID) {
      if (isEdit) {
        getExerciseDetail(exerciseID).then((re) => {
          this.exercise_title = re.data.title;
          this.$refs.myEditor.content = re.data.content;
          this.exercise_type = re.data.type;
          // this.exercise_answer = re.data.answer;
          this.exercise_tags = re.data.tags;
          this.exercise_choices = re.data.data.choices;
          this.arrStrToChoices(re.data.data.choices, re.data.type);
          this.transAnswerStr2Choice(re.data.answer, re.data.type);
        });
      }
    }
  },
  watch: {
    exercise_type: {
      handler(item1) {
        // item1为新值
        Array.from(document.getElementsByClassName('t1')).map((item) => (item.style.display = 'none'));
        Array.from(document.getElementsByClassName('t2')).map((item) => (item.style.display = 'none'));
        Array.from(document.getElementsByClassName('t3')).map((item) => (item.style.display = 'none'));
        if (item1 === 'fill_blank') {
          Array.from(document.getElementsByClassName('t1')).map((item) => (item.style.display = 'block'));
        } else if (item1 === 'single_choice') {
          Array.from(document.getElementsByClassName('t2')).map((item) => (item.style.display = 'block'));
        } else if (item1 === 'multi_choice') {
          Array.from(document.getElementsByClassName('t3')).map((item) => (item.style.display = 'block'));
        }
      }
    }
  },
  mounted() {
    Array.from(document.getElementsByClassName('t1')).map((item) => (item.style.display = 'none'));
    Array.from(document.getElementsByClassName('t2')).map((item) => (item.style.display = 'none'));
    Array.from(document.getElementsByClassName('t3')).map((item) => (item.style.display = 'none'));
  }
};
</script>

<style scoped>
.t1,
.t2,
.t3 {
  display: block;
}
</style>
