<template>
  <div id="app2">
    <div id="nav">
      <router-link to="/">客户端</router-link> |
      <router-link to="/server">服务端</router-link>
    </div>
    <router-view/>
  </div>
</template>

<style lang="scss">
#app2 {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
