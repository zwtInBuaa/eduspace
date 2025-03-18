<template>
  <v-container>
    <v-container>
      <v-container>
        <v-btn icon @click="$router.go(-1)">
          <v-icon>mdi-arrow-left</v-icon>
        </v-btn>
        <post-card
          :id="id"
          :title="blog.title"
          :content="blog.content"
          :avatar="blog.poster_avatar"
          :user_name="blog.poster_name"
          :created_at="blog.created_at"
          :updated_at="blog.updated_at"
          :showEditAndDelete="hasBlogAccess()"
          :isBlog="true"
          @reply="openReply"
          @edit="openEdit"
          @del="openDel"
        />
        <!-- 帖子回复卡片 -->
        <v-container v-for="reply in replies" :key="reply.id" class="mt-5">
          <post-card
            :id="reply.id"
            :title="'回复'"
            :content="reply.content"
            :avatar="reply.poster_avatar"
            :user_name="reply.poster_name"
            :created_at="reply.created_at"
            :updated_at="reply.updated_at"
            :showEditAndDelete="hasCommentAccess(reply.poster_id)"
            @reply="openReply"
            @edit="openEdit"
            @del="openDel"
          />
        </v-container>
      </v-container>

      <!-- 创建回复弹窗 -->
      <v-dialog v-model="replyEdit" persistent width="1024">
        <v-app-bar color="cyan">
          <v-app-bar-title>回复该帖</v-app-bar-title>
        </v-app-bar>
        <v-card>
          <v-card-text>
            <markdown-input ref="replyContent" :prop-info="''" />
          </v-card-text>
          <!--卡片操作-->
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="cyan darken-2" text @click="replyEdit = false"> 取消创建 </v-btn>
            <v-btn color="cyan darken-2" text @click="postReply"> 确认创建 </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

      <!-- 修改帖子/回复弹窗 -->
      <v-dialog v-model="editEdit" persistent width="1024">
        <v-app-bar color="cyan">
          <v-app-bar-title>修改帖子/回复</v-app-bar-title>
        </v-app-bar>
        <v-card>
          <v-card-text>
            <markdown-input ref="editContent" :prop-info="editInfo" />
          </v-card-text>
          <!--卡片操作-->
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="cyan darken-2" text @click="editEdit = false"> 取消修改 </v-btn>
            <v-btn color="cyan darken-2" text @click="postEdit"> 确认修改 </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

      <!-- 删除帖子/回复弹窗 -->
      <v-dialog v-model="delEdit" persistent max-width="600px">
        <v-card>
          <v-card-title> 确定要删除该帖子/回复吗？ </v-card-title>
          <!--卡片操作-->
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="cyan darken-2" text @click="delEdit = false"> 取消删除 </v-btn>
            <v-btn color="cyan darken-2" text @click="postDel"> 确认删除 </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-container>
  </v-container>
</template>

<script>
import {
  getBlogByID,
  getBlogRemark,
  postBlogRemark,
  postBlogEdit,
  postBlogRemarkEdit,
  postBlogDelete,
  postBlogRemarkDelete
} from '@/api/forum';
import PostCard from '@/components/forum/PostCard';
import MarkdownInput from '@/components/MarkdownInput';

export default {
  name: 'ForumArticleView',
  components: {
    PostCard,
    MarkdownInput
  },
  data() {
    return {
      id: 999,
      blog: {
        title: '',
        content: '',
        poster_name: '小红红',
        poster_avatar: 'https://cdn.vuetifyjs.com/images/lists/1.jpg',
        created_at: '2023-4-10 20:24',
        updated_at: '2023-4-10 20:24'
      },
      replies: [],
      replyEdit: false,
      editEdit: false,
      editId: 0,
      delEdit: false,
      delId: 0,
      isBlog: true,
      editInfo: '',
      replyInfo: ''
    };
  },
  methods: {
    // 获取权限
    hasBlogAccess() {
      return this.$store.state.user.userId === this.blog.poster_id || this.$store.state.user.role !== '学生';
    },
    hasCommentAccess(comment_id) {
      return this.$store.state.user.userId === comment_id || this.$store.state.user.role !== '学生';
    },
    // 获取帖子详情
    getPostDetail() {
      this.id = parseInt(this.$route.params.id);
      if (this.$route.params.id) {
        getBlogByID(this.$route.params.id).then((response) => {
          this.blog = {
            title: response.data.title,
            content: response.data.content,
            poster_id: response.data.poster_id,
            poster_name: response.data.poster_name,
            poster_avatar: response.data.poster_avatar,
            created_at: response.data.created_at,
            updated_at: response.data.updated_at
          };
        });
      }
    },
    getAllRemark() {
      this.replies = [];
      getBlogRemark(this.$route.params.id).then((response) => {
        if (response.data.comments !== null) {
          response.data.comments.forEach((comment) => {
            this.replies.push({
              id: comment.id,
              content: comment.content,
              poster_id: comment.user_id,
              poster_name: comment.user_name,
              poster_avatar: comment.user_avatar,
              created_at: comment.created_at,
              updated_at: comment.updated_at
            });
          });
        }
      });
    },
    // 回复帖子
    openReply() {
      this.replyEdit = true;
      if (this.$refs.replyContent) {
        this.$refs.replyContent.clearInfo();
      }
    },
    postReply() {
      postBlogRemark({ post_id: this.id, content: this.$refs.replyContent._data.content })
        .then(() => {
          this.$store.dispatch('snackbar/success', '回复成功');
          // alert('回复成功');
          this.getAllRemark();
        })
        .catch(() => {});
      this.replyEdit = false;
    },
    // 修改帖子/回复
    openEdit(id, content, isBlog) {
      // // console.log('openEditCalled');
      this.editInfo = content;
      this.editId = id;
      this.isBlog = isBlog;
      this.editEdit = true;
    },
    async postEdit() {
      // // console.log(this.$refs);
      if (this.isBlog) {
        // 修改帖子
        await postBlogEdit(this.editId, this.blog.title, this.$refs.editContent._data.content)
          .then(() => {
            this.$store.dispatch('snackbar/success', '修改成功');
            this.getPostDetail();
          })
          .catch(() => {
            this.editInfo = '';
          });
      } else {
        // 修改回复
        await postBlogRemarkEdit(this.editId, this.$refs.editContent._data.content)
          .then(() => {
            this.$store.dispatch('snackbar/success', '修改成功');
            this.getAllRemark();
          })
          .catch(() => {
            this.editInfo = '';
          });
      }
      this.editEdit = false;
    },
    // 删除帖子/回复
    openDel(id, isBlog) {
      this.delEdit = true;
      this.delId = id;
      this.isBlog = isBlog;
    },
    postDel() {
      this.delEdit = false;
      if (this.isBlog) {
        // 删除帖子
        postBlogDelete(this.delId)
          .then(() => {
            this.$store.dispatch('snackbar/success', '删除成功');
            this.$router.replace('/forum');
          })
          .catch(() => {});
      } else {
        // 删除回复
        postBlogRemarkDelete(this.delId)
          .then(() => {
            this.$store.dispatch('snackbar/success', '删除成功');
            this.getAllRemark();
          })
          .catch(() => {});
      }
    }
  },
  created() {
    this.getPostDetail();
    this.getAllRemark();
  }
};
</script>

<style scoped>
.v-note-wrapper {
  position: static !important;
}
</style>
