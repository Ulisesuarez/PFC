<template>
  <v-layout row wrap justify-center>
    <v-flex xs12 md12 xl10>
      <Card class="task-queue">
        <template slot="title">
          {{ $t('Tasks') }}
        </template>
        <template slot="content">
          <DataTable
            :headers="headers"
            service="taskSearch"
            :reload="false"
            :paginated-service="true"
            item-key="id"
          >
            <template slot="items" slot-scope="{ item }">
              <tr :class="{ 'deep-orange lighten-4': item.error && Object.keys(item.error).length > 0 }" class="task-row">
                <td>{{ studyName(item) }}</td>
                <td>{{date(item)}}</td>
                <td>{{ item.studyUID }}</td>
                <td>
                  {{ getStep(item) }}
                  <strong v-if="item.error && Object.keys(item.error).length > 0" class="red--text text--darken-4 d-block" >
                    {{ $t(item.error.code)  }}
                  </strong>
                </td>
                <td class="justify-end layout px-0">
                  <v-tooltip
                    bottom>
                    <v-btn
                      slot="activator"
                      color="grey darken-3"
                      icon
                      flat
                      @click="showDetails(item)">
                      <v-icon>info</v-icon>
                    </v-btn>
                  </v-tooltip>
                </td>
              </tr>
              <tr class="pa-0 progress-row" >
                <td colspan="100%" class="pa-0" >
                  <v-progress-linear
                    :value="100 * item.steps.currentStep / item.steps.totalSteps"
                    :color="progressColor(item)"
                    class="pa-0 ma-0"
                    height="3"
                  />
                </td>
              </tr>
            </template>
          </DataTable>
        </template>
      </Card>
    </v-flex>
  </v-layout>
</template>

<script>
import DataTable from '../components/common/DataTable'
import Card from '../components/common/Card'
import bus from '../assets/js/bus'

export default {
  name: 'Task',
  components: {
    DataTable,
    Card
  },
  data: () => {
    return {
      dialog: false,
      dialogTitle: 'output',
      dialogWidth: '1000px',
      output: {},
      showFrame: false
    }
  },
  computed: {
    headers () {
      return [
        {
          text: this.$t('study_name'),
          sortable: false,
          value: 'studyName',
          class: 'subheading'
        },
        {
          text: this.$t('date'),
          sortable: false,
          value: 'date',
          class: 'subheading'
        },
        {
          text: this.$t('study_UID'),
          sortable: false,
          value: 'studyUID',
          class: 'subheading'
        },
        {
          text: this.$t('step/Status'),
          sortable: false,
          value: 'step',
          class: 'subheading'
        },
        {
          text: this.$t('see_details'),
          sortable: false,
          value: 'details',
          class: 'subheading'
        }
      ]
    }
  },
  mounted () {
    this.$store.dispatch('navigation/setActions', [
      {
        text: 'navigation.back.files_selection',
        attrs: { color: 'grey lighten-2' },
        listeners: {
          click: () => {
            this.$router.push({ name: 'files-selection' })
          }
        }
      }
    ])
  },
  beforeMount () {
    bus.$on('taskSearch', () => {
      console.log(`IDONIA ${this.$store.state.session.token}`)
      this.$idoniapacs.defaults.headers.common['Authorization'] = `IDONIA ${this.$store.state.session.token}`
      this.$idoniapacs.get('tasks').then(
        (apiresponse) => {
          console.log(apiresponse)
          bus.$emit('searchDone', apiresponse)
        }).catch(e => {
        bus.$emit('searchError', e)
      })
    })
  },
  methods: {
    studyName (item) {
      return JSON.parse(item.aditionalFields).StudyName
    },
    date (item) {
      return new Date(item.createdAt).toLocaleDateString()
    },
    openInNewTab (output) {
      console.log('closing')
      this.dialog = false
      window.open(output.srcFrame, '_blank')
    },
    getStep (item) {
      if (item.steps.currentStepId) {
        if (item.steps.currentStep === item.steps.totalSteps && item.status === 'completed') {
          return this.$t('completed')
        } else {
          return this.$t(item.steps.currentStepId)
        }
      } else if (item.status) {
        return this.$t(item.status)
      }
    },
    progressColor (task) {
      if (task.error && Object.keys(task.error).length > 0) {
        return 'error'
      }

      if (task.status === 'completed') {
        return 'success'
      }

      return 'primary'
    },
    showDetails (item) {

    }
  }
}
</script>

<style lang="scss" scoped>
  .task-queue {
    /deep/ {
      .task-row {
        border-bottom: none !important;
      }
      .progress-row {
        border-top: none !important;
        td {
          height: 3px;
        }
      }
    }
  }

  .output {
    width: 100%;

  }
</style>
