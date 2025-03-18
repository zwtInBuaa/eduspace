<template>
  <div>
    <user-table
      :course-manage="false"
      :users="users"
      @deleteConfirm="deleteConfirm"
      @createConfirm="createConfirm"
      @editConfirm="editConfirm"
      @resetPassword="resetPassword"
      :create-button-name="cteateButtonName"
      :create-multi-button-name="createMultiButtonName"
    ></user-table>
  </div>
</template>

<script>
import UserTable from '@/components/manage/UserTable.vue';
import { deleteUser, getAllUsers, postUsers, putUser, resetPassword } from '@/api/manage';

export default {
  name: 'ManageUserView',
  components: {
    UserTable
  },

  data: () => ({
    users: []
  }),

  computed: {
    cteateButtonName() {
      return '新建用户';
    },
    createMultiButtonName() {
      return '批量新建用户';
    }
  },

  watch: {},

  created() {
    this.myinitialize(); /*initial数据*/
  },

  methods: {
    /*初始化用户数据数组*/
    myinitialize() {
      getAllUsers()
        .then((response) => {
          this.users = response.data;
        })
        .catch((error) => {
          this.$store.dispatch('snackbar/error', error);
          // alert(error);
        });
    },

    /*确认删除*/
    async deleteConfirm(id) {
      if (id === this.$store.state.user.userId) {
        await this.$store.dispatch('snackbar/error', '不能删除自己');
        return;
      }
      const vueThis = this;
      deleteUser(id)
        .then(function () {
          vueThis.myinitialize();
        })
        .catch(function () {
          // // console.log(error);
        });
    },

    /*确认编辑保存*/
    async editConfirm(editedItem) {
      if (editedItem.buaa_id === this.$store.state.user.buaaId) {
        await this.$store.dispatch('snackbar/success', '您已修改个人信息，下次登录后生效');
      }
      let b = {
        username: editedItem.username,
        buaa_id: editedItem.buaa_id,
        role: parseInt(editedItem.role)
      };
      const vueThis = this;
      putUser(editedItem.id, b)
        .then(function () {
          vueThis.myinitialize();
        })
        .catch(function () {
          // console.log(error);
        });
    },

    createConfirm(json) {
      let b = json.map(({ username, buaa_id, password, role }) => ({ username, buaa_id, password, role }));
      const vueThis = this;
      postUsers(b)
        .then(function () {
          vueThis.myinitialize();
        })
        .catch(function () {
          // console.log(error);
        });
    },

    async resetPassword(data) {
      if (data.buaa_id === this.$store.state.user.buaaId) {
        await this.$store.dispatch('snackbar/success', '您已修改个人信息，下次登录后生效');
      }
      resetPassword(data)
        .then(function () {
          this.$store.dispatch('snackbar/success', '密码重置成功');
          // alert('密码重置成功');
        })
        .catch(function () {
          // console.log(error);
        });
    }
  }
};
</script>

<style scoped></style>
