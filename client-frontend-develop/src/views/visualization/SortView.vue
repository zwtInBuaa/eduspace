<template>
  <v-container>
    <v-card class="mx-auto">
      <v-card-text>
        <v-row class="text-center">
          <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 6" class="align-self-center">
            <!-- 图 SVG 区 -->
            <sort-graphic ref="sort_grpah" @running="setDisableButtons"></sort-graphic>
            <v-col class="text-center">
              <v-btn text color="primary" :disabled="buttonDisabled" @click="codeSubmit">提交代码</v-btn>
              <v-btn text color="success" :disabled="buttonDisabled" @click="graphNextStep">单步执行</v-btn>
              <v-btn text color="success" :disabled="buttonDisabled" @click="runCode">全部执行</v-btn>
              <v-btn text color="error" @click="resetStep">重新执行</v-btn>
            </v-col>
          </v-col>
          <!-- 按钮、代码块区 -->
          <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 6">
            <v-select v-model="selectedSort" :items="sortOptions" label="选择排序方式"></v-select>
            <!-- <variable-table ref="varTable"></variable-table> -->
            <code-tracer ref="code_tracer" :code="code" :input-default="inputDefault"></code-tracer>
          </v-col>
        </v-row>

        <v-row> </v-row>
      </v-card-text>
    </v-card>
    <v-overlay :value="loading">
      <v-progress-circular :size="70" :width="7" color="purple" indeterminate></v-progress-circular>
    </v-overlay>
  </v-container>
</template>

<script>
import sortGraphic from '@/components/visualization/SortGraphic.vue';
import codeTracer from '@/components/visualization/CodeTracer.vue';
import { getRequest, postRequest } from '@/api/request';
// import variableTable from '@/components/visualization/VariableTable.vue';

export default {
  name: 'SortView',
  components: {
    sortGraphic,
    codeTracer
    // variableTable
  },
  data() {
    return {
      loading: false,
      buttonDisabled: false,
      selectedSort: null,
      sortOptions: ['冒泡排序', '选择排序', '快速排序'],
      code: '11111111',
      inputDefault: [],
      choice: 'bubbleSort',
      dataReady: false,
      highLightSteps: ['0', '1', '2', '3', '4', '5', '6', '7', '8'],
      timeoutIDs: []
    };
  },
  mounted() {
    // set code
    // set inputDefault
    this.$store.commit('userGuide/open', this.$route.path);
    // const stid = setTimeout(() => {
    //   this.$store.commit('userGuide/close');
    // }, 5000);
    //this.timeoutIDs.push(stid);
    this.choice = 'bubbleSort';
    if (this.$route.params.defaultAlgo) {
      this.choice = this.$route.params.defaultAlgo;
    }
    this.loadSortInfo(this.choice);
    this.$refs.code_tracer.steps = this.highLightSteps;
  },
  beforeDestroy() {
    for (let i = 0; i < this.timeoutIDs.length; i++) {
      clearTimeout(this.timeoutIDs[i]);
    }
  },
  methods: {
    async codeSubmit() {
      this.loading = true;
      await postRequest(`/visualizations/sort/submitCode/${this.choice}`, {
        codeBlanks: this.$refs.code_tracer.codeSubmit()
      })
        .then((response) => {
          let res = response.data;
          this.$refs.sort_grpah.setValues(res['heights'], res['values']);
          this.$refs.sort_grpah.setSteps(res['steps'], res['detail']);
          if (res['highlightSteps']) {
            this.highLightSteps = res['highlightSteps'];
            this.$refs.code_tracer.steps = res['highlightSteps'];
          }
          // if (res['varSteps']) {
          //   this.$refs.varTable.steps = res['varSteps'];
          // }
        })
        .catch((error) => {
          this.$store.dispatch('snackbar/error', error);
        });
      this.loading = false;
    },
    graphNextStep() {
      if (!this.buttonDisabled) {
        this.$refs.sort_grpah.nextStep();
        this.$refs.code_tracer.nextStep();
        // this.$refs.varTable.nextStep();
      }
    },
    loadSortInfo(type) {
      let req_url = `visualizations/sort/${type}`;
      // 发送get请求
      this.loading = true;
      getRequest(req_url)
        .then((response) => {
          let res = response.data;
          this.code = response.data['code'];
          this.inputDefault = response.data['inputDefault'];
          this.dataReady = true;
          if (res['heights']) {
            this.$refs.sort_grpah.setValues(res['heights'], res['values']);
            this.$refs.sort_grpah.setSteps(res['steps'], res['detail']);
          }
          if (res['highlightSteps']) {
            this.highLightSteps = res['highlightSteps'];
            this.$refs.code_tracer.steps = res['highlightSteps'];
          }
          // if (res['varSteps']) {
          //   this.$refs.varTable.setSteps(res['varSteps']);
          // }
          this.$refs.code_tracer.setCode(this.code, this.inputDefault);
          this.loading = false;
        })
        .catch((error) => {
          this.$store.dispatch('snackbar/error', error);
          this.loading = false;
        });
    },
    runCode() {
      if (!this.buttonDisabled) {
        this.$refs.sort_grpah.runSteps();
        this.$refs.code_tracer.runSteps(200);
        // this.$refs.varTable.runSteps(200);
      }
    },
    resetStep() {
      this.$refs.sort_grpah.resetStep();
      this.$refs.code_tracer.resetStep();
      // this.$refs.varTable.resetStep();
    },
    setDisableButtons(b) {
      this.buttonDisabled = b;
    }
  },
  watch: {
    async selectedSort(newValue) {
      let type = '';
      if (newValue === '冒泡排序') {
        type = 'bubbleSort';
      } else if (newValue === '选择排序') {
        type = 'selectionSort';
      } else if (newValue === '插入排序') {
        type = 'insertionSort'; // 没有插入排序！
      } else if (newValue === '快速排序') {
        type = 'quickSort';
      }
      if (type !== this.choice) {
        this.choice = type;
        this.dataReady = false;
        this.$refs.sort_grpah.endCurrentRun();
        await this.loadSortInfo(type);
      }
    }
  }
};
</script>

<style scoped></style>
