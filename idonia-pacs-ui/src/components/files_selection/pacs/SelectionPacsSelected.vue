<template>
  <div>
    <v-list two-line subheader>
      <v-subheader>{{ $t('selected_studies') }}</v-subheader>
      <transition-group name="slide-y-transition" >
        <template v-for="(study, index) in studies">
          <v-list-tile
            :key="study.StudyInstanceUID"
            xs12>
            <v-list-tile-content>
              <v-list-tile-title>{{ study.PatientID }} | {{ study.PatientName }}</v-list-tile-title>
              <v-list-tile-sub-title>{{ $moment(study.dateTime).format(localeData.longDateFormat('L')) }} - {{ study.StudyDescription }}</v-list-tile-sub-title>
            </v-list-tile-content>
            <v-list-tile-action>
              <v-btn icon ripple @click="$store.dispatch('payload/removeStudy', study)">
                <v-icon color="grey lighten-1">close</v-icon>
              </v-btn>
            </v-list-tile-action>
          </v-list-tile>
          <v-divider
            v-if="index + 1 < studies.length"
            :key="study.StudyInstanceUID + '-divider'"
          />
        </template>
      </transition-group>
    </v-list>
  </div>
</template>

<script>
import bus from '../../../assets/js/bus'
export default {
  name: 'SelectionPacsSelected',
  computed: {
    studies () {
      return Object.keys(this.$store.state.payload.data.studies)
        .map(key => this.$store.state.payload.data.studies[key])
    },
    hasStudiesOfMultiplePatients () {
      return Object.keys(
        this.studies.reduce((patientIds, study) => {
          patientIds[study.PatientID] = true
          return patientIds
        }, {})
      ).length > 1
    },
    localeData () {
      return this.$moment.localeData(this.$store.state.session.locale.shortCode)
    }
  },
  watch: {
    hasStudiesOfMultiplePatients (val) {
      bus.$emit('alert', {
        id: 'alertOfMultiplePatients',
        value: val,
        type: 'warning',
        text: this.$tname('studies_of_multiple_patients_warning')
      })
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
