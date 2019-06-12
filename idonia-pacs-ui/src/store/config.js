import Vue from 'vue'

export default {
  namespaced: true,
  state: {
    loading: false,
    loaded: false,
    data: {
      welcomeView: 'files-selection',
      login: {
        loginProvider: '',
        forgotPasswordUrl: ''
      },
      outputTemplate: {},
      pacs: [{ id: 'dcm4che', name: 'DCM4CHEE' }],
      consentId: [],
      destinations: [
      //   {
      //     id: 'another',
      //     name: 'Another'
      //   },
      //   {
      //     id: 'not_dkv',
      //     name: 'NOT DKV',
      //     imageUrl: 'https://i.redditmedia.com/zl8Z9CZ4PhQkx8xIbkUa1As66c8BZdkHGfQOQ-y6394.jpg?s=0f45065c34d01062ba7cd692db1811df'
      //   },
      //   {
      //     id: 'dkv',
      //     name: 'DKV',
      //     imageUrl: 'https://upload.wikimedia.org/wikipedia/commons/2/2c/DKV_%28Versicherung%29_logo.svg',
      //     additionalData: [
      //       {
      //         label: 'Contact info',
      //         id: 'contact_info',
      //         type: 'section',
      //         items: [
      //           { label: 'Phone', id: 'phone', type: 'phone', validation: { required: true, regex: /^(\+)?(?=.*\d)[\d|\s]+$/ } },
      //           { label: 'Email', id: 'email', type: 'email', validation: { required: true, regex: /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/ } }
      //         ]
      //       },
      //       {
      //         label: 'Other data',
      //         id: 'other_data',
      //         type: 'section',
      //         items: [
      //           {
      //             label: 'Other data 1',
      //             id: 'other_data_one',
      //             type: 'section',
      //             items: [
      //               { label: 'About the life', id: 'about_life', type: 'textarea' }
      //             ]
      //           },
      //           {
      //             label: 'Now the date.',
      //             id: 'the_date',
      //             type: 'section',
      //             items: [
      //               { label: 'The date', id: 'the_date_field', type: 'date' }
      //             ]
      //           }
      //         ]
      //       }
      //     ]
      //   }
      ]
    }
  },
  mutations: {
    data (state, data) {
      for (const key in data) {
        Vue.set(state.data, key, data[key])
      }
    },
    loaded (state, value) {
      state.loaded = value
    },
    loading (state, value) {
      state.loading = value
    }
  },
  actions: {
    async load ({ commit, state }) {
      if (state.loading) return
      commit('loading', true)

      try {
        const response = await Vue.services.config.get()
        if (!response.data) {
          throw new Error('data not found on response')
        }
        commit('data', response.data)
        commit('loaded', true)
      } catch (e) {
        commit('loaded', false)
        throw e
      } finally {
        commit('loading', false)
      }

      return true
    },
    async loadDestinations ({ commit, state, dispatch }) {
      const response = await Vue.services.config.destinations.get()
      if (!response.data) {
        throw new Error('data not found on response')
      }
      commit('data', { destinations: response.data || [] })
      if (state.data.destinations.length === 1) {
        dispatch('payload/setDestination', state.data.destinations[0].id, { root: true })
      }
    }
  },
  getters: {
    hasLoginProvider: state => {
      return state.data.login && state.data.login.loginProvider
    },
    destinationNameById: state => destinationId => (
      (state.data.destinations || [])
        .find(destination => destination.id === destinationId) || { name: destinationId }
    ).name
  }
}
