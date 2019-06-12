<template>
  <Card icon="search" >
    <template slot="title">
      {{ $t('Select studies') }}
    </template>
    <template slot="content">
      <SelectionPacsSearch
        ref="pacsSearch"
        :searching="searching"
        @search="catchPayload"
      />

      <v-layout row wrap>
        <v-flex xs12>
          <transition name="slide-y-transition" tag="div">
            <SelectionPacsItems
              v-if="showItems"
              :selected-pacs="$refs.pacsSearch.selectedPacs.name"
              :pacs-search-payload="selectionPacsSearchPayload"
              :hidden="!hasSearched"
              @has-searched="hasSearched = $event"
              @searching="searching = $event"
            />
          </transition>
        </v-flex>
      </v-layout>
      <template v-if="$store.getters['payload/hasStudiesSelected']" >
        <br>
        <v-divider />
        <br>
        <SelectionPacsSelected />
      </template>
    </template>
  </Card>
</template>

<script>
import SelectionPacsSelected from './SelectionPacsSelected'
import SelectionPacsSearch from './SelectionPacsSearch'
import SelectionPacsItems from './SelectionPacsItems'
import Card from '../../common/Card'

export default {
  name: 'SelectionPacs',
  components: {
    SelectionPacsSelected,
    SelectionPacsSearch,
    SelectionPacsItems,
    Card
  },
  computed: {
    showItems () {
      return this.pacsSearch && this.pacsSearch.selectedPacs !== null && this.selectionPacsSearchPayload !== null
    }
  },
  watch: {
    showItems (val) {
      console.log('showItems')
      console.log(val)
    },
    pacsSearch (val) {
      console.log('refs pacsSearch')
      console.log(val)
    }
  },
  data: () => ({
    selectionPacsSearchPayload: null,
    searching: false,
    hasSearched: false,
    pacsSearch: null
  }),
  mounted () {
    this.pacsSearch = this.$refs.pacsSearch
  },
  methods: {
    catchPayload (event) {
      console.log(event)
      console.log(this.pacsSearch)
      console.log(this.$refs)
      this.selectionPacsSearchPayload = event
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
