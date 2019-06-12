import Vue from 'vue'
import Router from 'vue-router'
import Tasks from './views/Tasks.vue'
import Login from './views/Login.vue'
import StudiesSelection from './views/StudiesSelection'
import store from './store/store'
Vue.use(Router)

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/tasks',
      name: 'tasks',
      component: Tasks
    },
    {
      path: '/',
      name: 'login',
      component: Login
    },
    {
      path: '/studies',
      name: 'studies',
      component: StudiesSelection,
      meta: { requiresAuth: true }
    }
  ]
})

router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth && !store.state.session.loggedIn) {
    next({ name: 'login' })
  }
  console.log(store.state.session.loggedIn)
  if (to.name === 'login' && store.state.session.loggedIn) {
    next({ name: 'studies' })
  }
  next()
})

export default router
