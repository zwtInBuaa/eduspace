<!--可编辑的习题列表-->
<template>
  <v-container>
    <v-toolbar flat>
      <v-toolbar-title>习题列表</v-toolbar-title>
      <v-divider class="mx-4" inset vertical></v-divider>
      <v-spacer></v-spacer>
      <v-btn color="primary" dark class="mb-2" @click="newItem">创建新习题</v-btn>
    </v-toolbar>

    <v-card class="mx-auto">
      <v-data-table :headers="headers" :items="exerciseList">
        <template v-slot:item.tags="{ item }">
          <v-chip color="green" text-color="white" class="mr-1" v-for="tag in item.tags" :key="tag" small>
            {{ tag }}
          </v-chip>
        </template>
        <template v-slot:item.actions="{ item }">
          <v-icon small class="mr-2" @click="editItem(item)"> mdi-pencil</v-icon>
          <v-icon small @click="deleteItem(item)"> mdi-delete</v-icon>
        </template>
      </v-data-table>
    </v-card>
    <v-dialog v-model="dialogDelete" max-width="500px">
      <v-card>
        <v-card-title class="text-h5">你确定要删除这项吗？</v-card-title>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="closeDelete">取消</v-btn>
          <v-btn color="blue darken-1" text @click="deleteItemConfirm">确定</v-btn>
          <v-spacer></v-spacer>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
export default {
  name: 'EditableExerciseList',
  props: {
    exerciseList: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      dialogDelete: false,
      headers: [
        {
          text: '题目ID',
          align: 'start',
          sortable: false,
          value: 'id'
        },
        { text: '题目名', value: 'title' },
        { text: 'Actions', value: 'actions', sortable: false }
      ],
      editedIndex: -1,
      editedItem: {
        id: 1,
        name: '题目1',
        description: 'akf'
      }
    };
  },

  watch: {
    dialog(val) {
      val || this.close();
    },
    dialogDelete(val) {
      val || this.closeDelete();
    }
  },

  created() {},

  methods: {
    newItem() {
      this.editedIndex = -1;
      this.$emit('invokeExerciseEdit', this.editedIndex);
    },

    editItem(item) {
      this.editedIndex = item.id;
      this.$emit('invokeExerciseEdit', this.editedIndex);
    },

    deleteItem(item) {
      this.editedIndex = this.exerciseList.indexOf(item);
      this.editedItem = Object.assign({}, item);
      this.dialogDelete = true;
    },

    deleteItemConfirm() {
      this.$emit('deleteExerciseItemConfirm', this.editedItem.id);
      this.closeDelete();
    },

    close() {
      this.dialog = false;
      this.$nextTick(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
        this.editedIndex = -1;
      });
    },

    closeDelete() {
      this.dialogDelete = false;
      this.$nextTick(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
        this.editedIndex = -1;
      });
    }
  }
};
</script>
