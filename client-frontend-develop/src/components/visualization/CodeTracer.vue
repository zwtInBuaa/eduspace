<template>
  <!--  <div class="code-block-wrapper">-->
  <v-card>
    <v-card-text>
      <v-container class="code-block" style="overflow-x: auto; overflow-y: auto; max-height: 600px">
        <v-row v-for="(line, index) in lines" :key="index" :ref="'line-' + index" :class="lineClass(index)" no-gutters>
          <div v-for="(part, partIndex) in line.parts" :key="partIndex">
            <span v-if="part.editable" class="code-block-line-editable">
              <custom-input-container
                v-model="inputValues[index]"
                @input="updateValue(index, $event)"
                :width="(7 * inputValues[index].length + 5).toString()"
                :ref="'input-' + index"
              ></custom-input-container>
            </span>
            <span v-else>{{ part.text }}</span>
          </div>
        </v-row>
      </v-container>
    </v-card-text>
  </v-card>
  <!--  </div>-->
</template>

<script>
import customInputContainer from '@/components/visualization/CustomInputContainer.vue';

export default {
  computed: {},
  components: { customInputContainer },
  props: {
    code: {
      type: String,
      required: true
    },
    inputDefault: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      lines: [],
      inputValues: ['', '', ''],
      activeLine: [0, 0, 0],
      steps: [[0], [1], [0, 1], [0], [], [1], []],
      curStep: 0,
      isEnding: false,
      timeoutIDs: []
    };
  },
  mounted() {
    this.lines = this.code.split('\n').map((line) => {
      const parts = line.split(/(_+)/).map((part) => ({
        text: part.startsWith('_') ? '' : part,
        editable: part.startsWith('_')
      }));
      return { parts };
    });
    for (let i = 0, j = 0; i < this.lines.length; i++) {
      this.inputValues.push('');
      if (this.lines[i].parts.some((part) => part.editable)) {
        this.inputValues[i] = this.inputDefault[j++];
      }
    }
  },
  methods: {
    codeSubmit() {
      let res = [];
      for (let i = 0; i < this.inputValues.length; i++) {
        if (this.inputValues[i] !== '') {
          let input = this.$refs['input-' + i.toString()][0];
          if (input && input._data.isSafe) {
            res.push(this.inputValues[i]);
          }
        }
      }
      return res;
    },
    updateValue(index, value) {
      this.$set(this.inputValues, index, value);
    },
    lineClass(index) {
      return {
        'code-block-line-editable': this.lines[index].parts.some((part) => part.editable)
      };
    },
    highlightActiveLines() {
      for (let i = 0; i < this.lines.length; i++) {
        let t = false;
        for (let j = 0; j < this.activeLine.length; j++) {
          if (i === this.activeLine[j]) {
            t = true;
            break;
          }
        }
        if (t) {
          const selectedLine = this.$refs['line-' + i][0];
          selectedLine.style.transaction = 'background-color 0.3s ease';
          selectedLine.style.backgroundColor = 'yellow'; // 将背景颜色修改为黄色
        } else {
          const selectedLine = this.$refs['line-' + i][0];
          selectedLine.style.transaction = 'background-color 0.3s ease';
          selectedLine.style.backgroundColor = '#f5f5f5';
        }
      }
      // this.activeLine[0] = (this.activeLine[0] + 1) % this.lines.length;
    },
    nextStep() {
      if (this.curStep < this.steps.length) {
        let highlights = this.steps[this.curStep];
        for (let i = 0; i < this.lines.length; i++) {
          if (highlights.includes(i)) {
            const selectedLine = this.$refs['line-' + i][0];
            selectedLine.style.transaction = 'background-color 0.3s ease';
            selectedLine.style.backgroundColor = 'yellow'; // 将背景颜色修改为黄色
          } else {
            const selectedLine = this.$refs['line-' + i][0];
            selectedLine.style.transaction = 'background-color 0.3s ease';
            selectedLine.style.backgroundColor = '#f5f5f5';
          }
        }
        this.curStep++;
      } else {
        this.curStep = 0;
      }
    },
    runSteps(tInterval) {
      if (this.curStep < this.steps.length) {
        if (!this.isEnding) {
          this.nextStep();
          const stid = setTimeout(() => {
            this.runSteps(tInterval);
          }, tInterval);
          this.timeoutIDs.push(stid);
        }
      }
      // else {
      //   const stid = setTimeout(() => {
      //     this.isRunning = false;
      //     this.buttonDisabled = false;
      //   }, 100);
      //   this.timeoutIDs.push(stid);
      // }
    },
    resetStep() {
      this.isEnding = true;
      for (let i = 0; i < this.timeoutIDs.length; i++) {
        clearTimeout(this.timeoutIDs[i]);
      }
      for (let i = 0; i < this.lines.length; i++) {
        const selectedLine = this.$refs['line-' + i][0];
        selectedLine.style.transaction = 'background-color 0.3s ease';
        selectedLine.style.backgroundColor = '#f5f5f5';
      }
      this.timeoutIDs = [];
      const stid = setTimeout(() => {
        this.curStep = 0;
        this.isEnding = false;
      }, 300);
      this.timeoutIDs.push(stid);
    },
    setCode(newCode, newInputDefault) {
      this.lines = newCode.split('\n').map((line) => {
        const parts = line.split(/(_+)/).map((part) => ({
          text: part.startsWith('_') ? '' : part,
          editable: part.startsWith('_')
        }));
        return { parts };
      });
      this.inputValues = [];
      for (let i = 0, j = 0; i < this.lines.length; i++) {
        this.inputValues.push('');
        if (this.lines[i].parts.some((part) => part.editable)) {
          this.inputValues[i] = newInputDefault[j++];
        }
      }
    },
    setCodeNoEdit(newCode) {
      this.lines = newCode.split('\n').map((line) => {
        const parts = [line].map((part) => ({
          text: part,
          editable: false
        }));
        return { parts };
      });
    }
  },
  beforeDestroy() {
    for (let i = 0; i < this.timeoutIDs.length; i++) {
      clearTimeout(this.timeoutIDs[i]);
    }
  }
};
</script>

<style>
.code-block {
  background-color: #f5f5f5;
  font-family: Consolas, Menlo, Courier, monospace;
  padding: 10px;
  line-height: 20px;
  font-size: 16px;
  white-space: pre;
  text-align: left;
}

.code-block-line-editable {
  /*background-color: #fff;*/
  display: flex;
  flex-wrap: nowrap;
  overflow: hidden;
}

.code-block-wrapper {
  overflow-x: scroll;
}

/*.code-block-line-editable-part {*/
/*    !*background-color: white;*!*/
/*    white-space: nowrap;*/
/*    !*width: 60px;*!*/
/*    !*height: 10px;*!*/
/*    !*padding: -10px;*!*/
/*    height:30px;*/
/*}*/

.code-block span {
  display: inline-block;
}
</style>
