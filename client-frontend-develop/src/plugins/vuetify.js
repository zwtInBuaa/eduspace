import Vue from 'vue';
import Vuetify from 'vuetify/lib/framework';
import VuetifyToast from 'vuetify-toast-snackbar';

Vue.use(VuetifyToast, {
  x: 'right', // default
  y: 'bottom', // default
  color: 'info', // default
  icon: 'info',
  timeout: 3000, // default
  dismissable: true, // default
  autoHeight: false, // default
  multiLine: false, // default
  vertical: false, // default
  shorts: {
    custom: {
      color: 'purple'
    }
  },
  property: '$toast' // default
});

Vue.use(Vuetify);

export default new Vuetify({
  breakpoint: {
    mobileBreakpoint: 'xs'
  }
});
