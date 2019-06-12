import Vue from 'vue'
import Vuex from 'vuex'
import session from './session'
import payload from './payload'
import config from './config'

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    session,
    payload,
    config
  }
})

export default store
