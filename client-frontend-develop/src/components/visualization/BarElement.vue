<template>
  <div class="bar" :style="{ height: `${height}px`, backgroundColor: color }">
    <div class="bar-value">{{ value }}</div>
    <div class="bar-content">
      <div class="bar-index">{{ barIndex }}</div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'BarElement',
  props: {
    height: {
      type: Number,
      required: true
    },
    color: {
      type: String,
      required: true
    },
    value: {
      type: Number,
      required: true
    },
    barIndex: {
      type: Number,
      required: true
    }
  }
};
</script>

<style scoped>
.bar {
  margin: 0 2px;
  flex: 1;
  border-radius: 5px;
  /*transition: height 0.3s ease-in-out;*/
  transition: background-color 0.3s ease-in-out;
  max-width: 60px;
}

.bar-value {
  font-size: 12px;
  text-align: center;
  margin-top: -15px;
}

.bar-index {
  font-size: 12px;
  text-align: center;
  margin-bottom: -15px;
}

.bar-content {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  align-items: center;
  height: 100%;
}
</style>
