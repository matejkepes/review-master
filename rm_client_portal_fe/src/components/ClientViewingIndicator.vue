<template>
  <div v-if="selectedClient" class="q-mb-md">
    <div class="row items-center">
      <q-icon name="business" color="grey-6" size="16px" class="q-mr-xs" />
      <span class="text-caption text-grey-6">
        Currently viewing: <span class="text-weight-medium text-grey-8">{{ selectedClient.name }}</span>
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useStore } from 'stores/store';

const store = useStore();

// Get the selected client from store
const selectedClient = computed(() =>
  store.clients.find(client => client.id === store.selectedClientId)
);
</script>