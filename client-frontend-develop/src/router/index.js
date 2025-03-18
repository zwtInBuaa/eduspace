import Vue from 'vue';
import VueRouter from 'vue-router';
import store from '@/store';

import Login from '@/views/LoginView.vue';
import Home from '@/views/HomeView.vue';
import Manage from '@/views/manage/ManageView.vue';

import DashBoard from '@/views/DashBoardView.vue';
import UserCenter from '@/views/UserCenterView.vue';

import ExerciseCenter from '@/views/exercise/ExerciseCenterView.vue';
import ExerciseInfo from '@/views/exercise/ExerciseInfoView.vue';
import ExerciseSetInfo from '@/views/exercise/ExerciseSetInfoView.vue';

import Forum from '@/views/forum/ForumView.vue';
import ForumArticle from '@/views/forum/ForumArticleView.vue';

import VisualizationView from '@/views/visualization/VisualizationView.vue';
import BinarySearchView from '@/views/visualization/BinarySearchView.vue';
import SortView from '@/views/visualization/SortView.vue';
import GraphAlgoView from '@/views/visualization/GraphAlgoView.vue';

import ManageUserView from '@/views/manage/ManageUserView.vue';
import ManageCourseView from '@/views/manage/ManageCourseView.vue';
import TeacherCourseView from '@/views/manage/TeacherCourseView.vue';
import ManageCourseUserView from '@/views/manage/ManageCourseUserView.vue';
import SelectUserView from '@/views/manage/SelectUserView.vue';
import CustomVisualizationView from '@/views/visualization/CustomVisualizationView.vue';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/login',
    name: 'login',
    meta: {
      title: 'EduCodingSpace - 用户登录'
    },
    component: Login
  },
  {
    path: '/home',
    name: 'home',
    meta: {
      title: 'EduCodingSpace'
    },
    component: Home,
    redirect: '/dashboard',
    children: [
      {
        path: '/dashboard',
        name: 'dashboard',
        meta: {
          title: 'EduCodingSpace - 仪表盘'
        },
        component: DashBoard
      },
      {
        path: '/visualization',
        name: 'visualization',
        meta: {
          title: 'EduCodingSpace - 交互学习'
        },
        component: VisualizationView
      },
      {
        path: '/visualization/binarySearch/:defaultAlgo?',
        name: 'binarySearch',
        meta: {
          title: 'EduCodingSpace - 二分查找'
        },
        component: BinarySearchView
      },
      {
        path: '/visualization/sort/:defaultAlgo?',
        name: 'sort',
        meta: {
          title: 'EduCodingSpace - 排序'
        },
        component: SortView
      },
      {
        path: '/visualization/graph/:defaultAlgo?',
        name: 'graph',
        meta: {
          title: 'EduCodingSpace - 图与树算法'
        },
        component: GraphAlgoView
      },
      {
        path: '/visualization/custom',
        name: 'custom',
        meta: {
          title: 'EduCodingSpace - 自定义算法'
        },
        component: CustomVisualizationView
      },
      {
        path: '/exerciseCenter',
        name: 'exerciseCenter',
        meta: {
          title: 'EduCodingSpace - 习题社区'
        },
        component: ExerciseCenter
      },
      {
        path: '/exercise/:id',
        name: 'exerciseInfo',
        meta: {
          title: 'EduCodingSpace - 习题详情'
        },
        component: ExerciseInfo
      },
      {
        path: '/exerciseSet/:id',
        name: 'exerciseSetInfo',
        meta: {
          title: 'EduCodingSpace - 题组详情'
        },
        component: ExerciseSetInfo
      },

      {
        path: '/forum',
        name: 'forum',
        meta: {
          title: 'EduCodingSpace - 交流社区'
        },
        component: Forum
      },
      {
        path: '/forumArticle/:id',
        name: 'forumArticle',
        meta: {
          title: 'EduCodingSpace - 交流社区'
        },
        component: ForumArticle
      },
      {
        path: '/userCenter',
        name: 'userCenter',
        meta: {
          title: 'EduCodingSpace - 个人中心'
        },
        component: UserCenter
      }
    ]
  },
  {
    path: '/manage',
    name: 'manage',
    meta: {
      title: 'EduCodingSpace - 管理端'
    },
    component: Manage,
    children: [
      {
        path: '/manageUser',
        name: 'manageUser',
        meta: {
          title: 'EduCodingSpace - 管理员用户管理'
        },
        component: ManageUserView
      },
      {
        path: '/manageCourse',
        name: 'manageCourse',
        meta: {
          title: 'EduCodingSpace - 管理员课程管理'
        },
        component: ManageCourseView
      },
      {
        path: '/teacherCourse',
        name: 'teacherCourse',
        meta: {
          title: 'EduCodingSpace - 教师课程管理'
        },
        component: TeacherCourseView
      },
      {
        path: '/manageCourseUser',
        name: 'manageCourseUser',
        meta: {
          title: 'EduCodingSpace - 课程成员管理'
        },
        component: ManageCourseUserView
      },
      {
        path: '/managerCenter',
        name: 'managerCenter',
        meta: {
          title: 'EduCodingSpace - 个人中心'
        },
        component: UserCenter
      },
      {
        path: '/selectUser',
        name: 'selectUser',
        meta: {
          title: 'EduCodingSpace - 课程成员管理'
        },
        component: SelectUserView
      }
    ]
  },
  // For Public
  // {
  //   path: '/public/visualization',
  //   name: 'visualization',
  //   meta: {
  //     title: 'EduCodingSpace - 交互学习'
  //   },
  //   component: VisualizationView
  // },
  // {
  //   path: '/public/visualization/binarySearch/:defaultAlgo?',
  //   name: 'binarySearch',
  //   meta: {
  //     title: 'EduCodingSpace - 二分查找'
  //   },
  //   component: BinarySearchView
  // },
  // {
  //   path: '/public/visualization/sort/:defaultAlgo?',
  //   name: 'sort',
  //   meta: {
  //     title: 'EduCodingSpace - 排序'
  //   },
  //   component: SortView
  // },
  // {
  //   path: '/public/visualization/graph/:defaultAlgo?',
  //   name: 'graph',
  //   meta: {
  //     title: 'EduCodingSpace - 图与树算法'
  //   },
  //   component: GraphAlgoView
  // },
  {
    path: '*',
    redirect: '/home'
  }
];

const router = new VueRouter({
  routes
});

router.beforeEach((to, from, next) => {
  let token = store.state.user.token;
  if (to.path.startsWith('/public')) {
    next();
  } else if (to.path === '/login') {
    if (token === null || token === '') {
      next();
    } else {
      next('/');
    }
  } else {
    if (token === null || token === '') {
      next({
        path: '/login',
        query: {
          redirect: to.fullPath
        }
      });
    } else {
      next();
    }
  }
});

router.afterEach((to) => {
  if (to.meta.title) {
    document.title = to.meta.title;
  }
});

export default router;
