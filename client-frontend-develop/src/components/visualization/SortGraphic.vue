<template>
  <div>
    <div class="bar-chart">
      <bar
        v-for="(bar, index) in bars"
        v-show="index < 15"
        :key="index"
        ref="bars"
        :height="bar.height"
        :color="bar.color"
        :value="bar.value"
        :bar-index="index"
        class="bar"
      />
    </div>
    <div class="btn-group">
      <!--        <button @click="shuffleBars">Shuffle Bars</button>-->
      <!--        <button @click="nextStep">Next Step</button>-->
    </div>
  </div>
</template>

<script>
import Bar from './BarElement.vue';
import Vue from 'vue';

export default {
  components: {
    Bar
  },
  data() {
    return {
      originBar: [],
      bars: [],
      barCount: 10,
      steps: [
        [0, 2],
        [1, 3],
        [1, 2],
        [0, 3],
        [0, 2]
      ],
      cutStep: 0,
      isRunning: false,
      isEnding: false,
      step_color: [[], [], [], [], []],
      timeoutIDs: []
    };
  },
  mounted() {
    this.initBars();
    this.cutStep = 0;
  },
  beforeDestroy() {
    for (let i = 0; i < this.timeoutIDs.length; i++) {
      clearTimeout(this.timeoutIDs[i]);
    }
  },
  methods: {
    initBars() {
      this.bars = Array.from({ length: this.barCount }, () => {
        return {
          height: Math.floor(Math.random() * 100 + 10),
          color: this.getRandomColor(),
          value: Math.floor(Math.random() * 100 + 10)
        };
      });
      this.copyOrigin();
    },
    getRandomColor() {
      const letters = '0123456789ABCDEF';
      let color = '#';
      for (let i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
      }
      return color;
    },
    shuffleBars() {
      this.bars = this.bars.map(() => {
        return {
          height: Math.floor(Math.random() * 100 + 10),
          color: this.getRandomColor()
        };
      });
    },
    swapBar(i, j) {
      // this.$emit('running', true);
      if (!this.$refs.bars[i]) return;
      const bar1 = this.$refs.bars[i].$el;
      const bar2 = this.$refs.bars[j].$el;
      let bar1RectLeft = bar1.getBoundingClientRect().left;
      let bar2RectLeft = bar2.getBoundingClientRect().left;
      if (i > 14) {
        bar1RectLeft = this.$refs.bars[14].$el.getBoundingClientRect().right;
      }
      if (j > 14) {
        bar2RectLeft = this.$refs.bars[14].$el.getBoundingClientRect().right;
      }

      const deltaX = bar2RectLeft - bar1RectLeft;
      const deltaY = 0; //bar2Rect.top - bar1Rect.top;
      const tmp = this.bars[i]; // 保存 Bar 数组中的元素
      bar1.animate([{ transform: 'translate(0, 0)' }, { transform: `translate(${deltaX}px, ${deltaY}px)` }], {
        duration: 300,
        easing: 'ease-out'
      }).onfinish = () => {};
      bar2.animate([{ transform: 'translate(0, 0)' }, { transform: `translate(${-deltaX}px, ${-deltaY}px)` }], {
        duration: 300,
        easing: 'ease-out'
      }).onfinish = () => {
        Vue.nextTick(() => {
          // 在下一次 DOM 更新后再交换元素
          this.$set(this.bars, i, this.bars[j]);
          this.$set(this.bars, j, tmp);
          this.$emit('running', this.isRunning);
        });
      };
    },
    setColor(randomColor) {
      let colors = this.step_color[this.cutStep];
      if (randomColor && this.cutStep === 0 && colors.length !== this.bars.length) {
        for (let i = 0; i < this.bars.length; i++) {
          this.bars[i].color = this.getRandomColor();
        }
      }
      if (colors.length > 0) {
        for (let i = 0; i < colors.length; i++) {
          this.bars[i].color = colors[i];
        }
      }
    },
    nextStep() {
      if (this.cutStep >= this.steps.length) {
        this.cutStep = 0;
        this.bars = this.originBar;
        this.copyOrigin();
      }
      if (!this.isEnding) {
        this.$emit('running', true);
        if (this.steps[this.cutStep].length !== 0) {
          this.swapBar(this.steps[this.cutStep][0], this.steps[this.cutStep][1]);
        } else {
          this.$emit('running', this.isRunning);
        }
        this.setColor(false);
        this.cutStep++;
      }
    },
    setSteps(new_steps, step_color) {
      this.cutStep = 0;
      this.steps = new_steps;
      this.step_color = step_color;
      this.setColor(true);
      this.copyOrigin();
    },
    setValues(new_heights, new_values) {
      let curBars = [];
      this.barCount = new_heights.length;
      for (let i = 0; i < new_heights.length; i++) {
        curBars.push({ height: new_heights[i], color: this.getRandomColor(), value: new_values[i] });
      }
      this.bars = curBars;
    },
    runSteps() {
      this.isRunning = true;
      if (this.cutStep >= this.steps.length || this.isEnding) {
        this.isRunning = false;
        this.$emit('running', false);
        return;
      } else {
        if (!this.isEnding) {
          this.nextStep();
          const stid = setTimeout(this.runSteps, 200);
          this.timeoutIDs.push(stid);
        }
      }
    },
    resetStep() {
      this.isEnding = true;
      const stid = setTimeout(() => {
        this.isEnding = false;
        this.bars = this.originBar; //this.steps[this.curStep++];
        this.copyOrigin();
        this.cutStep = 0;
        this.isRunning = false;
        this.$emit('running', false);
      }, 700);
      this.timeoutIDs.push(stid);
    },
    endCurrentRun() {
      this.isEnding = true;
      const stid = setTimeout(() => {
        this.isEnding = false;
        this.cutStep = 0;
        this.isRunning = false;
        this.$emit('running', false);
      }, 700);
      this.timeoutIDs.push(stid);
    },
    copyOrigin() {
      this.originBar = [];
      for (let i = 0; i < this.bars.length; i++) {
        this.originBar.push({ height: this.bars[i].height, color: this.bars[i].color, value: this.bars[i].value });
      }
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

.bar-chart ::v-deep .bar {
  transition: none;
}

.btn-group {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.btn-group button {
  margin: 0 10px;
  padding: 10px;
  border: none;
  border-radius: 5px;
  background-color: #ddd;
}

.btn-group button:hover {
  background-color: #bbb;
  cursor: pointer;
}
</style>
