<template>
  <div>
    <v-card width="90%" style="margin-left: 5%; margin-top: 5%">
      <v-data-table :headers="headers" :items="courses" sort-by="id" class="elevation-1">
        <template v-slot:top>
          <v-toolbar flat>
            <v-toolbar-title>{{ tableTitle }}</v-toolbar-title>

            <v-divider class="mx-4" inset vertical></v-divider>
            <v-spacer></v-spacer>

            <!--搜索栏-->
            <v-text-field
              class="search"
              v-model="searchItem.id"
              label="课程编号"
              append-icon="mdi-magnify"
            ></v-text-field>
            <v-text-field
              class="search"
              v-model="searchItem.name"
              label="课程名称"
              append-icon="mdi-magnify"
            ></v-text-field>
            <v-btn text color="blue" @click="showAll">全部</v-btn>
            <v-btn v-if="manage" text color="blue" @click="openEdit(defaultItem)">新建</v-btn>
          </v-toolbar>
        </template>

        <!--某一行列表最后的操作-->
        <template v-slot:[`item.actions`]="{ item }">
          <v-btn text color="#e39c04" @click="openEdit(item)">
            <v-icon color="#e39c04" left>mdi-pencil</v-icon>
            编辑
          </v-btn>
          <v-btn v-if="manage" text color="red" @click="openDelete(item)">
            <v-icon color="red" left> mdi-delete</v-icon>
            删除
          </v-btn>
          <v-btn text color="blue" @click="pushUserManage(item)">
            <v-icon color="blue" left>mdi-account</v-icon>
            用户管理
          </v-btn>
        </template>

        <!--没有数据时显示内容-->
        <template v-slot:no-data>
          <p class="v-text-field">这里没有数据</p>
        </template>
      </v-data-table>
    </v-card>

    <!--新建表项和修改表项  共用一个dialog-->
    <v-dialog v-model="dialogEdit" max-width="500px">
      <v-card>
        <v-card-title>
          <span class="text-h5">{{ formTitle }}</span>
        </v-card-title>

        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12" sm="12" md="12">
                <v-text-field v-model="editedItem.name" label="课程名称(必填)" :rules="notNullRule"></v-text-field>
              </v-col>
              <v-col cols="12" sm="12" md="12">
                <v-text-field
                  v-model="editedItem.description"
                  :rules="notNullRule"
                  label="课程介绍(必填)"
                ></v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="closeEdit">取消</v-btn>
          <v-btn color="blue darken-1" text @click="saveEdit">确定</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!--删除表项-->
    <v-dialog v-model="dialogDelete" max-width="500px">
      <v-card>
        <v-card-title class="text-h5">确定删除?</v-card-title>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="deleteClose">取消</v-btn>
          <v-btn color="blue darken-1" text @click="deleteConfirm">删除</v-btn>
          <v-spacer></v-spacer>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
export default {
  name: 'CourseTable',
  props: {
    /*存放的课程数据*/
    courses: []
  },
  data: () => ({
    /*用于控制dialog和dialogDelete窗口的打开和关闭*/
    dialogEdit: false,
    dialogDelete: false,

    /*操作表项的下标*/
    editedIndex: -1,

    /*新建，删除，编辑表项使用*/
    searchItem: {
      id: '',
      name: '',
      description: '',
      created_at: '',
      updated_at: ''
    },

    editedItem: {
      id: '',
      name: '',
      description: '',
      created_at: '',
      updated_at: ''
    },

    /*用于清空editedItem*/
    defaultItem: {
      id: '',
      name: '',
      description: '',
      created_at: '',
      updated_at: ''
    },
    notNullRule: [(v) => !!v || '此选项不能为空']
  }),

  computed: {
    manage() {
      return this.$store.state.user.role === '管理员';
    },
    tableTitle() {
      if (this.manage) {
        return '管理员课程管理';
      } else {
        return '教师端课程管理';
      }
    },

    /*新建和编辑同用一个dialog，选择不同dialog的标题*/
    formTitle() {
      return this.editedIndex === -1 ? '新建' : '编辑';
    },

    /*选择table可以有的列*/
    headers() {
      return [
        {
          text: '课程编号',
          value: 'id',
          filter: (value) => {
            if (!this.searchItem.id) return true;
            return value.toString().indexOf(this.searchItem.id.toString()) !== -1;
          }
        },
        {
          text: '课程名称',
          value: 'name',
          sortable: false,
          filter: (value) => {
            if (!this.searchItem.name) return true;
            return value.toString().indexOf(this.searchItem.name) !== -1;
          }
        },
        {
          text: '课程介绍',
          value: 'description',
          sortable: false,
          width: '20%'
        },
        {
          text: '创建时间',
          value: 'created_at'
        },
        {
          text: '最后更新时间',
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
  watch: {
    /*如果dialog为false，则调用close*/
    dialog(val) {
      val || this.closeEdit();
    },
    /*如果dialogDelete为false，则调用closeDelete()*/
    dialogDelete(val) {
      val || this.deleteClose();
    }
  },

  methods: {
    /*初始化课程数组*/
    showAll() {
      this.searchItem = Object.assign({}, this.defaultItem);
    },

    pushUserManage(item) {
      this.$router.push({
        path: '/manageCourseUser',
        query: {
          id: item.id,
          courseName: item.name
        }
      });
    },

    /*打开编辑表*/
    openEdit(item) {
      this.editedIndex = this.courses.indexOf(item);
      this.editedItem = Object.assign({}, item);
      this.dialogEdit = true;
    },
    /*关闭新建|编辑页*/
    closeEdit() {
      this.dialogEdit = false;
      /*把下次操作要再次使用的元素reset*/
      this.$nextTick(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
        this.editedIndex = -1;
      });
    },
    /*编辑|新建保存 + 关闭编辑新建页*/
    saveEdit() {
      if (!this.editedItem.name) {
        this.$store.dispatch('snackbar/warning', '请输入课程名');
        // alert('请输入课程名');

        return;
      } else if (!this.editedItem.description) {
        this.$store.dispatch('snackbar/warning', '请输入课程描述');
        // alert('请输入课程描述');
        return;
      }
      if (this.editedIndex > -1) {
        this.$emit('editItem', this.editedItem);
      } else {
        this.$emit('createItem', this.editedItem);
      }
      this.closeEdit();
    },

    /*打开删除表*/
    openDelete(item) {
      this.editedIndex = this.courses.indexOf(item);
      this.editedItem = Object.assign({}, item);
      this.dialogDelete = true;
    },
    /*关闭删除页*/
    deleteClose() {
      this.dialogDelete = false;
      this.$nextTick(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
        this.editedIndex = -1;
      });
    },
    /*删除行+关闭删除页*/
    deleteConfirm() {
      this.$emit('deleteItem', this.editedItem.id);
      this.deleteClose();
    }
  }
};
</script>

<style scoped>
.search {
  margin-left: 20px;
  margin-right: 20px;
  margin-top: 20px;
  max-width: 100px;
}

.no {
  color: #e39c04;
}
</style>
