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
        <div class="text-h6 q-mb-md">Reviews</div>

        <div v-if="isLoadingReviews">
          <SmartLoadingSpinner loadingType="dashboard" />
        </div>
        <div v-else-if="hasReviewsError" class="text-center q-pa-lg">
          <q-icon name="error_outline" size="48px" color="negative" class="q-mb-md" />
          <div class="text-subtitle1 text-grey-7">Failed to load reviews</div>
          <q-btn color="primary" label="Retry" @click="loadReviews" class="q-mt-md" />
        </div>
        <transition name="fade" appear>
          <div v-if="!isLoadingReviews && !hasReviewsError" class="row q-col-gutter-md">
            <!-- Ratings Distribution -->
            <div class="col-12">
              <apexchart type="bar" height="350" :options="reviewChartOptions" :series="reviewChartSeries" />
            </div>

          <!-- Insights -->
          <div class="col-12 q-mt-md">
            <apexchart type="bar" height="350" :options="insightsChartOptions" :series="insightsChartSeries" />
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
    color: '#11189E'
  },
  {
    name: 'Requested',
    data: [],
    color: '#F6CC53'
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
      color: '#11189E'
    },
    {
      name: 'Requested',
      data: mappedStats.map(stat => ({
        x: new Date(stat.date).getTime(),
        y: stat.requested
      })).reverse(),
      color: '#F6CC53'
    }
  ];
};

// Reviews chart configuration
const reviewChartOptions = ref({
  chart: {
    type: 'bar',
  },
  plotOptions: {
    bar: {
      horizontal: false,
      columnWidth: '55%',
      distributed: true
    },
  },
  dataLabels: {
    enabled: false
  },
  xaxis: {
    categories: ['⭐', '⭐⭐', '⭐⭐⭐', '⭐⭐⭐⭐', '⭐⭐⭐⭐⭐'],
    labels: {
      show: false
    }
  },
  title: {
    text: 'Rating Distribution',
    align: 'center'
  },
  colors: ['#05062C', '#05062C', '#05062C', '#05062C', '#05062C']
});

// Update reviewChartSeries to use the specified color
const reviewChartSeries = ref([
  {
    name: 'Reviews',
    data: [0, 0, 0, 0, 0],
  }
]);

// Insights chart configuration
const insightsChartOptions = ref({
  chart: {
    type: 'bar',
  },
  plotOptions: {
    bar: {
      horizontal: true,
    },
  },
  dataLabels: {
    enabled: true
  },
  xaxis: {
    categories: ['Website Clicks', 'Call Button Clicks'],
  },
  title: {
    text: 'Profile Interactions',
    align: 'center'
  }
});

// Update insightsChartSeries to use the specified color
const insightsChartSeries = ref([
  {
    name: 'Count',
    data: [0, 0],
    color: '#EC9714' // Set color for Insights
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
    color: '#EC9714'
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
</style>
