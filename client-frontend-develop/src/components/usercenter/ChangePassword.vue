<template>
  <span>
    <v-dialog v-model="dialog" persistent max-width="500">
      <v-card>
        <v-card-title class="text-h5">修改密码</v-card-title>
        <v-divider></v-divider>
        <v-card-text style="margin-top: 5%">
          <v-row align="center">
            <v-col class="v-col-test" cols="12" md="12" sm="12">
              <v-text-field
                label="旧密码"
                outlined
                v-model="checkOldPassword"
                :rules="checkOldPassordRule"
                style="margin-top: 5%"
                :append-icon="show1 ? 'mdi-eye' : 'mdi-eye-off'"
                :type="show1 ? 'text' : 'password'"
                @click:append="show1 = !show1"
              ></v-text-field>
            </v-col>
          </v-row>
          <v-row align="center">
            <v-col class="d-flex" cols="12" md="12" sm="12">
              <v-text-field
                label="新密码"
                outlined
                v-model="newPassword"
                :rules="inputNewPassordRule"
                :append-icon="show2 ? 'mdi-eye' : 'mdi-eye-off'"
                :type="show2 ? 'text' : 'password'"
                @click:append="show2 = !show2"
              ></v-text-field>
            </v-col>
          </v-row>
          <v-row align="center">
            <v-col class="d-flex" cols="12" md="12" sm="12">
              <v-text-field
                label="请再次输入新密码"
                outlined
                v-model="newPasswordConfirm"
                :rules="passwordConfirmRule"
                :append-icon="show3 ? 'mdi-eye' : 'mdi-eye-off'"
                :type="show3 ? 'text' : 'password'"
                @click:append="show3 = !show3"
              ></v-text-field>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="green darken-1" text @click="changePassword"> 确定</v-btn>
          <v-btn color="green darken-1" text @click="closeDialog"> 取消</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </span>
</template>

<script>
export default {
  name: 'ChangePassword',
  props: {
    passwordDialog: Boolean
  },
  data() {
    return {
      show1: false,
      show2: false,
      show3: false,
      /*输入的新密码*/
      newPassword: '',
      newPasswordConfirm: '',
      /*输入的旧密码，需验证和旧密码是否相同*/
      checkOldPassword: '',
      checkOldPassordRule: [(v) => !!v || '请输入旧密码'],
      inputNewPassordRule: [(v) => !!v || '新设置的密码不能为空'],
      passwordConfirmRule: [(v) => !!v || '请再次确认新密码', (v) => v === this.newPassword || '两次输入的密码不一致']
    };
  },
  computed: {
    dialog() {
      return this.passwordDialog;
    }
  },
  methods: {
    changePassword() {
      if (!this.checkOldPassword) {
        this.$store.dispatch('snackbar/error', '旧密码不能为空');
        // alert('旧密码不能为空');
        return;
      } else if (!this.newPassword) {
        this.$store.dispatch('snackbar/error', '新密码不能为空');
        // alert('新密码不能为空');
        return;
      } else if (!this.newPasswordConfirm) {
        this.$store.dispatch('snackbar/warning', '请再次确认新密码');
        // alert('请再次确认新密码');
        return;
      } else if (this.newPassword !== this.newPasswordConfirm) {
        this.$store.dispatch('snackbar/error', '两次输入的密码不一致');
        // alert('两次输入的密码不一致');
        return;
      }
      this.$emit('changePassword', this.checkOldPassword, this.newPassword);
      this.closeDialog();
    },
    closeDialog() {
      this.$emit('closePassword');
      this.newPassword = '';
      this.checkOldPassword = '';
      this.newPasswordConfirm = '';
    }
  }
};
</script>

<style scoped>
.v-col-test {
  color: black;
  font-size: 120%;
}
</style>
