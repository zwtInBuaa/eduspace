<template>
  <span>
    <v-dialog v-model="dialog" persistent max-width="500">
      <v-card>
        <v-card-title class="text-h5">修改头像</v-card-title>
        <v-divider></v-divider>

        <v-card-text style="margin-top: 5%" class="v-col-test">
          <v-file-input
            :rules="rules"
            accept="image/png"
            prepend-icon="mdi-camera"
            v-model="newAvatar"
            label="上传头像(.png格式)"
          ></v-file-input>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="green darken-1" text @click="avatarConfirm">确定修改</v-btn>
          <v-btn color="green darken-1" text @click="closeDialog">取消</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </span>
</template>

<script>
export default {
  name: 'ChangeAvatar',
  props: {
    changeAvatarDialog: Boolean
  },
  data() {
    return {
      newAvatar: [],
      rules: [(v) => !!v || '请选择头像上传']
    };
  },

  computed: {
    dialog() {
      return this.changeAvatarDialog;
    }
  },

  watch: {},
  methods: {
    async avatarConfirm() {
      if (!this.newAvatar || this.newAvatar.length === 0) {
        this.$store.dispatch('snackbar/warning', '请先上传一个头像(目前仅支持.png格式)');
        // alert('请先上传一个头像(目前仅支持.png格式)');
        return;
      }
      const fileReader = new FileReader(); // 内置方法new FileReader()读取文件
      fileReader.readAsDataURL(this.newAvatar);
      const vueThis = this;
      fileReader.onload = async function () {
        await vueThis.$emit('avatarConfirm', fileReader.result);
        vueThis.closeDialog();
        // window.setTimeout(function () {
        //   window.location.reload();
        // }, 300);
      };
    },

    closeDialog() {
      this.newAvatar = [];
      this.$emit('avatarCancel');
    }
  }
};
</script>

<style scoped></style>
