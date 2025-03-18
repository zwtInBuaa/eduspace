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
        <v-btn text color="blue" @click="openSelect">
          {{ createButtonName }}
        </v-btn>
        <v-btn text color="blue" @click="openMultiSelect">
          {{ createMultiButtonName }}
        </v-btn>
      </v-container>
      <v-data-table :headers="headers" :items="users" sort-by="id" class="elevation-1">
        <template v-slot:[`item.role`]="{ item }">
          {{ mapNumToRole(item.role) }}
        </template>
        <!--某一行列表最后的操作-->
        <template v-slot:[`item.actions`]="{ item }">
          <v-btn v-if="manage" text color="#e39c04" @click="openEdit(item)">
            <v-icon color="#e39c04" left>mdi-pencil</v-icon>
            编辑
          </v-btn>
          <v-btn v-if="manage" text color="#e39c04" @click="openResetPassword(item)">
            <v-icon color="#e39c04" left>mdi-pencil</v-icon>
            重置密码
          </v-btn>
          <!--          <v-btn v-if="manage" text color="#e39c04" @click="openResetRole(item)">-->
          <!--            <v-icon color="#e39c04" left>mdi-pencil</v-icon>-->
          <!--            重置权限-->
          <!--          </v-btn>-->
          <v-btn text color="red" @click="openDelete(item)">
            <v-icon color="red" left> mdi-delete</v-icon>
            删除
          </v-btn>
        </template>
        <!--没有数据时显示内容-->
        <template v-slot:no-data>
          <p class="v-text-field">这里没有数据</p>
        </template>
      </v-data-table>
    </v-card>
    <!--以下是本组件的Dialog模块,Id代表课程-用户管理模块，通过用户id进行课程的用户管理-->

    <!-- 新建多个表项，文件导入 id-->
    <v-dialog v-model="dialogAddMultiId" max-width="500px">
      <v-card>
        <v-card-title>批量向课程导入学生</v-card-title>
        <v-card-text>
          <p>您需要上传一个满足要求的CSV表格文件</p>
          <p>除CSV文件外，暂时不支持其他格式文件</p>
          <p>表格存在共计1列</p>
          <p>表头是:"id",下面为数字,即用户的编号</p>
          <p>请注意，若该id不存在，或者该id对应的用户已经在课程里，将导入失败</p>
          <p>请注意，仅允许导入非管理员用户</p>
          <v-file-input v-model="uploadUserIds" truncate-length="15"></v-file-input>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="addMultiIdConfirm">确定添加</v-btn>
          <v-btn color="blue darken-1" text @click="closeAddMultiId">取消</v-btn>
          <v-spacer></v-spacer>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!--新建表项-->
    <v-dialog v-model="dialogCreate" max-width="500px">
      <v-card>
        <v-card-title>
          <span class="text-h5">新建</span>
        </v-card-title>

        <v-card-text>
          <v-container>
            <v-text-field
              class="mt-5 create-bar"
              v-model="editedItem.username"
              label="用户名(必填)"
              :rules="notNullRule"
              required
            ></v-text-field>
            <v-text-field
              class="mt-5 create-bar"
              v-model="editedItem.buaa_id"
              label="学号(必填)"
              required
            ></v-text-field>
            <v-text-field
              class="mt-5 create-bar"
              v-model="editedItem.password"
              label="用户密码(必填)"
              :rules="notNullRule"
              required
            ></v-text-field>
            <v-select
              class="mt-5 create-bar"
              dense
              background-color="rgba(255, 255, 255, 0)"
              v-model="editedItem.role"
              :items="identities"
              filled
              label="用户权限(必填)"
              :rules="roleRule"
              required
            ></v-select>
          </v-container>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="closeCreate()">取消</v-btn>
          <v-btn color="blue darken-1" text @click="createConfirm()">确定</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 新建多个表项，文件导入-->
    <v-dialog v-model="dialogCreateMulti" max-width="500px">
      <v-card>
        <v-card-title>批量导入学生</v-card-title>
        <v-card-text>
          <p>您需要上传一个满足要求的CSV表格文件</p>
          <p>除CSV文件外，暂不支持其他格式文件</p>
          <p>表格存在共计4列</p>
          <p>表头依次是:"buaa_id"，"username","password","role"</p>
          <p>意义依次是:"学工号"，"用户名"，"密码"，"用户权限”</p>
          <p>用户权限栏应为数字，数字对应的权限为：</p>
          <p>0(管理员),1(教师),2(助教),3(学生)</p>
          <p>请注意，若该学工号对应的账号已存在，则导入失败</p>
          <v-file-input v-model="uploadUsers" truncate-length="15"></v-file-input>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="createMultiConfirm()">确定创建</v-btn>
          <v-btn color="blue darken-1" text @click="closeCreateMulti()">取消</v-btn>
          <v-spacer></v-spacer>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 编辑表项-->
    <v-dialog v-model="dialogEdit" max-width="500px">
      <v-card>
        <v-card-title>
          <span class="text-h5">编辑</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12" sm="12" md="12">
                <v-text-field v-model="editedItem.username" label="用户名(必填)" :rules="notNullRule"></v-text-field>
              </v-col>
              <v-col cols="12" sm="12" md="12">
                <v-text-field v-model="editedItem.buaa_id" label="学号(必填)"></v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="closeEdit">取消</v-btn>
          <v-btn color="blue darken-1" text @click="editConfirm">确定</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!--重置密码-->
    <v-dialog v-model="dialogResetPassword" max-width="500px">
      <v-card>
        <v-card-title>
          <span class="text-h5">重置密码</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12" sm="12" md="12">
                <v-text-field v-model="editedItem.password" label="重置密码(必填)"></v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="closeResetPassword">取消</v-btn>
          <v-btn color="blue darken-1" text @click="resetPasswordConfirm">确定</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!--重置用户权限-->
    <v-dialog v-model="dialogResetRole" max-width="500px">
      <v-card>
        <v-card-title>
          <span class="text-h5">重置用户权限</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12" sm="12" md="12">
                <v-select
                  background-color="rgba(255, 255, 255, 0)"
                  v-model="editedItem.role"
                  :items="identities"
                  filled
                  label="用户权限(必填)"
                ></v-select>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="closeResetRole">取消</v-btn>
          <v-btn color="blue darken-1" text @click="resetRoleConfirm">确定</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 删除表项 -->
    <v-dialog v-model="dialogDelete" max-width="500px">
      <v-card>
        <v-card-title class="text-h5">确定删除?</v-card-title>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="closeDelete()">取消</v-btn>
          <v-btn color="blue darken-1" text @click="deleteConfirm()">删除</v-btn>
          <v-spacer></v-spacer>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import csvToJson from 'csvtojson';

export default {
  /*此组件显示用户列表，供管理端用户管理+管理端课程管理的管理课程用户使用*/
  name: 'UserTable',
  props: {
    /*是否是课程用户管理*/
    courseManage: Boolean,
    courseName: String,
    /*传入的用户数据*/
    users: [],
    /*按键名字，UserManageView:"新建"
     * ManageCourseStudentView:"添加用户"*/
    createButtonName: String,
    createMultiButtonName: String
  },

  data: () => ({
    /*用于csv文件导入*/
    uploadUsers: [],
    uploadUserIds: [],
    /*用于控住Dialog窗口，只需输入id即可*/
    dialogAddId: false,
    /*判断要加入课程的一个用户是否存在(输入id，判断是否存在)*/
    addError: false,
    dialogAddMultiId: false,
    /*用于控制dialog窗口的打开和关闭*/
    dialogCreate: false,
    dialogCreateMulti: false,
    dialogEdit: false,
    dialogResetPassword: false,
    dialogResetRole: false,
    dialogDelete: false,

    /*身份下拉选项*/
    identities: ['管理员', '教师', '助教', '学生'],

    /*操作的表项的下标*/
    editedIndex: -1,

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

    /*新建，删除，编辑表项使用*/
    editedItem: {
      id: '',
      username: '',
      password: '',
      buaa_id: '',
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

    /*表单填入规则 2个*/
    notNullRule: [(v) => !!v || '此选项不能为空'],
    roleRule: [(v) => !!v || '此选项必填']
  }),

  computed: {
    manage() {
      return this.$store.state.user.role === '管理员';
    },
    tableTitle() {
      if (this.courseManage) {
        return this.courseName + '  用户管理';
      } else {
        return '用户管理';
      }
    },
    /*table表头*/
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
        },
        {
          text: '操作',
          value: 'actions',
          sortable: false
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
    openSelect() {
      if (this.courseManage) {
        this.go2SelectUser();
      } else {
        this.openCreate();
      }
    },
    openMultiSelect() {
      if (this.courseManage) {
        this.openAddMultiId();
      } else {
        this.openCreateMulti();
      }
    },

    go2SelectUser() {
      this.$emit('go2SelectUser');
    },

    /*打开创建多个用户Dialog*/
    openAddMultiId() {
      this.dialogAddMultiId = true;
    },
    /*关闭创建多个用户Dialog*/
    closeAddMultiId() {
      this.dialogAddMultiId = false;
      this.$nextTick(() => {
        this.uploadUserIds = [];
      });
    },
    /*确认创建多个用户Dialog*/
    async addMultiIdConfirm() {
      // bug:这一部分其实我不太会，目前只是能用
      if (this.uploadUserIds.length === 0) {
        this.$store.dispatch('snackbar/warning', '请先上传一个文件');
        // alert('请先上传一个文件');
        return;
      }
      const file = new FileReader();
      file.readAsBinaryString(this.uploadUserIds);
      const vueThis = this;
      file.onload = async function () {
        let json = await csvToJson().fromString(file.result);
        let data = [];
        for (let i = 0; i < json.length; i++) {
          data.push(parseInt(json[i].id));
        }
        vueThis.$emit('addIdsConfirm', data);
        vueThis.closeAddMultiId();
      };
    },

    /*打开新建Dialog*/
    openCreate() {
      this.dialogCreate = true;
    },
    /*关闭新建Dialog*/
    closeCreate() {
      this.dialogCreate = false;
      this.$nextTick(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
      });
    },
    /*新建保存+关闭编辑新建Dialog*/
    createConfirm() {
      if (!this.editedItem.username) {
        this.$store.dispatch('snackbar/warning', '请填写用户名');
        // alert('请填写用户名');
        return;
      } else if (!this.editedItem.buaa_id) {
        this.$store.dispatch('snackbar/warning', '请填写用户学号');
        // alert('请填写用户学号');
        return;
      } else if (!this.editedItem.password) {
        this.$store.dispatch('snackbar/warning', '请填写用户密码');
        // alert('请填写用户密码');
        return;
      } else if (!this.editedItem.role) {
        this.$store.dispatch('snackbar/warning', '请选择用户权限');
        // alert('请选择用户权限');
        return;
      }
      this.editedItem.role = this.mapRoleToNum(this.editedItem.role);
      const data = [this.editedItem];
      this.$emit('createConfirm', data);
      this.closeCreate();
    },

    /*打开创建多个用户Dialog*/
    openCreateMulti() {
      this.dialogCreateMulti = true;
    },
    /*关闭创建多个用户Dialog*/
    closeCreateMulti() {
      this.dialogCreateMulti = false;
      this.$nextTick(() => {
        this.uploadUsers = [];
      });
    },
    /*确认创建多个用户Dialog*/
    async createMultiConfirm() {
      if (this.uploadUsers.length === 0) {
        this.$store.dispatch('snackbar/warning', '请先上传文件');
        // alert('请先上传文件');
        return;
      }
      const file = new FileReader();
      file.readAsBinaryString(this.uploadUsers);
      const vueThis = this;
      file.onload = async function () {
        let json = await csvToJson().fromString(file.result);
        for (let i = 0; i < json.length; i++) {
          json[i].role = parseInt(json[i].role);
        }
        vueThis.$emit('createConfirm', json);
        vueThis.dialogMulti = false;
        vueThis.closeCreateMulti();
      };
    },

    /*打开编辑Dialog*/
    openEdit(item) {
      this.editedItem = Object.assign({}, item);
      this.dialogEdit = true;
    },
    /*关闭编辑Dialog*/
    closeEdit() {
      this.dialogEdit = false;
      this.$nextTick(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
      });
    },
    /*编辑保存+关闭编辑编辑Dialog*/
    editConfirm() {
      if (!this.editedItem.username) {
        this.$store.dispatch('snackbar/warning', '请填写用户名');
        // alert('请填写用户名');
        return;
      } else if (!this.editedItem.buaa_id) {
        this.$store.dispatch('snackbar/warning', '请填写用户学号');
        // alert('请填写用户学号');
        return;
      }
      //TODO :暂时不能编辑身份
      this.$emit('editConfirm', this.editedItem);
      this.closeEdit();
    },

    /*打开编辑密码Dialog*/
    openResetPassword(item) {
      this.editedItem = Object.assign({}, item);
      this.dialogResetPassword = true;
    },
    /*关闭编辑密码Dialog*/
    closeResetPassword() {
      this.dialogResetPassword = false;
      this.$nextTick(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
      });
    },
    /*编辑保存+关闭编辑密码Dialog*/
    resetPasswordConfirm() {
      if (!this.editedItem.password) {
        this.$store.dispatch('snackbar/warning', '请填写密码');
        // alert('请填写密码');
        return;
      }
      let data = {
        buaa_id: this.editedItem.buaa_id,
        password: this.editedItem.password
      };
      this.$emit('resetPassword', data);
      this.closeResetPassword();
    },

    /*打开编辑用户权限Dialog*/
    openResetRole(item) {
      this.editedItem = Object.assign({}, item);
      this.dialogResetRole = true;
    },
    /*关闭编辑用户权限Dialog*/
    closeResetRole() {
      this.dialogResetRole = false;
      this.$nextTick(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
      });
    },
    /*编辑保存+关闭编辑用户权限Dialog*/
    resetRoleConfirm() {
      if (
        !this.editedItem.role ||
        !(this.editedItem.role === '0' || this.editedItem.role === '1' || this.editedItem.role === '2')
      ) {
        this.$store.dispatch('snackbar/warning', '请选择一个用户权限');
        // alert('请选择一个用户权限');
        return;
      }
      this.$emit('editConfirm', this.editedItem);
      this.closeResetRole();
    },

    /*打开删除Dialog*/
    openDelete(item) {
      this.editedItem = Object.assign({}, item);
      this.dialogDelete = true;
    },
    /*关闭删除Dialog*/
    closeDelete() {
      this.dialogDelete = false;
      this.$nextTick(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
      });
    },
    /*删除行+关闭删除Dialog*/
    deleteConfirm() {
      this.$emit('deleteConfirm', this.editedItem.id);
      this.closeDelete();
    }
  }
};
</script>

<style scoped>
.search-bar {
  display: inline-block;
  max-width: 15%;
}

.create-bar {
  display: block;
  max-width: 100%;
}
</style>
