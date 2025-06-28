package server

import (
	"fmt"
	"log"
	"rm_client_portal/database"
	"rm_client_portal/google_my_business_api"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

// ClientsHandler - retrieve a list of clients
func ClientsHandler(c *gin.Context) {
	clients := database.GetClientsForUserEmail(getEmailFromJWT(c))
	c.JSON(200, gin.H{
		"clients": clients,
	})
}

// StatsUserHandler - retrieve a list of stats for a user using stats database table for user
func StatsUserHandler(c *gin.Context) {
	success := true
	var errStr string
	var stats []database.StatsNewResult
	var statsRequest database.StatsRequest
	startDay := c.Query("start_day")
	endDay := c.Query("end_day")
	timeGrouping := c.Query("time_grouping")
	// log.Printf("start_day: %s\n", startDay)
	statsRequest.StartDay = startDay
	statsRequest.EndDay = endDay
	statsRequest.TimeGrouping = timeGrouping
	// log.Printf("statsRequest: %+v\n", statsRequest)
	email := getEmailFromJWT(c)
	clients := database.GetClientsForUserEmail(email)
	for _, v := range clients {
		s, err := database.ClientStats(statsRequest, int(v.ID))
		if err != nil {
			log.Printf("error retrieving stats list, err: %+v\n", err)
			errStr = fmt.Sprintf("error retrieving stats list, error: %+v", err)
			success = false
		} else {
			stats = append(stats, s...)
			success = true
		}
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
		"stats":   stats,
	})
}

// ReportOnReviewsAndInsights - retrieve reviews and insights from Google
func ReportOnReviewsAndInsights(c *gin.Context) {
	email := getEmailFromJWT(c)
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	accounts := google_my_business_api.GetAccounts()
	clientIDs := database.GetClientIDsForUserEmail(email)
	// log.Printf("clientIDs = %v\n", clientIDs)
	var locs []database.GoogleReviewsConfigAndGoogleMyBusinessLocation
	for _, a := range accounts {
		l := google_my_business_api.GetLocationsCheckClientID(a, clientIDs)
		for _, g := range l {
			g.GoogleReviewRatings = google_my_business_api.ReportOnReviewsWeb(g, startTime, endTime)
			// insights
			g.GoogleInsights = google_my_business_api.ReportOnInsightsWeb(g, startTime, endTime)
			locs = append(locs, g)
		}
	}
	// log.Printf("locs: %+v\n", locs)

	c.JSON(200, gin.H{
		"locations": locs,
	})
}

// Get the email (which is the identity key) from the JWT claims
func getEmailFromJWT(c *gin.Context) string {
	// log.Printf("claims: %s\n", jwt.ExtractClaims(c))
	claims := jwt.ExtractClaims(c)
	email := claims[identityKey].(string)
	return email
}
