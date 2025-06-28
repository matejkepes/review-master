<template>
  <div id="user">
    <q-page class="flex flex-center">
      <h4>User</h4>
      <q-layout class="flex flex-center">
        <q-page padding class="absolute-full">
          <!-- <q-form ref="form" v-model="valid" lazy-validation autofocus autocomplete="off"> -->
          <q-form ref="form" lazy-validation autofocus autocomplete="off">
            <q-separator />
            <h5>User</h5>
            <q-input v-model="userId" label="User ID" required disable />
            <q-input v-model="email" :rules="emailRules" label="Email" required />
            <q-input v-model="password" :rules="passwordRules" label="Password" required />
            <q-separator />
            <div><br/>Clients: <strong>{{ clientList }}</strong></div>
            <div class="row">
              <div class="col">
                <q-input v-model="clientName" label="Client Name" @onchange="clientFilterMethod()" />
              </div>
              <div class="col">
                <div class="q-pa-md">
                  <q-btn label="reset filter" color="primary" @click="resetFilter"></q-btn>
                </div>
              </div>
            </div>
            <q-table
              title="Clients"
              :rows="clients"
              :columns="columns"
              row-key="id"
              :selected-rows-label="getSelectedClient"
              selection="multiple"
              v-model:selected="selectedClients"
              v-model:pagination="pagination"
              :filter="filter"
              :filter-method="clientFilterMethod"
            />
            <div><br/>Clients: <strong>{{ clientList }}</strong></div>

            <q-layout>
              <div class="q-pa-md q-gutter-sm">
                <q-btn color="primary" @click="validate">Submit</q-btn>
                <q-btn color="red" @click="reset">Reset Form</q-btn>
                <q-btn :disabled="!userIdProp" color="orange" @click="deleteUser">Delete</q-btn>
                <q-btn color="secondary" @click="cancel">Cancel</q-btn>
              </div>
            </q-layout>

          </q-form>
        </q-page>
      </q-layout>
    </q-page>
  </div>
</template>

<script>
import { api } from 'boot/axios'

export default {
  name: 'user',
  props: ['userIdProp'],
  provide () {
    return {
      // so can use validator in components
      $validator: this.$validator
    }
  },
  data () {
    return {
      userId: '',
      // userIdRules: [v => !!v || 'User ID is required'],
      email: '',
      emailRules: [
        v =>
          !!v ||
          'Email is required',
        v =>
          /^(?=[a-zA-Z0-9@._%+-]{6,254}$)[a-zA-Z0-9._%+-]{1,64}@(?:[a-zA-Z0-9-]{1,63}\.){1,8}[a-zA-Z]{2,63}$/.test(
            v
          ) || 'email must be valid'
      ],
      password: '',
      passwordRules: [
        v => !!v || 'Password is required',
        v => (v && v.length > 2) || 'Password must be greater than 3 characters'
      ],

      user: null,

      columns: [
        {
          name: 'id',
          required: true,
          label: 'ID',
          align: 'right',
          field: client => client.id,
          format: val => `${val}`,
          sortable: true
        },
        { name: 'name', required: true, label: 'Name', align: 'left', field: 'name', sortable: true },
        { name: 'country', required: true, label: 'Country', align: 'left', field: 'country', sortable: true }
      ],
      pagination: {
        sortBy: 'name',
        rowsPerPage: 10
      },
      clients: [],
      selectedClients: [],

      clientList: '',

      clientName: '',
      filter: {}
    }
  },
  mounted () {
    // this.fetchClients()
    // this.fetchUser()
    this.fetchInitialData()
  },
  methods: {
    fetchInitialData: function () {
      this.fetchClients()
      this.fetchUser()
    },
    fetchClients: function () {
      api.get('/auth/clients').then(result => {
        // console.log('result.data.clients = %O', result.data.clients)
        this.clients = result.data.clients
        // console.log(this.clients)
      })
    },
    fetchUser: function () {
      if (this.userIdProp) {
        api.get('/auth/user?id=' + this.userIdProp).then(result => {
          // console.log(result.data.user)
          // console.log(result.data.user)
          // console.log('read clients: ', result.data.user.clients)
          this.user = result.data.user.user
          if (this.user != null) {
            const user = this.user
            this.userId = user.id
            this.email = user.email
            this.password = user.password
          }
          // set initial selections
          const initialClientsSelected = result.data.user.clients
          if (initialClientsSelected != null) {
            this.selectedClients = initialClientsSelected
          }
        })
      } else {
        // Initialise user id
        this.userId = 0
      }
    },
    getSelectedClient () {
      // const sct = JSON.parse(JSON.stringify(this.selectedClients))
      // console.log(sct)
      const sc = this.selectedClients
      const o = []
      sc.forEach((values, index) => {
        o.push(values.name + ' (' + values.id + ')')
      })
      this.clientList = o
    },

    resetFilter () {
      this.clientName = ''
    },
    clientFilterMethod () {
      // console.log('In filter method')
      if (this.clientName.length > 0) {
        return this.clients.filter(row => row.name.toLowerCase().startsWith(this.clientName.toLowerCase()))
      } else {
        return this.clients
      }
    },

    validate () {
      // console.log(this.user)
      this.$refs.form.validate()
        .then(v => {
          if (v) {
            // console.log('valid')
            // user
            const u = {
              id: this.userId,
              email: this.email,
              password: this.password
            }
            // selected clients
            const sc = JSON.parse(JSON.stringify(this.selectedClients))
            if (sc.length < 1) {
              this.$q.notify({
                message: 'You must select at least one client',
                icon: 'warning',
                color: 'red',
                timeout: 5000,
                closeBtn: 'Close'
              })
              return
            }
            // console.log(u)
            // console.log(sc)
            const userClients = {
              user: u,
              clients: sc
            }
            // console.log(userClients)
            if (this.userIdProp) {
              api.put('/auth/user', userClients).then(result => {
                // console.log(result.data.success)
                const success = result.data.success
                const err = result.data.err
                if (!success) {
                  // console.log('not successful')
                  this.$q.notify({
                    message: 'Failed to update: ' + err,
                    icon: 'warning',
                    color: 'red'
                  })
                } else {
                  // console.log('successful')
                  this.$router.push('/users')
                }
              })
            } else {
              api.post('/auth/user', userClients).then(result => {
                // console.log(result.data.success)
                const success = result.data.success
                const err = result.data.err
                if (!success) {
                  // console.log('not successful')
                  this.$q.notify({
                    message: 'Failed to add: ' + err,
                    icon: 'warning',
                    color: 'red'
                  })
                } else {
                  // console.log('successful')
                  this.$router.push('/users')
                }
              })
            }
          }
        })
    },
    reset () {
      window.location.reload()
    },
    deleteUser () {
      this.$q.notify({
        message: 'Are you sure you want to delete the user?',
        color: 'red',
        icon: 'warning',
        actions: [
          // { label: 'Yes', color: 'yellow', handler: () => { this.deleteUserExec() } },
          {
            label: 'Yes',
            color: 'yellow',
            handler: () => {
              const userId = this.userId
              // console.log('userId: ', userId)
              api.delete('/auth/user?id=' + userId).then(result => {
                // console.log(result.data.success)
                const success = result.data.success
                const err = result.data.err
                if (!success) {
                  // console.log('not successful')
                  this.$q.notify({
                    message: 'Failed to delete: ' + err,
                    icon: 'warning',
                    color: 'red'
                  })
                } else {
                  // console.log('successful')
                  this.$router.push('/users')
                }
              })
            }
          },
          { label: 'No', color: 'white', handler: () => { /* ... */ } }
        ]
      })
    },
    cancel () {
      this.$router.push('/users')
    }
  }
}
</script>
