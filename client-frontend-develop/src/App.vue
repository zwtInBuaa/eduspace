<template>
  <v-app>
    <SnackBar />
    <UserGuide />
    <router-view />
  </v-app>
</template>

<script>
import SnackBar from '@/components/SnackBar';
import UserGuide from '@/components/UserGuide.vue';

export default {
  name: 'App',
  components: {
    SnackBar,
    UserGuide
  },
  data: () => ({
    //
  })
};
</script>

<style>
html,
body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
}

.hljs {
  background-color: transparent !important;
}

.hljs > code {
  background-color: transparent !important;
}

pre > code {
  background-color: transparent !important;
}
</style>
