<template>
  <div id="config">
    <q-page class="flex flex-center">
      <h4>Client Config</h4>
      <q-layout class="flex flex-center">
        <q-page padding class="absolute-full">
          <!-- <q-form ref="form" v-model="valid" lazy-validation autofocus autocomplete="off"> -->
          <q-form ref="form" lazy-validation autofocus autocomplete="off">
            <q-separator />
            <h5>Client</h5>
            <q-input v-model="clientId" label="Client ID" required disable />
            <q-checkbox v-model="clientEnabled" label="Client Enabled" />
            <q-input v-model="clientName" :rules="clientNameRules" label="Client Name" required />
            <q-input v-model="clientNote" :rules="clientNoteRules" label="Client Note" type="textarea" rows="6" required />
            <q-input v-model="clientCountry" :rules="clientCountryRules" label="Client Country" required />
            <q-input v-model="reportEmailAddress" :rules="reportEmailAddressRules" label="Report Email Address" />

            <!-- <q-separator /> -->
            <h5>Configs</h5>
            <q-list bordered>
              <g-r-config
                v-for="(grc, index) in grConfigs"
                :key="index"
                v-bind:grc="grc"
                v-bind:index="index"
                v-bind:clientId="clientId"
                v-on:updateGRC="updateGRC"
              />
            </q-list>
            <!-- <q-separator /> -->

            <q-layout>
              <div class="q-pa-md q-gutter-sm">
                <!-- <q-btn :disabled="!valid" color="primary" @click="validate">Submit</q-btn> -->
                <q-btn color="primary" @click="validate">Submit</q-btn>
                <q-btn color="red" @click="reset">Reset Form</q-btn>
                <!-- <q-btn color="orange" @click="resetValidation">Reset Validation</q-btn> -->
                <q-btn color="secondary" @click="cancel">Cancel</q-btn>
                <q-btn color="orange" @click="addConfig">Add Config</q-btn>
              </div>
            </q-layout>
          </q-form>
        </q-page>
      </q-layout>
    </q-page>
  </div>
</template>

<script>
import GRConfig from '../components/GRConfig'
import { api } from 'boot/axios'

export default {
  name: 'clientConfig',
  props: ['clientIdProp'],
  components: {
    GRConfig
  },
  provide () {
    return {
      // so can use validator in components
      $validator: this.$validator
    }
  },
  data () {
    return {
      // valid: false,
      // remoteErrors: false,
      // errStr: 'Unknown Error',
      clientId: '',
      // clientIdRules: [v => !!v || 'Client ID is required'],
      clientEnabled: false,
      clientName: '',
      clientNameRules: [
        v => !!v || 'Client Name is required',
        v =>
          (v && v.length > 2) || 'Client Name must be greater than 3 characters'
      ],
      clientNote: '',
      clientNoteRules: [],
      clientCountry: 'GB',
      clientCountryRules: [
        v => !!v || 'Client Country is required',
        v =>
          (v && v.length === 2) ||
          'Client Country must be 2 characters. It conforms to ISO 3166-1 alpha-2 see: https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2'
      ],
      reportEmailAddress: '',
      reportEmailAddressRules: [
        v => (!v || /^(?=[a-zA-Z0-9@._%+-]{6,254}$)[a-zA-Z0-9._%+-]{1,64}@(?:[a-zA-Z0-9-]{1,63}\.){1,8}[a-zA-Z]{2,63}$/.test(v)) || 'Report Email Address must be valid'
      ],

      clientConfig: null,
      grConfigs: null
    }
  },
  mounted () {
    this.fetchClient()
  },
  methods: {
    fetchClient: function () {
      if (this.clientIdProp) {
        // this.$http.get('/auth/config?id=' + this.clientIdProp).then(result => {
        api.get('/auth/config?id=' + this.clientIdProp).then(result => {
          // console.log(result.data.config)
          this.clientConfig = result.data.config
          if (this.clientConfig != null) {
            const client = this.clientConfig.client
            this.clientId = client.id
            this.clientEnabled = client.enabled
            this.clientName = client.name
            this.clientNote = client.note
            this.clientCountry = client.country
            this.reportEmailAddress = client.report_email_address

            const grcs = this.clientConfig.configs
            this.grConfigs = grcs
          }
        })
      } else {
        // TODO: Initialise clientConfig to have blank entry
      }
    },
    validate () {
      // console.log(this.clientConfig.configs)
      // if (this.$refs.form.validate()) {
      this.$refs.form.validate()
        .then(v => {
          if (v) {
            // console.log('valid')
            // console.log(this.grConfigs)
            // console.log(this.grConfigs.google_reviews_config.max_daily_send_count)
            // console.log(
            //   this.grConfigs[0].google_reviews_config.max_daily_send_count
            // )
            // var c = {
            //   client: {
            //     id: this.clientId,
            //     enabled: this.clientEnabled,
            //     name: this.clientName,
            //     note: this.clientNote,
            //     country: this.clientCountry
            //   },
            //   configs: [
            //     {
            //       google_reviews_config: {
            //         id: this.googleReviewsConfigId,
            //         enabled: this.googleReviewsConfigEnabled,
            //         min_send_frequency: this.googleReviewsConfigMinSendFrequency,
            //         max_send_count: this.googleReviewsConfigMaxSendCount,
            //         max_daily_send_count: 2,
            //         token: this.googleReviewsConfigToken,
            //         telephone_parameter: this.googleReviewsConfigTelephoneParameter,
            //         send_url: this.googleReviewsConfigSendURL,
            //         send_success_response: this
            //           .googleReviewsConfigSendSuccessResponse,
            //         time_zone: this.googleReviewsConfigTimeZone,
            //         client_id: this.clientId
            //       },
            //       google_reviews_config_times: [
            //         {
            //           id: this.googleReviewsConfigTimeId,
            //           enabled: this.googleReviewsConfigTimeEnabled,
            //           start: this.googleReviewsConfigTimeStart,
            //           end: this.googleReviewsConfigTimeEnd,
            //           google_reviews_config_id: this.googleReviewsConfigId
            //         }
            //       ]
            //     }
            //   ]
            // }
            const c = {
              client: {
                id: this.clientId,
                enabled: this.clientEnabled,
                name: this.clientName,
                note: this.clientNote,
                country: this.clientCountry,
                report_email_address: this.reportEmailAddress
              },
              configs: this.clientConfig.configs
            }

            // console.log(c)
            if (this.clientIdProp) {
              // this.$http.put('/auth/config', c).then(result => {
              api.put('/auth/config', c).then(result => {
                // console.log(result.data.success)
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
                  this.$router.push('/clients')
                }
              })
            } else {
              // this.$http.post('/auth/config', c).then(result => {
              api.post('/auth/config', c).then(result => {
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
                    message: 'Failed to update: ' + err,
                    icon: 'warning',
                    color: 'red'
                  })
                } else {
                  // console.log('successful')
                  // this.remoteErrors = false
                  this.$router.push('/clients')
                }
              })
            }
          }
        })
    },
    updateGRC: function (index, grc) {
      // console.log('index')
      // console.log(index)
      // console.log(grc)
      // console.log(this.clientConfig.configs[index])
      this.clientConfig.configs[index] = grc
      // console.log(this.clientConfig.configs[index])
    },
    reset () {
      // this.$refs.form.reset()
      // this.fetchClient()
      window.location.reload()
    },
    // resetValidation() {
    //   this.$refs.form.resetValidation()
    // }
    cancel () {
      this.$router.push('/clients')
    },
    addConfig () {
      // console.log('Add Config pushed')
      this.$router.push(
        '/grconfig/' + this.clientId + '/add'
      )
    }
  }
}
</script>
