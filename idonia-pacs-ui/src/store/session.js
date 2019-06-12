import { parseLocale, defaultLocale } from '../locales/'
import i18n from '../plugins/i18n'
import router from '../router'
import Vue from 'vue'
import cookies from 'vue-cookies'
import app from '../main'

export default {
  namespaced: true,
  state: {
    locale: defaultLocale,
    loggedIn: false,
    loggingIn: false,
    token: null
  },
  mutations: {
    locale (state, value) {
      state.locale = value
      cookies.set('locale', value.code)
    },
    token (state, value) {
      state.token = value
      cookies.set('token', value)
      // Vue.$cookies.set('token', value)
    },
    loggedIn (state, value) {
      state.loggedIn = value
    },
    loggingIn (state, value) {
      state.loggingIn = value
    }
  },
  actions: {
    setLocale ({ commit }, value) {
      let locale = parseLocale(value)
      if (!locale) {
        Vue.logger.warn(`Invalid locale: ${value}`)
        locale = defaultLocale
      }
      commit('locale', locale)
      i18n.locale = locale.code
      Vue.prototype.$vuetify.lang.current = locale.code
    },
    setToken ({ commit }, value) {
      commit('token', value)
      if (value) {
        console.log(value)
        Vue.axios.defaults.headers.common['Authorization'] = `Bearer ${value}`
        // Vue.pacs.defaults.headers.common, 'Authorization', `Bearer ${value}`
        console.log(Vue.pacs.defaults.headers.common)
      }
    },
    async login ({ dispatch, commit }, credentials = {}) {
      let email = credentials.email
      let password = credentials.password
      if (email && password) {
        dispatch('setToken', null)
      }
      commit('loggingIn', true)
      try {
        await app.$idoniapacs.post('auth/login', { username: email, password }, { params: { _: Date.now() } }).then(
          response => {
            if (response.data && response.code === 200) {
              dispatch('setToken', response.data.token)
              commit('loggedIn', true)
            }
            Vue.notify({
              title: i18n.t('notify.logged_in.title'),
              type: 'primary'
            })
          }, error => { console.log(error) }
        ).catch(error => {
          console.log(error)
        })
      } catch (e) {
        commit('loggedIn', false)
        throw e
      } finally {
        commit('loggingIn', false)
      }
    },
    async logout ({ commit, dispatch, getters }) {
      commit('token', null)
      commit('loggedIn', false)

      if (getters['config/hasLoginProvider']) {
        await Vue.services.auth.logout()
      }

      router.push({ name: 'init' })
    }
  }
}
