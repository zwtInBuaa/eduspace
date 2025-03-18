<template>
  <v-container>
    <v-card class="mx-auto">
      <v-card-text>
        <v-row class="text-center">
          <!-- 图 SVG 区 -->
          <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 6" class="align-self-center">
            <binary-search-graphic ref="binarySearchGraphicRef" style="margin-bottom: 20px"></binary-search-graphic>
            <v-col class="text-center">
              <v-btn text color="primary" @click="codeSubmit">提交代码</v-btn>
              <v-btn text color="success" @click="graphNextStep">单步执行</v-btn>
              <v-btn text color="success" @click="runCode">全部执行</v-btn>
              <v-btn text color="error" @click="resetStep">重新执行</v-btn>
            </v-col>
          </v-col>
          <!-- 按钮、代码块区 -->
          <v-col>
            <variable-table ref="varTable"></variable-table>
            <code-tracer
              ref="codeTracerRef"
              :code="code"
              :input-default="inputDefault"
              @input="updateCode"
            ></code-tracer>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
    <v-overlay :value="loading">
      <v-progress-circular :size="70" :width="7" color="purple" indeterminate></v-progress-circular>
    </v-overlay>
  </v-container>
</template>

<script>
import CodeTracer from '@/components/visualization/CodeTracer.vue';
import BinarySearchGraphic from '@/components/visualization/BinarySearchGraphic.vue';
// import axios from "axios";
import { postRequest } from '@/api/request';
import VariableTable from '@/components/visualization/VariableTable.vue';
export default {
  components: {
    VariableTable,
    CodeTracer,
    BinarySearchGraphic
  },
  data() {
    return {
      loading: false,
      code: `def binarySearch(array, target):
    left = 0
    right = len(array) - 1

    while ___:
        mid = (left + right) // 2
        if array[mid] == target:
            return mid
        elif array[mid] < target:
            left = ___
        else:
            right = ___
    return -1

array = [___]
target = ___
result = binarySearch(array, target)
print(result)`,
      inputDefault: ['left <= right', 'mid + 1', 'mid - 1', '1,2,3,4,5,7,8,9,10,11,12,13,14', '13'],
      highLightSteps: ['0', '1', '2', '3', '4', '5', '6', '7', '8']
    };
  },
  mounted() {
    this.$store.commit('userGuide/open', this.$route.path);
    this.executeBinarySearch(this.inputDefault);
    this.$refs.codeTracerRef.steps = this.highLightSteps;
  },
  methods: {
    updateCode(newCode) {
      this.code = newCode;
    },
    codeSubmit() {
      this.executeBinarySearch(this.$refs.codeTracerRef.codeSubmit());
    },
    graphNextStep() {
      this.$refs.binarySearchGraphicRef.nextStep();
      //this.$refs.codeTracerRef.highlightActiveLines();
      this.$refs.codeTracerRef.nextStep();
      this.$refs.varTable.nextStep();
    },
    runCode() {
      this.$refs.codeTracerRef.steps = this.highLightSteps;
      this.$refs.binarySearchGraphicRef.runSteps();
      this.$refs.codeTracerRef.runSteps(200);
      this.$refs.varTable.runSteps(200);
    },
    resetStep() {
      this.$refs.binarySearchGraphicRef.resetStep();
      this.$refs.codeTracerRef.resetStep();
      this.$refs.varTable.resetStep();
    },
    async executeBinarySearch(codeBlanks) {
      let res;
      this.loading = true;
      await postRequest('/visualizations/binarySearch/submitCode', { codeBlanks: codeBlanks })
        .then((response) => {
          res = response.data;
          this.loading = false;
        })
        .catch((error) => {
          this.$store.dispatch('snackbar/error', error);
          this.loading = false;
          res = 'error';
        });
      if (res === 'error') return;
      let jsonData = res; //JSON.parse(res)
      let new_steps = [],
        barNum = jsonData['barNum'],
        stepNum = jsonData['steps'].length;
      let barValue = jsonData['barValue'],
        barHeight = jsonData['barHeight'],
        color_info = jsonData['steps'];
      for (let i = 0; i < stepNum; i++) {
        let cur_step_info = [];
        for (let j = 0; j < barNum; j++) {
          cur_step_info.push({ height: barHeight[j], color: color_info[i][j], value: barValue[j] });
        }
        new_steps.push(cur_step_info);
      }
      if (res['highlightSteps']) {
        this.highLightSteps = res['highlightSteps'];
        this.$refs.codeTracerRef.steps = res['highlightSteps'];
      }
      if (res['varSteps']) {
        this.$refs.varTable.setSteps(res['varSteps']);
      }
      this.$refs.binarySearchGraphicRef.setSteps(new_steps);
      // this.$store.commit('userGuide/close');
    }
  }
};
</script>
