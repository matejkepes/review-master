<template>
  <div id="g-r-config-time">
    <q-page class="flex flex-center">
      <h4>Add Time</h4>
      <q-layout class="flex flex-center">
        <q-page padding class="absolute-full">
          <!-- <q-form ref="form" v-model="valid" lazy-validation autofocus autocomplete="off"> -->
          <q-form ref="form" lazy-validation autofocus autocomplete="off">
            <q-checkbox v-model="googleReviewsConfigTimeEnabled" label="Google Reviews Config Time Enabled" />
            <q-input v-model="googleReviewsConfigTimeStart" :rules="googleReviewsConfigTimeStartRules" label="Google Reviews Config Time Start time" required />
            <q-input v-model="googleReviewsConfigTimeEnd" :rules="googleReviewsConfigTimeEndRules" label="Google Reviews Config Time End time" required />
            <q-checkbox v-model="googleReviewsConfigTimeSunday" label="Sunday" />
            <q-checkbox v-model="googleReviewsConfigTimeMonday" label="Monday" />
            <q-checkbox v-model="googleReviewsConfigTimeTuesday" label="Tuesday" />
            <q-checkbox v-model="googleReviewsConfigTimeWednesday" label="Wednesday" />
            <q-checkbox v-model="googleReviewsConfigTimeThursday" label="Thursday" />
            <q-checkbox v-model="googleReviewsConfigTimeFriday" label="Friday" />
            <q-checkbox v-model="googleReviewsConfigTimeSaturday" label="Saturday" />

            <q-input v-model.number="googleReviewsConfigId" label="Google Reviews Config ID" required disable v-show="false" />

            <q-layout>
              <div class="q-pa-md q-gutter-sm">
                <!-- <q-btn :disabled="!valid" color="primary" @click="validate">Submit</q-btn> -->
                <q-btn color="primary" @click="validate">Submit</q-btn>
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
  name: 'GRConfigTime',
  props: ['clientIdProp', 'grctIdProp'],
  data () {
    return {
      // valid: false,
      // remoteErrors: false,
      // errStr: 'Unknown Error',
      googleReviewsConfigTimeEnabled: false,
      googleReviewsConfigTimeStart: '10:00',
      googleReviewsConfigTimeStartRules: [
        v =>
          !!v ||
          'Google Reviews Config Time Start time is required. Examples: 09:00, 21:00 this is 24 hour clock.',
        v =>
          // /^(0[0-9]|1[0-9]|2[0-3]|[0-9]):[0-5][0-9]$/.test(v) ||
          /^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$/.test(v) ||
          'Google Reviews Config Time Start time must be valid. Examples: 09:00, 21:00 this is 24 hour clock.'
      ],
      googleReviewsConfigTimeEnd: '14:00',
      googleReviewsConfigTimeEndRules: [
        v =>
          !!v ||
          'Google Reviews Config Time End time is required. Examples: 10:00, 22:00 this is 24 hour clock.',
        v =>
          // /^(0[0-9]|1[0-9]|2[0-3]|[0-9]):[0-5][0-9]$/.test(v) ||
          /^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$/.test(v) ||
          'Google Reviews Config Time End time must be valid. Examples: 09:00, 21:00 this is 24 hour clock.'
      ],
      googleReviewsConfigTimeSunday: true,
      googleReviewsConfigTimeMonday: true,
      googleReviewsConfigTimeTuesday: true,
      googleReviewsConfigTimeWednesday: true,
      googleReviewsConfigTimeThursday: true,
      googleReviewsConfigTimeFriday: true,
      googleReviewsConfigTimeSaturday: true,
      googleReviewsConfigId: ''
    }
  },
  mounted () {
    this.googleReviewsConfigId = this.grctIdProp
  },
  methods: {
    validate () {
      // if (this.$refs.form.validate()) {
      this.$refs.form.validate()
        .then(v => {
          if (v) {
            // console.log('valid')
            // console.log(this.googleReviewsConfigId)
            // console.log(this.grctIdProp)
            const grct = {
              enabled: this.googleReviewsConfigTimeEnabled,
              start: this.googleReviewsConfigTimeStart,
              end: this.googleReviewsConfigTimeEnd,
              sunday: this.googleReviewsConfigTimeSunday,
              monday: this.googleReviewsConfigTimeMonday,
              tuesday: this.googleReviewsConfigTimeTuesday,
              wednesday: this.googleReviewsConfigTimeWednesday,
              thursday: this.googleReviewsConfigTimeThursday,
              friday: this.googleReviewsConfigTimeFriday,
              saturday: this.googleReviewsConfigTimeSaturday,
              google_reviews_config_id: this.googleReviewsConfigId
              // google_reviews_config_id: this.grctIdProp
            }

            // console.log(grct)
            // this.$http.post('/auth/configtime', grct).then(result => {
            api.post('/auth/configtime', grct).then(result => {
              // console.log(result.data.success)
              const success = result.data.success
              const err = result.data.err
              // console.log(success)
              // console.log(err)
              if (!success) {
                // console.log('not successful')
                // this.remoteErrors = true
                // this.errStr = err
                this.$q.notify({
                  message: 'Failed to add: ' + err,
                  icon: 'warning',
                  color: 'red'
                })
              } else {
                // console.log('successful')
                // this.remoteErrors = false
                this.$router.push('/client/' + this.clientIdProp + '/edit')
              }
            })
          }
        })
    },
    cancel () {
      this.$router.push('/client/' + this.clientIdProp + '/edit')
    }
  }
}
</script>
