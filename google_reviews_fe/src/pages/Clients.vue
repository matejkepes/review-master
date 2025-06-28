<template>
  <q-page class="flex flex-center">
    <h4>Clients</h4>
    <q-layout class="flex flex-center">
      <q-page padding class="absolute-full">

        <div class="row">
          <div class="col">
            <q-toggle
              v-for="(option, index) in foptions"
              :key="index"
              v-model="option.value"
              :label="option.label"
              :color="option.color"
              :keep-color="option.keepColor"
              @update:model-value="toggleFilter(index)"
            ></q-toggle>
          </div>
          <div class="col">
            <q-input v-model="clientName" label="Client Name" @change="clientFilterMethod()" />
          </div>
          <div class="col">
            <div class="q-pa-md">
              <q-btn label="reset filter" color="primary" @click="resetFilter"></q-btn>
            </div>
          </div>
        </div>

        <q-btn color="primary" label="New" @click="newClient"/>
        <q-table
          title="Clients"
          :rows="clients"
          :columns="columns"
          row-key="id"
          :filter="filter"
          :filter-method="clientFilterMethod"
          v-model:pagination="pagination"
        >
          <!-- <template v-slot:body="props" :props="props"> -->
          <template v-slot:body="props">
            <q-tr :props="props" class="cursor-pointer" @click="fetchConfig(props.row)">
              <q-td key="id" :props="props" class="text-xs-right">{{ props.row.id }}</q-td>
              <q-td key="enabled" :props="props">{{ props.row.enabled }}</q-td>
              <q-td key="name" :props="props">{{ props.row.name }}</q-td>
              <q-td key="country" :props="props">{{ props.row.country }}</q-td>
              <q-td key="report_email_address" :props="props">{{ props.row.report_email_address }}</q-td>
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
  name: 'clients',
  data () {
    return {

      foptions: [
        {
          label: 'Enabled',
          value: false,
          color: 'positive',
          keepColor: true
        },
        {
          label: 'Disabled',
          value: false,
          color: 'negative',
          keepColor: true
        }
      ],

      clientName: '',

      columns: [
        { name: 'id', required: true, label: 'ID', align: 'right', field: 'id', sortable: true },
        { name: 'enabled', required: true, label: 'Enabled', align: 'left', field: 'enabled', sortable: true },
        { name: 'name', required: true, label: 'Name', align: 'left', field: 'name', sortable: true },
        { name: 'country', required: true, label: 'Country', align: 'left', field: 'country', sortable: true },
        { name: 'report_email_address', required: true, label: 'Report Email', align: 'left', field: 'report_email_address', sortable: true }
      ],
      pagination: {
        sortBy: 'name',
        rowsPerPage: 10
      },
      filter: { value: 'none' },
      clients: []
    }
  },
  mounted () {
    this.fetchClients()
  },
  methods: {
    fetchClients: function () {
      api.get('/auth/clients').then(result => {
        // console.log('result.data.clients = %O', result.data.clients)
        this.clients = result.data.clients
      })
    },
    fetchConfig (a) {
      // if (event.target.classList.contains('btn-content')) return
      // alert('id: ' + a.id)
      // api.get('/auth/configsimple?id=' + a.id).then(result => {
      //   console.log(result.data)
      //   // this.clients = result.data.clients
      // })
      // this.$router.push('/simpleclient/' + a.id + '/edit')
      this.$router.push('/client/' + a.id + '/edit')
    },
    newClient () {
      this.$router.push('/simpleclient')
      // this.$router.push('/client')
    },

    toggleFilter (index) {
      if (index === 0 && this.foptions[0].value) {
        this.foptions[1].value = false
        this.filter.value = 'enabled'
        return
      } else {
        this.filter.value = 'none'
      }

      if (index === 1 && this.foptions[1].value) {
        this.foptions[0].value = false
        this.filter.value = 'disabled'
        // return
      } else {
        this.filter.value = 'none'
      }
    },

    resetFilter () {
      this.foptions.forEach(option => { option.value = false })
      this.clientName = ''
      this.filter.value = 'none'
    },

    clientFilterMethod () {
      if (this.clients.length >= 1) {
        // if (this.filter.value === 'enabled') {
        //   return this.clients.filter(row => row.enabled)
        // }
        // if (this.filter.value === 'disabled') {
        //   return this.clients.filter(row => !row.enabled)
        // }
        if (this.clientName.length > 0) {
          // return this.clients.filter(row => row.name.toLowerCase().startsWith(this.clientName.toLowerCase()))
          if (this.filter.value === 'enabled') {
            return this.clients.filter(row => row.name.toLowerCase().startsWith(this.clientName.toLowerCase()) &&
              row.enabled)
          }
          if (this.filter.value === 'disabled') {
            return this.clients.filter(row => row.name.toLowerCase().startsWith(this.clientName.toLowerCase()) &&
              !row.enabled)
          }
          return this.clients.filter(row => row.name.toLowerCase().startsWith(this.clientName.toLowerCase()))
        } else {
          if (this.filter.value === 'enabled') {
            return this.clients.filter(row => row.enabled)
          }
          if (this.filter.value === 'disabled') {
            return this.clients.filter(row => !row.enabled)
          }
        }
        return this.clients
      }
    }
  }
}
</script>
