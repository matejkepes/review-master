<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="toggleLeftDrawer"
        />

        <q-toolbar-title>
          Google Reviews Front End
        </q-toolbar-title>

        <div>v0.2</div>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
      class="bg-grey-1"
    >
      <q-list>
        <q-item-label
          header
          class="text-grey-8"
        >
          Essential Links
        </q-item-label>

        <span v-if="isLoggedIn">
          <q-item clickable @click="logout">
            <q-item-section avatar>
              <q-icon name="person" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ displayCurrentUser() }}</q-item-label>
              <q-item-label caption>is logged in, select to Logout</q-item-label>
            </q-item-section>
          </q-item>
          <span v-if="isAdmin">
            <q-item clickable to="/clients">
              <q-item-section avatar>
                <q-icon name="people" />
              </q-item-section>
              <q-item-section>
                <q-item-label>Clients</q-item-label>
                <q-item-label caption>client list</q-item-label>
              </q-item-section>
            </q-item>
            <q-item clickable to="/sendtest">
              <q-item-section avatar>
                <q-icon name="send" />
              </q-item-section>
              <q-item-section>
                <q-item-label>Send Test</q-item-label>
                <q-item-label caption>send a test message to get response</q-item-label>
              </q-item-section>
            </q-item>
            <q-item clickable to="/stats">
              <q-item-section avatar>
                <q-icon name="bar_chart" />
              </q-item-section>
              <q-item-section>
                <q-item-label>Stats</q-item-label>
                <q-item-label caption>some statistics</q-item-label>
              </q-item-section>
            </q-item>
            <q-item clickable to="/statsnew">
              <q-item-section avatar>
                <q-icon name="bar_chart" />
              </q-item-section>
              <q-item-section>
                <q-item-label>Stats (New)</q-item-label>
                <q-item-label caption>some new statistics</q-item-label>
              </q-item-section>
            </q-item>
            <q-item clickable to="/checklog">
              <q-item-section avatar>
                <q-icon name="error" />
              </q-item-section>
              <q-item-section>
                <q-item-label>Errors</q-item-label>
                <q-item-label caption>Error checks</q-item-label>
              </q-item-section>
            </q-item>
            <q-item clickable to="/checknothingsent">
              <q-item-section avatar>
                <q-icon name="clear" />
              </q-item-section>
              <q-item-section>
                <q-item-label>Nothing Sent Check</q-item-label>
                <q-item-label caption>No messages sent check</q-item-label>
              </q-item-section>
            </q-item>
            <q-item clickable to="/gmbreport">
              <q-item-section avatar>
                <q-icon name="description" />
              </q-item-section>
              <q-item-section>
                <q-item-label>My Business Report</q-item-label>
                <q-item-label caption>Google My Business report</q-item-label>
              </q-item-section>
            </q-item>
            <q-item clickable to="/users">
              <q-item-section avatar>
                <q-icon name="people" />
              </q-item-section>
              <q-item-section>
                <q-item-label>Users</q-item-label>
                <q-item-label caption>user list</q-item-label>
              </q-item-section>
            </q-item>
          </span>
        </span>
        <span v-else>
          <q-item to='/login' exact>
            <q-item-section avatar>
              <q-icon name="input" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Login</q-item-label>
              <q-item-label caption>login to the system</q-item-label>
            </q-item-section>
          </q-item>
        </span>

      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script>
import { defineComponent, ref } from 'vue'

export default defineComponent({
  name: 'MainLayout',

  computed: {
    isLoggedIn: function () {
      // console.log(this.$store.getters)
      // return this.$store.getters.isLoggedIn
      return this.$store.getters['login/isLoggedIn']
    },
    isAdmin: function () {
      // console.log(this.$store.getters['login/roleCurrentUser'])
      return this.$store.getters['login/roleCurrentUser'] === 'admin'
    }
  },

  methods: {
    // openURL,
    logout: function () {
      // this.$store.dispatch('logout').then(() => {
      this.$store.dispatch('login/logout').then(() => {
        this.$router.push('/login')
      })
    },
    displayCurrentUser: function () {
      // return this.$store.getters.displayCurrentUser
      return this.$store.getters['login/displayCurrentUser']
    }
  },

  setup () {
    const leftDrawerOpen = ref(false)

    return {
      leftDrawerOpen,
      toggleLeftDrawer () {
        leftDrawerOpen.value = !leftDrawerOpen.value
      }
    }
  }
})
</script>
