const userGuide = {
  namespaced: true,
  state: {
    msg: '',
    visible: false
  },
  getters: {
    msg: (state) => state.msg,
    visible: (state) => state.visible
  },
  mutations: {
    open(state, message) {
      state.visible = true;
      state.msg = message;
    },
    close(state) {
      state.visible = false;
    }
  }
};
export default userGuide;
