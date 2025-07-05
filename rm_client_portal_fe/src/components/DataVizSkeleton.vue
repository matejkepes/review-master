<template>
  <div class="dataviz-skeleton">
    <!-- Profile Interactions Section -->
    <div class="skeleton-section">
      <q-skeleton type="text" width="200px" height="24px" class="section-header" />
      <div class="row q-col-gutter-lg q-mt-md">
        <div class="col-6">
          <div class="metric-card-skeleton">
            <q-skeleton type="rect" width="120px" height="60px" class="metric-number" />
            <q-skeleton type="text" width="140px" height="16px" class="metric-label" />
          </div>
        </div>
        <div class="col-6">
          <div class="metric-card-skeleton">
            <q-skeleton type="rect" width="120px" height="60px" class="metric-number" />
            <q-skeleton type="text" width="160px" height="16px" class="metric-label" />
          </div>
        </div>
      </div>
    </div>
    
    <!-- Rating Distribution Section -->
    <div class="skeleton-section">
      <div class="row items-center q-mb-md justify-between">
        <q-skeleton type="text" width="180px" height="24px" class="section-header" />
        <q-skeleton type="rect" width="32px" height="32px" class="menu-button" />
      </div>
      <div class="chart-skeleton">
        <div class="bar-chart-skeleton">
          <div class="bar-skeleton bar-1"></div>
          <div class="bar-skeleton bar-2"></div>
          <div class="bar-skeleton bar-3"></div>
          <div class="bar-skeleton bar-4"></div>
          <div class="bar-skeleton bar-5"></div>
        </div>
        <div class="chart-labels">
          <q-skeleton type="text" width="20px" height="12px" />
          <q-skeleton type="text" width="20px" height="12px" />
          <q-skeleton type="text" width="20px" height="12px" />
          <q-skeleton type="text" width="20px" height="12px" />
          <q-skeleton type="text" width="20px" height="12px" />
        </div>
      </div>
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
  gap: 3rem;
}

.skeleton-section {
  width: 100%;
}

.section-header {
  margin-bottom: 1rem;
}

/* Profile Interactions Metric Cards */
.metric-card-skeleton {
  text-align: center;
  padding: 2rem 1rem;
  background: white;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  height: 140px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 1rem;
}

.metric-number {
  border-radius: 8px;
}

.metric-label {
  border-radius: 4px;
}

/* Rating Distribution Chart */
.chart-skeleton {
  background: white;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  padding: 1.5rem;
  height: 350px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.bar-chart-skeleton {
  display: flex;
  align-items: end;
  justify-content: space-around;
  height: 280px;
  gap: 1rem;
  padding: 0 2rem;
}

.bar-skeleton {
  background: #f0f0f0;
  border-radius: 4px 4px 0 0;
  width: 40px;
  animation: skeleton-pulse 1.5s ease-in-out infinite alternate;
}

.bar-1 { height: 20px; animation-delay: 0s; }
.bar-2 { height: 10px; animation-delay: 0.2s; }
.bar-3 { height: 15px; animation-delay: 0.4s; }
.bar-4 { height: 50px; animation-delay: 0.6s; }
.bar-5 { height: 200px; animation-delay: 0.8s; }

.chart-labels {
  display: flex;
  justify-content: space-around;
  padding: 0 2rem;
  margin-top: 1rem;
}

.menu-button {
  border-radius: 50%;
}

/* Enhanced pulsating animation */
.q-skeleton {
  animation: skeleton-pulse 1.5s ease-in-out infinite alternate;
}

.section-header {
  animation-delay: 0s;
}

.metric-number {
  animation-delay: 0.2s;
}

.metric-label {
  animation-delay: 0.4s;
}

.menu-button {
  animation-delay: 0.6s;
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