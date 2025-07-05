<template>
  <div class="smart-loading-container">
    <q-spinner-dots 
      size="40px" 
      color="primary" 
      class="q-mb-md"
    />
    <div class="loading-message">
      {{ currentMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

interface Props {
  loadingType?: 'dashboard' | 'reports' | 'general'
}

const props = withDefaults(defineProps<Props>(), {
  loadingType: 'general'
});

const currentMessage = ref('');
const messageTimer = ref<NodeJS.Timeout | null>(null);

const messages = {
  dashboard: [
    { time: 0, text: "Loading your dashboard..." },
    { time: 3000, text: "Analyzing your Google Business data..." },
    { time: 8000, text: "Processing performance metrics..." },
    { time: 15000, text: "Almost ready - finalizing your reports..." }
  ],
  reports: [
    { time: 0, text: "Loading your reports..." },
    { time: 3000, text: "Gathering report data..." },
    { time: 8000, text: "Processing analytics..." },
    { time: 15000, text: "Almost ready - preparing final report..." }
  ],
  general: [
    { time: 0, text: "Loading..." },
    { time: 3000, text: "Processing data..." },
    { time: 8000, text: "Almost ready..." },
    { time: 15000, text: "Finalizing..." }
  ]
};

const startMessageProgression = () => {
  const messageSet = messages[props.loadingType];
  let messageIndex = 0;
  
  // Set initial message
  const initialMessage = messageSet[0];
  if (initialMessage) {
    currentMessage.value = initialMessage.text;
  }
  
  const scheduleNextMessage = () => {
    if (messageIndex < messageSet.length - 1) {
      const nextMessage = messageSet[messageIndex + 1];
      const currentMessageObj = messageSet[messageIndex];
      
      if (nextMessage && currentMessageObj) {
        const delay = nextMessage.time - currentMessageObj.time;
        
        messageTimer.value = setTimeout(() => {
          messageIndex++;
          const newMessage = messageSet[messageIndex];
          if (newMessage) {
            currentMessage.value = newMessage.text;
          }
          scheduleNextMessage();
        }, delay);
      }
    }
  };
  
  scheduleNextMessage();
};

const clearMessageTimer = () => {
  if (messageTimer.value) {
    clearTimeout(messageTimer.value);
    messageTimer.value = null;
  }
};

onMounted(() => {
  startMessageProgression();
});

onUnmounted(() => {
  clearMessageTimer();
});
</script>

<style scoped>
.smart-loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.loading-message {
  color: #666;
  font-size: 14px;
  text-align: center;
  min-height: 20px;
  font-weight: 500;
  transition: opacity 0.3s ease;
}
</style>