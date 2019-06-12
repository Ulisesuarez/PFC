import Vue from 'vue'

export default {
  namespaced: true,
  state: {
    sending: false,
    data: {
      studies: {},
      dirfiles: [],
      files: [],
      reports: [],
      destination: null,
      additionalData: {}
    }
  },
  mutations: {
    addStudy (state, study) {
      Vue.set(state.data.studies, study.StudyInstanceUID, study)
    },
    removeStudy (state, study) {
      Vue.delete(state.data.studies, study.StudyInstanceUID)
    },
    studies (state, studies) {
      state.data.studies = studies
    },
    files (state, files) {
      state.data.files = files
    },
    dirfiles (state, files) {
      state.data.dirfiles = files
    },
    addFiles (state, files) {
      for (let i = 0; i < files.length; i++) {
        Vue.set(state.data.files, state.data.files.length, files[i])
      }
    },
    deleteFile (state, uuid) {
      Vue.delete(state.data.files, state.data.files.findIndex(file => file.uuid === uuid))
    },
    addDirFiles (state, files) {
      for (let i = 0; i < files.length; i++) {
        Vue.set(state.data.dirfiles, state.data.dirfiles.length, files[i])
      }
    },
    deleteDirFile (state, uuid) {
      Vue.delete(state.data.dirfiles, state.data.dirfiles.findIndex(file => file.uuid === uuid))
    },
    reports (state, reports) {
      state.data.reports = reports
    },
    destination (state, destination) {
      state.data.destination = destination
    },
    additionalData (state, additionalData) {
      state.data.additionalData = additionalData
    },
    sending (state, value) {
      state.sending = value
    }
  },
  actions: {
    addStudy ({ commit }, study) {
      commit('addStudy', study)
    },
    removeStudy ({ commit }, study) {
      commit('removeStudy', study)
    },
    setStudies ({ commit }, studies) {
      commit('studies', studies)
    },
    setFiles ({ commit }, files) {
      commit('files', files)
    },
    addFiles ({ commit }, files) {
      commit('addFiles', files)
    },
    deleteFile ({ commit }, uuid) {
      commit('deleteFile', uuid)
    },
    setDirFiles ({ commit }, files) {
      commit('dirfiles', files)
    },
    addDirFiles ({ commit }, files) {
      commit('addDirFiles', files)
    },
    deleteDirFile ({ commit }, uuid) {
      commit('deleteDirFile', uuid)
    },
    setReports ({ commit }, reports) {
      commit('reports', reports)
    },
    setDestination ({ commit }, destination) {
      commit('destination', destination)
    },
    addAdditionalData ({ commit, state }, { field, value }) {
      commit('additionalData', {
        ...state.data.additionalData,
        [field]: value
      })
    },
    isSending ({ commit }, value) {
      commit('sending', value)
    },
    clear ({ dispatch, commit, rootState }) {
      dispatch('setStudies', {})
      dispatch('setFiles', [])
      dispatch('setReports', [])
      if (rootState.config.data.destinations.length > 1) {
        dispatch('setDestination', null)
      }
      commit('additionalData', {})
    },
    async send ({ dispatch, state }) {
      dispatch('isSending', true)
      try {
        const payload = new FormData()
        Object.keys(state.data.studies).forEach(key => {
          let description = {
            patientName: state.data.studies[key].PatientName,
            patientId: state.data.studies[key].PatientID,
            modalitiesInStudy: state.data.studies[key].ModalitiesInStudy,
            institutionName: state.data.studies[key].InstitutionName || '',
            studyDate: state.data.studies[key].StudyDate,
            studyTime: state.data.studies[key].StudyTime,
            studyDescription: state.data.studies[key].StudyDescription || ''
          }
          payload.append('studies[]', JSON.stringify({
            id: state.data.studies[key].StudyInstanceUID,
            pacs: state.data.studies[key].pacs,
            description: description
          }))
        })

        state.data.dirfiles.forEach(file => {
          if (file.isDicom) {
            payload.append('dicoms[]', file.uuid)
          }
          payload.append(`files[${file.uuid}]`, file)
        })
        state.data.files.forEach(file => {
          payload.append(`files[${file.uuid}]`, file)
        })

        state.data.reports.forEach(report => {
          payload.append('reports[]', report)
        })

        if (Object.keys(state.data.additionalData).length) {
          payload.append('additionalData', JSON.stringify(state.data.additionalData))
        }

        payload.append('destination', state.data.destination)
        console.log(payload)
        await Vue.services.tasks.add(payload)
        dispatch('clear')
      } finally {
        dispatch('isSending', false)
      }
    }
  },
  getters: {
    hasStudiesSelected: state => Object.keys(state.data.studies).length > 0,
    destination: (state, _, rootState) => state.data.destination
      ? rootState.config.data.destinations.find(item => item.id === state.data.destination)
      : null
  }
}
