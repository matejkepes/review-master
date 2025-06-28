<template>
  <q-page class="flex flex-center">
    <h4>Login Page</h4>
    <q-layout class="flex flex-center">
      <q-page>
        <!-- <q-form ref="form" v-model="valid" lazy-validation autofocus autocomplete="off"> -->
        <q-form ref="form" lazy-validation autofocus autocomplete="off">
          <q-input v-model="username" :rules="usernameRules" label="Username" required ref="usernameInput"/>
          <q-input v-model="password" :rules="passwordRules" type="password" label="Password" required ref="passwordInput" v-on:keyup.enter="validate"/>
          <!-- <q-btn :disabled="!valid" color="primary" @click="validate">Login</q-btn> -->
          <q-btn color="primary" @click="validate">Login</q-btn>
        </q-form>
      </q-page>
    </q-layout>
  </q-page>
</template>

<script>
export default {
  data () {
    return {
      // valid: false,
      username: '',
      usernameRules: [
        v => !!v || 'Username is required',
        v => (v && v.length > 2) || 'Username must be greater than 3 characters'
      ],
      password: '',
      passwordRules: [
        v => !!v || 'Password is required',
        v => (v && v.length > 2) || 'Password must be greater than 3 characters'
      ]
    }
  },
  mounted () {
    this.$refs.usernameInput.focus()
  },
  computed: {
    isLoggedIn: function () {
      return this.$store.getters.isLoggedIn
    }
  },
  methods: {
    validate: function () {
      // console.log('this.$refs.form.validate() = %O', this.$refs.form.validate())
      // if (this.$refs.form.validate()) {
      //   this.login()
      // }
      this.$refs.form.validate()
        .then(v => {
          if (v) {
            this.login()
          }
        })
    },
    login: function () {
      const username = this.username
      const password = this.password
      // console.log('this.$store = %O', this.$store)
      this.$store
        .dispatch('login/login', { username: username, password: password })
        .then(() => this.$router.push('/'))
        .catch(err => {
          console.log(err)
          this.$q.notify({
            message: 'Invalid credentials',
            icon: 'warning',
            color: 'red'
          })
        })
    }
  }
}
</script>
