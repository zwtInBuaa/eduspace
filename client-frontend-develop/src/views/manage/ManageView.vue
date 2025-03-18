<template>
  <v-app>
    <!-- 整体界面左侧：导航栏 -->
    <v-navigation-drawer permanent mini-variant expand-on-hover app>
      <!-- 导航栏上方：用户头像、用户名、课程名称 -->
      <v-list>
        <v-list-item class="home-head">
          <v-list-item-avatar>
            <v-img :src="avatar"></v-img>
          </v-list-item-avatar>
          <v-list-item-content>
            <v-list-item-title>{{ userName }}</v-list-item-title>
            <v-list-item-subtitle>{{ courseName }}</v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
      </v-list>
      <v-divider v-if="!$vuetify.breakpoint.mobile"></v-divider>

      <!-- 导航栏下方：功能路由 -->
      <v-list>
        <v-list-item v-for="item in items" :key="item.title" :to="item.index">
          <v-list-item-icon>
            <v-icon color="cyan darken-1">{{ item.icon }}</v-icon>
          </v-list-item-icon>
          <v-list-item-title class="text--green">
            {{ item.title }}
          </v-list-item-title>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <!-- 整体界面顶部：工具栏 -->
    <v-app-bar color="white" app>
      <v-toolbar-title>{{ $route.meta['title'] }}</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-toolbar-title> {{ currentDateTime }}</v-toolbar-title>
    </v-app-bar>

    <!-- 整体界面主体部分 -->
    <v-main>
      <router-view></router-view>
    </v-main>
  </v-app>
</template>

<script>
export default {
  data() {
    return {
      currentDateTime: null
    };
  },
  computed: {
    items() {
      let ret = [];
      if (this.$store.state.user.role === '管理员') {
        ret.push(
          {
            title: '管理端用户管理',
            index: '/manageUser',
            icon: 'mdi-account-edit'
          },
          {
            title: '管理端课程管理',
            index: '/manageCourse',
            icon: 'mdi-table-account'
          }
        );
      } else if (this.$store.state.user.role === '老师' || this.$store.state.user.role === '助教') {
        ret.push({
          title: '课程管理',
          index: '/teacherCourse',
          icon: 'mdi-table-account'
        });
      }
      ret.push(
        {
          title: '个人中心',
          index: '/managerCenter',
          icon: 'mdi-account'
        },
        {
          title: '返回用户端',
          index: '/home',
          icon: 'mdi-home'
        }
      );
      return ret;
    },
    userName() {
      return this.$store.state.user.userName;
    },
    courseName() {
      return this.$store.state.user.curCourseName;
    },
    avatar() {
      return this.$store.state.user.avatar;
    }
  },
  beforeMount() {
    setInterval(this.updateTime, 1000);
  },
  methods: {
    updateTime() {
      this.currentTime = new Date().toLocaleTimeString();
    }
  }
};
</script>

<style scoped>
.home-head {
  height: 48px;
  padding-left: 8px;
}
</style>
