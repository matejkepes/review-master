<template>
  <q-page padding>
    <!-- max 800px width, centered -->
    <div style="max-width: 800px; margin: 0 auto;">

      <!-- Currently Viewing Client Indicator -->
      <ClientViewingIndicator />

      <div class="text-h6 q-mb-md">Monthly Reports</div>

      <!-- Loading spinner -->
      <div v-if="isLoading" class="row items-center justify-center q-pa-lg">
        <SmartLoadingSpinner loadingType="reports" />
      </div>

      <!-- Reports table -->
      <div v-else-if="reports.length > 0">
        <q-table
          flat
          bordered
          :rows="reports"
          :columns="columns"
          row-key="report_id"
          class="q-mt-md"
        >
          <template v-slot:body-cell-actions="props">
            <q-td :props="props">
              <q-btn
                color="primary"
                label="View Report"
                size="sm"
                @click="viewReport(props.row.report_id)"
              />
            </q-td>
          </template>
        </q-table>
      </div>

      <!-- No reports message -->
      <div v-else class="text-center q-pa-lg">
        <q-icon name="description" size="48px" color="grey-5" class="q-mb-md" />
        <div class="text-subtitle1 text-grey-7">No reports found</div>
        <div class="text-body2 text-grey-6">
          Monthly reports will appear here once they are generated.
        </div>
      </div>

    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useStore } from 'stores/store';
import { apiService } from 'src/services/api-service';
import ClientViewingIndicator from 'src/components/ClientViewingIndicator.vue';
import SmartLoadingSpinner from 'src/components/SmartLoadingSpinner.vue';

const router = useRouter();
const store = useStore();
const isLoading = ref(true);

// Define the report type
interface Report {
  report_id: number;
  period_start: string;
  period_end: string;
  generated_at: string;
  client_name: string;
}

const reports = ref<Report[]>([]);

// Table columns configuration
const columns = [
  {
    name: 'period',
    label: 'Report Period',
    align: 'left' as const,
    field: (row: Report) => `${formatDate(row.period_start)} - ${formatDate(row.period_end)}`,
    sortable: true,
  },
  {
    name: 'generated_at',
    label: 'Generated',
    align: 'left' as const,
    field: 'generated_at',
    sortable: true,
  },
  {
    name: 'client_name',
    label: 'Client',
    align: 'left' as const,
    field: 'client_name',
    sortable: true,
  },
  {
    name: 'actions',
    label: 'Actions',
    align: 'center' as const,
    field: '',
  },
];

// Format date helper
const formatDate = (dateStr: string) => {
  const date = new Date(dateStr);
  return date.toLocaleDateString('en-US', { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric' 
  });
};

// Load reports from API
const loadReports = async () => {
  try {
    isLoading.value = true;
    const clientId = store.selectedClientId;
    const response = await apiService.getReports(clientId || undefined);
    reports.value = response.reports || [];
  } catch (error) {
    console.error('Failed to load reports:', error);
    // Could add a notification here
  } finally {
    isLoading.value = false;
  }
};

// Navigate to report view
const viewReport = (reportId: number) => {
  router.push(`/reports/${reportId}`);
};

// Watch for client changes and reload reports
watch(() => store.selectedClientId, () => {
  loadReports();
});

// Load reports on mount
onMounted(() => {
  loadReports();
});

</script>

<style scoped>
.q-table {
  border-radius: 8px;
  border: 1px solid #e0e0e0;
}
</style>