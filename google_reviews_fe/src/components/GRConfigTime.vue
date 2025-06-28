<template>
  <div id="g-r-config-time">
    <!-- <q-input v-model="googleReviewsConfigTimeId" label="Google Reviews Config Time ID" required disable v-show="false" /> -->
    <q-checkbox v-model="googleReviewsConfigTimeEnabled" label="Google Reviews Config Time Enabled" @update:model-value="updateTime" />
    <q-input v-model="googleReviewsConfigTimeStart" :rules="googleReviewsConfigTimeStartRules" label="Google Reviews Config Time Start time" required @update:model-value="updateTime" />
    <q-input v-model="googleReviewsConfigTimeEnd" :rules="googleReviewsConfigTimeEndRules" label="Google Reviews Config Time End time" required @update:model-value="updateTime" />
    <q-checkbox v-model="googleReviewsConfigTimeSunday" label="Sunday" @update:model-value="updateTime" />
    <q-checkbox v-model="googleReviewsConfigTimeMonday" label="Monday" @update:model-value="updateTime" />
    <q-checkbox v-model="googleReviewsConfigTimeTuesday" label="Tuesday" @update:model-value="updateTime" />
    <q-checkbox v-model="googleReviewsConfigTimeWednesday" label="Wednesday" @update:model-value="updateTime" />
    <q-checkbox v-model="googleReviewsConfigTimeThursday" label="Thursday" @update:model-value="updateTime" />
    <q-checkbox v-model="googleReviewsConfigTimeFriday" label="Friday" @update:model-value="updateTime" />
    <q-checkbox v-model="googleReviewsConfigTimeSaturday" label="Saturday" @update:model-value="updateTime" />

    <q-separator />
  </div>
</template>

<script>
export default {
  name: 'GRConfigTime',
  props: ['grct', 'index'],
  // so this component can use parent validator
  inject: ['$validator'],
  data () {
    return {
      googleReviewsConfigTimeId: '',
      // googleReviewsConfigTimeIdRules: [
      //   v => !!v || 'Google Reviews Config Time ID is required'
      // ],
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
    this.fetchTime()
  },
  methods: {
    fetchTime: function () {
      if (this.grct) {
        this.googleReviewsConfigTimeId = this.grct.id
        this.googleReviewsConfigTimeEnabled = this.grct.enabled
        this.googleReviewsConfigTimeStart = this.grct.start
        this.googleReviewsConfigTimeEnd = this.grct.end
        this.googleReviewsConfigTimeSunday = this.grct.sunday
        this.googleReviewsConfigTimeMonday = this.grct.monday
        this.googleReviewsConfigTimeTuesday = this.grct.tuesday
        this.googleReviewsConfigTimeWednesday = this.grct.wednesday
        this.googleReviewsConfigTimeThursday = this.grct.thursday
        this.googleReviewsConfigTimeFriday = this.grct.friday
        this.googleReviewsConfigTimeSaturday = this.grct.saturday
        this.googleReviewsConfigId = this.grct.google_reviews_config_id
      } else {
        // TODO: initialise with empty data
      }
    },
    updateTime () {
      const tm = {
        id: this.googleReviewsConfigTimeId,
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
      }
      this.$emit('updateGRCT', this.index, tm)
    }
  },
  emits: ['updateGRCT']
}
</script>
