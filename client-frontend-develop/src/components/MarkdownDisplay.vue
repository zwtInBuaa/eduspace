<!--用于展示已知内容的mavon-editor，可调整内容和最小高度-->
<template>
  <mavon-editor
    :value="content"
    :subfield="false"
    :defaultOpen="'preview'"
    :toolbarsFlag="false"
    :boxShadow="false"
    :boxShadowStyle="'0'"
    :preview-background="'#ffffff'"
    :code-style="'github'"
    :transition="false"
    class="v-note-wrapper mav"
    :style="styleVar"
  >
  </mavon-editor>
</template>

<script>
export default {
  name: 'MarkdownDisplay',
  props: {
    content: {
      type: String,
      required: true
    },
    minHeight: {
      type: Number,
      required: false
    }
  },
  data() {
    return {};
  },
  computed: {
    styleVar() {
      return {
        '--minHeight': this.minHeight + 'px'
      };
    }
  }
};
</script>

<style scoped>
.v-note-wrapper {
  display: block !important;
  position: static !important;
}

.mav {
  min-width: 30px;
}
</style>
