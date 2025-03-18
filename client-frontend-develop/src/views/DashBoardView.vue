<template>
  <v-card class="mx-auto my-6" :width="'95%'">
    <v-row dense>
      <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 7">
        <v-card class="ml-4" flat>
          <v-card-title> 欢迎，{{ this.userName }}！ </v-card-title>

          <v-card-subtitle> 在使用平台过程中，遇到问题可以及时向教师反馈，祝你学习顺利！ </v-card-subtitle>
          <v-spacer></v-spacer>
          <v-card-title> 为你推荐的习题 </v-card-title>
          <v-card-text>
            <exercise-list ref="exerciseList" :exercise-list="exerciseList" @clickRow="rowOnClick"></exercise-list>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 5">
        <v-card flat class="mt-4">
          <v-row>
            <v-col>
              <v-card>
                <v-card-title>题目正确率分析</v-card-title>
                <v-card-text>
                  <radar-chart
                    :data-arr="dataArr"
                    :indicator-arr="indicatorArr"
                    :style-string="'width: 100%; height: 240px;'"
                    :r="'70%'"
                    ref="radarChart"
                  />
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-card>
                <v-card-title>做题情况总览</v-card-title>
                <circle-pie-chart
                  ref="pieChart"
                  :data-arr="overviewData"
                  :style-string="'width: 100%; height: 280px;'"
                  :line-length="$vuetify.breakpoint.mobile ? 0 : 15"
                />
              </v-card>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>
  </v-card>
</template>

<script>
import radarChart from '@/components/weaknessAnalysis/radarChart.vue';
import circlePieChart from '@/components/weaknessAnalysis/circlePieChart.vue';
import ExerciseList from '@/components/exercise/ExerciseList.vue';
import { loadFailureRate, loadQuestionOverview, loadRecommendQuestions } from '@/api/weakness';
export default {
  components: {
    radarChart,
    circlePieChart,
    ExerciseList
  },
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
      ],
      exerciseList: [
        // {
        //   id: 1,
        //   title: '题目1',
        //   tags: ['tag1', 'tag2']
        // },
        // {
        //   id: 2,
        //   title: '题目2',
        //   tags: ['tag1', 'tag2']
        // }
      ],
      overviewData: [
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
    rowOnClick(item) {
      // // console.log('rowOnClick' + item.problem_id);
      this.$router.push('/exercise/' + item.id);
    }
  },
  mounted() {
    loadFailureRate(this.userId).then((res) => {
      if (res.data['data']) {
        this.dataArr = res.data['data'];
        this.indicatorArr = res.data['labels'];
        this.$refs.radarChart.indicatorArr = this.indicatorArr;
        this.$refs.radarChart.dataArr = this.dataArr;
        this.$refs.radarChart.renderChart();
      }
    });
    loadRecommendQuestions(this.userId).then((res) => {
      if (res.data && res.data['questionList']) {
        // // console.log(res.data['questionList']);
        this.exerciseList = [];
        let idPool = [];
        for (let i = 0; i < res.data['questionList'].length; i++) {
          if (idPool.includes(res.data['questionList'][i]['problem_id'])) continue;
          this.exerciseList.push({
            id: res.data['questionList'][i]['problem_id'],
            title: res.data['questionList'][i]['problem_title'],
            tags: res.data['questionList'][i]['tags']
          });
          idPool.push(res.data['questionList'][i]['problem_id']);
        }
      }
    });
    loadQuestionOverview(this.userId).then((res) => {
      if (res.data['overview']) {
        // // console.log(res.data['overview']);
        this.overviewData = [
          {
            name: '正确题目',
            value: res.data['overview']['correctNum']
          },
          {
            name: '错误题目',
            value: res.data['overview']['wrongNum']
          },
          {
            name: '未做题目',
            value: res.data['overview']['notDoneNum']
          }
        ];
        this.$refs.pieChart.dataArr = this.overviewData;
        this.$refs.pieChart.renderChart();
      }
    });
  },
  computed: {
    userId() {
      return this.$store.state.user.userId;
    },
    buaaId() {
      return this.$store.state.user.buaaId;
    },
    userName() {
      return this.$store.state.user.userName;
    },
    role() {
      return this.$store.state.user.role;
    },
    courseName() {
      return this.$store.state.user.curCourseName;
    },
    avatar() {
      return this.$store.state.user.avatar;
    },
    courses() {
      return this.$store.state.user.courses;
    },
    items() {
      return [
        //TODO
        { icon: 'mdi-account', type: '个人姓名', value: this.userName },
        { icon: 'mdi-format-list-numbered', type: '个人学号', value: this.buaaId },
        { icon: 'mdi-account-multiple', type: '用户身份', value: this.role },
        { icon: 'mdi-book-open-page-variant-outline', type: '当前课程', value: this.courseName }
      ];
    }
  }
};
</script>

<style scoped></style>
