<template>
  <q-layout view="lHh Lpr lFf">
    <!-- Unified Header -->
    <q-header elevated class="bg-primary">
      <q-toolbar class="q-px-md">
        <!-- Left Section: Logo + Navigation -->
        <div class="row items-center no-wrap">
          <!-- Logo -->
          <img src="~src/assets/logo-nav.png" alt="logo" height="40px" class="q-mr-lg gt-sm" />
          
          <!-- Navigation Tabs -->
          <q-tabs 
            v-model="activeTab" 
            class="text-white navigation-tabs" 
            active-color="white" 
            indicator-color="white"
            dense
            no-caps
          >
            <q-tab name="dashboard" label="Dashboard" @click="navigateToTab('dashboard')" />
            <q-tab name="reports" label="Reports" @click="navigateToTab('reports')" />
          </q-tabs>
        </div>

        <!-- Spacer -->
        <q-space />

        <!-- Right Section: Client + User Context -->
        <div class="row items-center no-wrap q-gutter-sm">
          <!-- Client Selector -->
          <q-btn
            v-if="store.clients.length > 1"
            flat
            no-caps
            class="client-selector-btn"
            aria-label="Select Client"
          >
            <q-avatar size="28px" class="q-mr-sm app-avatar-base app-avatar-inverse">
              {{ currentClientName.charAt(0).toUpperCase() }}
            </q-avatar>
            <span class="client-name">{{ currentClientName }}</span>
            <q-icon name="expand_more" size="16px" class="q-ml-xs" />

            <!-- Client Selection Dropdown -->
            <q-menu 
              anchor="bottom right" 
              self="top right"
              :offset="[0, 8]"
              class="client-dropdown"
            >
              <q-list style="min-width: 280px; max-width: 400px;">
                <!-- Header -->
                <q-item-label header class="text-weight-bold text-grey-8 q-py-sm">
                  <q-icon name="business" class="q-mr-sm" />
                  Select Client
                </q-item-label>
                
                <q-separator />

                <!-- Client List -->
                <q-item
                  v-for="client in store.clients"
                  :key="client.id"
                  clickable
                  v-ripple
                  @click="selectClient(client.id)"
                  :active="client.id === store.selectedClientId"
                  active-class="bg-primary text-white"
                  class="client-dropdown-item"
                >
                  <q-item-section avatar>
                    <q-avatar 
                      size="36px"
                      :class="[
                        'app-avatar-base',
                        client.id === store.selectedClientId ? 'app-avatar-inverse' : 'app-avatar-primary'
                      ]"
                    >
                      {{ client.name.charAt(0).toUpperCase() }}
                    </q-avatar>
                  </q-item-section>
                  
                  <q-item-section>
                    <q-item-label class="text-weight-medium">{{ client.name }}</q-item-label>
                    <q-item-label caption v-if="client.id === store.selectedClientId">
                      Current client
                    </q-item-label>
                  </q-item-section>
                  
                  <q-item-section side v-if="client.id === store.selectedClientId">
                    <q-icon name="check_circle" color="white" size="20px" />
                  </q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>

          <!-- User Menu -->
          <q-btn flat round dense class="user-menu-btn">
            <q-avatar size="32px" class="app-avatar-base" style="background: rgba(255,255,255,0.2); color: white;">
              {{ userEmail.charAt(0).toUpperCase() }}
            </q-avatar>

            <q-menu>
              <q-list style="min-width: 220px">
                <!-- User Email -->
                <q-item>
                  <q-item-section avatar>
                    <q-avatar size="40px" class="app-avatar-base app-avatar-primary">
                      {{ userEmail.charAt(0).toUpperCase() }}
                    </q-avatar>
                  </q-item-section>
                  <q-item-section>
                    <q-item-label caption>Signed in as</q-item-label>
                    <q-item-label class="text-weight-medium">{{ userEmail }}</q-item-label>
                  </q-item-section>
                </q-item>

                <q-separator class="q-my-sm" />

                <!-- Sign Out Button -->
                <q-item clickable v-close-popup @click="handleSignOut">
                  <q-item-section avatar>
                    <q-icon name="logout" />
                  </q-item-section>
                  <q-item-section>Sign Out</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </div>
      </q-toolbar>
    </q-header>


    <q-page-container class="q-px-lg q-py-md">
      <!-- Show loading spinner while clients are being fetched -->
      <div v-if="isLoading" class="row items-center justify-center" style="height: 100vh">
        <SmartLoadingSpinner loadingType="general" />
      </div>

      <!-- Show router view only after clients are loaded -->
      <router-view v-else />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue';
import { useStore } from 'stores/store';
import { apiService } from 'src/services/api-service';
import { useRouter, useRoute } from 'vue-router';

const router = useRouter();
const route = useRoute();
const store = useStore();
const isLoading = ref(true);
const activeTab = ref('dashboard');

// Get user email from localStorage
const userEmail = computed(() => {
  const user = localStorage.getItem('USER');
  return user ? JSON.parse(user).email : '';
});

// Get current client name
const currentClientName = computed(() => {
  const client = store.clients.find(c => c.id === store.selectedClientId);
  return client ? client.name : '';
});

// Handle sign out
const handleSignOut = () => {
  localStorage.removeItem('USER');
  router.push('/login');
};

// Select client from dropdown
const selectClient = (clientId: number) => {
  store.setSelectedClient(clientId);
};

// Handle navigation tab change
const navigateToTab = (tabName: string) => {
  activeTab.value = tabName;
  if (tabName === 'dashboard') {
    router.push('/dashboard');
  } else if (tabName === 'reports') {
    router.push('/reports');
  }
};

// Set active tab based on current route
const updateActiveTab = () => {
  if (route.path.startsWith('/reports')) {
    activeTab.value = 'reports';
  } else {
    activeTab.value = 'dashboard';
  }
};

onMounted(async () => {
  try {
    // Check if clients are already in store
    if (!store.clients || store.clients.length === 0) {
      const response = await apiService.getClients();

      if (response.clients) {
        store.setClients(response.clients);
      } else {
        console.error('No clients found in response');
        localStorage.removeItem('USER');
        router.push('/login');
      }
    }
    isLoading.value = false;
    updateActiveTab(); // Set the active tab based on current route
  } catch (error) {
    console.error('Failed to load clients:', error);
    isLoading.value = false;
  }
});
</script>

<style scoped>
/* Shared Avatar Utility Classes */
.app-avatar-base {
  font-weight: 600;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

.app-avatar-primary {
  background: var(--q-primary);
  color: white;
}

.app-avatar-inverse {
  background: white;
  color: var(--q-primary);
}

/* Header Layout */
.q-toolbar {
  min-height: 64px;
  padding: 0 16px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

/* Navigation Tabs */
.navigation-tabs {
  margin-left: 0;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 6px;
  padding: 2px;
  flex-shrink: 0;
}

.navigation-tabs .q-tab {
  padding: 0 16px;
  min-height: 36px;
  font-weight: 500;
  border-radius: 4px;
  transition: all 0.2s ease;
  font-size: 14px;
}

.navigation-tabs .q-tab--active {
  background: rgba(255, 255, 255, 0.15);
}

/* Client selector button styles - Desktop (pill shape) */
.client-selector-btn {
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 20px;
  padding: 4px 12px 4px 4px;
  transition: all 0.2s ease;
  height: 40px;
  margin-right: 12px;
  display: flex;
  align-items: center;
  min-width: fit-content;
}

.client-selector-btn:hover {
  background: rgba(255, 255, 255, 0.18);
  border-color: rgba(255, 255, 255, 0.3);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.client-selector-btn .client-name {
  color: white;
  font-weight: 500;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 140px;
  margin-left: 10px;
}

.client-selector-btn .q-avatar {
  width: 28px !important;
  height: 28px !important;
  flex-shrink: 0;
}

.client-selector-btn .q-icon {
  color: white;
  opacity: 0.9;
  margin-left: 6px;
}

/* User menu button */
.user-menu-btn {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 50%;
  padding: 0;
  transition: all 0.2s ease;
  height: 40px;
  width: 40px;
  min-width: 40px;
  overflow: hidden;
}

.user-menu-btn:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.25);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.user-menu-btn .q-avatar {
  background: rgba(255, 255, 255, 0.2) !important;
  width: 32px !important;
  height: 32px !important;
  margin: 0 auto;
}

.user-menu-btn:hover .q-avatar {
  background: rgba(255, 255, 255, 0.3) !important;
}

/* Simple responsive - just hide client name on small screens */
@media (max-width: 640px) {
  /* Hide pill elements to make it a circle */
  .client-selector-btn .client-name,
  .client-selector-btn .q-icon {
    display: none;
  }
  
  /* Transform client selector to circle like user avatar */
  .client-selector-btn {
    background: rgba(255, 255, 255, 0.08) !important;
    border: 1px solid rgba(255, 255, 255, 0.15) !important;
    border-radius: 50% !important;
    padding: 0 !important;
    height: 40px !important;
    width: 40px !important;
    min-width: 40px !important;
    max-width: 40px !important;
    margin-right: 12px;
  }
  
  .client-selector-btn:hover {
    background: rgba(255, 255, 255, 0.15) !important;
    border-color: rgba(255, 255, 255, 0.25) !important;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }
  
  .client-selector-btn .q-avatar {
    width: 32px !important;
    height: 32px !important;
    margin: 0 auto;
  }
  
  .client-selector-btn:hover .q-avatar {
    background: white !important;
  }
}

/* Client dropdown styles */
.client-dropdown {
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.15);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.client-dropdown .q-list {
  padding: 8px;
  background: white;
}

.client-dropdown-item {
  padding: 12px 16px;
  margin: 4px 0;
  border-radius: 8px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.client-dropdown-item:hover {
  background-color: rgba(25, 118, 210, 0.04);
  transform: translateX(2px);
}

.client-dropdown-item.q-item--active {
  background-color: var(--q-primary) !important;
  margin: 4px 0;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(25, 118, 210, 0.3);
}

.client-dropdown .q-separator {
  margin: 4px 0 8px 0;
}

.client-dropdown .q-avatar {
  border: 2px solid rgba(0, 0, 0, 0.08);
  transition: all 0.2s ease;
}

.client-dropdown-item.q-item--active .q-avatar {
  border-color: rgba(255, 255, 255, 0.3);
}

.client-dropdown-item:hover .q-avatar {
  transform: scale(1.05);
}

/* User menu styles */
.user-menu-btn .q-menu {
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.15);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.user-menu-btn .q-list {
  padding: 8px;
  background: white;
}

.user-menu-btn .q-item {
  padding: 12px 16px;
  margin: 4px 0;
  border-radius: 8px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.user-menu-btn .q-item:hover {
  background-color: rgba(25, 118, 210, 0.04);
  transform: translateX(2px);
}

.user-menu-btn .q-item__label--caption {
  color: #666;
  font-size: 12px;
}

.user-menu-btn .q-separator {
  margin: 8px 0;
  background-color: rgba(0, 0, 0, 0.08);
}

/* Mobile dropdown adjustments */
@media (max-width: 640px) {
  .client-dropdown .q-list {
    min-width: 240px;
  }
  
  .user-menu-btn .q-list {
    min-width: 180px;
  }
}
</style>
