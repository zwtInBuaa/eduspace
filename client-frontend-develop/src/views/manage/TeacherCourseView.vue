<template>
  <div>
    <course-table @editItem="editItem" :courses="courses"></course-table>
  </div>
</template>

<script>
import CourseTable from '@/components/manage/CourseTable.vue';
import { getStudentsAllCourse, putCourse, userAllCourse } from '@/api/manage';

export default {
  name: 'TeacherCourseView',
  components: {
    CourseTable
  },
  data: () => ({
    courses: []
  }),
  computed: {},

  created() {
    this.initial(); /*initial数据*/
  },

  methods: {
    /*初始化课程数组*/
    initial() {
      const vueThis = this;
      getStudentsAllCourse(this.$store.state.user.userId)
        .then(function (response) {
          vueThis.courses = response.data;
        })
        .catch(function () {
          // console.log(error);
        });
    },

    async editItem(editedItem) {
      const data = {
        name: editedItem.name,
        description: editedItem.description
      };
      const vueThis = this;
      putCourse(editedItem.id, data)
        .then(function () {
          vueThis.initial();
        })
        .catch(function () {
          // console.log(error);
        });
      /*在管理人员修改自己所在课程时及时更新*/
      let courses = await userAllCourse(vueThis.$store.state.user.userId);
      vueThis.$store.commit('user/setCourses', courses);
    }
  }
};
</script>

<style scoped></style>
