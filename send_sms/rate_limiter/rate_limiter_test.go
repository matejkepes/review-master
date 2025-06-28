package rate_limiter

import (
	"log"
	"testing"

	"send_sms/config"
)

func TestRateLimiterEnabled(t *testing.T) {
	// read config file
	config := config.ReadProperties()

	rateLimiterEnabled, rateLimiterRestrictor := RateLimiter(config)
	log.Printf("rateLimiterEnabled: %t, rateLimiterRestrictor: %+v\n", rateLimiterEnabled, rateLimiterRestrictor)
}

func TestRateLimiterDisabled(t *testing.T) {
	// read config file
	config := config.ReadProperties()
	config.RateLimiterEnabled = false

	rateLimiterEnabled, rateLimiterRestrictor := RateLimiter(config)
	log.Printf("rateLimiterEnabled: %t, rateLimiterRestrictor: %+v\n", rateLimiterEnabled, rateLimiterRestrictor)
}
