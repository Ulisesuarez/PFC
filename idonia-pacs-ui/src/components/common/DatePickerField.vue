<template>
  <v-menu
    ref="menu"
    v-model="menu"
    :close-on-content-click="false"
    :nudge-right="40"
    :z-index="500"
    :return-value.sync="returnable"
    content-class="date-picker-field"
    transition="scale-transition"
    min-width="290px"
    full-width
    offset-y
    lazy>
    <v-text-field
      slot="activator"
      :value="localeStringDate"
      :label="label"
      :append-icon="value === null ? null : 'close'"
      v-bind="$attrs"
      type="text"
      prepend-icon="event"
      readonly
      @input="$emit('input', this.$moment($event))"
      @click:append="() => { $emit('input', null) }"
    />
    <v-date-picker
      :value="isoString"
      :locale="$store.state.session.locale.code"
      :first-day-of-week="1"
      no-title
      scrollable
      @input="$refs.menu.save($moment($event))"
    />
  </v-menu>
</template>

<script>
export default {
  name: 'DatePickerField',
  props: {
    label: {
      type: String,
      default: 'Date'
    },
    value: {
      type: [Object, Date],
      default: () => null
    }
  },
  data: () => ({
    menu: false
  }),
  computed: {
    date () {
      return this.value === null ? null : this.$moment(this.value)
    },
    localeStringDate () {
      if (this.value === null) {
        return null
      }
      return this.date.format('DD/MM/YYYY')
    },
    isoString () {
      return this.value === null
        ? null
        : this.value.toISOString(true).split('T')[0]
    },
    returnable: {
      get () {
        return this.date
      },
      set (value) {
        this.$emit('input', value)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
  .date-picker-field {
    /deep/ .v-picker__body {
      padding-bottom: 12px;
    }
  }
</style>
