const snackbar = {
  namespaced: true,
  state: {
    msg: '',
    color: '',
    visible: false,
    showClose: true,
    timeout: 5000
  },
  getters: {
    msg: (state) => state.msg,
    color: (state) => state.color,
    visible: (state) => state.visible,
    showClose: (state) => state.showClose,
    timeout: (state) => state.timeout
  },
  mutations: {
    open(state, options) {
      state.visible = true;
      state.msg = options.msg;
      state.color = options.color;
    },
    close(state) {
      state.visible = false;
    },
    setShowClose(state, isShow) {
      state.showClose = isShow;
    },
    setTimeout(state, timeout) {
      state.timeout = timeout;
    }
  },
  actions: {
    openSnackbar(context, options) {
      let timeout = context.state.timeout;
      context.commit('open', {
        msg: options.msg,
        color: options.color
      });
      setTimeout(() => {
        context.commit('close');
      }, timeout);
    },
    success(context, msg) {
      context.dispatch('openSnackbar', {
        msg: msg,
        color: 'success'
      });
    },
    error(context, msg) {
      context.dispatch('openSnackbar', {
        msg: msg,
        color: 'error'
      });
    },
    warning(context, msg) {
      context.dispatch('openSnackbar', {
        msg: msg,
        color: 'warning'
      });
    }
  }
};
export default snackbar;
