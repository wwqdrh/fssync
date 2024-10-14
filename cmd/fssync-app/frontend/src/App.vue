<template>
  <div id="app">
    <main>
      <keep-alive>
        <router-view v-if="$route.meta.keepAlive"></router-view>
      </keep-alive>
      <router-view v-if="!$route.meta.keepAlive"></router-view>
      <!-- <router-view></router-view> -->
    </main>
  </div>
</template>

<style>
body {
  margin: 0;
}

#app {
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB',
    'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  line-height: 1.7;
  background-image: linear-gradient(to top, #fbc2eb 0%, #a6c1ee 100%);
  background-blend-mode: screen, overlay, hard-light, color-burn, color-dodge, normal;
  background-attachment: fixed;
  background-repeat: no-repeat;
  background-position: 0 100%;
  color: #ffffff;
}

main {
  position: fixed;
  top: 0;
  bottom: 0;
  right: 0;
  left: 0;
  overflow: auto;
  box-sizing: border-box;
  scrollbar-width: none;
}
</style>
