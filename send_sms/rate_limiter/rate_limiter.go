package rate_limiter

import (
	"log"

	"send_sms/config"

	"github.com/EagleChen/restrictor"
)

// RateLimiter - rate limiter to be used to restrict the number of requests per telephone number.
// Returns a boolean as to whether to use rate limiting and the restrictor for use in further prcessing.
// Pass this to the send_sms_handler and do check:
// if r.LimitReached(tel) {
// limit is reached, notify user
// } else {
// continue with remaining operations
// }
//
func RateLimiter(config config.Config) (bool, restrictor.Restrictor) {
	rateLimitEnabled := config.RateLimiterEnabled

	var r restrictor.Restrictor
	if rateLimitEnabled {
		// create a 'store' to store various limiters for each key
		store, err := restrictor.NewMemoryStore()
		// if there are many backend servers using one central limiter store use 'NewRedisStore' instead (TODO: may later)
		if err != nil {
			log.Println("Error creating rate limiter store")
			rateLimitEnabled = false
		} else {
			// create restrictor
			// e.g. 5 request every 5 minutes per telephone number
			// usually pick some number from 60 to 100 for number of buckets. This number will affect the deviation
			r = restrictor.NewRestrictor(config.RateLimiterWindowMinutes, uint32(config.RateLimiterUpperLimit), uint32(config.RateLimiterBucketSpan), store)
		}
	}

	return rateLimitEnabled, r
}
