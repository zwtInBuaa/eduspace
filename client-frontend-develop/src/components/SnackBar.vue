<template>
  <v-snackbar top text v-model="visible" :color="color" :timeout="this.$store.state.snackbar.timeout">
    {{ this.$store.state.snackbar.msg }}
    <!-- 关闭按钮 -->
    <template v-slot:action="{ attrs }">
      <v-btn v-bind="attrs" v-if="showClose" icon @click="close" :color="color"><v-icon>mdi-close</v-icon> </v-btn>
    </template>
  </v-snackbar>
</template>

<script>
export default {
  name: 'SnackBar',
  computed: {
    visible() {
      return this.$store.state.snackbar.visible;
    },
    showClose() {
      return this.$store.state.snackbar.showClose;
    },
    color() {
      return this.$store.state.snackbar.color;
    }
  },
  methods: {
    close() {
      this.$store.commit('snackbar/close');
    }
  }
};
</script>
