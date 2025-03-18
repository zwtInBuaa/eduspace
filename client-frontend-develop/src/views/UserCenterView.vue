<template>
  <v-container>
    <v-container>
      <v-card class="mx-auto" :width="$vuetify.breakpoint.mobile ? '100%' : '65%'">
        <!-- 顶部标注 -->
        <v-card-title> 个人信息</v-card-title>
        <v-divider />
        <!-- 中部内容区 -->
        <v-card-text>
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
        </v-card-text>
        <v-divider></v-divider>
        <!-- 底部按钮区 -->
        <v-card-text>
          <v-row>
            <v-col class="text-left" cols="4">
              <v-btn text color="blue" @click="openAvatar">修改头像</v-btn>
            </v-col>
            <v-col class="text-right" cols="8">
              <v-btn text color="blue" @click="openChangeCourse">切换课程</v-btn>
              <v-btn text color="error" @click="openPassword">修改密码</v-btn>
              <v-btn text color="error" @click="logout">退出登录</v-btn>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-container>
    <choose-course
      :choose-course-dialog="chooseCourseDialog"
      :courses="courses"
      @chooseCourse="changeCourse"
      :cur-course-name="courseName"
    ></choose-course>
    <change-avatar
      :change-avatar-dialog="avatarDialog"
      @avatarConfirm="updateAvatar"
      @avatarCancel="closeAvatar"
    ></change-avatar>
    <change-password
      :passwordDialog="passwordDialog"
      @changePassword="changePassword"
      @closePassword="closePassword"
    ></change-password>
  </v-container>
</template>

<script>
import { logout, changeAvatar, changePassword } from '@/api/users';

import ChooseCourse from '@/components/login/ChooseCourse.vue';
import ChangeAvatar from '@/components/usercenter/ChangeAvatar.vue';
import ChangePassword from '@/components/usercenter/ChangePassword.vue';

export default {
  components: {
    ChooseCourse,
    ChangeAvatar,
    ChangePassword
  },
  data() {
    return {
      passwordDialog: false,
      avatarDialog: false,
      chooseCourseDialog: false
    };
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
  },
  methods: {
    openAvatar() {
      this.avatarDialog = true;
    },
    closeAvatar() {
      this.avatarDialog = false;
    },
    updateAvatar(avatar) {
      let data = { image: avatar };
      changeAvatar(this.$store.state.user.userId, data)
        .then((response) => {
          if (response.status === 200) {
            this.$store.commit('user/setAvatar', response.data.path);
            this.$store.dispatch('snackbar/success', '头像设置成功');
          }
        })
        .catch(() => {});
    },
    openPassword() {
      this.passwordDialog = true;
    },
    closePassword() {
      this.passwordDialog = false;
    },
    changePassword(oldPassword, newPassword) {
      let data = {
        old_password: oldPassword,
        new_password: newPassword
      };
      changePassword(this.$store.state.user.userId, data)
        .then((response) => {
          if (response.status === 200) {
            logout().then(() => {
              this.$store.dispatch('snackbar/success', '密码修改成功，请重新登录');
              this.$store.commit('user/setToken', '');
              localStorage.clear();
              this.$router.push('/login');
            });
          } else {
            this.$store.dispatch('snackbar/error', '密码设置失败（原密码错误）');
          }
        })
        .catch(() => {
          this.$store.dispatch('snackbar/error', '密码设置失败（原密码错误）');
        });
    },

    openChangeCourse() {
      this.chooseCourseDialog = true;
    },
    changeCourse(course) {
      // 保存课程信息
      this.chooseCourseDialog = false;
      if (course !== null) {
        this.$store.commit('user/setCurCourseId', course.id);
        this.$store.commit('user/setCurCourseName', course.name);
        // 跳转到首页
        this.$router.push('/');
        this.$store.dispatch('snackbar/success', '切换课程成功');
      } else {
        this.$store.dispatch('snackbar/success', '未切换课程');
      }
    },
    async logout() {
      await logout();
      this.$store.commit('user/setToken', '');
      this.$store.dispatch('snackbar/success', '退出登录成功');
      localStorage.clear();
      this.$router.push('/login');
    }
  }
};
</script>
