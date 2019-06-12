<template>
  <v-form ref="form" lazy-validation>
    <v-container>
      <v-layout row wrap>
        <v-flex xs12>
          <v-select
            v-model="selectedPacs"
            :items="$store.state.config.data.pacs"
            item-value="id"
            item-text="name"
            return-object
            label="PACS"
          />
        </v-flex>
        <v-flex xs12 sm6>
          <v-text-field
            v-model="studyUID"
            :label="$t('study_id') "
            clearable
          />
        </v-flex>
        <v-flex xs12 sm6>
          <v-text-field
            v-model="patientId"
            :label="$t('patient_id') "
            clearable
          />
        </v-flex>
        <v-flex xs12 sm6>
          <v-text-field
            v-model="patientSurname"
            :label="$t('patient_surname') "
            clearable
          />
        </v-flex>

        <v-flex xs12 sm6>
          <v-text-field
            v-model="patientName"
            :label="$t('patient_name') "
            clearable
          />
        </v-flex>
        <v-flex xs12 sm6>
          <DatePickerField
            v-model="studyDateFrom"
            :label="$t('from')"
          />
        </v-flex>

        <v-flex xs12 sm6>
          <DatePickerField
            v-model="studyDateUntil"
            :label="$t('until')"
          />
        </v-flex>

        <v-flex xs12 class="text-xs-right">
          <v-btn
            :disabled="isValid"
            :loading="searching"
            class="primary"
            @click="$emit('search', getPayload())">
            {{ $t('search') }}
          </v-btn>
        </v-flex>

      </v-layout>
    </v-container>
  </v-form>
</template>

<script>
import DatePickerField from '../../common/DatePickerField'

export default {
  name: 'SelectionPacsSearch',
  components: { DatePickerField },
  props: {
    searching: {
      type: Boolean,
      required: true
    }
  },
  data () {
    return {
      selectedPacs: this.$store.state.config.data.pacs.length === 1
        ? this.$store.state.config.data.pacs[0]
        : null,
      patientId: '',
      patientName: '',
      patientSurname: '',
      studyUID: '',
      studyDateFrom: this.$moment().subtract(1, 'week'),
      studyDateUntil: this.$moment()
    }
  },
  computed: {
    fullName () {
      let name = this.patientSurname
      if (this.patientName !== '') {
        name = name + '^' + this.patientName
      }
      return name
    },
    studyDate () {
      if (this.studyDateUntil && this.studyDateFrom) {
        return this.studyDateFrom.format('YYYYMMDD') + '-' + this.studyDateUntil.format('YYYYMMDD')
      } else if (this.studyDateUntil) {
        return '-' + this.studyDateUntil.format('YYYYMMDD')
      } else if (this.studyDateFrom) {
        return this.studyDateFrom.format('YYYYMMDD') + '-'
      }
      return ''
    },
    isValid () {
      return this.$refs.form ? this.$refs.form.validate() : false
    }
  },
  validations: {
    selectedPacs: {
    // required
    },
    studyDateFrom: {
      // isBefore: date.isBefore('studyDateUntil')
    },
    studyDateUntil: {
      // isAfter: date.isAfter('studyDateFrom')
    }
  },
  methods: {
    getPayload () {
      const filters = {}
      const fields = {
        '00100020': this.patientId,
        '00100010': this.fullName,
        '0020000D': this.studyUID,
        '00080020': this.studyDate
      }

      for (const fieldName in fields) {
        let filterName = fieldName
        let value = fields[fieldName]
        if (value) {
          filters[filterName] = value
        }
      }

      return {
        aet: this.selectedPacs.name,
        filters
      }
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
