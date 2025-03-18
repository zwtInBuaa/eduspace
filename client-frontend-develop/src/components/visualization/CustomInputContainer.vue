<template>
  <div class="custom-input-container">
    <input
      v-model="inputValue"
      @input="$emit('input', inputValue)"
      :style="{ width: width + 'px' }"
      class="custom-input"
      maxlength="53"
    />
  </div>
</template>

<script>
export default {
  name: 'custom-input-container',
  props: {
    value: String,
    width: String
  },
  data() {
    return {
      inputValue: this.value,
      isSafe: true
    };
  },
  methods: {
    detectSQLInjection(sql) {
      const sqlKeywords = [
        // SELECT
        'SELECT',
        'FROM',
        'WHERE',
        'GROUP BY',
        'HAVING',
        'ORDER BY',
        'LIMIT',
        'OFFSET',
        // INSERT
        'INSERT',
        'INTO',
        'VALUES',
        // UPDATE
        'UPDATE',
        'SET',
        // DELETE
        'DELETE',
        // DROP
        'DROP',
        'TABLE',
        'DATABASE',
        'INDEX',
        // TRUNCATE
        'TRUNCATE'
      ];

      // 匹配注释符号（-- 和 /*）并删除
      sql = sql.replace(/(--.*)|(\s*\/\*[\s\S]*?\*\/)/g, '');

      // 匹配关键字
      for (let i = 0; i < sqlKeywords.length; i++) {
        const keyword = sqlKeywords[i];

        // 匹配关键字（不区分大小写），但排除掉作为标识符的情况
        const pattern = new RegExp(
          `(^|[\\s\\(\\)\\+\\-\\*/!><=\\|&,'"\`])${keyword}($|[\\s\\(\\)\\+\\-\\*/!><=\\|&,'"\`])`,
          'i'
        );
        if (pattern.test(sql)) {
          return true;
        }
      }

      return false;
    }
  },
  watch: {
    inputValue(newVal) {
      if (
        newVal.includes('\n') ||
        newVal.includes(';') ||
        newVal.includes('\r') ||
        newVal.includes('\\r') ||
        newVal.includes('\\n') ||
        newVal.includes('#') ||
        newVal.includes('//') ||
        newVal.includes('/*') ||
        newVal.includes('*/') ||
        newVal.includes('\\')
      ) {
        this.isSafe = false;
        this.$store.dispatch('snackbar/error', '输入中存在违禁字符，バカ!');
      } else if (newVal.toLowerCase().includes('import') || this.detectSQLInjection(newVal.toLowerCase())) {
        this.isSafe = false;
        this.$store.dispatch('snackbar/error', '输入中存在不允许的关键词，バカ!');
      } else {
        this.isSafe = true;
      }
    }
  }
};
</script>

<style scoped>
.custom-input-container {
  display: inline-block;
  border: 1px solid gray;
  border-radius: 4px;
  padding: 0;
  height: 20px;
  background-color: white;
}

.custom-input {
  border: none;
  outline: none;
  font-family: monospace;
  font-size: 14px;
  line-height: 1;
  padding: 0;
  margin: 0;
  width: 100%;
  box-sizing: border-box;
}
</style>
