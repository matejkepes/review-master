<template>
  <div 
    class="smart-loading-container" 
    :class="{ 'section-header': sectionLevel }"
    role="status"
    aria-live="polite"
    :aria-label="`Loading: ${currentMessage}`"
  >
    <q-spinner-dots 
      size="40px" 
      color="primary" 
      class="loading-spinner"
      :class="{ 'section-spinner': sectionLevel }"
      aria-hidden="true"
    />
    <div 
      class="loading-message" 
      :class="{ 'section-message': sectionLevel }"
    >
      {{ currentMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { loadingMessages, isValidLoadingType } from 'src/utils/loadingMessages';

import type { LoadingType } from 'src/utils/loadingMessages';

interface Props {
  loadingType?: LoadingType
  sectionLevel?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loadingType: 'general',
  sectionLevel: false
});

const currentMessage = ref('');
const messageTimer = ref<ReturnType<typeof setTimeout> | null>(null);

const startMessageProgression = () => {
  const loadingType = isValidLoadingType(props.loadingType) ? props.loadingType : 'general';
  const messageSet = loadingMessages[loadingType];
  
  if (!messageSet || messageSet.length === 0) {
    currentMessage.value = 'Loading...';
    return;
  }
  
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

/* Section-level loading styling */
.section-header {
  flex-direction: row;
  align-items: center;
  justify-content: flex-start;
  padding: 1rem 0;
  gap: 1rem;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 1.5rem;
  background: linear-gradient(135deg, #fafafa 0%, #f5f5f5 100%);
  border-radius: 8px;
  padding: 1.5rem;
}

.loading-spinner {
  margin-bottom: 1rem;
}

.section-spinner {
  margin-bottom: 0;
  animation: section-pulse 2s ease-in-out infinite;
}

.loading-message {
  color: #666;
  font-size: 14px;
  text-align: center;
  min-height: 20px;
  font-weight: 500;
  transition: opacity 0.3s ease;
}

.section-message {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  text-align: left;
  animation: section-text-pulse 2s ease-in-out infinite;
}

/* Enhanced animations for section-level loading */
@keyframes section-pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.05);
    opacity: 0.8;
  }
}

@keyframes section-text-pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}


</style>