<template>
  <v-layout>
  <v-dialog
    v-model="dialog"
    max-width="600"
  >
    <v-card>
      <v-card-title class="headline">{{$t('Transfer File')}}</v-card-title>

      <v-card-text>
        {{$t('If you transfer the file, you can only access it if the new owner shares it with you')}}
      </v-card-text>
      <v-card-text>
        <v-form>
          <v-text-field
            ref="email"
            v-model="email"
            :label="$t('email')"
            :rules="emailRules"
            prepend-icon="person"
            name="transfer"
            type="text"
          />
          <v-text-field
            ref="phone"
            v-model="phone"
            :label="$t('phone')"
            :rules="phoneRules"
            prepend-icon="phone"
            name="transfer"
            type="text"
          />
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>

        <v-btn
          color="primary"
          flat="flat"
          @click="bus.$emit('closeTransfer')"
        >
          {{$t('Cancel')}}
        </v-btn>

        <v-btn
          color="primary"
          flat="flat"
          @click="addTask"
        >
          {{$t('Transfer')}}
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
    studyUid: {
      default: '',
      type: String
    }
  },
  data () {
    return {
      bus: bus,
      phone: '',
      email: ''
    }
  },
  methods: {
    addTask () {
      this.$idoniapacs.post('task', {
        studyUID: this.studyUid,
        sopObjects: [],
        studyName: 'testnumero1001',
        containerName: 'toDO',
        steps: ['TransferFile'],
        transferRQ: {
          phone: this.phone,
          email: this.email
        }

      }).then(response => {
        console.log(response)
      }).catch((error) => {
        console.log('error')
        console.log(error)
      })
      bus.$emit('closeTransfer')
    }
  }
}
</script>

<style scoped>

</style>
