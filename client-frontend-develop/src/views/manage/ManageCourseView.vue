<template>
  <div>
    <course-table
      :courses="courses"
      @deleteItem="deleteItem"
      @editItem="editItem"
      @createItem="createItem"
    ></course-table>
  </div>
</template>

<script>
import CourseTable from '@/components/manage/CourseTable.vue';
import { getAllCourses, postCourse, deleteCourse, putCourse, userAllCourse } from '@/api/manage';

export default {
  name: 'ManageCourseView',
  components: {
    CourseTable
  },
  data: () => ({
    courses: [
      {
        id: 84,
        name: 'OO',
        description: '描述',
        created_at: '1984-07-15 09:10:53',
        updated_at: '1985-12-10 16:47:08'
      }
    ]
  }),
  computed: {},

  created() {
    this.myInitial(); /*initial数据*/
  },

  methods: {
    /*初始化课程数组*/
    myInitial() {
      const vueThis = this;
      getAllCourses()
        .then(function (response) {
          vueThis.courses = response.data;
        })
        .catch(function () {});
    },

    /*删除表项*/
    async deleteItem(id) {
      const vueThis = this;
      await deleteCourse(id)
        .then(function () {
          vueThis.myInitial();
        })
        .catch(function () {});
      /*管理人员修改自己的课程时再次刷新*/
      let courses = await userAllCourse(vueThis.$store.state.user.userId);
      vueThis.$store.commit('user/setCourses', courses);
    },

    async editItem(editedItem) {
      const data = {
        name: editedItem.name,
        description: editedItem.description
      };
      const vueThis = this;
      await putCourse(editedItem.id, data)
        .then(function () {
          vueThis.myInitial();
        })
        .catch(function () {});
      /*管理人员修改自己的课程时再次刷新*/
      let courses = await userAllCourse(vueThis.$store.state.user.userId);
      vueThis.$store.commit('user/setCourses', courses);
    },

    createItem(editedItem) {
      let data = [editedItem];
      data = data.map(({ name, description }) => ({ name, description }));
      const vueThis = this;
      postCourse(data)
        .then(function () {
          vueThis.myInitial();
        })
        .catch(function () {
          // console.log(error);
        });
    }
  }
};
</script>

<style scoped></style>
