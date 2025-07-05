<template>
  <div class="dataviz-skeleton">
    <!-- Main content area that fills the space properly -->
    <div class="skeleton-main">
      <q-skeleton type="rect" width="100%" height="250px" class="main-block" />
    </div>
    
    <!-- Secondary content area -->
    <div class="skeleton-secondary">
      <q-skeleton type="rect" width="100%" height="200px" class="secondary-block" />
    </div>
  </div>
</template>

<script setup lang="ts">
// No script needed for this simple skeleton
</script>

<style scoped>
.dataviz-skeleton {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.skeleton-main {
  width: 100%;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  background: #fafafa;
  padding: 1rem;
}

.skeleton-secondary {
  width: 100%;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  background: #fafafa;
  padding: 1rem;
}

.main-block,
.secondary-block {
  border-radius: 6px;
}

/* Enhanced pulsating animation */
.q-skeleton {
  animation: skeleton-pulse 1.5s ease-in-out infinite alternate;
}

.main-block {
  animation-delay: 0s;
}

.secondary-block {
  animation-delay: 0.3s;
}

@keyframes skeleton-pulse {
  0% {
    opacity: 1;
  }
  100% {
    opacity: 0.4;
  }
}
</style>