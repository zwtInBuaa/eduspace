<!--题组列表-->
<template>
  <v-container>
    <v-card class="mx-auto">
      <v-data-table :headers="headers" :items="exerciseSetList" @click:row="rowOnClick">
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
  name: 'ExerciseSetList',
  props: {
    exerciseSetList: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      headers: [
        {
          text: '题组ID',
          align: 'start',
          sortable: false,
          value: 'id'
        },
        { text: '题组名', value: 'name' }
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
