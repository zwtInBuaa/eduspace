const user = {
  namespaced: true,
  state: {
    // 用户登录状态
    token: '',
    // 用户信息
    buaaId: '',
    userId: '',
    userName: '',
    role: '',
    avatar: '',
    // 用户课程
    courses: [],
    curCourseId: '',
    curCourseName: ''
  },
  getters: {
    // 获取用户登录状态
    token: (state) => state.token,
    // 获取用户信息
    buaaId: (state) => state.buaaId,
    userId: (state) => state.userId,
    userName: (state) => state.userName,
    role: (state) => state.role,
    avatar: (state) => state.avatar,
    // 获取用户课程
    courses: (state) => state.courses,
    curCourseId: (state) => state.cur_course_id,
    curCourseName: (state) => state.cur_course_name
  },
  mutations: {
    // 设置用户登录状态
    setToken: (state, token) => {
      state.token = token;
    },
    // 设置用户信息
    setBuaaId: (state, buaaId) => {
      state.buaaId = buaaId;
    },
    setUserId: (state, userId) => {
      state.userId = userId;
    },
    setUserName: (state, userName) => {
      state.userName = userName;
    },
    setRole: (state, role) => {
      state.role = role;
    },
    setAvatar: (state, avatar) => {
      state.avatar = avatar;
    },
    // 设置用户课程
    setCourses: (state, courses) => {
      state.courses = courses;
    },
    setCurCourseId: (state, curCourseId) => {
      state.curCourseId = curCourseId;
    },
    setCurCourseName: (state, curCourseName) => {
      state.curCourseName = curCourseName;
    }
  },
  actions: {}
};

export default user;
