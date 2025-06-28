<template>
  <div id="g-r-config-new">
    <q-page class="flex flex-center">
      <h4>Add Config</h4>
      <q-layout class="flex flex-center">
        <q-page padding class="absolute-full">
          <!-- <q-form ref="form" v-model="valid" lazy-validation autofocus autocomplete="off"> -->
          <q-form ref="form" lazy-validation autofocus autocomplete="off">

            <!-- <q-input v-model="googleReviewsConfigId" label="Google Reviews Config ID" required disable v-show="false" /> -->
            <q-checkbox v-model="googleReviewsConfigEnabled" label="Google Reviews Config Enabled" />
            <q-checkbox v-model="ai_responses_enabled" label="AI Responses Enabled" />
            <q-checkbox v-model="monthly_review_analysis_enabled" label="Monthly Review Analysis Enabled" />
            <q-input v-model="contact_method" label="Contact method for AI responses, e.g. email us at x@y.z or call us XXX" />
            <q-input v-model="googleReviewsConfigMinSendFrequency" :rules="googleReviewsConfigMinSendFrequencyRules"
              label="Google Reviews Config Min Send Frequency" type="number" min="0" required />
            <q-input v-model="googleReviewsConfigMaxSendCount" :rules="googleReviewsConfigMaxSendCountRules"
              label="Google Reviews Config Max Send Count" type="number" min="0" required />
            <q-input v-model="googleReviewsConfigMaxDailySendCount" :rules="googleReviewsConfigMaxDailySendCountRules"
              label="Google Reviews Config Max Daily Send Count" type="number" min="0" required />
            <q-input v-model="googleReviewsConfigToken" :rules="googleReviewsConfigTokenRules"
              label="Google Reviews Config Token (use the Generate Token button for a new one)" required />
            <div class="q-pa-md q-gutter-sm">
              <q-btn color="primary" @click="generateToken">Generate Token</q-btn>
            </div>
            <q-input v-model="googleReviewsConfigTelephoneParameter" :rules="googleReviewsConfigTelephoneParameterRules"
              label="Google Reviews Config Telephone Parameter" required />
            <q-input v-model="googleReviewsConfigSendURL" :rules="googleReviewsConfigSendURLRules"
              label="Google Reviews Config Send URL" required />
            <q-checkbox v-model="googleReviewsConfigHttpGet"
              label="Google Reviews Config Send SMS request as an HTTP GET (default POST)" />
            <q-input v-model="googleReviewsConfigSendSuccessResponse"
              :rules="googleReviewsConfigSendSuccessResponseRules"
              label="Google Reviews Config Send Success Response (enter EMPTY if no response) (if using iCabbi or Review Master SMS Gateway APP set to anything e.g. ok)"
              required />
            <q-input v-model="googleReviewsConfigTimeZone" :rules="googleReviewsConfigTimeZoneRules"
              label="Google Reviews Config Time Zone" required />
            <q-checkbox v-model="googleReviewsConfigMultiMessageEnabled"
              label="Google Reviews Config Multi Message Enabled" />
            <q-input v-model="googleReviewsConfigMessageParameter" :rules="googleReviewsConfigMessageParameterRules"
              label="Google Reviews Config Message Parameter" required />
            <q-input v-model="googleReviewsConfigMultiMessageSeparator"
              :rules="googleReviewsConfigMultiMessageSeparatorRules"
              label="Google Reviews Config Multi Message Separator" required />
            <q-checkbox v-model="googleReviewsConfigUseDatabaseMessage"
              label="Google Reviews Config Use Database Message" />
            <q-input v-model="googleReviewsConfigMessage" :rules="googleReviewsConfigMessageRules"
              label="Google Reviews Message" type="textarea" rows="3" required />

            <q-checkbox v-model="googleReviewsConfigSendFromIcabbiApp"
              label="Google Reviews Config Send From iCabbi App" />
            <q-input v-model="googleReviewsConfigAppKey"
              label="Google Reviews Config App Key (iCabbi, Autocab V1) / Username (Autocab) (REQUIRED if sending from iCabbi App or Dispatcher Checks Enabled)" />
            <q-input v-model="googleReviewsConfigSecretKey"
              label="Google Reviews Config Secret Key (iCabbi) / Password (Autocab) (Autocab V1 enter anything) (REQUIRED if sending from iCabbi App or Dispatcher Checks Enabled)" />
            <q-input v-model="googleReviewsConfigDispatcherURL" :rules="googleReviewsConfigDispatcherURLRules"
              label="Google Reviews Config Dispatcher URL used when Dispatcher Checks Enabled (REQUIRES iCabbi App Key and Secrect Key) MUST be included for Autocab and Username and Password"
              required />
            <q-select v-model="googleReviewsConfigDispatcherType" :options="selectDispatcherType"
              label="Dispatcher Type" />

            <q-checkbox v-model="googleReviewsConfigSendDelayEnabled"
              label="Google Reviews Config Send Delay Enabled" />
            <q-input v-model="googleReviewsConfigSendDelay" :rules="googleReviewsConfigSendDelayRules"
              label="Google Reviews Config Send Delay in Minutes" type="number" min="0" required />

            <q-checkbox v-model="googleReviewsConfigDispatcherChecksEnabled"
              label="Google Reviews Config Dispatcher Checks Enabled" />
            <q-input v-model="googleReviewsConfigBookingIdParameter" :rules="googleReviewsConfigBookingIdParameterRules"
              label="Google Reviews Config Booking ID Parameter" required />
            <q-input v-model="googleReviewsConfigIsBookingForNowDiffMinutes"
              :rules="googleReviewsConfigIsBookingForNowDiffMinutesRules"
              label="Google Reviews Config Is Booking For Now Difference in Minutes" type="number" min="0" required />
            <q-input v-model="googleReviewsConfigBookingNowPickupToContactMinutes"
              :rules="googleReviewsConfigBookingNowPickupToContactMinutesRules"
              label="Google Reviews Config Booking for Now Pickup To Contact Minutes" type="number" min="0" required />
            <q-input v-model="googleReviewsConfigPreBookingPickupToContactMinutes"
              :rules="googleReviewsConfigPreBookingPickupToContactMinutesRules"
              label="Google Reviews Config Pre Booking Pickup To Contact Minutes" type="number" min="0" required />

            <q-checkbox v-model="googleReviewsConfigReplaceTelephoneCountryCode"
              label="Google Reviews Config Replace Telephone Country Code (for when SIM will not send international numbers, normally leave disabled)" />
            <q-input v-model="googleReviewsConfigReplaceTelephoneCountryCodeWith"
              :rules="googleReviewsConfigReplaceTelephoneCountryCodeWithRules"
              label="Google Reviews Config Replace Telephone Country Code With this value is normally a number zero"
              required />

            <q-separator />
            <h5>Review Master SMS Gateway</h5>
            <p>
              Use this for companies using the Review Master SMS Gateway App. If this is enabled make sure that the Send
              From iCabbi App is NOT also enabled.
            </p>
            <q-checkbox v-model="googleReviewsConfigReviewMasterSMSGatewayEnabled"
              label="Google Reviews Config Review Master SMS Gateway Enabled" />
            <q-checkbox v-model="googleReviewsConfigReviewMasterSMSGatewayUseMasterQueue"
              label="Google Reviews Config Review Master SMS Gateway Use Master Queue" />
            <q-input v-model="googleReviewsConfigReviewMasterSMSGatewayPairCode"
              label="Google Reviews Config Review Master SMS Gateway Pair Code (used for pairing customers app)" />
            <div class="q-pa-md q-gutter-sm">
              <q-btn color="primary" @click="generatePairCode">Generate Pairing Code</q-btn>
            </div>
            <!-- <vue-qr v-bind:text="googleReviewsConfigReviewMasterSMSGatewayPairCodeQrCode" :size="200"></vue-qr> -->
            <qrcode-vue v-bind:value="googleReviewsConfigReviewMasterSMSGatewayPairCodeQrCode" :size="200" />

            <q-separator />
            <h5>Alternate Message Service</h5>
            <p>
              Use this for companies using an alternate message service. These normally use alternate APIs that do not
              use standard HTTP POST (e.g. require the use of JSON).
            </p>
            <ul>
              <li>For Message Media set the Send URL above to https://api.messagemedia.com/v1/messages</li>
              <li>For Message Media set the Send Success Response to anything e.g. ok</li>
            </ul>
            <ul>
              <li>For Veezu set the Send URL above to https://messages.veezu.com/api/messages</li>
              <li>For Veezu set the Send Success Response to anything e.g. ok</li>
              <li>For Veezu set the Secret1 to the auth token</li>
            </ul>
            <ul>
              <li>For Autocab set the Send URL above to https://autocab-api.azure-api.net/ (i.e. normal URL for Autocab,
                should not need changing)</li>
              <li>For Autocab set the Send Success Response to anything e.g. ok</li>
              <li>For Autocab set the Secret1 to the taxi companies subscription key</li>
            </ul>
            <q-checkbox v-model="googleReviewsConfigAlternateMessageServiceEnabled"
              label="Google Reviews Config Alternate Message Service Enabled" />
            <q-select v-model="googleReviewsConfigAlternateMessageService" :options="selectAlternateMessageServiceType"
              label="Google Reviews Config Alternate Message Service Type" />
            <q-input v-model="googleReviewsConfigAlternateMessageServiceSecret1"
              label="Google Reviews Config Alternate Message Service Secret1" />

            <q-separator />
            <h5>Autocab specific filtering</h5>
            <p>
              The companies list is used to restrict the companies processed for Autocab (in a multi-company setup).
              <br />
              It is a comma separated list of company ID's to process. The ID's can be got from the dispatcher.
              <br />
              If left empty it will not do any filtering on the company ID's.
              <br />
              e.g. 1,3,7
            </p>
            <q-input v-model="googleReviewsConfigCompanies"
              label="Google Reviews Config Companies used to restrict companies for processing with Autocab (multi company setup)" />
            Booking Source Mobile App filtering (select All will ignore this filter):
            <q-option-group v-model="googleReviewsConfigBookingSourceMobileAppState"
              :options="googleReviewsConfigBookingSourceMobileAppStateOptions" color="primary" />

            <q-separator />
            <h5>Google My Business</h5>
            <p>
              The following section is for automatically replying to or reporting on Google My Business reviews.
            </p>
            <p>
              If you want to use this feature the company (client) needs to give permission. They need to add
              your agency account as a manager to their individual locations or their location group.
              Once a client invites you to an existing location group, you can access their account
              through your generated OAuth2.0 token.
            </p>
            <ul>
              <li>
                The Location Name and Postal Code must match the clients Location Name and Postal Code Address in their
                Google My Business.
              </li>
              <li>
                To automatically reply the Review Reply must be enabled and at least one of the star rating replies
                needs to be checked / enabled and the corresponding reply message.
              </li>
              <li>
                You can use &lt;name&gt; in the message and it will be replaced by the name of the reviewer.
              </li>
              <li>
                Multiple reply messages can be added by separating them using the multi message separator above.
              </li>
              <li>
                To report on the replies the Review Report must be enabled and at least one email where to send the
                report must be added. For reporting only it is not necessary to have the Reply Automatically to be
                enabled.
              </li>
            </ul>
            <q-checkbox v-model="googleReviewsConfigGoogleMyBusinessReviewReplyEnabled"
              label="Google Reviews Config Google My Business Reply to Review Automatically" />
            <q-input v-model="googleReviewsConfigGoogleMyBusinessLocationName"
              :rules="googleReviewsConfigGoogleMyBusinessLocationNameRules"
              label="Google Reviews Config Google My Business Location Name (This must match the Google My Business Location Name for this to work)"
              required />
            <q-input v-model="googleReviewsConfigGoogleMyBusinessPostalCode"
              :rules="googleReviewsConfigGoogleMyBusinessPostalCodeRules"
              label="Google Reviews Config Google My Business Postal Code (This must match the Google My Business Postal Code Address for this to work)"
              required />
            <q-file v-model="googleReviewsConfigGoogleMyBusinessFile" bg-color="primary" label-color="white"
              label="Optionally select a file with Star Rating text to auto populate ratings (NOTE: file must use .txt extension)"
              accept=".txt, text/plain" filled @update:model-value="readStarRatingFile()" />
            <q-checkbox v-model="googleReviewsConfigGoogleMyBusinessReplyToUnspecfifiedStarRating"
              label="Google Reviews Config Google My Business Reply To Unspecfified Star Rating" />
            <q-input v-model="googleReviewsConfigGoogleMyBusinessUnspecfifiedStarRatingReply"
              :rules="googleReviewsConfigGoogleMyBusinessUnspecfifiedStarRatingReplyRules"
              label="Google Reviews Config Google My Business Unspecfified Star Rating Reply" type="textarea" rows="6"
              required />
            <q-checkbox v-model="googleReviewsConfigGoogleMyBusinessReplyToOneStarRating"
              label="Google Reviews Config Google My Business Reply To One Star Rating" />
            <q-input v-model="googleReviewsConfigGoogleMyBusinessOneStarRatingReply"
              :rules="googleReviewsConfigGoogleMyBusinessOneStarRatingReplyRules"
              label="Google Reviews Config Google My Business One Star Rating Reply" type="textarea" rows="6"
              required />
            <q-checkbox v-model="googleReviewsConfigGoogleMyBusinessReplyToTwoStarRating"
              label="Google Reviews Config Google My Business Reply To Two Star Rating" />
            <q-input v-model="googleReviewsConfigGoogleMyBusinessTwoStarRatingReply"
              :rules="googleReviewsConfigGoogleMyBusinessTwoStarRatingReplyRules"
              label="Google Reviews Config Google My Business Two Star Rating Reply" type="textarea" rows="6"
              required />
            <q-checkbox v-model="googleReviewsConfigGoogleMyBusinessReplyToThreeStarRating"
              label="Google Reviews Config Google My Business Reply To Three Star Rating" />
            <q-input v-model="googleReviewsConfigGoogleMyBusinessThreeStarRatingReply"
              :rules="googleReviewsConfigGoogleMyBusinessThreeStarRatingReplyRules"
              label="Google Reviews Config Google My Business Three Star Rating Reply" type="textarea" rows="6"
              required />
            <q-checkbox v-model="googleReviewsConfigGoogleMyBusinessReplyToFourStarRating"
              label="Google Reviews Config Google My Business Reply To Four Star Rating" />
            <q-input v-model="googleReviewsConfigGoogleMyBusinessFourStarRatingReply"
              :rules="googleReviewsConfigGoogleMyBusinessFourStarRatingReplyRules"
              label="Google Reviews Config Google My Business Four Star Rating Reply" type="textarea" rows="6"
              required />
            <q-checkbox v-model="googleReviewsConfigGoogleMyBusinessReplyToFiveStarRating"
              label="Google Reviews Config Google My Business Reply To Five Star Rating" />
            <q-input v-model="googleReviewsConfigGoogleMyBusinessFiveStarRatingReply"
              :rules="googleReviewsConfigGoogleMyBusinessFiveStarRatingReplyRules"
              label="Google Reviews Config Google My Business Five Star Rating Reply" type="textarea" rows="6"
              required />
            <q-checkbox v-model="googleReviewsConfigGoogleMyBusinessReportEnabled"
              label="Google Reviews Config Google My Business Report Enabled" />
            <q-input v-model="emailAddress" :rules="emailAddressRules"
              label="Email Address used for sending information for this configuration e.g. reports. To send individual (NOT CSV) report use email that is different to CSV report email"
              required />

            <q-input v-model.number="googleReviewsClientId" label="Google Reviews Client ID" required disable
              v-show="false" />

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
// import VueQr from 'vue-qr'
import QrcodeVue from 'qrcode.vue'
import { api } from 'boot/axios'
import { ref } from 'vue'

export default {
  name: 'GRConfigNew',
  props: ['clientIdProp'],
  components: {
    // VueQr
    QrcodeVue
  },
  data () {
    return {
      // valid: false,
      // remoteErrors: false,
      // errStr: 'Unknown Error',
      // googleReviewsConfigId: '',
      // googleReviewsConfigIdRules: [
      //   v => !!v || 'Google Reviews Config ID is required'
      // ],
      googleReviewsConfigEnabled: false,
      googleReviewsConfigMinSendFrequency: '90',
      googleReviewsConfigMinSendFrequencyRules: [
        v => !!v || 'Google Reviews Config Min Send Frequency is required',
        v =>
          (v && v > -1) ||
          'Google Reviews Config Min Send Frequency must be greater than 0. It is the minimum number of days before a message will be sent to an individual telephone number.'
      ],
      googleReviewsConfigMaxSendCount: '2000',
      googleReviewsConfigMaxSendCountRules: [
        v => !!v || 'Google Reviews Config Max Send Count is required',
        v =>
          (v && v > 0) ||
          'Google Reviews Config Max Send Count must be greater than  or equal to 0. It is the maximum number of messages sent to an individual telephone number.'
      ],
      googleReviewsConfigMaxDailySendCount: '200',
      googleReviewsConfigMaxDailySendCountRules: [
        v => !!v || 'Google Reviews Config Max Daily Send Count is required',
        v =>
          (v && v > 0) ||
          'Google Reviews Config Max Daily Send Count must be greater than  or equal to 0. It is the maximum number of messages sent each day in total.'
      ],
      googleReviewsConfigToken: '',
      googleReviewsConfigTokenRules: [
        v => !!v || 'Google Reviews Config Token is required',
        v =>
          (v && v.length > 30) ||
          'Google Reviews Config Max Send Count must at least 30 characters.'
      ],
      googleReviewsConfigTelephoneParameter: '',
      googleReviewsConfigTelephoneParameterRules: [
        v =>
          !!v ||
          'Google Reviews Config Telephone Parameter is required. Get this from a hook on iCabbi standard send message.'
      ],
      googleReviewsConfigSendFromIcabbiApp: false,
      googleReviewsConfigAppKey: '',
      googleReviewsConfigSecretKey: '',
      googleReviewsConfigSendURL: '',
      googleReviewsConfigSendURLRules: [
        v =>
          !!v ||
          'Google Reviews Config Send URL is required. Get this from a hook on iCabbi standard send message.',
        v =>
          /[-a-zA-Z0-9@:%._+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_+.~#?&//=]*)/.test(
            v
          ) || 'Google Reviews Config Send URL must be valid'
      ],
      googleReviewsConfigHttpGet: false,
      googleReviewsConfigSendSuccessResponse: '',
      googleReviewsConfigSendSuccessResponseRules: [
        v =>
          !!v ||
          'Google Reviews Config Send Success Response is required. Get this from making a manual request to the send SMS service e.g. use CURL.'
      ],
      googleReviewsConfigTimeZone: 'Europe/London',
      googleReviewsConfigTimeZoneRules: [
        v =>
          !!v ||
          'Google Reviews Config Time Zone is required. Example Europe/London for a list see https://en.wikipedia.org/wiki/List_of_tz_database_time_zones.',
        v => /.+\/.+/.test(v) || 'Google Reviews Config Time Zone must be valid'
      ],
      googleReviewsConfigMultiMessageEnabled: false,
      googleReviewsConfigMessageParameter: 'm',
      googleReviewsConfigMessageParameterRules: [
        v =>
          !!v ||
          'Google Reviews Config Message Parameter is required. If multi message is NOT enabled this is ignored so can put anything. Get this from a hook on iCabbi standard send message.'
      ],
      googleReviewsConfigMultiMessageSeparator: 'SSSSS',
      googleReviewsConfigMultiMessageSeparatorRules: [
        v =>
          !!v ||
          'Google Reviews Config Multi Message Separator is required. If multi message is NOT enabled this is ignored so can put anything. This should be something unlikely to appear in a message.'
      ],
      googleReviewsConfigUseDatabaseMessage: false,
      googleReviewsConfigMessage: 'change me',
      googleReviewsConfigMessageRules: [
        v =>
          !!v ||
          'Google Reviews Config Message is required. If use database message is NOT enabled this is ignored so can put anything.'
      ],
      googleReviewsConfigSendDelayEnabled: false,
      googleReviewsConfigSendDelay: '10',
      googleReviewsConfigSendDelayRules: [
        v => !!v || 'Google Reviews Config Send Delay is required',
        v =>
          (v && v > -1) ||
          'Google Reviews Config Send Delay must be greater than 0. It is the minimum number of minutes the message is delayed before sending.'
      ],
      googleReviewsConfigDispatcherChecksEnabled: false,
      googleReviewsConfigDispatcherURL: '',
      googleReviewsConfigDispatcherURLRules: [
        v =>
          !!v ||
          'Google Reviews Config Dispatcher URL is required. Get this from iCabbi it is likely to be the same as the send URL if using iCabbi to send message requires the app key and secret key above to be set if doing checks.'
      ],
      // NOTE: There is a bug in the q-select in that it only updates if another value changes (e.g. a checkbox) (input change does not work)
      googleReviewsConfigDispatcherType: 'ICABBI',
      selectDispatcherType: ['ICABBI', 'AUTOCAB', 'AUTOCAB_V1', 'AUTOCAB_V2', 'CORDIC', 'CAB 9'],
      googleReviewsConfigBookingIdParameter: 'b',
      googleReviewsConfigBookingIdParameterRules: [
        v =>
          !!v ||
          'Google Reviews Config Booking ID Parameter is required. Get this from a hook on iCabbi standard send message (it is #trip_id).'
      ],
      googleReviewsConfigIsBookingForNowDiffMinutes: '10',
      googleReviewsConfigIsBookingForNowDiffMinutesRules: [
        v => !!v || 'Google Reviews Config Is Booking For Now Difference in Minutes is required',
        v =>
          (v && v > 0) ||
          'Google Reviews Config Is Booking For Now Difference in Minutes must be greater than  or equal to 0. It is the number of minutes that determines whether the booking is for now or a pre booking.'
      ],
      googleReviewsConfigBookingNowPickupToContactMinutes: '10',
      googleReviewsConfigBookingNowPickupToContactMinutesRules: [
        v => !!v || 'Google Reviews Config Booking for Now Pickup To Contac Minutes is required',
        v =>
          (v && v > 0) ||
          'Google Reviews Config Booking for Now Pickup To Contac Minutes must be greater than  or equal to 0. It is the allowed difference in minutes of a booking for now between pickup time and contact time.'
      ],
      googleReviewsConfigPreBookingPickupToContactMinutes: '3',
      googleReviewsConfigPreBookingPickupToContactMinutesRules: [
        v => !!v || 'Google Reviews Config Pre Booking Pickup To Contac Minutes is required',
        v =>
          (v && v > 0) ||
          'Google Reviews Config Pre Booking Pickup To Contac Minutes must be greater than  or equal to 0. It is the allowed difference in minutes of a pre booking between pickup time and contact time.'
      ],
      googleReviewsConfigReplaceTelephoneCountryCode: false,
      googleReviewsConfigReplaceTelephoneCountryCodeWith: '0',
      googleReviewsConfigReplaceTelephoneCountryCodeWithRules: [
        v =>
          !!v ||
          'Google Reviews Config Replace Telephone Country Code With is required. If replace telephone country code is NOT enabled this is ignored so can put anything normally a zero.'
      ],
      googleReviewsConfigReviewMasterSMSGatewayEnabled: false,
      googleReviewsConfigReviewMasterSMSGatewayUseMasterQueue: false,
      googleReviewsConfigReviewMasterSMSGatewayPairCode: '{code}',
      googleReviewsConfigReviewMasterSMSGatewayPairCodeQrCode: 'rmsg-pair://',

      googleReviewsConfigAlternateMessageServiceEnabled: false,
      googleReviewsConfigAlternateMessageService: '',
      selectAlternateMessageServiceType: ['', 'Message Media', 'Veezu', 'AUTOCAB_V1'],
      googleReviewsConfigAlternateMessageServiceSecret1: '',

      googleReviewsConfigCompanies: '',
      googleReviewsConfigBookingSourceMobileAppState: -1,
      googleReviewsConfigBookingSourceMobileAppStateOptions: [
        {
          label: 'All',
          value: -1
        },
        {
          label: 'Mobile App Only',
          value: 1
        },
        {
          label: 'NOT Mobile App',
          value: 0
        }
      ],

      googleReviewsConfigGoogleMyBusinessReviewReplyEnabled: false,
      googleReviewsConfigGoogleMyBusinessLocationName: '',
      googleReviewsConfigGoogleMyBusinessLocationNameRules: [],
      googleReviewsConfigGoogleMyBusinessPostalCode: '',
      googleReviewsConfigGoogleMyBusinessPostalCodeRules: [],
      googleReviewsConfigGoogleMyBusinessReplyToUnspecfifiedStarRating: false,
      googleReviewsConfigGoogleMyBusinessUnspecfifiedStarRatingReply: '',
      googleReviewsConfigGoogleMyBusinessUnspecfifiedStarRatingReplyRules: [],
      googleReviewsConfigGoogleMyBusinessReplyToOneStarRating: false,
      googleReviewsConfigGoogleMyBusinessOneStarRatingReply: '',
      googleReviewsConfigGoogleMyBusinessOneStarRatingReplyRules: [],
      googleReviewsConfigGoogleMyBusinessReplyToTwoStarRating: false,
      googleReviewsConfigGoogleMyBusinessTwoStarRatingReply: '',
      googleReviewsConfigGoogleMyBusinessTwoStarRatingReplyRules: [],
      googleReviewsConfigGoogleMyBusinessReplyToThreeStarRating: false,
      googleReviewsConfigGoogleMyBusinessThreeStarRatingReply: '',
      googleReviewsConfigGoogleMyBusinessThreeStarRatingReplyRules: [],
      googleReviewsConfigGoogleMyBusinessReplyToFourStarRating: false,
      googleReviewsConfigGoogleMyBusinessFourStarRatingReply: '',
      googleReviewsConfigGoogleMyBusinessFourStarRatingReplyRules: [],
      googleReviewsConfigGoogleMyBusinessReplyToFiveStarRating: false,
      googleReviewsConfigGoogleMyBusinessFiveStarRatingReply: '',
      googleReviewsConfigGoogleMyBusinessFiveStarRatingReplyRules: [],
      googleReviewsConfigGoogleMyBusinessFile: ref(null),
      googleReviewsConfigGoogleMyBusinessReportEnabled: false,
      emailAddress: '',
      emailAddressRules: [
        v =>
          !!v ||
          'Email Address is required',
        v =>
          /^(?=[a-zA-Z0-9@._%+-]{6,254}$)[a-zA-Z0-9._%+-]{1,64}@(?:[a-zA-Z0-9-]{1,63}\.){1,8}[a-zA-Z]{2,63}$/.test(
            v
          ) || 'email Address must be valid'
      ],
      googleReviewsClientId: '',
      ai_responses_enabled: false,
      monthly_review_analysis_enabled: false,
      contact_method: null
    }
  },
  mounted () {
    this.googleReviewsClientId = this.clientIdProp
    this.generateToken()
    this.generatePairCode()
    this.updateQrCode()
  },
  methods: {
    validate () {
      this.$refs.form.validate()
        .then(v => {
          if (v) {
            // console.log('valid')
            // console.log(this.googleReviewsClientId)
            // console.log(this.clientIdProp)
            const googleReviewsConfig = {
              id: this.googleReviewsConfigId,
              enabled: this.googleReviewsConfigEnabled,
              min_send_frequency: this.googleReviewsConfigMinSendFrequency,
              max_send_count: this.googleReviewsConfigMaxSendCount,
              max_daily_send_count: this.googleReviewsConfigMaxDailySendCount,
              token: this.googleReviewsConfigToken,
              telephone_parameter: this.googleReviewsConfigTelephoneParameter,
              send_from_icabbi_app: this.googleReviewsConfigSendFromIcabbiApp,
              app_key: this.googleReviewsConfigAppKey,
              secret_key: this.googleReviewsConfigSecretKey,
              send_url: this.googleReviewsConfigSendURL,
              http_get: this.googleReviewsConfigHttpGet,
              send_success_response: this.googleReviewsConfigSendSuccessResponse,
              time_zone: this.googleReviewsConfigTimeZone,
              multi_message_enabled: this.googleReviewsConfigMultiMessageEnabled,
              message_parameter: this.googleReviewsConfigMessageParameter,
              multi_message_separator: this.googleReviewsConfigMultiMessageSeparator,
              use_database_message: this.googleReviewsConfigUseDatabaseMessage,
              message: this.googleReviewsConfigMessage,
              send_delay_enabled: this.googleReviewsConfigSendDelayEnabled,
              send_delay: this.googleReviewsConfigSendDelay,
              dispatcher_checks_enabled: this.googleReviewsConfigDispatcherChecksEnabled,
              dispatcher_url: this.googleReviewsConfigDispatcherURL,
              dispatcher_type: this.googleReviewsConfigDispatcherType,
              booking_id_parameter: this.googleReviewsConfigBookingIdParameter,
              is_booking_for_now_diff_minutes: this.googleReviewsConfigIsBookingForNowDiffMinutes,
              booking_now_pickup_to_contact_minutes: this.googleReviewsConfigBookingNowPickupToContactMinutes,
              pre_booking_pickup_to_contact_minutes: this.googleReviewsConfigPreBookingPickupToContactMinutes,
              replace_telephone_country_code: this.googleReviewsConfigReplaceTelephoneCountryCode,
              replace_telephone_country_code_with: this.googleReviewsConfigReplaceTelephoneCountryCodeWith,
              review_master_sms_gateway_enabled: this.googleReviewsConfigReviewMasterSMSGatewayEnabled,
              review_master_sms_gateway_use_master_queue: this.googleReviewsConfigReviewMasterSMSGatewayUseMasterQueue,
              review_master_sms_gateway_pair_code: this.googleReviewsConfigReviewMasterSMSGatewayPairCode,
              alternate_message_service_enabled: this.googleReviewsConfigAlternateMessageServiceEnabled,
              alternate_message_service: this.googleReviewsConfigAlternateMessageService,
              alternate_message_service_secret1: this.googleReviewsConfigAlternateMessageServiceSecret1,
              companies: this.googleReviewsConfigCompanies,
              booking_source_mobile_app_state: this.googleReviewsConfigBookingSourceMobileAppState,
              google_my_business_review_reply_enabled: this.googleReviewsConfigGoogleMyBusinessReviewReplyEnabled,
              google_my_business_location_name: this.googleReviewsConfigGoogleMyBusinessLocationName,
              google_my_business_postal_code: this.googleReviewsConfigGoogleMyBusinessPostalCode,
              google_my_business_reply_to_unspecfified_star_rating: this.googleReviewsConfigGoogleMyBusinessReplyToUnspecfifiedStarRating,
              google_my_business_unspecfified_star_rating_reply: this.googleReviewsConfigGoogleMyBusinessUnspecfifiedStarRatingReply,
              google_my_business_reply_to_one_star_rating: this.googleReviewsConfigGoogleMyBusinessReplyToOneStarRating,
              google_my_business_one_star_rating_reply: this.googleReviewsConfigGoogleMyBusinessOneStarRatingReply,
              google_my_business_reply_to_two_star_rating: this.googleReviewsConfigGoogleMyBusinessReplyToTwoStarRating,
              google_my_business_two_star_rating_reply: this.googleReviewsConfigGoogleMyBusinessTwoStarRatingReply,
              google_my_business_reply_to_three_star_rating: this.googleReviewsConfigGoogleMyBusinessReplyToThreeStarRating,
              google_my_business_three_star_rating_reply: this.googleReviewsConfigGoogleMyBusinessThreeStarRatingReply,
              google_my_business_reply_to_four_star_rating: this.googleReviewsConfigGoogleMyBusinessReplyToFourStarRating,
              google_my_business_four_star_rating_reply: this.googleReviewsConfigGoogleMyBusinessFourStarRatingReply,
              google_my_business_reply_to_five_star_rating: this.googleReviewsConfigGoogleMyBusinessReplyToFiveStarRating,
              google_my_business_five_star_rating_reply: this.googleReviewsConfigGoogleMyBusinessFiveStarRatingReply,
              google_my_business_report_enabled: this.googleReviewsConfigGoogleMyBusinessReportEnabled,
              email_address: this.emailAddress,
              client_id: parseInt(this.googleReviewsClientId, 10),
              ai_responses_enabled: this.ai_responses_enabled,
              monthly_review_analysis_enabled: this.monthly_review_analysis_enabled,
              contact_method: this.contact_method
            }

            this.updateQrCode()

            // console.log(googleReviewsConfig)
            // this.$http.post('/auth/configadd', googleReviewsConfig).then(result => {
            api.post('/auth/configadd', googleReviewsConfig).then(result => {
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
    },
    generateToken () {
      const token = this.generateCode(64)
      this.googleReviewsConfigToken = token
      // Need this to make the change recongnised
      // this.updateConfig()
    },
    generatePairCode () {
      const pairCode = this.generateCode(12).toLowerCase()
      this.googleReviewsConfigReviewMasterSMSGatewayPairCode = pairCode
      // Need this to make the change recongnised (this will also update the QR code)
      this.updateQrCode()
    },
    generateCode (length) {
      let code = ''
      const possible =
        'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_'

      for (let i = 0; i < length; i++) {
        code += possible.charAt(Math.floor(Math.random() * possible.length))
      }
      return code
    },
    updateQrCode () {
      // console.log('rmsg-pair://' + this.googleReviewsConfigReviewMasterSMSGatewayPairCode)
      this.googleReviewsConfigReviewMasterSMSGatewayPairCodeQrCode = 'rmsg-pair://' + this.googleReviewsConfigReviewMasterSMSGatewayPairCode
    },
    readStarRatingFile () {
      this.file = this.googleReviewsConfigGoogleMyBusinessFile
      // console.log(this.file)
      const reader = new FileReader()
      if (this.file.name.includes('.txt')) {
        reader.onload = (res) => {
          // this.googleReviewsConfigGoogleMyBusinessUnspecfifiedStarRatingReply = res.target.result
          const txt = res.target.result
          // this.googleReviewsConfigGoogleMyBusinessUnspecfifiedStarRatingReply = txt
          const lines = txt.split('\n')
          const starLine = '^\\d Star\\s*'
          const startCheck = new RegExp(starLine)
          let whichStar = '-1'
          let words = ''
          for (const line of lines) {
            // console.log(line)
            let starLine = false
            if (startCheck.test(line)) {
              const ln = line
              words = ln.split(' ')
              // console.log('star line words[0]: ' + words[0])
              whichStar = words[0]
              starLine = true
            }
            // console.log('star line: ' + whichStar)
            // console.log('line: ' + line)
            // if (!starLine) {
            switch (whichStar) {
              case '1':
                this.googleReviewsConfigGoogleMyBusinessReplyToOneStarRating = true
                if (starLine) {
                  this.googleReviewsConfigGoogleMyBusinessOneStarRatingReply = ''
                } else {
                  this.googleReviewsConfigGoogleMyBusinessOneStarRatingReply += line
                }
                break
              case '2':
                this.googleReviewsConfigGoogleMyBusinessReplyToTwoStarRating = true
                if (starLine) {
                  this.googleReviewsConfigGoogleMyBusinessTwoStarRatingReply = ''
                } else {
                  this.googleReviewsConfigGoogleMyBusinessTwoStarRatingReply += line
                }
                break
              case '3':
                this.googleReviewsConfigGoogleMyBusinessReplyToThreeStarRating = true
                if (starLine) {
                  this.googleReviewsConfigGoogleMyBusinessThreeStarRatingReply = ''
                } else {
                  this.googleReviewsConfigGoogleMyBusinessThreeStarRatingReply += line
                }
                break
              case '4':
                this.googleReviewsConfigGoogleMyBusinessReplyToFourStarRating = true
                if (starLine) {
                  this.googleReviewsConfigGoogleMyBusinessFourStarRatingReply = ''
                } else {
                  this.googleReviewsConfigGoogleMyBusinessFourStarRatingReply += line
                }
                break
              case '5':
                this.googleReviewsConfigGoogleMyBusinessReplyToFiveStarRating = true
                if (starLine) {
                  this.googleReviewsConfigGoogleMyBusinessFiveStarRatingReply = ''
                } else {
                  this.googleReviewsConfigGoogleMyBusinessFiveStarRatingReply += line
                }
                break
            }
            // }
          }
        }
        reader.onerror = (err) => console.log(err)
        reader.readAsText(this.file)
      } else {
        // this.googleReviewsConfigGoogleMyBusinessUnspecfifiedStarRatingReply = 'check the console for file output'
        alert('check the console for file output')
        reader.onload = (res) => {
          console.log(res.target.result)
        }
        reader.onerror = (err) => console.log(err)
        reader.readAsText(this.file)
      }
    }
  }
}
</script>
