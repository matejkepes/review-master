<template>
  <!-- centered login form with email and password fields, take up 100% width on mobile, 50% on rest, text fields full width -->
  <q-page class="row items-center justify-center">
    <!-- width 90% on mobile, 60% on rest -->
    <q-card class="col-xs-10 col-md-4 q-pa-md q-pb-lg">

      <!-- review master logo h1 -->
      <q-card-section class="text-h6 text-center">
        <!-- logo -->
        <img src="~src/assets/logo-login.png" alt="logo" class="q-mb-md" width="100%" />
      </q-card-section>


      <q-card-section>
        <q-input v-model="email" label="Email" outlined dense clearable type="email" />
      </q-card-section>
      <q-card-section>
        <q-input v-model="password" label="Password" outlined dense clearable type="password" />
      </q-card-section>
      <q-card-section>
        <q-btn color="primary" label="Login" @click="login" class="full-width" />
      </q-card-section>

    </q-card>
  </q-page>

</template>
<script setup lang="ts">
import { ref } from 'vue';
import { ApiError, useApiService } from 'src/services/api-service';
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';

const router = useRouter();
const $q = useQuasar();
const email = ref('');
const password = ref('');

const { apiService } = useApiService();

const login = async () => {
  try {
    const response = await apiService.login(email.value, password.value);
    // if response is not an instance of ApiError

    // we know it's a LoginResponse
    localStorage.setItem('USER', JSON.stringify({
      email: email.value,
      token: response.token
    }));

    // Navigate to dashboard
    router.push('/dashboard');

  } catch (error) {
    let errorMessage = 'An error occurred. Please try again later.';

    if (error instanceof ApiError) {
      if (error.statusCode == 401) {
        errorMessage = 'Invalid email or password';
      }
    }
    $q.notify({
      message: errorMessage,
      color: 'negative',
      position: 'top',
      timeout: 2000
    });

  }
};

</script>
