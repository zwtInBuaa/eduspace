<template>
  <v-container>
    <div @click="handleGlobalClick">
      <v-row class="my-auto">
        <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 5">
          <v-card
            v-if="pageIndex === 0 || sb"
            :height="sb ? '500px' : '700px'"
            style="display: flex; align-items: center"
          >
            <v-row class="text-center">
              <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 5" class="align-self-center">
                <v-avatar :size="$vuetify.breakpoint.mobile ? 100 : 150">
                  <v-img :src="avatar" />
                </v-avatar>
              </v-col>
              <v-col>
                <v-col v-for="item in items" :key="item.type" class="text-left">
                  <v-icon style="margin-bottom: 1.5%">{{ item.icon }}</v-icon>
                  <span class="ml-3">
                    {{ item.type }}
                    <span class="mx-3"></span>
                    {{ item.value }}
                  </span>
                </v-col>
              </v-col>
            </v-row>
          </v-card>
        </v-col>
        <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 5">
          <v-card v-if="pageIndex === 1 || sb" :height="sb ? '500px' : '700px'">
            <v-card-title>题目错误率分析</v-card-title>
            <v-card-text>
              <radar-chart
                :data-arr="dataArr"
                :indicator-arr="indicatorArr"
                :width="'100%'"
                :r="$vuetify.breakpoint.mobile ? '65%' : '70%'"
              ></radar-chart>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      <v-row class="my-auto">
        <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 5">
          <v-card v-if="pageIndex === 2 || sb" class="mx-auto" :height="sb ? '500px' : '700px'">
            <v-card-title>做题情况总览</v-card-title>
            <circle-pie-chart
              :data-arr="overviewData"
              :line-length="$vuetify.breakpoint.mobile ? 0 : 15"
            ></circle-pie-chart>
          </v-card>
        </v-col>
        <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 5">
          <v-card v-if="pageIndex === 3 || sb" class="mx-auto" :height="sb ? '500px' : '700px'">
            <v-card-title>推荐习题</v-card-title>
            <!--        <exercise-list-->
            <!--          :exercise-list="exerciseList"-->
            <!--          :row-on-click="rowOnClick"-->
            <!--          style="overflow-y: auto; height: 600px"-->
            <!--        ></exercise-list>-->
          </v-card>
        </v-col>
      </v-row>
    </div>
    <v-container style="text-align: center" v-if="!sb">
      <v-btn @click="prevPage" :disabled="prevDisable">上一页</v-btn>
      <v-btn @click="nextPage" :disabled="nextDisable">下一页</v-btn>
    </v-container>
  </v-container>
</template>

<script>
import radarChart from '@/components/weaknessAnalysis/radarChart.vue';
// import exerciseList from '@/components/ExerciseList.vue';
import circlePieChart from '@/components/weaknessAnalysis/circlePieChart.vue';
import { loadFailureRate, loadQuestionOverview, loadRecommendQuestions } from '@/api/weakness';
export default {
  components: {
    radarChart,
    // exerciseList,
    circlePieChart
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
        {
          problem_id: 1,
          problem_title: '题目1',
          tags: ['tag1', 'tag2']
        },
        {
          problem_id: 2,
          problem_title: '题目2',
          tags: ['tag1', 'tag2']
        },
        {
          problem_id: 3,
          problem_title: '题目3',
          tags: ['tag1', 'tag2']
        },
        {
          problem_id: 1,
          problem_title: '题目1',
          tags: ['tag1', 'tag2']
        },
        {
          problem_id: 2,
          problem_title: '题目2',
          tags: ['tag1', 'tag2']
        },
        {
          problem_id: 3,
          problem_title: '题目3',
          tags: ['tag1', 'tag2']
        }
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
      ],
      pageIndex: 0,
      maxPageIndex: 2,
      nextDisable: false,
      prevDisable: true,
      sb: true
    };
  },
  methods: {
    rowOnClick(item) {
      this.$router.push('/exercise/' + item.problem_id);
    },
    nextPage() {
      if (this.nextDisable) return;
      this.prevDisable = false;
      this.pageIndex++;
      if (this.pageIndex === this.maxPageIndex) {
        this.nextDisable = true;
      }
    },
    prevPage() {
      if (this.prevDisable) return;
      this.nextDisable = false;
      this.pageIndex--;
      if (this.pageIndex === 0) {
        this.prevDisable = true;
      }
    },
    handleGlobalClick(event) {
      const screenCenterX = window.innerWidth / 2; // 获取屏幕中心点的 x 坐标

      if (event.clientX < screenCenterX) {
        // 点击屏幕左侧
        this.prevPage();
      } else {
        // 点击屏幕右侧
        this.nextPage();
      }
    }
  },
  mounted() {
    loadFailureRate(this.buaaId).then((res) => {
      if (res.data['data']) {
        this.dataArr = res.data['data'];
        this.indicatorArr = res.data['labels'];
      }
    });
    loadRecommendQuestions(this.buaaId).then((res) => {
      if (res.data['questionList']) {
        this.exerciseList = res.data['questionList'];
      }
    });
    loadQuestionOverview(this.buaaId).then((res) => {
      if (res.data['overview']) {
        this.overviewData = res.data['overview'];
      }
    });
  },
  computed: {
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
