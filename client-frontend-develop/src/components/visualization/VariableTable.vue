<template>
  <div style="overflow-x: auto">
    <table class="list-table">
      <thead>
        <tr class="table-title">
          <th>变量名</th>
          <th v-for="(item, i) in varNames" :key="i">{{ item }}</th>
        </tr>
      </thead>
      <tbody>
        <tr class="table-content">
          <th>变量值</th>
          <td v-for="(item, i) in values" :key="i" :class="{ flash: updateData[i] }">
            {{ item }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  data() {
    return {
      values: [-122, 'null', 'false'],
      varNames: ['a', 'b', 'c'],
      steps: [
        {
          varNames: ['a', 'b', 'c', 'd', 'e'],
          values: [1, 2, 3, 4, 5]
        },
        {
          varNames: ['a', 'b', 'c', 'd', 'e'],
          values: [11, 2, 3, 41, 51]
        },
        {
          varNames: ['a', 'b', 'c', 'd', 'e'],
          values: [11, 21, 31, 41, 51]
        },
        {
          varNames: ['str1', 'bool3', 'bool2', 'd', 'e'],
          values: ['aaaa', 'true', 'false', 4, 5]
        }
      ],
      updateData: new Array(100).fill(true),
      curStep: 0,
      isEnding: false,
      timeoutIDs: []
    };
  },
  methods: {
    nextStep() {
      if (this.curStep < this.steps.length) {
        // 空步骤return
        if (!(this.steps[this.curStep] && this.steps[this.curStep]['varNames'])) return;
        for (let i = 0; i < this.steps[this.curStep]['values'].length; i++) {
          if (i < this.values.length && this.values[i] === this.steps[this.curStep]['values'][i]) {
            this.$set(this.updateData, i, false);
          } else {
            this.$set(this.updateData, i, true);
            setTimeout(() => {
              this.$set(this.updateData, i, false);
            }, 500);
          }
        }
        // // console.log(this.updateData);
        this.varNames = this.steps[this.curStep]['varNames'];
        this.values = this.steps[this.curStep]['values'];
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
    },
    resetStep() {
      this.isEnding = true;
      for (let i = 0; i < this.timeoutIDs.length; i++) {
        clearTimeout(this.timeoutIDs[i]);
      }
      this.varNames = this.steps[0]['varNames'];
      this.values = this.steps[0]['values'];
      this.timeoutIDs = [];
      const stid = setTimeout(() => {
        this.curStep = 0;
        this.isEnding = false;
      }, 300);
      this.timeoutIDs.push(stid);
    },
    beforeDestroy() {
      for (let i = 0; i < this.timeoutIDs.length; i++) {
        clearTimeout(this.timeoutIDs[i]);
      }
    },
    setSteps(newSteps) {
      this.curStep = 0;
      this.steps = newSteps;
      this.varNames = this.steps[this.curStep]['varNames'];
      this.values = this.steps[this.curStep]['values'];
    }
  }
};
</script>

<style>
.list-table {
  border-collapse: collapse;
  border-radius: 5px;
  overflow: hidden;
  font-family: Arial, sans-serif;
}

.table-title th,
.table-content td {
  padding: 10px;
  text-align: center;
  border: 1px solid #fcfcff;
}

.table-title th {
  background-color: #fafaff;
  /*color: #6d8fb7;*/
}

.table-content td {
  background-color: #fcfcff;
  /*color: #506784;*/
}

.table-content th {
  background-color: #fcfcff;
  /*color: #506784;*/
}
.flash {
  animation-name: flash;
  animation-duration: 0.5s;
}
@keyframes flash {
  0% {
    background-color: transparent;
  }
  50% {
    background-color: yellow;
  }
  100% {
    background-color: transparent;
  }
}
</style>
