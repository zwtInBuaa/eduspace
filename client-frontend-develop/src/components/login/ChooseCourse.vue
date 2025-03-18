<template>
  <span>
    <v-dialog v-model="dialog" @click:outside="cancel" :width="$vuetify.breakpoint.mobile ? '100%' : '40%'">
      <v-card>
        <v-card-title class="text-h5">当前课程：{{ curCourseName }}</v-card-title>
        <!-- 描述框 -->
        <v-card-text class="mx-auto"> 您可以在下方进行切换课程，选择一个课程后点击确定即可切换到该课程。 </v-card-text>
        <!-- 选择课程栏 -->
        <v-card-text>
          <v-select
            :items="courses"
            item-text="name"
            item-value="id"
            label="选择课程"
            v-model="changeTo"
            :rules="notNullRule"
            return-object
          ></v-select>
        </v-card-text>
        <!-- 下方按钮 -->
        <v-card-actions>
          <v-col class="text-right">
            <v-btn color="primary" text @click="submit">确定</v-btn>
          </v-col>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </span>
</template>

<script>
export default {
  name: 'ChooseCourse',
  props: {
    courses: {
      type: Array,
      required: true
    },
    chooseCourseDialog: {
      type: Boolean,
      required: true
    },
    curCourseName: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      changeTo: null,
      notNullRule: [(v) => !!v || '请选择一个课程']
    };
  },
  computed: {
    dialog() {
      return this.chooseCourseDialog;
    }
  },
  methods: {
    submit() {
      if (this.changeTo) {
        this.$emit('chooseCourse', this.changeTo);
      } else {
        this.$emit('chooseCourse', null);
      }
    },
    cancel() {
      this.$emit('chooseCourse', null);
    }
  }
};
</script>

<style scoped></style>
