<!--用于编辑的mavon-edoitor，注意父子组件传参问题-->
<template>
  <mavon-editor :ishljs="true" ref="md" @imgAdd="$imgAdd" @imgDel="$imgDel" v-model="content" class="v-note-wrapper" />
</template>

<script>
import { postRequest } from '@/api/request';

export default {
  name: 'MarkdownInput',
  props: {
    propInfo: {
      type: String,
      required: false
    }
  },
  data() {
    return {
      content: '',
      img_file: {},
      toolbars: {
        bold: true, // 粗体
        italic: true, // 斜体
        header: true, // 标题
        underline: true, // 下划线
        strikethrough: true, // 中划线
        mark: true, // 标记
        superscript: true, // 上角标
        subscript: true, // 下角标
        quote: true, // 引用
        ol: true, // 有序列表
        ul: true, // 无序列表
        link: true, // 链接
        // imagelink: true, // 图片链接,暂不支持
        code: true, // code
        table: true, // 表格
        fullscreen: false, // 全屏编辑
        readmodel: false, // 沉浸式阅读
        htmlcode: true, // 展示html源码
        help: true, // 帮助
        /* 1.3.5 */
        undo: true, // 上一步
        redo: true, // 下一步
        trash: true, // 清空
        save: false, // 保存（触发events中的save事件）
        /* 1.4.2 */
        navigation: true, // 导航目录
        /* 2.1.8 */
        alignleft: true, // 左对齐
        aligncenter: true, // 居中
        alignright: true, // 右对齐
        /* 2.2.1 */
        subfield: true, // 单双栏模式
        preview: false // 预览
      }
    };
  },
  methods: {
    // 绑定@imgAdd event
    $imgAdd(pos, $file) {
      this.img_file[pos] = $file; // 缓存图片信息
      // 第一步.将图片上传到服务器.
      var formdata = new FormData();
      formdata.append('data', $file);
      postRequest('/utils/imgdbs', formdata).then((re) => {
        // 第二步.将返回的url替换到文本原位置![...](0) -> ![...](url)
        /**
         * $vm 指为mavonEditor实例，可以通过如下两种方式获取
         * 1. 通过引入对象获取: `import {mavonEditor} from ...` 等方式引入后，`$vm`为`mavonEditor`
         * 2. 通过$refs获取: html声明ref : `<mavon-editor ref=md ></mavon-editor>，`$vm`为 `this.$refs.md`
         */
        this.$refs.md.$img2Url(pos, re.data.url);
      });
    },
    $imgDel(pos) {
      delete this.img_file[pos];
    },
    clearInfo() {
      this.content = '';
    }
  },
  mounted() {
    if (this.propInfo) {
      this.content = this.propInfo;
    } else if (this.propInfo === '') {
      this.content = '';
    }
  },
  watch: {
    propInfo: function (newVal) {
      this.content = newVal;
    }
  }
};
</script>

<style scoped>
.v-note-wrapper {
  display: block !important;
  position: static !important;
}
</style>
