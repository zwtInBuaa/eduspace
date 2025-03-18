<!--习题列表-->
<template>
  <v-container>
    <v-card class="mx-auto">
      <v-data-table :headers="headers" :items="exerciseList" @click:row="rowOnClick">
        <template v-slot:item.tags="{ item }">
          <v-chip color="green" text-color="white" class="mr-1" v-for="tag in item.tags" :key="tag" small>
            {{ tag }}
          </v-chip>
        </template>
      </v-data-table>
    </v-card>
  </v-container>
</template>

<script>
export default {
  name: 'ExerciseList',
  props: {
    exerciseList: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      dialog: false,
      dialogDelete: false,
      headers: [
        {
          text: '题目ID',
          align: 'start',
          sortable: false,
          value: 'id'
        },
        { text: '题目名', value: 'title' },
        { text: '标签', value: 'tags' }
      ]
    };
  },
  methods: {
    rowOnClick(item) {
      this.$emit('clickRow', item);
    }
  }
};
</script>
