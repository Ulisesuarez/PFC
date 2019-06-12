import Vue from 'vue'
import Plugin from './plugins/axios'
import './plugins/vuetify'
import App from './App.vue'
import router from './router'
import store from './store/store'
import './registerServiceWorker'
import 'roboto-fontface/css/roboto/roboto-fontface.css'
import 'material-design-icons-iconfont/dist/material-design-icons.css'
import i18n from './plugins/i18n'
import VueCookies from 'vue-cookies'
import Notifications from 'vue-notification'
import './plugins/moment'
Vue.use(VueCookies)
Vue.config.productionTip = false
Vue.use(Notifications)
Vue.prototype.$services = Plugin.services
const app = new Vue({
  router,
  store,
  i18n,
  render: h => h(App)
}).$mount('#app')

export default app
