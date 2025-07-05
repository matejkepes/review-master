<template>
  <q-page padding>

    <!-- max 800px width, centered -->
    <div style="max-width: 800px; margin: 0 auto;">

      <!-- Currently Viewing Client Indicator -->
      <ClientViewingIndicator />

      <!-- Performance Overview Section -->
      <div class="q-mt-lg">
        <div class="row items-center q-mb-md justify-between">
          <div class="text-h6">Performance Overview</div>
          <q-select v-model="selectedPeriod" :options="periodOptions" outlined dense options-dense class="q-ml-md"
            style="width: 200px" color="primary" />
        </div>

        <!-- Stats content will go here -->
        <div v-if="isLoading">
          <SmartLoadingSpinner loadingType="dashboard" />
        </div>
        <div v-else-if="hasStatsError" class="text-center q-pa-lg">
          <q-icon name="error_outline" size="48px" color="negative" class="q-mb-md" />
          <div class="text-subtitle1 text-grey-7">Failed to load statistics</div>
          <q-btn color="primary" label="Retry" @click="loadStats" class="q-mt-md" />
        </div>
        <transition name="fade" appear>
          <div v-if="!isLoading && !hasStatsError" class="row q-col-gutter-md">
            <div class="col-12">
              <apexchart type="area" height="350" :options="chartOptions" :series="chartSeries" />
            </div>
          </div>
        </transition>
      </div>

      <!-- Reviews Section -->
      <div class="q-mt-xl">

        <div v-if="isLoadingReviews">
          <!-- Section-level loading message -->
          <SmartLoadingSpinner loadingType="reviews" :sectionLevel="true" />
          
          <!-- Single unified content area -->
          <div class="unified-skeleton-container">
            <DataVizSkeleton />
          </div>
        </div>
        <div v-else-if="hasReviewsError" class="text-center q-pa-lg">
          <q-icon name="error_outline" size="48px" color="negative" class="q-mb-md" />
          <div class="text-subtitle1 text-grey-7">Failed to load reviews</div>
          <q-btn color="primary" label="Retry" @click="loadReviews" class="q-mt-md" />
        </div>
        <transition name="fade" appear>
          <div v-if="!isLoadingReviews && !hasReviewsError">
            <!-- Profile Interactions -->
            <div class="q-mt-xl">
              <div class="text-h6 q-mb-lg">Profile Interactions</div>
              <div class="row q-col-gutter-lg justify-center">
                <div class="col-6">
                  <div class="interaction-metric-card">
                    <div class="interaction-number">{{ insightsChartSeries[0]?.data[0] || 0 }}</div>
                    <div class="interaction-label text-grey-6">WEBSITE CLICKS</div>
                  </div>
                </div>
                <div class="col-6">
                  <div class="interaction-metric-card">
                    <div class="interaction-number">{{ insightsChartSeries[0]?.data[1] || 0 }}</div>
                    <div class="interaction-label text-grey-6">CALL BUTTON CLICKS</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Rating Distribution -->
            <div class="q-mt-xl">
              <div class="row items-center q-mb-md justify-between">
                <div class="text-h6">Rating Distribution</div>
                <q-btn flat round icon="menu" color="grey-6" />
              </div>
              <apexchart type="bar" height="350" :options="reviewChartOptions" :series="reviewChartSeries" />
            </div>
          </div>
        </transition>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useStore } from 'stores/store';
import { useApiService, type UserStatsResponse } from 'src/services/api-service';
import ClientViewingIndicator from 'src/components/ClientViewingIndicator.vue';
import SmartLoadingSpinner from 'src/components/SmartLoadingSpinner.vue';
import DataVizSkeleton from 'src/components/DataVizSkeleton.vue';
import { apiCache } from 'src/utils/apiCache';

const { apiService } = useApiService();


const store = useStore();
const isLoading = ref(false);
const isLoadingReviews = ref(false);
const hasStatsError = ref(false);
const hasReviewsError = ref(false);

// Add this type near the top of the script
type ChartDataPoint = {
  x: number;
  y: number;
};

// Chart data
const chartOptions = ref({
  chart: {
    type: 'area',
    toolbar: {
      show: false
    }
  },
  dataLabels: {
    enabled: false
  },
  stroke: {
    curve: 'smooth',
    width: 2
  },
  xaxis: {
    type: 'datetime'
  },
  yaxis: [
    {
      title: {
        text: 'Count'
      }
    }
  ],
  tooltip: {
    x: {
      format: 'dd MMM yyyy'
    }
  },
  legend: {
    position: 'top'
  }
});

// Update the chartSeries ref definition
const chartSeries = ref<{ name: string; data: ChartDataPoint[]; color?: string }[]>([
  {
    name: 'Sent',
    data: [],
    color: '#05062c'
  },
  {
    name: 'Requested',
    data: [],
    color: '#F2C037'
  }
]);

// Get the selected client from store (for chart filtering)
const selectedClient = computed(() =>
  store.clients.find(client => client.id === store.selectedClientId)
);

// Period selector options
const periodOptions = [
  { label: 'Last 7 days', value: { days: 7 } },
  { label: 'Last 14 days', value: { days: 14 } },
  { label: 'Last 30 days', value: { days: 30 } },
  { label: 'Last 3 months', value: { months: 3 } },
  { label: 'Last 6 months', value: { months: 6 } },
];

const selectedPeriod = ref(periodOptions[2]);

// Calculate date range based on selected period
const getDateRange = (period: { days?: number; months?: number }) => {
  const endDate = new Date();
  const startDate = new Date();

  if (period.days) {
    startDate.setDate(startDate.getDate() - period.days);
  } else if (period.months) {
    startDate.setMonth(startDate.getMonth() - period.months);
  }

  // Format dates as YYYY-MM-DD
  return {
    startDate: startDate.toISOString().split('T')[0],
    endDate: endDate.toISOString().split('T')[0],
  };
};

// Load stats
const loadStats = async () => {
  if (!selectedClient.value) return;

  const { startDate, endDate } = getDateRange(selectedPeriod.value!.value);
  const timeGrouping = selectedPeriod.value!.value.months ? 'Month' : 'Day';
  
  // Check cache first
  const cacheKey = apiCache.createStatsKey(selectedClient.value.id, startDate!, endDate!, timeGrouping);
  const cachedResponse = apiCache.get<UserStatsResponse>(cacheKey);
  
  if (cachedResponse) {
    parseStatsResponse(cachedResponse);
    return;
  }

  isLoading.value = true;
  hasStatsError.value = false;
  try {
    const response = await apiService.getUserStats(startDate!, endDate!, timeGrouping);
    
    // Cache the response
    apiCache.set(cacheKey, response);
    
    parseStatsResponse(response);
  } catch (error) {
    console.error('Failed to load stats:', error);
    hasStatsError.value = true;
  } finally {
    isLoading.value = false;
  }
};


const parseStatsResponse = (response: UserStatsResponse) => {
  // generate an array of days or months based on the selected period and number of periods
  const period = selectedPeriod.value!.value.months ? 'Month' : 'Day';
  const periodCount = selectedPeriod.value!.value.months ?? selectedPeriod.value!.value.days;

  // generate an array of DateTimes based on the selected period and number of periods
  const periodArray = Array.from({ length: periodCount }, (_, i) => {
    const date = period === 'Month'
      ? new Date(new Date().getFullYear(), new Date().getMonth() - i, 1)
      : new Date(new Date().getFullYear(), new Date().getMonth(), new Date().getDate() - i);

    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}T00:00:00Z`;
  });

  const filteredStats = (response.stats ?? []).filter(stat => stat.client_id === selectedClient.value?.id);

  const mappedStats = periodArray.map(period => {
    const stat = filteredStats.find(stat => stat.group_period === period);
    return {
      date: period,
      sent: stat?.sent ?? 0,
      requested: stat?.requested ?? 0,
    };
  });

  chartSeries.value = [
    {
      name: 'Sent',
      data: mappedStats.map(stat => ({
        x: new Date(stat.date).getTime(),
        y: stat.sent
      })).reverse(),
      color: '#05062c'
    },
    {
      name: 'Requested',
      data: mappedStats.map(stat => ({
        x: new Date(stat.date).getTime(),
        y: stat.requested
      })).reverse(),
      color: '#F2C037'
    }
  ];
};

// Reviews chart configuration
const reviewChartOptions = ref({
  chart: {
    type: 'bar',
    toolbar: {
      show: false
    }
  },
  plotOptions: {
    bar: {
      horizontal: false,
      columnWidth: '60%',
      distributed: true
    },
  },
  dataLabels: {
    enabled: false
  },
  xaxis: {
    categories: ['1⭐', '2⭐', '3⭐', '4⭐', '5⭐'],
    labels: {
      style: {
        fontSize: '12px',
        colors: '#666'
      }
    }
  },
  yaxis: {
    title: {
      text: 'Count'
    },
    min: 0,
    max: 20,
    tickAmount: 4,
    labels: {
      formatter: function(val: number) {
        return val.toString();
      }
    }
  },
  grid: {
    show: true,
    strokeDashArray: 3,
    borderColor: '#e0e0e0'
  },
  colors: ['#05062c'],
  legend: {
    show: false
  }
});

// Update reviewChartSeries to use the specified color
const reviewChartSeries = ref([
  {
    name: 'Reviews',
    data: [0, 0, 0, 0, 0],
  }
]);

// Update insightsChartSeries to use the specified color (keeping for data)
const insightsChartSeries = ref([
  {
    name: 'Count',
    data: [0, 0],
    color: '#F2C037' // Set color for Insights
  }
]);

// Load reviews
const loadReviews = async () => {
  if (!selectedClient.value) return;

  const { startDate, endDate } = getDateRange(selectedPeriod.value!.value);
  const startTime = `${startDate}T00:00:00Z`;
  const endTime = `${endDate}T23:59:59Z`;
  
  // Check cache first
  const cacheKey = apiCache.createReviewsKey(selectedClient.value.id, startTime, endTime);
  const cachedResponse = apiCache.get(cacheKey);
  
  if (cachedResponse) {
    processReviewsData(cachedResponse);
    return;
  }

  isLoadingReviews.value = true;
  hasReviewsError.value = false;
  try {
    const response = await apiService.getReviews(startTime, endTime, selectedClient.value.id);
    
    // Cache the response
    apiCache.set(cacheKey, response);
    
    processReviewsData(response);
  } catch (error) {
    console.error('Failed to load reviews:', error);
    hasReviewsError.value = true;
  } finally {
    isLoadingReviews.value = false;
  }
};

// Extract reviews processing into separate function for reuse
const processReviewsData = (response: any) => {
  // Aggregate ratings across all locations (now already filtered by backend)
  const totalRatings = {
    one: 0,
    two: 0,
    three: 0,
    four: 0,
    five: 0
  };

  let totalWebsiteClicks = 0;
  let totalCallButtonClicks = 0;

  (response.locations ?? []).forEach((location: any) => {
    totalRatings.one += location.review_ratings.one;
    totalRatings.two += location.review_ratings.two;
    totalRatings.three += location.review_ratings.three;
    totalRatings.four += location.review_ratings.four;
    totalRatings.five += location.review_ratings.five;

    totalWebsiteClicks += location.insights.number_of_business_profile_website_clicked;
    totalCallButtonClicks += location.insights.number_of_business_profile_call_button_clicked;
  });

  // Update review chart
  reviewChartSeries.value = [{
    name: 'Reviews',
    data: [
      totalRatings.one,
      totalRatings.two,
      totalRatings.three,
      totalRatings.four,
      totalRatings.five
    ],
  }];

  // Update insights chart
  insightsChartSeries.value = [{
    name: 'Count',
    data: [totalWebsiteClicks, totalCallButtonClicks],
    color: '#F2C037'
  }];
};

// Watch for changes in selected client or period
watch([selectedClient, selectedPeriod], () => {
  loadStats();
  loadReviews();
}, { immediate: true });

</script>

<style scoped>
/* Ensure charts take full width */
.apexcharts-canvas {
  width: 100% !important;
}

/* Fade transition */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.5s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

/* Unified skeleton container styling */
.unified-skeleton-container {
  width: 100%;
}

/* Profile Interactions styling */
.interaction-metric-card {
  text-align: center;
  padding: 2rem 1rem;
  background: white;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  transition: box-shadow 0.2s ease;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.interaction-metric-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.interaction-number {
  font-size: 4rem;
  font-weight: bold;
  line-height: 1;
  color: #F2C037;
  margin-bottom: 0.5rem;
}

.interaction-label {
  font-size: 0.875rem;
  font-weight: 500;
  letter-spacing: 0.1em;
  margin-top: 0.5rem;
}
</style>
