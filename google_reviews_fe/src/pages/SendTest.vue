<template>
  <div id="send">
    <q-page class="flex flex-center">
      <h4>Send Message Test</h4>
      <q-layout class="flex flex-center">
        <q-page padding class="absolute-full">
          <p>
            <u><b>iCabbi APP</b></u><br/>
            If sending from iCabbi APP the Send URL should be: <b>&lt;taxi company specific iCabbi API URL&gt;/sms/add</b><br/>
            and should have the following paramters:<br/>
            app_key=<b>&lt;app key from iCabbi&gt;</b><br/>
            secret_key=<b>&lt;secret key from iCabbi&gt;</b><br/>
            recipient=<b>&lt;telephone number&gt;</b><br/>
            body=<b>&lt;message&gt;</b><br/>
            <br/>
            <u><b>Testing iCabbi Hook parameters</b></u><br/>
            If testing the iCabbi hook parameters use:
            URL: <b>https://gr.taxi-magic.com/googlereviews</b><br/>
            and may have the following parameters (these should be the same as those entered on iCabbi but replacing dynamic values with static ones).<br/>
            gr_token=<b>&lt;token from the config&gt;</b><br/>
            t=<b>&lt;telephone number&gt;</b><br/>
            m=<b>&lt;message You do NOT need to include this parameter if set in database it is not necessary&gt;</b><br/>
            <b>NOTE: You may not receive a SMS if the attempt is outside of the configured times.</b><br/>
            <b>NOTE: If you want to ignore the telephone checks (to facilitate for sending multiple tests to the same telephone number) add the following parameter:</b><br />ignore_telephone_checks=1<br/>
            <b>NOTE: If you want to ignore the dispatcher checks (to facilitate for testing and not having a booking) add the following parameter:</b><br />ignore_dispatcher_checks=1<br/>
            <b>NOTE: If you want to ignore the time and sent count checks (to facilitate for testing and not wanting to change the configured times and max daily sent count) add the following parameter:</b><br />ignore_time_and_sent_count_checks=1<br/>
            </p>
            <!-- <q-form ref="form" v-model="valid" lazy-validation autofocus autocomplete="off"> -->
            <q-form ref="form" lazy-validation autofocus autocomplete="off">
            <q-input v-model="sendURL" :rules="sendURLRules" label="Send URL" required ref="sendURLInput" />
            <q-checkbox v-model="httpGet" label="Send request as an HTTP GET (default POST)" />
            <q-input type="textarea" v-model="parameters" :rules="parametersRules" label="Enter each Parameter on a separate line." rows="10" required />

            <div class="q-pa-md q-gutter-sm">
              <!-- <q-btn :disabled="!valid" color="primary" @click="validate">Send</q-btn> -->
              <q-btn color="primary" @click="validate">Send</q-btn>
            </div>
          </q-form>
        </q-page>
      </q-layout>
    </q-page>
  </div>
</template>

<script>
import { api } from 'boot/axios'

export default {
  name: 'sendMsg',
  data () {
    return {
      // valid: false,
      // err: false,
      // errStr: 'Unknown Error',
      // results: false,
      // resultStr: 'Unknown Result',
      sendURL: '',
      sendURLRules: [
        v =>
          !!v ||
          'Send URL is required. Get this from a hook on iCabbi standard send message.',
        v =>
          /[-a-zA-Z0-9@:%._+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_+.~#?&//=]*)/.test(
            v
          ) || 'Google Reviews Config Send URL must be valid'
      ],
      httpGet: false,
      parameters: '',
      parametersRules: [
        v =>
          !!v ||
          'Parameters is required. Get these from a hook on iCabbi standard send message. Can probably just copy them here.'
      ]
    }
  },
  mounted () {
    this.$refs.sendURLInput.focus()
  },
  methods: {
    validate () {
      this.$refs.form.validate()
        .then(v => {
          if (v) {
            console.log(this.parameters)
            // var parameters = this.parameters.replace(/=/g, ": ").replace(/\r?\n/g, ",")
            // var parameters = encodeURI(this.parameters.replace(/\r?\n/g, '&'))
            // var parameters = this.parameters.split(/\r?\n/g)
            // console.log(parameters)
            // var pStr = ""
            // for (var p of parameters) {
            //   console.log(p)
            //   var ps = p.split(/=/)
            //   if (pStr.length > 0) {
            //     pStr += "&"
            //   }
            //   pStr += ps[0] + "=" + encodeURI(ps[1])
            // }
            // console.log(pStr)
            // Get CORS error need to send to server which then sends and get response from there
            // this.$http.post(this.sendURL, parameters).then(result => {
            //   console.log(result.data)
            //   this.results = true
            //   this.resultStr = result.data
            //   this.resultStr = result
            // })

            // example request parameters:
            // https://www.mobitexi.co.uk/texitxt/private/sendsms/send

            // u=ambercars
            // p=chrislovesamber
            // t=07123456789
            // m=Your Ride arrived at  somewhere and couldn't find you. Please contact our office if you still require a taxi.
            // this.$http
            api
              .post('/auth/sendtest', {
                send_url: this.sendURL,
                http_get: this.httpGet,
                parameters: this.parameters
              })
              .then(result => {
                // console.log(result.data.success)
                const success = result.data.success
                const err = result.data.err
                const response = result.data.response
                console.log(success)
                console.log(err)
                console.log(response)
                if (err !== '') {
                  // this.err = true
                  // this.results = false
                  // this.errStr = 'Error getting response from remote server: ' + err
                  this.$q.notify({
                    message: 'Error getting response from remote server: ' + err,
                    icon: 'warning',
                    color: 'red',
                    timeout: 600000,
                    closeBtn: 'Close'
                  })
                } else if (response !== '') {
                  // this.err = false
                  // this.results = true
                  // this.resultStr =
                  //   "If you receive a SMS (use characters between single quotes): '" +
                  //   response +
                  //   "'"
                  this.$q.notify({
                    message: "If you receive a SMS (use characters between single quotes EXCEPT for iCabbi where you use ok): '" + response + "'",
                    icon: 'notification_important',
                    color: 'green',
                    timeout: 600000,
                    closeBtn: 'Close'
                  })
                } else {
                  // this.err = false
                  // this.results = true
                  // this.resultStr =
                  //   "If you receive a SMS (use characters between single quotes): 'EMPTY'"
                  this.$q.notify({
                    message: "If you receive a SMS (use characters between single quotes): 'EMPTY'",
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
