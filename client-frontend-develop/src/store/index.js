import Vue from 'vue';
import Vuex from 'vuex';
import VuexPersistence from 'vuex-persist';

import user from './user';
import snackbar from './snackbar';
import userGuide from '@/store/uesrGuide';

Vue.use(Vuex);

const vuexLocal = new VuexPersistence({
  storage: window.localStorage,
  modules: ['user']
});

export default new Vuex.Store({
  modules: {
    user,
    snackbar,
    userGuide
  },
  plugins: [vuexLocal.plugin]
});
