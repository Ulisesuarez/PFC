<template>
  <div class="data-table">
    <DataTable
      v-show="!hidden"
      :headers="headers"
      service="pacsSearch"
      :paginated-service="true"
      :payload="pacsSearchPayload"
      sort-default="StudyDate"
      :rows-per-page-items="[ 25, 50, 100, { 'text': '$vuetify.dataIterator.rowsPerPageAll', 'value': -1 } ]"
      item-key="id"
      class="elevation-4 pb-4"
      select-all
      @searching="$emit('searching', $event)"
      @has-searched="$emit('has-searched', $event)">
      <template slot="items" slot-scope="{ item }">
        <tr>
          <td class="px-1">
            <v-checkbox
              :input-value="isStudyActive(item)"
              hide-details
              primary
            />
          </td>
          <td class="justify-center layout align-center px-0">
            <upload-btn icon
                        flat
                        :custom-prop="item"
                        :hover="false"
                        :ripple="false"
                        small
                        maxWidth="15px"
                        class="mr-2"
                        v-if="!item.isReported"
                        style="padding: 0px; margin: 0px; color:#646464"
                        @file-update="reportStudy">
              <template slot="icon">
                  <v-icon>attach_file</v-icon>
                </template>
            </upload-btn>
            <v-icon :key="action.ID" v-for="action in composeActions(item)"
              medium
              class="mr-2"
              @click="executeAction(action.ID, item)"
            >
              {{ action.icon }}
            </v-icon>
          </td>
          <td>{{ patientIds(item) }}</td>
          <td>{{ patientNames(item) }}</td>
          <td>{{ studyDescriptions(item)}}</td>
          <td>{{ studyDates(item) }}</td>
          <td>{{ numberOfSeriesRelatedToStudy(item)}}</td>
        </tr>
      </template>
    </DataTable>
    <transfer
      v-if="transferDialog"
      :dialog="transferDialog"
      :study-uid="studyUID(activeItem)"></transfer>
    <share
      v-if="shareDialog"
      :dialog="shareDialog"
      :study="activeItem"></share>
  </div>
</template>

<script>
import DataTable from '../../common/DataTable'
import bus from '../../../assets/js/bus'
import tags from '../../../assets/js/tags.json'
import UploadButton from '../../common/CustomUploadButton'
import TransferDialog from '../../common/TransferDialog'
import ShareDialog from '../../common/ShareDialog'
export default {
  name: 'SelectionPacsItems',
  components: { DataTable, 'upload-btn': UploadButton, 'transfer': TransferDialog, 'share': ShareDialog },
  props: {
    selectedPacs: {
      type: String,
      required: true
    },
    pacsSearchPayload: {
      type: Object,
      required: true
    },
    hidden: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    headers () {
      return [
        {
          text: this.$t('actions'),
          sortable: false,
          value: 'actions',
          class: 'subheading'
        },
        {
          text: this.$t('patient_id'),
          sortable: true,
          value: 'PatientID',
          class: 'subheading'
        },
        {
          text: this.$t('name'),
          sortable: true,
          value: 'PatientName',
          class: 'subheading'
        },
        {
          text: this.$t('description'),
          sortable: false,
          value: 'StudyDescription',
          class: 'subheading'
        },
        {
          text: this.$t('date'),
          sortable: true,
          value: 'StudyDate',
          class: 'subheading'
        },
        {
          text: this.$t('series'),
          sortable: false,
          value: 'NumberOfStudyRelatedSeries',
          class: 'subheading'
        }
      ]
    },
    localeData () {
      return this.$moment.localeData(this.$store.state.session.locale.shortCode)
    }
  },
  beforeCreate () {
    bus.$on('closeTransfer', () => {
      this.transferDialog = false
    })
    bus.$on('closeShare', () => {
      this.shareDialog = false
    })
    console.log('-----------SelectionPacsItems----------------------------------')
    bus.$on('pacsSearch', data => {
      console.log('datatable this pagination', data.pagination)
      let aet = data.payload.aet
      console.log('aets/' + aet + '/rs/studies?includefield=all' +
        this.makeQuery(data.payload.filters, data.pagination))
      this.$axios.get('aets/' + aet + '/rs/studies?includefield=all' +
        this.makeQuery(data.payload.filters, data.pagination)).then(response => {
        let result = {
          data: {
            items: {},
            pagination: {
              total: 0 }
          }
        }
        let completedData
        let studyUIDs = response.data.map(study => this.studyUID(study))
        console.log(this.$idoniapacs)
        this.$idoniapacs.post('studies', studyUIDs).then(
          (apiresponse) => {
            console.log('-----------------')
            if (apiresponse.data && apiresponse.data.length > 0) {
              completedData = response.data.map((study) => {
                for (let i = 0; i < apiresponse.data.length; i++) {
                  console.log('APIRESPONSE', apiresponse.data[i])
                  console.log('study', study)
                  console.log(apiresponse.data[i].studyUID, this.studyUID(study))
                  if (apiresponse.data[i].studyUID === this.studyUID(study)) {
                    console.log(study)
                    study.isUploaded = true
                    study.isReported = apiresponse.data[i].isReported
                    study.reportID = apiresponse.data[i].reportID
                    study.studyID = apiresponse.data[i].studyID
                    study.magicLink = apiresponse.data[i].magicLink
                    console.log(study)
                    return study
                  }
                }
                return study
              })
              console.log('CAMBIO¿?')
              console.log(response.data)
            } else {
              completedData = response.data
            }
            result.data.items = completedData
            result.data.pagination.total = result.data.items.length
            bus.$emit('searchDone', result)
          }
        ).catch((e) => {
          console.log(e)
        })
        console.log(studyUIDs)
      }).catch(e => {
        bus.$emit('searchError', e)
      })
    })
  },
  beforeDestroy () {
    // bus.$off('pacsSearch')
  },
  data () {
    return {
      tags: tags,
      file: null,
      transferDialog: false,
      shareDialog: false,
      activeItem: null
    }
  },
  methods: {
    showTransferDialog (item) {
      this.transferDialog = true
      this.activeItem = item
    },
    showShareDialog (item) {
      this.shareDialog = true
      this.activeItem = item
    },
    getMagicLink (item) {
      console.log(this.studyUID(item))
      this.$idoniapacs.post('task', {
        studyUID: this.studyUID(item),
        sopObjects: [],
        studyName: 'testnumero1002',
        containerName: 'toDO',
        steps: ['GetMagicLink']
      }).then(response => {
        console.log(response)
      }).catch((error) => {
        console.log('error')
        console.log(error)
      })
    },
    executeAction (action, item) {
      switch (action) {
        case 'viewStudy':
          this.viewStudy(item)
          break
        case 'sendToIdonia':
          this.sendToIdonia(item)
          break
        case 'openInIdonia':
          this.openInIdonia(item)
          break
        case 'GetMagicLink':
          console.log(item)
          if (item.magicLink) {
            this.showShareDialog(item)
          } else {
            this.getMagicLink(item)
          }
          break
        case 'transfer':
          this.showTransferDialog(item)
          break
      }
    },
    composeActions (item) {
      let actions = [{
        ID: 'viewStudy',
        icon: 'visibility'
      }]
      if (item.isUploaded) {
        actions.push({
          ID: 'openInIdonia',
          icon: 'folder_open'
        })
      } else {
        actions.push({
          ID: 'sendToIdonia',
          icon: 'cloud_upload'
        })
      }
      actions.push({
        ID: 'GetMagicLink',
        icon: 'link'
      })
      actions.push({
        ID: 'transfer',
        icon: 'send'
      })
      return actions
    },
    isStudyActive (study) {
      return typeof this.$store.state.payload.data.studies[study.StudyInstanceUID] !== 'undefined'
    },
    toggleStudy (study) {
      if (this.isStudyActive(study)) {
        this.$store.dispatch('payload/removeStudy', study)
      } else {
        study.pacs = this.selectedPacs
        this.$store.dispatch('payload/addStudy', study)
      }
    },
    makeQuery (fields, pagination) {
      let offset = '0'
      let limit = pagination.rowsPerPage ? pagination.rowsPerPage.toString() : '25'
      if (pagination.page && pagination.rowsPerPage) {
        offset = pagination.page === 0 ? 0 : (pagination.page - 1) * pagination.rowsPerPage
        limit = pagination.page * pagination.rowsPerPage
      }
      let order = 'StudyDate'
      if (pagination.sortBy && pagination.sortBy !== order) {
        order = pagination.sortBy + ',-StudyDate'
      }
      if (pagination.descending) {
        order = '-' + order
      }
      let query = '&offset=' + offset + '&limit=' + limit + '&orderby=' + order + '&'
      // let oldQuery = '&offset=0&limit=25&orderby=-StudyDate,-StudyTime&'
      for (let field in fields) {
        if (fields.hasOwnProperty(field)) {
          console.log(field)
          query = query + field + '=' + fields[field] + '&'
        }
      }
      return query.substring(0, query.length - 1)
    },
    patientIds (item) {
      if (item['00100020'] && item['00100020'].Value) {
        return item['00100020'].Value.join(', ')
      }
      return ''
    },
    patientNames (item) {
      if (item['00100010'] && item['00100010'].Value) {
        return item['00100010'].Value.map(a => a.Alphabetic).join(', ')
      }
      return ''
    },
    studyDescriptions (item) {
      if (item['00081030'] && item['00081030'].Value) {
        return item['00081030'].Value
      }
      return ''
    },
    studyDates (item) {
      if (item['00080020'] && item['00080020'].Value) {
        return item['00080020'].Value.map(a => {
          let rs = a.replace(/[^0-9]/g, '')
          return new Date(rs.substring(0, 4),
            rs.substring(4, 6) - 1,
            rs.substring(6, 8)).toLocaleDateString()
        }).join(', ')
      }
      return ''
    },
    numberOfSeriesRelatedToStudy (item) {
      if (item && item['00201206'] && item['00201206'].Value) {
        return item['00201206'].Value.join(', ')
      }
      return ''
    },
    studyUID (item) {
      if (item && item['0020000D'] && item['0020000D'].Value.length > 0) {
        return item['0020000D'].Value[0]
      }
      return false
    },
    seriesUID (item) {
      if (item && item['0020000E'] && item['0020000E'].Value.length > 0) {
        return item['0020000E'].Value[0]
      }
      return false
    },
    sopUID (item) {
      if (item && item['00080018'] && item['00080018'].Value.length > 0) {
        return item['00080018'].Value[0]
      }
      return false
    },
    viewStudy (item) {
      let studyUID = this.studyUID(item)
      if (studyUID) {
        // TODO change url config field
        window.open('http://localhost:3000/oviyam?serverName=arc&studyUID=' + studyUID, '_blank')
      } else {

      }
    },
    openInIdonia (item) {
      if (item.studyID) {
        // TODO change url config field
        window.open('https://staging.idonia.com/file/' + item.studyID, '_blank')
      } else {

      }
    },
    reportStudy (item) {
      let self = this
      let reader = new FileReader()
      reader.onload = function () {
        let arrayBuffer = this.result
        let array = new Uint8Array(arrayBuffer)
        let binaryString = String.fromCharCode.apply(null, array)
        console.log(binaryString)
        self.$idoniapacs.post('task', {
          studyUID: self.studyUID(item.customProp),
          file: binaryString,
          steps: ['ReportStudy']
        }).then(response => {
          console.log(response)
        }).catch((error) => {
          console.log('error')
          console.log(error)
        })
      }
      reader.readAsArrayBuffer(item.file)
    },
    sendToIdonia (item) {
      console.log('send to idonia')
      this.$axios.get('http://localhost:8080/dcm4chee-arc/aets/DCM4CHEE/rs/studies/' +
        this.studyUID(item) + '/metadata')
        .then(response => {
          console.log(response)
          let identifiers = []
          if (response.data) {
            response.data.forEach(instance => {
              if (this.studyUID(instance) && this.seriesUID(instance) && this.sopUID(instance)) {
                identifiers.push({
                  studyUID: this.studyUID(instance),
                  sopInstanceUID: this.sopUID(instance),
                  seriesInstanceUID: this.seriesUID(instance)
                })
              }
            })
            console.log(identifiers)
            console.log('despues de que se complete')
            this.$idoniapacs.post('task', {
              studyUID: this.studyUID(item),
              sopObjects: identifiers,
              studyName: 'testnumero1000',
              containerName: 'toDO',
              steps: ['uploadToIdonia']
            }).then(response => {
              console.log(response)
            }).catch((error) => {
              console.log('error')
              console.log(error)
            })
          }
        }, (reject) => {
          console.log('rejected')
          console.error(reject)
        }).catch((error) => {
          console.log('error')
          console.log(error)
        })
    }
  }
}
</script>

<style lang="scss" scoped>
  .data-table /deep/ {
    thead .v-input--checkbox {
      display: none;
    }
  }
</style>
