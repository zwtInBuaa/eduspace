<!--可选择习题列表-->
<template>
  <v-container>
    <v-card class="mx-auto">
      <v-data-table
        v-model="selected"
        :headers="headers"
        :items="exerciseList"
        item-key="id"
        @click:row="rowOnClick"
        show-select
        class="elevation-1"
      >
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
  name: 'SelectExerciseList',
  props: {
    exerciseList: {
      type: Array,
      required: true
    },
    rowOnClick: {
      type: Function,
      required: false
    }
    // exercises: {
    //   type: Array,
    //   required: false
    // }
  },
  data() {
    return {
      selected: [],
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
  }
  // watch: {
  //   selected: {
  //     handler() {
  //       // // console.log(item);
  //       var IdList = new Set();
  //       this.selected.forEach((e) => {
  //         IdList.add(e.id);
  //       });
  //       // console.log(Array.from(IdList));
  //       this.$emit('update:exercises', IdList);
  //     }
  //   }
  // }
};
</script>
