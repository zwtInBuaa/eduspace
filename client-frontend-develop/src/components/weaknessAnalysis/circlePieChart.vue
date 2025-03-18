<template>
  <div ref="chartDom" :style="styleString"></div>
</template>

<script>
import * as echarts from 'echarts';

export default {
  mounted() {
    // this.renderChart();
    // // console.log(this.dataArr);
  },
  props: {
    lineLength: {
      type: Number,
      required: true
    },
    styleString: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      dataArr: [
        {
          name: '正确题目',
          value: 666
        },
        {
          name: '错误题目',
          value: 777
        },
        {
          name: '未作题目',
          value: 888
        }
      ]
    };
  },
  methods: {
    renderChart() {
      const chartDom = this.$refs.chartDom;
      const myChart = echarts.init(chartDom);
      const option = {
        tooltip: {
          trigger: 'item'
        },
        legend: {
          top: '5%',
          left: 'center'
        },
        series: [
          {
            name: '',
            type: 'pie',
            radius: ['30%', '60%'],
            avoidLabelOverlap: false,
            itemStyle: {
              borderRadius: 10,
              // borderColor: '#fff',
              borderWidth: 2,
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            },
            // label: {
            //   show: false,
            //   position: 'center'
            // },
            emphasis: {
              label: {
                show: true,
                fontSize: 20,
                fontWeight: 'bold'
              }
            },
            labelLine: {
              show: true,
              length: this.lineLength
            },
            data: this.dataArr,
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

<style>
/* 可以根据需要设置样式 */
</style>
