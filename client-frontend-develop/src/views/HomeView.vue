<template>
  <v-app>
    <!-- 整体界面左侧：导航栏 -->
    <v-navigation-drawer v-if="!$vuetify.breakpoint.mobile" permanent mini-variant expand-on-hover app>
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

    <v-bottom-navigation v-else app scroll-threshold="99999">
      <v-btn v-for="item in items.slice(0, 5)" :key="item.title" :to="item.index" text>
        <small> {{ item.title }} </small>
        <v-icon color="cyan darken-1">{{ item.icon }}</v-icon>
      </v-btn>
    </v-bottom-navigation>

    <!-- 整体界面顶部：工具栏 -->
    <v-app-bar color="white" app>
      <v-toolbar-title>{{ $route.meta['title'] }}</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-toolbar-title> {{ currentDateTime }}</v-toolbar-title>
      <v-btn icon @click="showUserGuide">
        <v-icon>mdi-help-circle-outline</v-icon>
      </v-btn>
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
      let ret = [
        {
          title: '主页',
          index: '/dashboard',
          icon: 'mdi-view-dashboard'
        },
        {
          title: '交互学习',
          index: '/visualization',
          icon: 'mdi-graph'
        },
        {
          title: '习题社区',
          index: '/exerciseCenter',
          icon: 'mdi-book-open-page-variant-outline'
        },
        {
          title: '交流社区',
          index: '/forum',
          icon: 'mdi-forum'
        },
        {
          title: '个人中心',
          index: '/userCenter',
          icon: 'mdi-account'
        }
      ];
      if (this.$store.state.user.role !== '学生') {
        ret.push({
          title: '管理端',
          index: '/manage',
          icon: 'mdi-account-cog'
        });
      }
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
      this.currentDateTime = new Date().toLocaleTimeString();
    },
    showUserGuide() {
      this.$store.commit('userGuide/open', this.$route.path);
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
