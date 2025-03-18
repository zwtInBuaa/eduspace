<template>
  <!-- 主帖子卡片 -->
  <v-card width="100%">
    <!-- 帖子标题 -->
    <v-card-title>
      <v-icon left>mdi-forum-outline</v-icon>
      {{ title }}
    </v-card-title>
    <v-divider></v-divider>
    <!-- 帖子用户 -->
    <v-card-text>
      <!-- 帖子作者信息 -->
      <v-row class="ma-1">
        <v-col>
          <!-- 帖子头像 -->
          <v-avatar :size="50">
            <v-img :src="avatar" />
          </v-avatar>
          <!-- 帖子姓名 -->
          <strong class="ml-3">{{ user_name }}</strong>
        </v-col>
        <v-col alignSelf="center" class="text-right"> 创建于：{{ created_at }} </v-col>
      </v-row>
      <!-- 帖子正文 -->
      <v-row>
        <v-col>
          <MarkdownDisplay :content="content" class="v-note-wrapper" />
        </v-col>
      </v-row>
    </v-card-text>
    <!-- 帖子功能按钮 -->
    <v-card-actions>
      <v-row>
        <v-col class="ml-4">
          <p style="font-size: 12px">最近一次修改于：{{ updated_at }}</p>
        </v-col>
        <v-col class="text-right">
          <v-btn color="cyan darken-2" outlined @click="reply" class="mr-2">
            <v-icon left>mdi-forum-plus-outline</v-icon>
            回复
          </v-btn>
          <v-btn v-if="showEditAndDelete" color="warning" outlined @click="edit" class="mr-2">
            <v-icon left>mdi-pencil-outline</v-icon>
            编辑
          </v-btn>
          <v-btn v-if="showEditAndDelete" color="error" outlined @click="del" class="mr-2">
            <v-icon left>mdi-delete-outline</v-icon>
            删除
          </v-btn>
        </v-col>
      </v-row>
    </v-card-actions>
  </v-card>
</template>

<script>
import MarkdownDisplay from '../MarkdownDisplay.vue';

export default {
  name: 'PostCard',
  props: {
    id: {
      type: Number,
      required: true
    },
    title: {
      type: String,
      required: true
    },
    content: {
      type: String,
      required: true
    },
    avatar: {
      type: String,
      required: true
    },
    user_name: {
      type: String,
      required: true
    },
    created_at: {
      type: String,
      required: true
    },
    updated_at: {
      type: String,
      required: true
    },
    showEditAndDelete: {
      type: Boolean,
      default: false
    },
    isBlog: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    reply() {
      this.$emit('reply');
    },
    edit() {
      this.$emit('edit', this.id, this.content, this.isBlog);
    },
    del() {
      this.$emit('del', this.id, this.isBlog);
    }
  },
  components: { MarkdownDisplay }
};
</script>

<style scoped>
.v-note-wrapper {
  position: static !important;
}
</style>
