<template>
  <div>
    <v-btn icon @click="$router.go(-1)" class="mt-1">
      <v-icon>mdi-arrow-left</v-icon>
    </v-btn>
    <select-user :users="users" @addUsersToCourse="addUsersToCourse"></select-user>
  </div>
</template>

<script>
import SelectUser from '@/components/manage/SelectUser.vue';
import { getAllUsers, getCourseAllStudents, postStudentsToCourse, userAllCourse } from '@/api/manage';

export default {
  name: 'SelectUserView',
  data() {
    return {
      users: []
    };
  },
  components: {
    SelectUser
  },
  created() {
    this.initial();
  },
  methods: {
    async initial() {
      this.users = [];
      let allSet = [];
      let courseSet = [];
      await getAllUsers()
        .then(function (response) {
          allSet = response.data;
        })
        .catch(function () {
          // console.log(error);
        });
      await getCourseAllStudents(this.$route.query.courseID)
        .then(function (response) {
          courseSet = response.data;
        })
        .catch(function () {
          // console.log(error);
        });
      this.users = allSet.filter((v) => {
        return courseSet.every((e) => e.id !== v.id);
      });
    },

    async addUsersToCourse(data) {
      const vueThis = this;
      const courseId = vueThis.$route.query.courseID;
      if (data.length !== 0) {
        await postStudentsToCourse(courseId, data)
          .then()
          .catch(function () {
            // console.log(error);
          });
      }
      /*将管理员/教师/助教自己加入课程，前端及时刷新*/
      let courses = await userAllCourse(vueThis.$store.state.user.userId);
      vueThis.$store.commit('user/setCourses', courses);

      await this.$router.push({
        path: '/manageCourseUser',
        query: {
          id: this.$route.query.courseID,
          courseName: this.$route.query.courseName
        }
      });
    }
  }
};
</script>

<style scoped></style>
