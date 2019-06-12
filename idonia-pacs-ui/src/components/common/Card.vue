<template>
  <v-card class="card" >
    <v-toolbar
      v-ripple="!disabledActions && value === false"
      :class="{ clickable: !disabledActions && value === false }"
      :color="color"
      dark
      flat
      @click="disabledActions ? null : $emit('input', true)">
      <transition name="slide-y-reverse-transition">
        <slot name="icon">
          <v-icon v-if="icon" v-text="icon" />
        </slot>
      </transition>
      <v-toolbar-title>
        <slot name="title" />
      </v-toolbar-title>
      <v-spacer />
      <transition name="fade-transition">
        <v-btn
          v-if="!disabledActions && value === false"
          icon
          @click="$emit('input', true)" >
          <v-icon>edit</v-icon>
        </v-btn>
      </transition>
    </v-toolbar>
    <div
      ref="collapsable"
      :class="{
        collapsable: value !== null,
        closed: !active
      }"
      :style="{
        maxHeight: collapsableHeight,
        transition: value === null ? null : `max-height ${collapsableTransitionTime}ms ease-in-out`,
        willChange: value === null ? null : 'max-height'
      }"
    >
      <v-progress-linear
        :active="loading"
        height="5"
        class="pa-0 ma-0"
        indeterminate
      />
      <v-card-text>
        <slot name="content" />
      </v-card-text>
      <v-card-actions>
        <slot name="actions" />
      </v-card-actions>
    </div>
  </v-card>
</template>

<script>
import ResizeSensor from 'css-element-queries/src/ResizeSensor'

export default {
  name: 'Card',
  props: {
    color: {
      type: String,
      default: 'primary'
    },
    loading: {
      type: Boolean,
      default: false
    },
    icon: {
      type: String,
      default: ''
    },
    disabledActions: {
      type: Boolean,
      default: false
    },
    value: {
      type: [Boolean, Object, String],
      default: () => null
    }
  },
  data: () => ({
    resizeSensor: null,
    collapsableHeight: null,
    collapsableTransitionTime: 400,
    timeoutID: null
  }),
  computed: {
    active () {
      return this.value !== null && this.value
    }
  },
  watch: {
    value (value, oldValue) {
      if (value !== null && oldValue === null) {
        this.bindResizeEvents()
      }
    }
  },
  mounted () {
    this.bindResizeEvents()
  },
  destroyed () {
    if (this.resizeSensor !== null) {
      this.resizeSensor = null
      delete this.resizeSensor
      window.removeEventListener('resize', this.storeCollapsableHeight)
    }
  },
  methods: {
    bindResizeEvents () {
      if (this.resizeSensor === null && this.value !== null) {
        this.resizeSensor = new ResizeSensor(this.$refs.collapsable, this.storeCollapsableHeight)
        window.addEventListener('resize', this.storeCollapsableHeight)
        this.storeCollapsableHeight()
      }
    },
    async storeCollapsableHeight (event) {
      clearTimeout(this.timeoutID)
      if (this.value === true) {
        this.timeoutID = setTimeout(() => {
          if (this.$refs && this.$refs.collapsable) {
            let initialValue = 0
            let totalHeight = Array.prototype.slice.call(this.$refs.collapsable.childNodes)
              .slice(0, Array.prototype.slice.call(this.$refs.collapsable.childNodes).length - 1)
              .reduce((accumulator, currentNode) => {
                let style = window.getComputedStyle(currentNode)
                return accumulator + parseInt(currentNode.offsetHeight) + parseInt(style.marginBottom) + parseInt(style.marginTop)
              }, initialValue)
            this.collapsableHeight = (totalHeight || this.$refs.collapsable.offsetHeight) + 'px'
          }
        }, this.collapsableTransitionTime + 50)
      } else if (event && event.type === 'resize' && this.value === false) {
        this.collapsableHeight = null
      }
    }
  }
}
</script>

<style lang="scss">
  .v-dialog > .card {
    margin-bottom: 0 !important;
  }
</style>

<style lang="scss" scoped>
  .card {
    margin-bottom: 1rem;

    .slide-y-reverse-transition-leave-active {
      position: absolute;
    }

    .collapsable {
      overflow: hidden;

      &.closed {
        max-height: 0px !important;
      }
    }

    /deep/ {
      .v-toolbar__title {
        max-height: 100%;
      }
    }
  }
</style>
