<template>
  <router-view />
</template>

<script setup>
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
}

.el-menu {
  border-right: none;
}

.layout-container {
  height: 100vh;
}

.main-content {
  padding: 20px;
  background-color: #f0f2f5;
  min-height: calc(100vh - 60px);
}
</style>
