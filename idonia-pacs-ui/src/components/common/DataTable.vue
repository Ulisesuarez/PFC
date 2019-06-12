<template>
  <div class="data-table" >
    <div class="table">
      <v-data-table
        v-bind="$attrs"
        :items="items"
        :loading="searching"
        :pagination.sync="pagination"
        :total-items="paginatedService ? totalItems : null"
        v-on="$listeners">
        <template slot="headers" slot-scope="props">
          <tr>
<!--            <th>
              <v-checkbox
                :input-value="props.all"
                :indeterminate="props.indeterminate"
                primary
                hide-details
                @click.stop="toggleAll"
              ></v-checkbox>
            </th>-->
            <th
              v-for="header in props.headers"
              :key="header.text"
              :class="['column sortable', pagination.descending ? 'desc' : 'asc', header.value === pagination.sortBy ? 'active' : '']"
              @click="changeSort(header.value)"
            >
              <v-icon small>arrow_upward</v-icon>
              {{ header.text }}
            </th>
          </tr>
        </template>
        <v-progress-linear
          color="primary"
          size="3"
          indeterminate
          style="height:5px !important;"
        />
        <template slot="items" slot-scope="scope">
          <slot v-bind="scope" name="items" />
        </template>
      </v-data-table>
    </div>
    <div v-if="pages>1" class="pagination">
      <v-pagination
        v-model="pagination.page"
        :disabled="searching"
        :total-visible="6"
        :length="pages"
      />
    </div>
  </div>
</template>

<script>
import bus from '../../assets/js/bus'

export default {
  name: 'DataTable',
  props: {
    payload: {
      type: Object,
      default: () => ({})
    },
    service: {
      type: String,
      required: true
    },
    sortDefault: {
      type: String,
      default: ''
    },
    paginatedService: {
      type: Boolean,
      default: false
    },
    mapItem: {
      type: Function,
      default: item => item
    },
    reload: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      items: [],
      searching: false,
      totalItems: 0,
      pagination: {
        sortBy: this.sortDefault
      },
      id: null
    }
  },
  computed: {
    pages () {
      const totalItems = this.paginatedService ? this.totalItems : this.pagination
        ? this.pagination.totalItems : 0
      if (this.pagination.rowsPerPage == null || totalItems == null) {
        return 0
      }
      console.log('totalItems', totalItems)
      console.log(Math.ceil(totalItems / this.pagination.rowsPerPage))
      return Math.ceil(totalItems / this.pagination.rowsPerPage)
    },
    reloadItems () {
      return this.reload
    }
  },
  watch: {
    pagination: {
      deep: true,
      async handler (value) {
        if (this.paginatedService) {
          console.log('holaaa')
          this.search(value)
        }
      }
    },
    searching (value) {
      this.$emit('searching', value)
    },
    service (val) {
      console.log('watcher service', val)
      this.search({ ...this.pagination, page: 1 })
    },
    payload: {
      deep: true,
      handler (val) {
        console.log('watcher')
        console.log(val)
        this.search({ ...this.pagination, page: 1 })
      }
    },
    reloadItems (val) {
      if (!val) {
        clearInterval(this.id)
      }
    }
  },
  mounted () {
    console.log(typeof this.service === 'undefined')
    console.log('mounted', this.service)
    let self = this
    self.search({ ...self.pagination, page: self.pagination.page })
    bus.$on('searchDone', (response) => {
      console.log(response)
      if (this.paginatedService) {
        this.items = (response.data.items || []).map(this.mapItem)
        this.totalItems = response.data.pagination.total
      } else {
        console.log('yp vpy despues', response.data)
        this.items = (response.data || []).map(this.mapItem)
        this.pagination.totalItems = this.items.length
        this.pagination.page = 1
      }
      self.searching = false
      self.$emit('has-searched', true)
    })
    bus.$on('searchError', (e) => {
      console.log(e)
      if (e.response && e.response.status !== 401) {
        this.$notify({
          title: this.$t('notify.exception_search.title'),
          text: e.error || this.$t('notify.exception_search.text'),
          type: 'error'
        })
      }
      self.searching = false
      self.$emit('has-searched', true)
    })
    if (this.reloadItems) {
      this.id = setInterval(function () {
        self.search({ ...self.pagination, page: self.pagination.page })
      }
      , 2000)
    }
  },
  beforeDestroy () {
    bus.$off('searchDone')
    bus.$off('searchError')
    clearInterval(this.id)
  },
  methods: {
    search (pagination) {
      this.searching = true
      console.log(this.service, 'searching?¿?', pagination)
      bus.$emit(this.service, { pagination, payload: this.payload })
      console.log('se emite o no')
    },
    toggleAll () {
      if (this.selected.length) this.selected = []
      else this.selected = this.items.slice()
    },
    changeSort (column) {
      if (this.pagination.sortBy === column) {
        this.pagination.descending = !this.pagination.descending
      } else {
        this.pagination.sortBy = column
        this.pagination.descending = false
      }
    }
  }
}
</script>

<style lang="scss" scoped>
  .data-table {
    /deep/ {
      .v-datatable__actions__range-controls {
        display: none;
      }
      .v-input__control {
        justify-content: center;
        align-items: center;
        flex-flow: row nowrap;
      }
      .v-datatable__progress {
        height: 5px !important;
      }
    }

    .pagination {
      text-align: center;
    }
  }
  </style>
  <style>
    .v-progress-linear{
      height: 5px !important;
    }
  .v-progress-linear__bar{
    height: 5px !important;
  }
  .v-progress-linear__background{
    height: 5px !important;
  }
  .v-progress-linear__bar__indeterminate{
    height: 5px !important;
  }
</style>
