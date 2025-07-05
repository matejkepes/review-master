/**
 * Centralized loading messages for SmartLoadingSpinner component
 * Uses contextual messages throughout timer progression
 */

export interface LoadingMessage {
  time: number;
  text: string;
}

export type LoadingType = 'dashboard' | 'reports' | 'reviews' | 'general' | 'clientData' | 'stats' | 'report';

export type LoadingMessageSet = {
  [K in LoadingType]: LoadingMessage[];
}

// Progressive loading messages based on elapsed time (0-5s, 5-15s, 15-30s, 30+s)
export const loadingMessages: LoadingMessageSet = {
  dashboard: [
    { time: 0, text: "Loading your dashboard..." },
    { time: 5000, text: "Gathering your latest data..." },
    { time: 15000, text: "Processing your analytics..." },
    { time: 30000, text: "Almost there - finalizing your dashboard..." }
  ],
  
  reports: [
    { time: 0, text: "Loading your reports..." },
    { time: 5000, text: "Compiling your monthly insights..." },
    { time: 15000, text: "Analyzing performance trends..." },
    { time: 30000, text: "Generating detailed analytics..." }
  ],
  
  reviews: [
    { time: 0, text: "Gathering your latest customer reviews..." },
    { time: 5000, text: "Analyzing review sentiment..." },
    { time: 15000, text: "Processing customer feedback..." },
    { time: 30000, text: "Compiling review insights..." }
  ],
  
  general: [
    { time: 0, text: "Loading..." },
    { time: 5000, text: "Processing data..." },
    { time: 15000, text: "Almost ready..." },
    { time: 30000, text: "Finalizing..." }
  ],

  clientData: [
    { time: 0, text: "Connecting to your Google Business Profile..." },
    { time: 5000, text: "Authenticating with Google Services..." },
    { time: 15000, text: "Syncing your business data..." },
    { time: 30000, text: "Finalizing profile synchronization..." }
  ],

  stats: [
    { time: 0, text: "Analyzing your performance metrics..." },
    { time: 5000, text: "Calculating key performance indicators..." },
    { time: 15000, text: "Processing historical data..." },
    { time: 30000, text: "Generating statistical insights..." }
  ],

  report: [
    { time: 0, text: "Compiling your monthly insights report..." },
    { time: 5000, text: "Analyzing performance trends..." },
    { time: 15000, text: "Generating detailed analytics..." },
    { time: 30000, text: "Finalizing your comprehensive report..." }
  ]
};

// Validation helper
export const isValidLoadingType = (type: string): type is LoadingType => {
  const validTypes: LoadingType[] = ['dashboard', 'reports', 'reviews', 'general', 'clientData', 'stats', 'report'];
  return validTypes.includes(type as LoadingType);
};