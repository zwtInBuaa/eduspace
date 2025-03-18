<template>
  <div>
    <v-btn icon @click="go2ManageCourse" class="mt-6">
      <v-icon>mdi-arrow-left</v-icon>
    </v-btn>
    <user-table
      :users="users"
      :course-manage="Boolean(true)"
      :course-name="courseName"
      @go2SelectUser="go2SelectUser"
      @addIdsConfirm="addIdsConfirm"
      @editConfirm="editConfirm"
      @deleteConfirm="deleteConfirm"
      @resetPassword="resetPassword"
      :create-button-name="cteateButtonName"
      :create-multi-button-name="createMultiButtonName"
    ></user-table>
  </div>
</template>

<script>
import UserTable from '@/components/manage/UserTable.vue';
import {
  deleteStudentFromCourse,
  getCourseAllStudents,
  postStudentsToCourse,
  putUser,
  resetPassword,
  userAllCourse
} from '@/api/manage';

export default {
  name: 'ManageCourseStudentView',
  components: {
    UserTable
  },
  data() {
    return {
      courseID: '',
      courseName: '',
      users: [
        {
          id: 1
        }
      ]
    };
  },
  created() {
    this.courseID = this.$route.query.id;
    this.courseName = this.$route.query.courseName;
    this.initial();
  },
  computed: {
    cteateButtonName() {
      return '选择用户添加';
    },
    createMultiButtonName() {
      return '文件导入用户组';
    }
  },

  methods: {
    go2ManageCourse() {
      this.$router.push({
        path: '/manageCourse'
      });
    },

    go2SelectUser() {
      this.$router.push({
        path: '/selectUser',
        query: {
          courseID: this.courseID,
          courseName: this.courseName
        }
      });
    },

    async initial() {
      const vueThis = this;
      this.users = [];
      /*获取全部学生*/
      await getCourseAllStudents(this.courseID)
        .then(function (response) {
          vueThis.users = response.data;
        })
        .catch(function () {});
    },

    /*确认向课程增加用户*/
    async addIdsConfirm(data) {
      const vueThis = this;
      let toAddstudents = [];
      for (let i = 0; i < data.length; i++) {
        let id = parseInt(data[i]);
        toAddstudents.push(id);
      }
      if (toAddstudents.length !== 0) {
        await postStudentsToCourse(this.courseID, toAddstudents)
          .then(function () {
            vueThis.initial();
          })
          .catch(function () {});
      }
      /*管理人员修改自己的课程时再次刷新*/
      let courses = await userAllCourse(vueThis.$store.state.user.userId);
      vueThis.$store.commit('user/setCourses', courses);
    },

    /*确认编辑保存*/
    async editConfirm(editedItem) {
      let data = {
        username: editedItem.username,
        password: editedItem.password,
        buaa_id: editedItem.buaa_id,
        role: editedItem.role
      };
      if (editedItem.buaa_id === this.$store.state.user.buaaId) {
        await this.$store.dispatch('snackbar/success', '您已修改自己信息，下次登录后生效');
      }
      const vueThis = this;
      await putUser(editedItem.id, data)
        .then(function () {
          vueThis.initial();
        })
        .catch(function () {});
    },

    /*确认删除*/
    async deleteConfirm(id) {
      id = parseInt(id);
      const vueThis = this;
      await deleteStudentFromCourse(this.courseID, id)
        .then(function () {
          vueThis.initial();
        })
        .catch(function () {});
      /*管理人员修改自己的课程时再次刷新*/
      let courses = await userAllCourse(vueThis.$store.state.user.userId);
      vueThis.$store.commit('user/setCourses', courses);
    },

    async resetPassword(data) {
      if (data.buaa_id === this.$store.state.user.buaaId) {
        await this.$store.dispatch('snackbar/success', '您已修改个人信息，下次登录后生效');
      }
      resetPassword(data)
        .then(function () {
          this.$store.dispatch('snackbar/success', '密码重置成功');
          // alert('密码重置成功');
        })
        .catch(function () {});
    }
  }
};
</script>
<style scoped></style>
