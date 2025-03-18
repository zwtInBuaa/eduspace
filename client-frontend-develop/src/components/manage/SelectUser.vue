<template>
  <div>
    <v-card max-width="90%" class="mt-6 mx-auto">
      <v-container>
        <v-text-field
          class="mr-3 search-bar"
          v-model="searchItem.id"
          label="用户编码"
          append-icon="mdi-magnify"
        ></v-text-field>
        <v-text-field
          class="mr-3 search-bar"
          v-model="searchItem.buaa_id"
          label="学号"
          append-icon="mdi-magnify"
        ></v-text-field>
        <v-text-field
          class="mr-3 search-bar"
          v-model="searchItem.username"
          label="用户名"
          append-icon="mdi-magnify"
        ></v-text-field>
        <v-select class="mr-3 search-bar" v-model="searchItem.role" :items="identities" label="用户权限"></v-select>
        <v-btn text color="blue" @click="showAll">全部</v-btn>
        <v-btn text color="blue" @click="openConfirm">确定添加</v-btn>
        <v-btn text color="blue" @click="goBack">取消添加</v-btn>
      </v-container>
      <v-data-table v-model="selected" :headers="headers" :items="users" item-key="id" show-select class="elevation-1">
        <template v-slot:[`item.role`]="{ item }">
          {{ mapNumToRole(item.role) }}
        </template>
      </v-data-table>
    </v-card>

    <v-dialog v-model="dialogConfim" max-width="500px">
      <v-card>
        <v-card-title class="text-h5">确定添加?</v-card-title>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="closeConfirm">取消</v-btn>
          <v-btn color="blue darken-1" text @click="confirm">确定添加</v-btn>
          <v-spacer></v-spacer>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>
<script>
export default {
  name: 'SelectUser',
  props: {
    users: []
  },
  data() {
    return {
      dialogConfim: false,
      singleSelect: false,
      selected: [],
      /*搜索列表*/
      searchItem: {
        id: '',
        username: '',
        buaa_id: '',
        password: '',
        role: '',
        created_at: '',
        updated_at: ''
      },
      /*用于清空editedItem和searchItem*/
      defaultItem: {
        id: '',
        username: '',
        password: '',
        buaa_id: '',
        role: '',
        created_at: '',
        updated_at: ''
      },
      /*身份下拉选项*/
      identities: ['管理员', '教师', '助教', '学生']
    };
  },
  computed: {
    headers() {
      return [
        {
          text: '用户编号',
          value: 'id',
          filter: (value) => {
            if (!this.searchItem.id) return true;
            return value.toString().indexOf(this.searchItem.id.toString()) !== -1;
          }
        },
        {
          text: '学号',
          value: 'buaa_id',
          filter: (value) => {
            if (!this.searchItem.buaa_id) return true;
            return value.toString().indexOf(this.searchItem.buaa_id.toString()) !== -1;
          }
        },
        {
          text: '用户名',
          value: 'username',
          sortable: false,
          filter: (value) => {
            if (!this.searchItem.username) return true;
            return value.toString().indexOf(this.searchItem.username) !== -1;
          }
        },
        {
          text: '角色ID',
          value: 'role',
          filter: (value) => {
            if (!this.searchItem.role) return true;
            if (value < 0 || value > 2) return this.searchItem.role === '学生';
            else return this.identities[value] === this.searchItem.role;
          }
        },
        {
          text: '创建时间',
          value: 'created_at'
        },
        {
          text: '修改时间',
          value: 'updated_at'
        }
      ];
    }
  },
  methods: {
    mapRoleToNum(role) {
      for (let i = 0; i < this.identities.length; i++) {
        if (this.identities[i] === role) {
          return i;
        }
      }
      return 3;
    },
    mapNumToRole(role) {
      if (role >= 0 && role <= 3) return this.identities[role];
      else return '学生';
    },
    /*清空搜索条件，显示全部表项*/
    showAll() {
      this.searchItem = Object.assign({}, this.defaultItem);
    },

    /*选择完毕*/
    openConfirm() {
      this.dialogConfim = true;
    },
    closeConfirm() {
      this.dialogConfim = false;
      this.$nextTick(() => {
        this.selected = [];
      });
    },
    confirm() {
      let data = [];
      for (let i = 0; i < this.selected.length; i++) {
        data.push(this.selected[i].id);
      }
      this.$emit('addUsersToCourse', data);
      this.closeConfirm();
    },
    goBack() {
      let data = [];
      this.$emit('addUsersToCourse', data);
    }
  }
};
</script>

<style scoped>
.search-bar {
  display: inline-block;
  max-width: 15%;
}
</style>
