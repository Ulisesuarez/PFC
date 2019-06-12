<template>
  <v-layout>
  <v-dialog
    v-model="dialog"
    max-width="600"
  >
    <v-card>
      <v-card-title class="headline">{{$t('Share File')}}</v-card-title>

      <v-card-text>
        {{$t('Here is the link and the access pin, share both with the person with whom you want to share the study.')}}
      </v-card-text>
      <v-card-text>
      Link to your File:      <a :href="link">{{link}}</a>
      </v-card-text>
      <v-card-text>
      PIN :                    {{pin}}
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>

        <v-btn
          color="primary"
          flat="flat"
          @click="bus.$emit('closeShare')"
        >
          {{$t('Accept')}}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
  </v-layout>
</template>

<script>
import bus from '../../assets/js/bus'
export default {
  name: 'TransferDialog',
  props: {
    dialog: {
      default: false,
      type: Boolean
    },
    study: {
      default: null,
      type: Object
    }
  },
  data () {
    return {
      bus: bus
    }
  },
  computed: {
    link () {
      if (this.study && this.study.magicLink) {
        return this.study.magicLink.substring(0, this.study.magicLink.indexOf('#'))
      }
      return ''
    },
    pin () {
      if (this.study && this.study.magicLink) {
        return this.study.magicLink.substring(this.study.magicLink.indexOf('#') + 1, this.study.magicLink.length)
      }
      return ''
    }
  }
}
</script>

<style scoped>

</style>
