<template>
  <div class="bar-chart" style="text-align: center">
    <bar
      v-for="(bar, index) in bars"
      v-show="index < 15"
      :key="index"
      ref="bars"
      :height="bar.height"
      :color="bar.color"
      :value="bar.value"
      :bar-index="index"
    />
  </div>
</template>

<script>
import Bar from '@/components/visualization/BarElement.vue';

export default {
  components: { Bar },
  data() {
    return {
      timeoutIDs: [],
      curStep: 0,
      bars: [
        { height: 20, color: 'blue', value: 20 },
        { height: 30, color: 'blue', value: 120 },
        { height: 40, color: 'blue', value: 220 },
        { height: 50, color: 'blue', value: 320 },
        { height: 60, color: 'blue', value: 420 },
        { height: 70, color: 'blue', value: 520 },
        { height: 80, color: 'blue', value: 620 },
        { height: 90, color: 'blue', value: 720 }
      ],
      steps: [
        [
          { height: 20, color: 'red', value: 20 },
          { height: 30, color: 'blue', value: 20 },
          { height: 40, color: 'blue', value: 20 },
          { height: 50, color: 'yellow', value: 20 },
          { height: 60, color: 'blue', value: 20 },
          { height: 70, color: 'blue', value: 20 },
          { height: 80, color: 'blue', value: 20 },
          { height: 90, color: 'red', value: 20 }
        ],
        [
          { height: 20, color: 'blue', value: 20 },
          { height: 30, color: 'blue', value: 20 },
          { height: 40, color: 'blue', value: 20 },
          { height: 50, color: 'red', value: 20 },
          { height: 60, color: 'blue', value: 20 },
          { height: 70, color: 'yellow', value: 20 },
          { height: 80, color: 'blue', value: 20 },
          { height: 90, color: 'red', value: 20 }
        ],
        [
          { height: 20, color: 'blue', value: 20 },
          { height: 30, color: 'blue', value: 20 },
          { height: 40, color: 'blue', value: 20 },
          { height: 50, color: 'red', value: 20 },
          { height: 60, color: 'yellow', value: 20 },
          { height: 70, color: 'red', value: 20 },
          { height: 80, color: 'blue', value: 20 },
          { height: 90, color: 'blue', value: 20 }
        ],
        [
          { height: 20, color: 'blue', value: 20 },
          { height: 30, color: 'blue', value: 20 },
          { height: 40, color: 'blue', value: 20 },
          { height: 50, color: 'blue', value: 20 },
          { height: 60, color: 'yellow', value: 20 },
          { height: 70, color: 'blue', value: 20 },
          { height: 80, color: 'blue', value: 20 },
          { height: 90, color: 'blue', value: 20 }
        ]
      ]
    };
  },
  methods: {
    randomize() {},
    nextStep() {
      if (this.curStep >= this.steps.length) {
        this.curStep = 0;
      }
      this.bars = this.steps[this.curStep++];
    },
    setSteps(new_steps) {
      this.bars = new_steps[1];
      this.steps = new_steps;
      this.curStep = 2;
      //this.offset = (this.chartWidth - new_steps[1].length * this.barWidth)/2
    },
    runSteps() {
      if (this.curStep < this.steps.length) {
        this.nextStep();
        this.timeoutIDs.push(setTimeout(this.runSteps, 200));
      } else {
        return;
      }
    },
    resetStep() {
      for (let i = 0; i < this.timeoutIDs.length; i++) {
        clearTimeout(this.timeoutIDs[i]);
      }
      this.timeoutIDs = [];
      this.curStep = 0;
      this.bars = this.steps[this.curStep++];
    }
  },
  beforeDestroy() {
    for (let i = 0; i < this.timeoutIDs.length; i++) {
      clearTimeout(this.timeoutIDs[i]);
    }
  }
};
</script>

<style scoped>
.bar-chart {
  display: flex;
  align-items: flex-end;
  height: 200px;
  max-width: 800px;
}

button {
  margin-top: 20px;
}
</style>
