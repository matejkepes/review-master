<template>
  <div id="stats">
    <q-page class="flex flex-center">
      <h4>Stats (New)</h4>
      <q-layout class="flex flex-center">
        <q-page padding class="absolute-full">
          <!-- <q-form ref="form" v-model="valid" lazy-validation autocomplete="off"> -->
          <q-form ref="form" lazy-validation autocomplete="off">
            <div class="row">
              <div class="col">
                <q-date v-model="pickerStart" :landscape="false" title="Start" :reactive="true" mask="YYYY-MM-DD" />
              </div>
              <div class="col">
                <q-date v-model="pickerEnd" :landscape="false" title="End" :reactive="true" mask="YYYY-MM-DD" />
              </div>
            </div>
            <div class="row">
              <div class="col">
                <q-select v-model="timeGrouping" :options="selectTimeGroup" label="Time Grouping" />
              </div>
              <div class="col">
                <div class="q-pa-md q-gutter-sm">
                  <!-- <q-btn :disabled="!valid" color="primary" @click="validate">Send</q-btn> -->
                  <q-btn color="primary" @click="validate">Send</q-btn>
                  <q-btn :disabled="stats.length == 0" color="orange" @click="csv">CSV</q-btn>
                </div>
              </div>
            </div>
          </q-form>
          <div class="row">
            <div class="col">
              <q-table
                title="Stats"
                :rows="stats"
                :columns="columns"
                row-key="client_id"
                v-model:pagination="pagination"
              >
                <!-- <template v-slot:body="props" :props="props"> -->
                <template v-slot:body="props">
                  <q-tr :props="props">
                    <q-td key="client_id" :props="props" class="text-xs-right">{{ props.row.client_id }}</q-td>
                    <q-td key="client_name" :props="props">{{ props.row.client_name }}</q-td>
                    <q-td key="sent" :props="props">{{ props.row.sent }}</q-td>
                    <q-td key="requested" :props="props">{{ props.row.requested }}</q-td>
                    <q-td key="group_period" :props="props" class="text-xs-right">{{ formatGroupPeriod(props.row.group_period) }}</q-td>
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
import { copyToClipboard } from 'quasar'

export default {
  name: 'stats',
  data () {
    return {
      // valid: true,
      // remoteErrors: false,
      // errStr: 'Unknown Error',

      // pickerStart: new Date().toISOString().substr(0, 10),
      pickerStart: this.setStartDate(),
      // pickerEnd: new Date().toISOString().substr(0, 10),
      pickerEnd: this.setEndDate(),
      timeGrouping: 'Day',
      selectTimeGroup: ['Day', 'Week', 'Month', 'Year'],

      columns: [
        { name: 'client_id', required: true, label: 'Client ID', align: 'right', field: 'client_id', sortable: true },
        { name: 'client_name', required: true, label: 'Client Name', align: 'left', field: 'client_name', sortable: true },
        { name: 'sent', required: true, label: 'Sent', align: 'left', field: 'sent', sortable: true },
        { name: 'requested', required: true, label: 'Requested', align: 'left', field: 'requested', sortable: true },
        { name: 'group_period', required: true, label: 'Group Period', align: 'left', field: 'group_period', sortable: false }
      ],
      pagination: {
        rowsPerPage: 10
      },
      stats: []
    }
  },
  methods: {
    validate () {
      // if (this.$refs.form.validate()) {
      this.$refs.form.validate()
        .then(v => {
          if (v) {
            // console.log('valid')
            // console.log(this.pickerStart)
            // console.log(this.pickerEnd)
            // console.log(this.timeGrouping)
            // this.$http
            this.getStats('/auth/statsnew?' +
                    'start_day=' +
                    this.pickerStart +
                    '&end_day=' +
                    this.pickerEnd +
                    '&time_grouping=' +
                    this.timeGrouping
            )
          }
        })
    },
    csv () {
      if (this.stats.length > 0) {
        let st = 'ClientID,Client Name,Sent,Requested,Group Period'
        for (const s of this.stats) {
          st += '\n'
          st += s.client_id + ',' + s.client_name + ',' + s.sent + ',' + s.requested + ',' + this.formatGroupPeriod(s.group_period)
        }
        copyToClipboard(st).then(() => {
          // success
          let as = 'The CSV format has been copied to your clipboard open a file and paste the contents\n'
          as += 'If this fails the content can be copied below and pasted to a file.\n\n' + st
          alert(as)
        })
          .catch(() => {
            // fail
            alert(st)
          })
      }
    },
    getStats (url) {
      this.stats = []
      api
        .get(
          url
        )
        .then(result => {
          // console.log(result.data)
          // console.log(result.data.stats)
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
            if (result.data.stats !== null) {
              this.stats = result.data.stats
            } else {
              this.stats = []
            }
          }
        })
    },
    // formatDate (s) {
    //   var d = new Date(Date.parse(s))
    //   var a = d.toString()
    //   var b = a.split(/\s\d\d:\d\d:\d\d/)
    //   return b[0]
    // },
    formatGroupPeriod (s) {
      // This is for Day when the response is in ISO format YYYY-MM-ddThh:mm:ssZ with a T for time
      const a = s.split('T')
      return a[0]
    },
    setStartDate () {
      const n = new Date()
      // n.setDate(n.getDate() - 1)
      return n.toISOString().substr(0, 10)
    },
    setEndDate () {
      const n = new Date()
      // n.setDate(n.getDate() + 1)
      return n.toISOString().substr(0, 10)
    }
  }
}
</script>
