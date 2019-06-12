<template>
  <v-layout align-center justify-center fill-height>
    <v-flex xs12 sm10 md6>
      <Card :loading="$store.state.session.loggingIn" >
        <template slot="title">
          <template>
            {{ $t('login') }}
          </template>
        </template>
        <template slot="content">
          <v-form
            ref="form"
            lazy-validation
            @keydown.native.enter="login">
            <v-text-field
              ref="email"
              v-model="email"
              :label="$t('email')"
              :rules="emailRules"
              prepend-icon="person"
              name="login"
              type="text"
            />
            <v-text-field
              id="password"
              v-model="password"
              :rules="passRules"
              :label="$t('password')"
              :append-icon="passwordIsVisible ? 'visibility_off' : 'visibility'"
              :type="passwordIsVisible ? 'text' : 'password'"
              prepend-icon="lock"
              name="password"
              @click:append="passwordIsVisible = !passwordIsVisible"
            />
            <div class="text-xs-right">
              <a
                href="https://idonia.com/access/forgot-password"
                target="_blank"
                v-text="$t('Forgot password?')"
              />
            </div>
            <v-alert
              :value="error"
              type="error"
              transition="slide-x-transition">
              {{ error }}
            </v-alert>
          </v-form>
        </template>
        <template slot="actions">
          <v-spacer />
          <v-btn
            :loading="$store.state.session.loggingIn"
            color="primary"
            @click="login">
            {{ $t('login') }}
          </v-btn>
        </template>
      </Card>
    </v-flex>
  </v-layout>
</template>

<script>
import Card from '../components/common/Card'

export default {
  name: 'Login',
  components: { Card },
  data: () => ({
    email: '',
    password: '',
    passwordIsVisible: false,
    isValid: true,
    error: ''
  }),
  computed: {
    emailRules () {
      const rules = []
      const rule = () => !!this.email || this.$t('Email required')
      rules.push(rule)

      return rules
    },
    passRules () {
      const rules = []
      const rule = () => !!this.password || this.$t('Password required')
      rules.push(rule)

      return rules
    }
  },
  methods: {
    async login () {
      this.validateField()
      if (!this.isValid) {
        this.$notify({
          group: 'app',
          title: this.$t('Invalid credentials'),
          text: this.$t('try again or contact with the admin'),
          type: 'error'
        })
        return
      }

      try {
        await this.$store.dispatch('session/login', {
          email: this.email,
          password: this.password
        })
        this.$router.push({ name: 'studies' })
        this.$idoniapacs.defaults.headers.common['Authorization'] = `IDONIA ${this.$store.state.session.token}`
      } catch (e) {
        if (e.response && e.response.data) {
          this.error = this.$t(e.response.data.code)
        } else {
          console.log(e)
          console.log(JSON.stringify(e))
          this.$idoniapacs.get('2', { params: { _: Date.now() } }).then(response => {
            console.log(response)
          }, err => {
            console.log(err)
            console.log('melotemia')
          }
          ).catch(err => {
            console.log('nOOO')
            console.log(err)
          })
        }
      }
    },
    validateField () {
      if (this.$refs.form) {
        this.isValid = this.$refs.form.validate()
        console.log(this.isValid)
        return this.isValid
      }
      this.isValid = false
      return this.isValid
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
