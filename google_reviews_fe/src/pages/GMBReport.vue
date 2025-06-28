<template>
  <div id="gMBReport">
    <q-page class="flex flex-center">
      <h4>Run Google My Business Report</h4>
      <q-layout class="flex flex-center">
        <q-page padding class="absolute-full">
          <!-- <q-form ref="form" v-model="valid" lazy-validation autocomplete="off"> -->
          <q-form ref="form" lazy-validation autocomplete="off">
            <!-- <div class="row"> -->
              <div class="col">
                <q-input v-model="monthsBack" :rules="monthsBackRules" label="Number of months backwards to report 0 is current month" type="number" min="0" required />
              </div>
              <div class="col">
                <q-checkbox v-model="csvReportOnly" label="Only send the CSV report, do NOT send individual reports" />
              </div>
              <div class="col">
                <div class="q-pa-md q-gutter-sm">
                  <!-- <q-btn :disabled="!valid" color="primary" @click="validate">Send</q-btn> -->
                  <q-btn :disable="reportRunning" color="primary" @click="validate">Send</q-btn>
                </div>
              </div>
            <!-- </div> -->
          </q-form>
          <div v-if="reportRunning" class="absolute-center">
            <div class="row">
              <div class="col">
                <q-spinner
                color="primary"
                size="3em"
                />
                <span><b>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Please wait this could take sometime...</b></span>
              </div>
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
  name: 'gMBReport',
  data () {
    return {
      // valid: true,

      monthsBack: '0',
      monthsBackRules: [
        v => !!v || 'Number of months backwards to report is required',
        v =>
          (v && v > -1) ||
          'Number of months backwards to report must be at least 0. It is the number of months backwards to report for errors.'
      ],
      csvReportOnly: true,
      reportRunning: false
    }
  },
  methods: {
    validate () {
      // if (this.$refs.form.validate()) {
      this.$refs.form.validate()
        .then(v => {
          if (v) {
            // console.log('valid')
            this.reportRunning = true
            api
              .get(
                '/auth/gmybusinessreport?' +
              'months_back=' +
              this.monthsBack +
              '&cvs_report_only=' +
              this.csvReportOnly
              )
              .then(result => {
                // console.log(result.data)
                // console.log(result.data.nothingSentResults)
                this.reportRunning = false
                const success = result.data.success
                const err = result.data.err
                if (!success) {
                  // console.log('not successful')
                  this.$q.notify({
                    message: 'Failed to run report: ' + err,
                    icon: 'warning',
                    color: 'red',
                    timeout: 600000,
                    closeBtn: 'Close'
                  })
                } else {
                  // console.log('successful')
                  this.$q.notify({
                    message: 'Success check your email',
                    icon: 'notification_important',
                    color: 'green',
                    timeout: 600000,
                    closeBtn: 'Close'
                  })
                }
              })
          }
        })
    }
  }
}
</script>
