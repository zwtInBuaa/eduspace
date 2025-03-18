<template>
  <v-row justify="center">
    <v-dialog
      v-model="showUserGuide"
      width="600px"
      :persistent="true"
      @click:outside="$store.commit('userGuide/close')"
    >
      <v-card>
        <v-card-title>
          <span class="text-h5">用户指引</span>
        </v-card-title>
        <v-card-text style="white-space: pre-wrap">{{ pathToGuideInfo }}</v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="green darken-1" text @click="$store.commit('userGuide/close')"> 关闭 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-row>
</template>

<script>
import {
  binarySearchGuideInfo,
  customGuideInfo,
  defaultGuideInfo,
  graphGuideInfo,
  loginGuideInfo,
  sortGuideInfo
} from '@/components/constStringDef';

export default {
  name: 'UserGuide',
  computed: {
    showUserGuide() {
      return this.$store.state.userGuide.visible;
    },
    pathToGuideInfo() {
      let curPath = this.$store.state.userGuide.msg;
      //采用这种写法是考虑到public路由,交互式学习部分的路由可行
      //其它部分的路由要注意includes直接匹配可能会产生的bug
      if (curPath.includes('binarySearch')) {
        return binarySearchGuideInfo;
      } else if (curPath.includes('sort')) {
        return sortGuideInfo;
      } else if (curPath.includes('graph')) {
        return graphGuideInfo;
      } else if (curPath.includes('login')) {
        return loginGuideInfo;
      } else if (curPath.includes('custom')) {
        return customGuideInfo;
      }
      return defaultGuideInfo;
    }
  }
};
</script>

<style scoped></style>
