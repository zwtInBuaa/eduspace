<template>
  <v-container>
    <!-- 交流社区主界面 -->
    <v-container class="mt-10">
      <v-card class="mx-auto">
        <!-- 顶部搜索栏界面 -->
        <v-row>
          <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 7" class="ml-6 mt-1">
            <v-row v-on:keyup.enter="search">
              <v-menu offset-y>
                <template v-slot:activator="{ on }">
                  <v-text-field
                    outlined
                    dense
                    hide-details
                    class="ma-1"
                    prepend-inner-icon="mdi-magnify"
                    label="请输入关键词"
                    v-model="text"
                    v-on="on"
                    autocomplete="off"
                    ref="search"
                    clearable
                    style="max-width: 80%"
                  ></v-text-field>
                </template>
                <v-list v-if="items.length > 0" class="border-list" dense>
                  <v-list-item
                    v-for="(item, index) in items.slice(0, Math.min(6, items.length))"
                    :key="index"
                    @click="itemClick(item)"
                  >
                    <v-list-item-title>{{ item.title }}</v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-menu>
            </v-row>
          </v-col>
          <v-col>
            <v-btn color="cyan darken-1" outlined @click="search()" class="ml-4">
              <v-icon v-if="!$vuetify.breakpoint.mobile" left>mdi-magnify</v-icon>
              搜索
            </v-btn>
            <v-btn color="cyan darken-1" outlined @click="showAllArticle()" class="ml-2">
              <v-icon v-if="!$vuetify.breakpoint.mobile" left>mdi-eye</v-icon>
              全部
            </v-btn>
            <v-btn color="cyan darken-1" @click="dialog = true" outlined class="small-padding ml-2">
              <v-icon v-if="!$vuetify.breakpoint.mobile" left>mdi-plus</v-icon>
              创建讨论贴
            </v-btn>
          </v-col>
        </v-row>
        <!-- 帖子列表 -->
        <v-data-table
          :headers="[{ text: '', value: 'post' }]"
          :items="postList"
          sort-by="created_at"
          :sort-desc="true"
          hide-default-header
        >
          <template v-slot:item="{ item }">
            <v-container class="my-4 text-left">
              <router-link :to="{ name: 'forumArticle', params: { id: item.id } }">
                <v-card v-ripple>
                  <v-row class="mx-auto" align-content="center">
                    <v-col v-if="!$vuetify.breakpoint.mobile" cols="1">
                      <v-card-text>
                        <v-icon color="cyan darken-1" large> mdi-note-text-outline</v-icon>
                      </v-card-text>
                    </v-col>
                    <v-col :cols="$vuetify.breakpoint.mobile ? 12 : 6">
                      <!-- 直接使用slice 截取题目长度-->
                      <v-card-title
                        class="headline"
                        v-text="
                          item.title.length > maxTitleLen ? item.title.slice(0, maxTitleLen - 3) + '...' : item.title
                        "
                      ></v-card-title>
                      <v-card-subtitle>{{ item.poster_name }} | {{ item.intro }}</v-card-subtitle>
                    </v-col>
                    <v-col v-if="!$vuetify.breakpoint.mobile">
                      <v-card-text>
                        本帖创建于：{{ item.created_at }} <br />
                        本帖更新于：{{ item.updated_at }}
                      </v-card-text>
                    </v-col>
                  </v-row>
                </v-card>
              </router-link>
            </v-container>
          </template>
        </v-data-table>
      </v-card>
    </v-container>

    <!-- 创建新帖弹出框 -->
    <v-dialog v-model="dialog" persistent>
      <v-app-bar color="cyan">
        <v-app-bar-title>创建讨论帖</v-app-bar-title>
      </v-app-bar>
      <v-card>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field :rules="titleRules" v-model="newTitle" label="讨论贴标题" required></v-text-field>
              </v-col>
              <markdown-input ref="newContent" />
              <!--<mavon-editor :ishljs="true" v-model="newContent" class="v-note-wrapper" />-->
            </v-row>
          </v-container>
        </v-card-text>
        <!-- 卡片操作 -->
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="cyan darken-2"
            text
            @click="
              () => {
                this.$refs.newContent.content = '';
                this.newTitle = '';
                dialog = false;
              }
            "
          >
            取消创建</v-btn
          >
          <v-btn color="cyan darken-2" text @click="postNew()" :disabled="!checkTitleValid()"> 确认创建</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import { getAllBlogs, postBlog } from '@/api/forum';
import MarkdownInput from '@/components/MarkdownInput';

export default {
  name: 'ForumView',
  components: {
    MarkdownInput
  },
  data() {
    return {
      dialog: false,
      page: 1,
      newTitle: '',
      // newContent: '',
      postListAll: [
        {
          poster_name: '品东质变气层',
          intro: 'in eiusmod veniam',
          id: 64,
          title: '公系就求安太',
          created_at: '1994-12-18 19:29:31',
          updated_at: '1994-05-15 19:37:47'
        },
        {
          poster_name: '公的之',
          intro: 'eu Duis dolore',
          id: 12,
          title: '当时么性内样',
          created_at: '1997-10-31 02:01:56',
          updated_at: '1978-10-27 11:39:29'
        },
        {
          poster_name: '导线社布',
          intro: 'in mollit dolor esse laboris',
          id: 91,
          title: '分音细',
          created_at: '1972-10-25 18:06:06',
          updated_at: '2007-03-31 11:52:20'
        },
        {
          poster_name: '导线社布',
          intro: 'in mollit dolor esse laboris',
          id: 91,
          title: '分音细',
          created_at: '1972-10-25 18:06:06',
          updated_at: '2007-03-31 11:52:20'
        },
        {
          poster_name: '识更毛比般每建',
          intro: 'ullamco aliqua',
          id: 70,
          title: '最需着定志派',
          created_at: '2014-12-18 18:08:50',
          updated_at: '2004-09-09 12:31:07'
        }
      ],
      //search
      text: '',
      showSelect: true,
      postList: [],
      items: [],
      titleRules: [(v) => !!v || '必须填写标题', (v) => (v && v.length <= 32) || '标题不能超过32个字符']
    };
  },
  created() {
    this.getAllPosts();
  },
  // mounted() {
  //     this.postList = this.postListAll;
  // },
  watch: {
    text: 'inputHandle'
  },
  computed: {
    maxTitleLen() {
      return this.$vuetify.breakpoint.mobile ? 12 : 20;
    }
  },
  methods: {
    checkTitleValid() {
      return this.newTitle.length <= 32 && this.newTitle.length > 0;
    },
    getAllPosts() {
      getAllBlogs()
        .then((response) => {
          this.postListAll = response.data.posts;
          this.postList = this.postListAll;
        })
        .catch(function () {});
    },

    postNew: function () {
      postBlog({ title: this.newTitle, content: this.$refs.newContent.content }, this.$store.state.user.token)
        .then(() => {
          this.$store.dispatch('snackbar/success', '发帖成功');
          this.getAllPosts();
        })
        .catch(function () {
          this.$store.dispatch('snackbar/success', '发帖失败');
        });
      this.$refs.newContent.content = '';
      this.newTitle = '';
      this.dialog = false;
    },

    showAllArticle() {
      this.getAllPosts();
    },

    itemClick(item) {
      this.text = item.title;
      this.$refs.search.blur(); //输入框失去焦点
    },

    inputHandle() {
      let textCopy = this.text;
      if (this.text.trim() !== '') {
        //移除字符串两侧的空白字符
        this.items = this.postListAll.filter(function (value) {
          return value.title.match(new RegExp('.*' + textCopy + '.*'));
        });
      } else {
        this.items = this.postListAll;
      }
    },

    search() {
      this.$refs.search.blur(); //输入框失去焦点
      let textCopy = this.text;
      if (this.text.trim() !== '') {
        //移除字符串两侧的空白字符
        this.postList = this.postListAll.filter(function (value) {
          return value.title.match(new RegExp('.*' + textCopy + '.*'));
        });
      } else {
        this.postList = this.postListAll;
      }
    }
  }
};
</script>

<style scoped>
.v-note-wrapper {
  display: block !important;
  position: static !important;
}

.border-list {
  border: 1px solid #eee !important;
}

a {
  text-decoration: none;
}

.router-link-active {
  text-decoration: none;
}
</style>
