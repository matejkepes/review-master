<template>
  <div id="checkNothingSent">
    <q-page class="flex flex-center">
      <h4>Check No Messages Sent for Companies</h4>
      <q-layout class="flex flex-center">
        <q-page padding class="absolute-full">
          <!-- <q-form ref="form" v-model="valid" lazy-validation autocomplete="off"> -->
          <q-form ref="form" lazy-validation autocomplete="off">
            <div class="row">
              <div class="col">
                <!-- <q-input v-model="hoursBack" :rules="hoursBackRules" label="Number of hours backwards to search" type="number" min="1" required /> -->
                <q-input v-model="daysBack" :rules="daysBackRules" label="Number of days backwards to search" type="number" min="0" required />
              </div>
              <div class="col">
                <div class="q-pa-md q-gutter-sm">
                  <!-- <q-btn :disabled="!valid" color="primary" @click="validate">Send</q-btn> -->
                  <q-btn color="primary" @click="validate">Send</q-btn>
                </div>
              </div>
            </div>
          </q-form>
          <div class="row">
            <div class="col">
              <q-table
                title="No messages sent for companies for specified period"
                :rows="nothingSents"
                :columns="columns"
                row-key="client_id"
                v-model:pagination="pagination"
              >
                <template v-slot:body="props"  :props="props">
                  <q-tr :props="props">
                    <q-td key="client_id" :props="props" class="text-xs-right">{{ props.row.client_id }}</q-td>
                    <q-td key="client_name" :props="props">{{ props.row.client_name }}</q-td>
                  </q-tr>
                </template>
              </q-table>
            </div>
          </div>
        </q-page>
      </q-layout>
    </q-page>
  </div>
</template>

<script>
import { api } from 'boot/axios'

export default {
  name: 'checkNothingSent',
  data () {
    return {
      // valid: true,

      // hoursBack: '24',
      // hoursBackRules: [
      //   v => !!v || 'Number of hours backwards to search is required',
      //   v =>
      //     (v && v > -1) ||
      //     'Number of hours backwards to search must be greater than 0. It is the number of hours backwards to search for errors.'
      // ],
      daysBack: '0',
      daysBackRules: [
        v => !!v || 'Number of days backwards to search is required',
        v =>
          (v && v > -1) ||
          'Number of days backwards to search must be at least 0. It is the number of days backwards to search for errors.'
      ],

      columns: [
        { name: 'client_id', required: true, label: 'Client ID', align: 'right', field: 'client_id', sortable: true },
        { name: 'client_name', required: true, label: 'Client Name', align: 'left', field: 'client_name', sortable: true }
      ],
      pagination: {
        rowsPerPage: 10
      },
      nothingSents: []
    }
  },
  methods: {
    validate () {
      // if (this.$refs.form.validate()) {
      this.$refs.form.validate()
        .then(v => {
          if (v) {
            // console.log('valid')
            // this.$http
            // api
            //   .get(
            //     '/auth/checknothingsent?' +
            //   'hours_back=' +
            //   this.hoursBack
            //   )
            api
              .get(
                '/auth/checknothingsent?' +
              'days_back=' +
              this.daysBack
              )
              .then(result => {
                // console.log(result.data)
                // console.log(result.data.nothingSentResults)
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
                  if (result.data.nothingSentResults !== null) {
                    this.nothingSents = result.data.nothingSentResults
                  } else {
                    this.nothingSents = []
                  }
                }
              })
          }
        })
    }
  }
}
</script>
