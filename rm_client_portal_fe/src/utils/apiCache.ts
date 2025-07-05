// Simple session-based cache for API responses
// MVP approach - no persistence, just memory cache to avoid duplicate calls

interface CacheEntry<T> {
  data: T;
  timestamp: number;
  ttl: number; // Time to live in milliseconds
}

class SimpleCache {
  private cache = new Map<string, CacheEntry<any>>();
  private defaultTTL = 5 * 60 * 1000; // 5 minutes default TTL

  set<T>(key: string, data: T, ttl?: number): void {
    this.cache.set(key, {
      data,
      timestamp: Date.now(),
      ttl: ttl || this.defaultTTL
    });
  }

  get<T>(key: string): T | null {
    const entry = this.cache.get(key);
    
    if (!entry) {
      return null;
    }

    // Check if entry has expired
    if (Date.now() - entry.timestamp > entry.ttl) {
      this.cache.delete(key);
      return null;
    }

    return entry.data;
  }

  has(key: string): boolean {
    const entry = this.cache.get(key);
    
    if (!entry) {
      return false;
    }

    // Check if entry has expired
    if (Date.now() - entry.timestamp > entry.ttl) {
      this.cache.delete(key);
      return false;
    }

    return true;
  }

  clear(): void {
    this.cache.clear();
  }

  // Generate cache keys for common API patterns
  createStatsKey(clientId: number, startDate: string, endDate: string, timeGrouping: string): string {
    return `stats:${clientId}:${startDate}:${endDate}:${timeGrouping}`;
  }

  createReviewsKey(clientId: number, startTime: string, endTime: string): string {
    return `reviews:${clientId}:${startTime}:${endTime}`;
  }
}

// Export singleton instance for MVP simplicity
export const apiCache = new SimpleCache();