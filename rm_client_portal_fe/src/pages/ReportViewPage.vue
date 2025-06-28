<template>
  <q-page padding>

    <!-- Loading spinner -->
    <div v-if="isLoading" class="row items-center justify-center q-pa-lg">
      <q-spinner-dots size="40px" color="primary" />
      <div class="q-ml-md">Loading report...</div>
    </div>

    <!-- Error message -->
    <div v-else-if="error" class="text-center q-pa-lg">
      <q-icon name="error" size="48px" color="negative" class="q-mb-md" />
      <div class="text-subtitle1 text-negative">Failed to load report</div>
      <div class="text-body2 text-grey-6 q-mb-md">{{ error }}</div>
      <q-btn color="primary" label="Back to Reports" @click="router.push('/reports')" />
    </div>

    <!-- Report content -->
    <div v-else class="report-container">
      <!-- Back button -->
      <div class="q-mb-md">
        <q-btn flat color="primary" icon="arrow_back" label="Back to Reports" @click="router.push('/reports')" />
      </div>

      <!-- Report iframe -->
      <iframe ref="reportIframe" :src="reportHtml" class="report-iframe" scrolling="no" @load="onIframeLoad"></iframe>
    </div>

  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { apiService } from 'src/services/api-service';

const router = useRouter();
const route = useRoute();
const isLoading = ref(true);
const error = ref('');
const reportHtml = ref('');
const reportIframe = ref<HTMLIFrameElement>();

// Handle iframe load event
const onIframeLoad = () => {
  // With blob URLs, we should be able to access iframe content for height adjustment
  if (reportIframe.value && reportIframe.value.contentWindow) {
    try {
      const body = reportIframe.value.contentWindow.document.body;
      const html = reportIframe.value.contentWindow.document.documentElement;
      const height = Math.max(
        body.scrollHeight,
        body.offsetHeight,
        html.clientHeight,
        html.scrollHeight,
        html.offsetHeight
      );
      reportIframe.value.style.height = height + 'px';
      console.log('Iframe height adjusted to:', height + 'px');
    } catch (e) {
      // Fallback to a reasonable height if cross-origin restrictions still apply
      reportIframe.value.style.height = '1200px';
      console.log('Using fallback height due to:', e);
    }
  }
};

// Load report HTML from API
const loadReport = async () => {
  try {
    isLoading.value = true;
    error.value = '';
    const reportId = route.params.id as string;
    const htmlContent = await apiService.getReportHTML(reportId);

    // Create blob URL for better CSP compatibility and JavaScript execution
    const blob = new Blob([htmlContent], { type: 'text/html' });
    const blobUrl = URL.createObjectURL(blob);
    reportHtml.value = blobUrl;
  } catch (err: any) {
    console.error('Failed to load report:', err);
    error.value = err.message || 'An error occurred while loading the report';
  } finally {
    isLoading.value = false;
  }
};

// Cleanup blob URL when component unmounts to prevent memory leaks
onUnmounted(() => {
  if (reportHtml.value && reportHtml.value.startsWith('blob:')) {
    URL.revokeObjectURL(reportHtml.value);
  }
});

// Load report on mount
onMounted(() => {
  loadReport();
});

</script>

<style scoped>
.report-container {
  max-width: 100%;
  height: 100%;
}

.report-iframe {
  width: 100%;
  min-height: 1200px;
  border: none;
  display: block;
}
</style>