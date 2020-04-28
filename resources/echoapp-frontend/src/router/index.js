import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import Server from '../views/Server.vue'

Vue.use(VueRouter)

  const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
      path: '/server',
      name: 'server',
      component: Server
  },
  {
    path: '/about_old',
    name: 'server',
    component: () => import(/* webpackChunkName: "about" */ '../views/Server.vue')
  }
]

const router = new VueRouter({
  routes
})

export default router
