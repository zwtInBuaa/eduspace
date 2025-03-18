<template>
  <div id="loginPage">
    <v-app>
      <v-main>
        <v-container class="fill-height" fluid>
          <v-row>
            <v-col :cols="this.$vuetify.breakpoint.mobile ? 12 : 5" class="mx-auto">
              <v-card class="fill-height" flat>
                <v-container class="text-center text-h5">EduCodingSpace 登录</v-container>
                <v-card-text>
                  <v-form>
                    <v-text-field
                      label="学工号"
                      name="login"
                      prepend-icon="mdi-account"
                      v-model="buaa_id"
                      required
                    ></v-text-field>
                    <v-text-field
                      id="password"
                      label="密码"
                      name="password"
                      prepend-icon="mdi-lock"
                      type="password"
                      v-model="password"
                      required
                    ></v-text-field>
                  </v-form>
                </v-card-text>
                <v-card-actions>
                  <v-btn color="primary" @click="loginWithPwd" block>登录</v-btn>
                </v-card-actions>
              </v-card>
            </v-col>
          </v-row>
        </v-container>
        <choose-course
          :choose-course-dialog="chooseCourseDialog"
          :courses="courses"
          @chooseCourse="changeCourse"
          :cur-course-name="curCourseName"
        ></choose-course>
      </v-main>
    </v-app>
  </div>
</template>

<script>
import { login } from '@/api/users';
import ChooseCourse from '@/components/login/ChooseCourse.vue';

export default {
  name: 'LoginView',
  components: {
    ChooseCourse
  },
  data() {
    return {
      buaa_id: '',
      password: '',
      chooseCourseDialog: false
    };
  },
  computed: {
    curCourseName() {
      return this.$store.state.user.curCourseName;
    },
    courses() {
      return this.$store.state.user.courses;
    }
  },
  methods: {
    async loginWithPwd() {
      const vueThis = this;
      login({
        buaa_id: this.buaa_id,
        password: this.password
      })
        .then(function (response) {
          if (
            response !== null &&
            response.data !== null &&
            response.data.token !== null &&
            response.data.token !== ''
          ) {
            let data = response.data;
            // 判断是否有多个课程
            if (data.courses !== null && data.courses.length > 0) {
              // 保存用户信息
              vueThis.$store.commit('user/setToken', data.token);
              vueThis.$store.commit('user/setBuaaId', data.buaa_id);
              vueThis.$store.commit('user/setUserId', data.user_id);
              vueThis.$store.commit('user/setUserName', data.username);
              vueThis.$store.commit('user/setRole', data.role);
              vueThis.$store.commit('user/setAvatar', data.avatar);
              // 保存课程信息
              vueThis.$store.commit('user/setCourses', data.courses);
              vueThis.$store.commit('user/setCurCourseId', data.courses[0].id);
              vueThis.$store.commit('user/setCurCourseName', data.courses[0].name);
              if (data.courses.length > 1) {
                // 弹出选择课程的弹窗
                vueThis.chooseCourseDialog = true;
              } else {
                // 跳转到首页
                vueThis.$router.push('/');
                vueThis.$store.dispatch('snackbar/success', '登录成功');
              }
            } else {
              // 没有课程
              vueThis.$store.dispatch('snackbar/warning', '您还没有加入任何课程，请联系管理员');
            }
          }
        })
        .catch(function () {
          vueThis.password = '';
        });
    },
    async changeCourse(course) {
      // 保存课程信息
      if (course !== null) {
        this.$store.commit('user/setCurCourseId', course.id);
        this.$store.commit('user/setCurCourseName', course.name);
      }
      // 关闭弹窗
      this.chooseCourseDialog = false;
      // 跳转到首页
      this.$router.push('/');
      this.$store.dispatch('snackbar/success', '登录成功');
    }
  }
};
</script>

<style scoped></style>
