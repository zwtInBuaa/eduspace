<template>
  <div class="radarChart" id="radar-chart" :style="styleString"></div>
</template>

<script>
import * as echarts from 'echarts';
export default {
  data() {
    return {
      dataArr: [
        { label: '算法', weight: 0.88 },
        { label: '数据结构', weight: 0.73 },
        { label: '编程语言', weight: 0.52 },
        { label: '数据库', weight: 0.61 },
        { label: '前端框架', weight: 0.84 }
      ],
      indicatorArr: [
        { name: '算法', max: 1 },
        { name: '数据结构', max: 1 },
        { name: '编程语言', max: 1 },
        { name: '数据库', max: 1 },
        { name: '前端框架', max: 1 }
      ]
    };
  },
  props: {
    styleString: {
      type: String,
      required: true
    },
    r: {
      type: String,
      required: true
    }
  },
  mounted() {
    // this.renderChart();
  },
  methods: {
    renderChart() {
      const myChart = echarts.init(document.getElementById('radar-chart'));

      const option = {
        // title: {
        //   text: '知识点弱项报告'
        // },
        tooltip: {},
        radar: {
          axisName: {
            //修改indicator文字的颜色
            // textStyle: {
            color: 'steelblue'
            // }
          },
          // 雷达图的指示器（也就是每个维度的名称）
          indicator: this.indicatorArr.map((item) => {
            return {
              name: item.name,
              max: item.max
            };
          }),
          // 雷达图的半径大小
          radius: this.r
        },
        series: [
          {
            // 雷达图的数据
            type: 'radar',
            data: [
              {
                // 这里假设数据是从父组件通过 props 传递过来的
                name: '正确率',
                value: this.dataArr.map((item) => item.weight.toFixed(2)),
                label: {
                  show: false
                  // formatter: '{b}: {c}'
                },
                emphasis: {
                  label: {
                    show: true
                    // formatter: '{b}: {c}'
                  }
                },
                areaStyle: {
                  // 雷达图的背景色和透明度
                  color: 'rgba(255, 0, 0, 0.3)'
                },
                lineStyle: {
                  // 雷达图的线条颜色
                  color: 'red'
                }
              }
            ],
            animationDuration: 1500, // 动画时长，单位为毫秒
            animationEasing: 'quarticInOut', // 缓动函数
            animationDelay: 50 // 动画延迟时间，单位为毫秒
          }
        ]
      };

      myChart.setOption(option);
    }
  }
};
</script>

<style scoped></style>
