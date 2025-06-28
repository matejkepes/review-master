<template>
  <div id="checkLog">
    <q-page class="flex flex-center">
      <h4>Check Log for Errors</h4>
      <q-layout class="flex flex-center">
        <q-page padding class="absolute-full">
          <!-- <q-form ref="form" v-model="valid" lazy-validation autocomplete="off"> -->
          <q-form ref="form" lazy-validation autocomplete="off">
            <div class="row">
              <div class="col">
                <q-input v-model="hoursBack" :rules="hoursBackRules" label="Number of hours backwards to search" type="number" min="1" required />
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
                title="Checking Send SMS Errors Only"
                :rows="errors"
                :columns="columns"
                row-key="client_id"
                v-model:pagination="pagination"
              >
                <template v-slot:body="props"  :props="props">
                  <q-tr :props="props">
                    <q-td key="client_id" :props="props" class="text-xs-right">{{ props.row.client_id }}</q-td>
                    <q-td key="frequency" :props="props">{{ props.row.frequency }}</q-td>
                    <q-td key="last_error" :props="props">{{ props.row.last_error }}</q-td>
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
  name: 'checkLog',
  data () {
    return {
      // valid: true,
      // remoteErrors: false,
      // errStr: 'Unknown Error',

      hoursBack: '24',
      hoursBackRules: [
        v => !!v || 'Number of hours backwards to search is required',
        v =>
          (v && v > -1) ||
          'Number of hours backwards to search must be greater than 0. It is the number of hours backwards to search for errors.'
      ],

      columns: [
        { name: 'client_id', required: true, label: 'Client ID', align: 'right', field: 'client_id', sortable: true },
        { name: 'frequency', required: true, label: 'Frequency', align: 'right', field: 'frequency', sortable: true },
        { name: 'last_error', required: true, label: 'Last Error', align: 'left', field: 'last_error', sortable: true }
      ],
      pagination: {
        rowsPerPage: 10
      },
      errors: []
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
            api
              .get(
                '/auth/checklog?' +
              'hours_back=' +
              this.hoursBack
              )
              .then(result => {
                // console.log(result.data)
                // console.log(result.data.errors)
                const success = result.data.success
                const err = result.data.err
                if (!success) {
                  // console.log('not successful')
                  // this.remoteErrors = true
                  // this.errStr = err
                  this.$q.notify({
                    message: 'Failed to update: ' + err,
                    icon: 'warning',
                    color: 'red'
                  })
                } else {
                  // console.log('successful')
                  // this.remoteErrors = false
                  if (result.data.errors !== null) {
                    if (result.data.errors.length < 1) {
                      this.$q.notify({
                        message: 'There are no errors for this time period',
                        icon: 'notification_important',
                        color: 'green',
                        timeout: 6000,
                        closeBtn: 'Close'
                      })
                    }
                    this.errors = result.data.errors
                  } else {
                    this.errors = []
                  }
                }
              })
          }
        })
    }
    // formatDate (s) {
    //   var d = new Date(Date.parse(s))
    //   var a = d.toString()
    //   var b = a.split(/\s\d\d:\d\d:\d\d/)
    //   return b[0]
    // },
  }
}
</script>
