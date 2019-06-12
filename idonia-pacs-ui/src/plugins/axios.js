'use strict'

import Vue from 'vue'
import axios from 'axios'

// Full config:  https://github.com/axios/axios#request-config
// axios.defaults.baseURL = process.env.baseURL || process.env.apiUrl || '';
// axios.defaults.headers.common['Authorization'] = AUTH_TOKEN;
// axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';

let config = {
  baseURL: 'http://localhost:8080/dcm4chee-arc/' || process.env.baseURL || process.env.apiUrl || '',
  timeout: 60 * 1000, // Timeout
  withCredentials: false // Check cross-site Access-Control
}

const _axios = axios.create(config)

_axios.interceptors.request.use(
  function (config) {
    // Do something before request is sent
    return config
  },
  function (error) {
    // Do something with request error
    return Promise.reject(error)
  }
)

// Add a response interceptor
_axios.interceptors.response.use(
  function (response) {
    // Do something with response data
    return response
  },
  function (error) {
    // Do something with response error
    return Promise.reject(error)
  }
)

const Services = ({ axios }) => ({
  auth: {
    login: payload => axios.post('auth/login', payload),
    logout: () => axios.post('auth/logout')
  },
  pacs: {
    search: payload => {
      console.log(axios)
      axios.get('studies?…')
      console.log('me llaman studies')
    }
  }
})

Plugin.install = function (Vue, options) {
  Vue.axios = _axios
  window.axios = _axios
  Object.defineProperties(Vue.prototype, {
    axios: {
      get () {
        return _axios
      }
    },
    $axios: {
      get () {
        return _axios
      }
    }
  })
}
const authInstance = axios.create({
  baseURL: 'http://localhost:9001/api/',
  timeout: 60 * 1000, // Timeout
  withCredentials: false, // Check cross-site Access-Control
  headers: {
    'content-type': 'application/json'
  }
})
const servicesInitialized = Services({ axios: _axios })
Plugin.services = servicesInitialized
Vue.services = servicesInitialized
Vue.pacs = authInstance
Vue.prototype.$services = servicesInitialized
Vue.prototype.$idoniapacs = Vue.pacs
Vue.use(Plugin)

export default Plugin
