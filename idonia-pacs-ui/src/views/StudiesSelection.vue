<template>
  <v-layout
    row
    wrap
    justify-center>
    <v-flex xs12 lg8>
      <v-flex xs12>
        <SelectionPacs />
      </v-flex>
      <v-flex xs12>
      </v-flex>
    </v-flex>
  </v-layout>
</template>

<script>
import SelectionPacsItems from '../components/files_selection/pacs/SelectionPacs'
import bus from '../assets/js/bus'
export default {
  name: 'StudiesSelection',
  components: {
    SelectionPacs: SelectionPacsItems
  },
  computed: {
    hasPacsData () {
      return Array.isArray(this.$store.state.config.data.pacs) &&
        this.$store.state.config.data.pacs.length > 0
    },
    hasFiles () {
      return this.$store.state.payload.data.dirfiles.length !== 0 || this.$store.state.payload.data.files.length !== 0
    },
    noStudies () {
      return Object.keys(this.$store.state.payload.data.studies).length === 0 &&
      this.$store.state.payload.data.dirfiles.findIndex(file => file.isDicom) === -1 &&
      (this.$store.state.payload.data.dirfiles.length !== 0 || this.$store.state.payload.data.files.length !== 0)
    }
  },
  watch: {
    noStudies (val) {
      bus.$emit('alert', {
        id: 'NoStudiesinSeletion',
        value: val,
        type: 'warning',
        text: this.$tname('no_studies_warning')
      })
    }
  },
  mounted () {
    const nextIsDisabled = () =>
      Object.keys(this.$store.state.payload.data.studies).length === 0 &&
      Object.keys(this.$store.state.payload.data.dirfiles).length === 0 &&
      Object.keys(this.$store.state.payload.data.files).length === 0
    const justOneDestinationWithoutAdditionalData = (
      this.$store.state.config.data.destinations.length === 1 &&
      (this.$store.state.config.data.destinations[0].additionalData || []).length === 0
    )
    const theSingleDestination = this.$store.state.config.data.destinations[0]

    this.$store.dispatch('navigation/setActions', [
      {
        i18nParams: justOneDestinationWithoutAdditionalData
          ? { destinationName: theSingleDestination.name }
          : null,
        text: () => nextIsDisabled()
          ? this.hasPacsData
            ? 'navigation.disabled.no_study_selected'
            : 'navigation.disabled.no_file_added'
          : justOneDestinationWithoutAdditionalData
            ? 'navigation.action.send'
            : 'common.continue',
        icon: () => nextIsDisabled()
          ? null
          : justOneDestinationWithoutAdditionalData
            ? { name: 'send', attrs: { right: true } }
            : { name: 'navigate_next', attrs: { right: true } },
        attrs: {
          color: 'primary',
          disabled: () => nextIsDisabled(),
          loading: () => this.$store.state.payload.sending
        },
        listeners: {
          click: async () => {
            if (!nextIsDisabled()) {
              if (justOneDestinationWithoutAdditionalData) {
                try {
                  await this.$store.dispatch('payload/send')
                  this.$store.dispatch('navigation/clear')
                  this.$router.push({ name: 'task-queue' })
                } catch (e) {
                  if (e.response.status !== 401) {
                    this.$notify({
                      title: this.$t('notify.error_sending_payload.title'),
                      text: this.$t('notify.error_sending_payload.text'),
                      type: 'error'
                    })
                  }
                }
              } else {
                this.$router.push({ name: 'destination' })
              }
            }
          }
        }
      }
    ])
  }
}
</script>

<style lang="scss" scoped>
</style>
