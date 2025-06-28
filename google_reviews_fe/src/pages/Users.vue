<template>
  <q-page class="flex flex-center">
    <h4>Users</h4>
    <q-layout class="flex flex-center">
      <q-page padding class="absolute-full">

        <div class="row">
          <div class="col">
            <q-input v-model="emailAddress" label="Email" @onchange="userFilterMethod()" />
          </div>
          <div class="col">
            <div class="q-pa-md">
              <q-btn label="reset filter" color="primary" @click="resetFilter"></q-btn>
            </div>
          </div>
        </div>

        <q-btn color="primary" label="New" @click="newUser"/>
        <q-table
          title="Users"
          :rows="users"
          :columns="columns"
          row-key="user_id"
          :filter="filter"
          :filter-method="userFilterMethod"
          v-model:pagination="pagination"
        >
          <template v-slot:body="props">
            <q-tr :props="props" class="cursor-pointer" @click="fetchUser(props.row)">
              <q-td key="user_id" :props="props" class="text-xs-right">{{ props.row.user_id }}</q-td>
              <q-td key="email" :props="props">{{ props.row.email }}</q-td>
              <q-td key="clients" :props="props">{{ props.row.clients }}</q-td>
            </q-tr>
          </template>
        </q-table>
      </q-page>
    </q-layout>
  </q-page>
</template>

<script>
import { api } from 'boot/axios'

export default {
  name: 'users',
  data () {
    return {

      emailAddress: '',

      columns: [
        { name: 'user_id', required: true, label: 'ID', align: 'right', field: 'user_id', sortable: true },
        { name: 'email', required: true, label: 'Email', align: 'left', field: 'email', sortable: true },
        { name: 'clients', required: true, label: 'Clients', align: 'left', field: 'clients', sortable: true }
      ],
      pagination: {
        sortBy: 'email',
        rowsPerPage: 10
      },
      filter: {},
      users: []
    }
  },
  mounted () {
    this.fetchUsers()
  },
  methods: {
    fetchUsers: function () {
      api.get('/auth/users').then(result => {
        // console.log('result.data.users = %O', result.data.users)
        const users = result.data.users
        if (users == null) {
          this.users = []
        } else {
          this.users = users
        }
      })
    },

    fetchUser (a) {
      // console.log('fetch user ' + a.user_id)
      this.$router.push('/user/' + a.user_id + '/edit')
    },

    newUser () {
      this.$router.push('/user')
    },

    resetFilter () {
      this.emailAddress = ''
    },

    userFilterMethod () {
      // console.log('filter method')
      if (this.emailAddress.length > 0) {
        // console.log('emailAddress: ' + this.emailAddress)
        return this.users.filter(row => row.email.toLowerCase().startsWith(this.emailAddress.toLowerCase()))
      }
      return this.users
    }
  }
}
</script>
