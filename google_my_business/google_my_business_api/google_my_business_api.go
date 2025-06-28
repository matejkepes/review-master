package google_my_business_api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"google_my_business/ai_service"
	"google_my_business/ai_service/review"
	"google_my_business/config"
	"google_my_business/database"
)

const googleMyBusinessAPIURL = "https://mybusiness.googleapis.com/v4/"
const googleMyBusinessAccountManagementAPIURL = "https://mybusinessaccountmanagement.googleapis.com/v1/"
const googleMyBusinessBusinessInformationAPIURL = "https://mybusinessbusinessinformation.googleapis.com/v1/"
const googleMyBusinessBusinessProfilePerformanceAPIURL = "https://businessprofileperformance.googleapis.com/v1/"
const StarRatingUnspecified = "STAR_RATING_UNSPECIFIED"
const StarRatingOne = "ONE"
const StarRatingTwo = "TWO"
const StarRatingThree = "THREE"
const StarRatingFour = "FOUR"
const StarRatingFive = "FIVE"

// GetAccounts - get accounts
func GetAccounts(client *http.Client) []string {
	var accounts []string
	apiURL := googleMyBusinessAccountManagementAPIURL + "accounts"

	// Make the request
	resp, err := client.Get(apiURL)
	if err != nil {
		log.Printf("Error calling Google API: %v", err)
		return accounts
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-200 response status: %d", resp.StatusCode)
		return accounts
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return accounts
	}

	// Parse JSON response
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return accounts
	}

	// Extract accounts
	if accountsArray, ok := data["accounts"].([]interface{}); ok {
		for _, record := range accountsArray {
			if rec, ok := record.(map[string]interface{}); ok {
				if accountName, ok := rec["name"].(string); ok {
					accounts = append(accounts, accountName)
				}
			}
		}
	}

	return accounts
}

// GetAccountsFromJSON - get accounts from JSON
func GetAccountsFromJSON(in []byte) []string {
	f := map[string]interface{}{
		"key": "value",
	}
	err := json.Unmarshal([]byte(in), &f)
	if err != nil {
		log.Fatal(err)
	}
	var accounts []string
	for key, record := range f {
		// fmt.Println(key)
		switch key {
		case "accounts":
			// fmt.Println("interface type", reflect.TypeOf(record))
			// test whether a single account, which will be of type map[string]interface{}, if so create a slice and add
			b := record
			// fmt.Printf("b = %+v\n", b)
			if _, ok := b.(map[string]interface{}); ok {
				s := make([]interface{}, 1)
				b = append(s, record)
			}

			if r, ok := b.([]interface{}); ok {
				for _, a := range r {
					if rec, ok := a.(map[string]interface{}); ok {
						// fmt.Printf("rec = %+v\n", rec)
						// the name holds the path that is appended to the google api URL
						name := fmt.Sprintf("%s", rec["name"])
						accounts = append(accounts, name)
					}
				}
			}
			// default:
			// fmt.Println(key, ":", record)
		}
	}
	// fmt.Printf("accounts = %+v", accounts)
	return accounts
}

// GetLocations - get locations for account
// lookupMode parameter indicates which mode to use (reply, report, or analysis)
func GetLocations(client *http.Client, account string, db *sql.DB, lookupMode int) []database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode {
	var pageToken string
	var grcfgmblns []database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode

	for {
		grcfgmblnsp, nextPageToken := getLocationsPage(client, account, db, lookupMode, pageToken)
		grcfgmblns = append(grcfgmblns, grcfgmblnsp...)
		if nextPageToken == "" {
			break
		} else {
			pageToken = nextPageToken
		}
	}
	return grcfgmblns
}

// getLocationsPage - get locations for account page at a time (default page size is 100)
// lookupMode parameter indicates which mode to use (reply, report, or analysis)
// pageToken parameter indicates when the locations has already been called and there are more results
func getLocationsPage(client *http.Client, account string, db *sql.DB, lookupMode int, pageToken string) ([]database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, string) {
	locationsParams := ""
	if pageToken != "" {
		locationsParams += "?pageToken=" + pageToken
	}
	// resp, err := client.Get(googleMyBusinessAPIURL + account + "/locations" + locationsParams)
	// readMask required query parameter as of new API v1
	readMaskParams := "readMask=name,title,storefrontAddress"
	if locationsParams == "" {
		readMaskParams = "?" + readMaskParams
	} else {
		readMaskParams = "&" + readMaskParams
	}
	// change as of new API v1
	resp, err := client.Get(googleMyBusinessBusinessInformationAPIURL + account + "/locations" + locationsParams + readMaskParams)
	if err != nil {
		log.Fatalf("Unable to retrieve locations error: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading body err: %v", err)
	}
	// change as of new API v1 location name does not include the accounts/{accountId}/ part of the path
	grcfgmblns, nextPageToken := GetLocationsFromJSON(body, account, db, lookupMode)
	return grcfgmblns, nextPageToken
}

// GetLocationsFromJSON - get locations from JSON
// change as of new API v1 location name does not include the accounts/{accountId}/ part of the path
func GetLocationsFromJSON(in []byte, account string, db *sql.DB, lookupMode int) ([]database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, string) {
	f := map[string]interface{}{
		"key": "value",
	}
	err := json.Unmarshal([]byte(in), &f)
	if err != nil {
		log.Fatal(err)
	}
	var grcfgmblns []database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode
	var nextPageToken string
	for key, record := range f {
		// fmt.Println(key)
		switch key {
		case "locations":
			// fmt.Println("interface type", reflect.TypeOf(record))
			// test whether a single account, which will be of type map[string]interface{}, if so create a slice and add
			b := record
			// fmt.Printf("b = %+v\n", b)
			if _, ok := b.(map[string]interface{}); ok {
				s := make([]interface{}, 1)
				b = append(s, record)
			}

			if r, ok := b.([]interface{}); ok {
				for _, a := range r {
					if rec, ok := a.(map[string]interface{}); ok {
						// fmt.Printf("rec = %+v\n", rec)
						// check the location name against database to get configuration fields for review reply
						// locationName := fmt.Sprintf("%s", rec["locationName"])
						// change as of new API v1
						locationName := fmt.Sprintf("%s", rec["title"])
						// fmt.Printf("locationName = %s\n", locationName)

						// Debug raw title field and handle null/empty case
						log.Printf("DEBUG Location: Raw title field: %v, Type: %T", rec["title"], rec["title"])
						if rec["title"] == nil || locationName == "" || locationName == "<nil>" || locationName == "null" {
							log.Printf("WARNING: Location without title found - Using a fallback name")
							// Try to use the name field for the location as fallback
							if locName, ok := rec["name"].(string); ok && locName != "" {
								// Extract last part of the path
								parts := strings.Split(locName, "/")
								if len(parts) > 0 {
									locationName = fmt.Sprintf("Location %s", parts[len(parts)-1])
								} else {
									locationName = "Unknown Location"
								}
							} else {
								locationName = "Unknown Location"
							}
						}

						// address := rec["address"]
						address := rec["storefrontAddress"]
						// fmt.Printf("reviewer = %+v\n", reviewer)
						postalCode := ""
						if v, ok := address.(map[string]interface{}); ok {
							postalCode = fmt.Sprintf("%s", v["postalCode"])
							// fmt.Printf("postalCode = %s\n", postalCode)
						}
						// require a postal code
						if len(postalCode) == 0 {
							continue
						}

						// Load the location configuration from the database
						grcfgmbln := database.ConfigFromGoogleMyBusinessLocationNameAndPostalCode(db, locationName, postalCode, lookupMode)
						if grcfgmbln.ClientID == 0 {
							log.Printf("WARNING: No configuration found for location '%s' with postal code '%s'", locationName, postalCode)
							continue
						}

						// Set the location path
						grcfgmbln.GoogleMyBusinessLocationPath = fmt.Sprintf("%s/%s", account, rec["name"])

						log.Printf("DEBUG Location: Using database configuration - Name: '%s', Path: '%s', Postal: '%s', ClientID: %d",
							grcfgmbln.GoogleMyBusinessLocationName,
							grcfgmbln.GoogleMyBusinessLocationPath,
							grcfgmbln.GoogleMyBusinessPostalCode,
							grcfgmbln.ClientID)

						// Removed because some customers like to respond themselves but we always want the report for these
						// if !(grcfgmbln.GoogleMyBusinessReplyToUnspecfifiedStarRating ||
						// 	grcfgmbln.GoogleMyBusinessReplyToOneStarRating ||
						// 	grcfgmbln.GoogleMyBusinessReplyToTwoStarRating ||
						// 	grcfgmbln.GoogleMyBusinessReplyToThreeStarRating ||
						// 	grcfgmbln.GoogleMyBusinessReplyToFourStarRating ||
						// 	grcfgmbln.GoogleMyBusinessReplyToFiveStarRating) {
						// 	continue
						// }
						// the name holds the path to be appended to the google api URL
						// name := fmt.Sprintf("%s", rec["name"])
						// change as of new API v1 location name does not include the accounts/{accountId}/ part of the path
						name := fmt.Sprintf("%s/%s", account, rec["name"])
						grcfgmbln.GoogleMyBusinessLocationPath = name
						grcfgmblns = append(grcfgmblns, grcfgmbln)
					}
				}
			}
		case "nextPageToken":
			// If the number of locations exceeded the requested page size, this field is populated
			// with a token to fetch the next page of locations on a subsequent call to locations.
			// If there are no more locations, this field is not present in the response.
			nextPageToken = record.(string)
			// default:
			// fmt.Println(key, ":", record)
		}
	}
	// fmt.Printf("grcfgmblns = %+v", grcfgmblns)
	return grcfgmblns, nextPageToken
}

// // ProcessReviews - process reviews for location
// func ProcessReviews(client *http.Client, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, reviewsNotBeforeDays int, replySeparator string) {
// 	resp, err := client.Get(googleMyBusinessAPIURL + grcfgmbln.GoogleMyBusinessLocationPath + "/reviews")
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve reviews for %s (clientID: %d) error: %v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
// 	}
// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalf("Error reading body for %s (clientID: %d) err: %v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
// 	}
// 	// fmt.Println(string(body))
// 	processReviewsFromJSON(client, body, grcfgmbln, reviewsNotBeforeDays, replySeparator)
// }

// ProcessReviews - process reviews for location
func ProcessReviews(client *http.Client, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, reviewsNotBeforeDays int, replySeparator string) {
	var pageToken string
	for {
		nextPageToken := processReviewsPage(client, grcfgmbln, reviewsNotBeforeDays, replySeparator, pageToken)
		if nextPageToken == "" {
			break
		} else {
			pageToken = nextPageToken
		}
	}
}

// processReviewsPage - process reviews for location per page
func processReviewsPage(client *http.Client, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, reviewsNotBeforeDays int, replySeparator string, pageToken string) string {
	url := googleMyBusinessAPIURL + grcfgmbln.GoogleMyBusinessLocationPath + "/reviews"
	if len(strings.Trim(pageToken, " ")) > 0 {
		url += "?pageToken=" + strings.Trim(pageToken, " ")
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Unable to retrieve reviews for %s (clientID: %d) error: %v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading body for %s (clientID: %d) err: %v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
	}
	// fmt.Println(string(body))
	nextPageToken, _ := processReviewsFromJSON(client, body, grcfgmbln, reviewsNotBeforeDays, replySeparator, true)
	// fmt.Printf("nextPageToken = %s\n", nextPageToken)
	return nextPageToken
}

// processReviewsFromJSON - process reviews from JSON
func processReviewsFromJSON(client *http.Client, in []byte, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, reviewsNotBeforeDays int, replySeparator string, processReplies bool) (string, []map[string]interface{}) {
	// nextPageToken default to empty
	var nextPageToken string
	// check complete indicates that all the reviews have been checked for period and do not need to go to next page
	var checkComplete bool
	var reviews []map[string]interface{}

	// Get configuration from properties file
	props := config.ReadProperties()

	// Log database connection attempt
	log.Printf("Attempting to open database connection for %s. AI responses enabled: %v",
		grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.AIResponsesEnabled)

	db := database.OpenDB(props.DbName, props.DbAddress, props.DbPort, props.DbUsername, props.DbPassword)
	defer db.Close()

	log.Printf("Successfully connected to database. Processing reviews for %s. AI responses enabled: %v",
		grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.AIResponsesEnabled)

	f := map[string]interface{}{
		"key": "value",
	}
	err := json.Unmarshal([]byte(in), &f)
	if err != nil {
		log.Fatalf("Error unmarshalling jSON for %s (clientID: %d) err: %v \n", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
	}
	// Do not process reviews before check time
	checkTimeNotBefore := time.Now().AddDate(0, 0, -reviewsNotBeforeDays)
	// fmt.Printf("checkTimeNotBefore = %+v\n", checkTimeNotBefore)
	for key, record := range f {
		// fmt.Println(key)
		switch key {
		case "reviews":
			// fmt.Println("interface type", reflect.TypeOf(record))
			// test whether a single account, which will be of type map[string]interface{}, if so create a slice and add
			b := record
			// fmt.Printf("b = %+v\n", b)
			if _, ok := b.(map[string]interface{}); ok {
				s := make([]interface{}, 1)
				b = append(s, record)
			}

			// debug log, processing reviews
			log.Printf("Processing reviews for %s (clientID: %d)", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID)

			if r, ok := b.([]interface{}); ok {
				for _, a := range r {
					if rec, ok := a.(map[string]interface{}); ok {
						// check update time if exists, else check create time, after configured time
						checkTime := fmt.Sprintf("%s", rec["createTime"])
						// fmt.Printf("createTime = %s\n", createTime)
						if _, ok := rec["updateTime"]; ok {
							checkTime = fmt.Sprintf("%s", rec["updateTime"])
						}
						// tm, err := time.Parse("2006-01-02T15:04:05.000Z", createTime)
						tm, err := time.Parse(time.RFC3339, checkTime)
						if err != nil {
							log.Printf("Error converting checkTime %s string to a time for %s (clientID: %d)\n", checkTime, grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID)
							continue
						}
						// fmt.Printf("tm = %+v\n", tm)
						if tm.Before(checkTimeNotBefore) {
							// The reviews are in create time descending order, so safe to stop processing
							// set the check complete so do not fetch next page of reviews
							checkComplete = true
							break
						}

						// Add review to the list if we're just fetching
						if !processReplies {
							reviews = append(reviews, rec)
							continue
						}

						// if review reply exists then already responded to so skip
						reviewReply := rec["reviewReply"]
						if _, ok := reviewReply.(map[string]interface{}); ok {
							continue
						}
						// star rating can be anyone of these:
						// STAR_RATING_UNSPECIFIED - Not specified.
						// ONE - One star out of a maximum of five.
						// TWO - Two stars out of a maximum of five.
						// THREE - Three stars out of a maximum of five.
						// FOUR - Four stars out of a maximum of five.
						// FIVE - The maximum star rating.
						starRating := fmt.Sprintf("%s", rec["starRating"])
						// fmt.Printf("starRating = %s\n", starRating)
						// Initialize reply variable
						var reply string

						reviewer := rec["reviewer"]
						var displayName string
						if v, ok := reviewer.(map[string]interface{}); ok {
							displayName = extractReviewerName(v)
						}

						// Get review text for AI response
						reviewText := extractReviewText(rec)

						// Add logging in the if condition for AI responses
						if grcfgmbln.AIResponsesEnabled && len(reviewText) > 0 {
							log.Printf("AI responses are enabled and review text is not empty. Attempting to generate AI response...")
							log.Printf("Review details - Rating: %s, Author: %s, Text length: %d",
								starRating, displayName, len(reviewText))

							// Try to generate AI response
							aiResponse, err := getAIResponse(reviewText, starRating, displayName, grcfgmbln, db)
							if err == nil {
								reply = aiResponse
								log.Printf("Successfully generated AI response: %s", truncateString(reply, 50))
							} else {
								// Log the error and fall back to template response
								log.Printf("AI response generation failed with error: %v", err)
								log.Printf("Falling back to template response")
								templateReply, ok := getTemplateResponse(starRating, displayName, grcfgmbln, replySeparator)
								if ok {
									reply = templateReply
									log.Printf("Using template response as fallback: %s", truncateString(reply, 50))
								}
							}
						} else {
							// Log why AI response is not being used
							if !grcfgmbln.AIResponsesEnabled {
								log.Printf("AI responses are disabled for this client (AIResponsesEnabled=false). Using template response.")
							}
							if len(reviewText) == 0 {
								log.Printf("Review text is empty. Using template response.")
							}

							templateReply, ok := getTemplateResponse(starRating, displayName, grcfgmbln, replySeparator)
							if ok {
								reply = templateReply
								log.Printf("Using template response: %s", truncateString(reply, 50))
							}
						}

						// debug log, reply
						log.Printf("Reply: %s", reply)

						// If no reply was generated, skip this review
						if len(reply) == 0 {
							continue
						}

						// the name holds the path to be appended to the google api URL
						name := fmt.Sprintf("%s", rec["name"])
						// fmt.Printf("name: %s\n", name)
						// updateReply method (https://developers.google.com/my-business/reference/rest/v4/accounts.locations.reviews/updateReply)
						// updateReviewReply := fmt.Sprintf("{\"comment\":\"%s\",\"updateTime\":\"%s\"}", reply, time.Now().Format("2006-01-02T15:04:05.000000000Z"))
						updateReviewReply := fmt.Sprintf("{\"comment\":\"%s\",\"updateTime\":\"%s\"}", reply, time.Now().Format(time.RFC3339))
						log.Printf("update review for %s (clientID: %d) with star rating %s with reply: %s for review path: %s\n", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, starRating, updateReviewReply, name)
						req, err := http.NewRequest(http.MethodPut, googleMyBusinessAPIURL+name+"/reply", strings.NewReader(updateReviewReply))
						// log.Printf("req: %+v", req)
						if err != nil {
							log.Printf("Unable creating update review request for %s (clientID: %d) error: %v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
							continue
						}
						_, err = client.Do(req)
						if err != nil {
							log.Printf("Unable to update review for %s (clientID: %d) error: %v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
							continue
						}
					}
				}
			}
		case "nextPageToken":
			// If the number of reviews exceeded the requested page size, this field is populated
			// with a token to fetch the next page of reviews on a subsequent call to reviews.
			// If there are no more locations, this field is not present in the response.
			nextPageToken = record.(string)
			// default:
			// fmt.Println(key, ":", record)
		}
	}
	// do not need to fetch next page if check is complete (i.e. got all reviews for the specified period)
	if checkComplete {
		nextPageToken = ""
	}
	return nextPageToken, reviews
}

// TODO: NOTE see: https://developers.google.com/my-business/content/review-data and the Get reviews from multiple locations section may be more efficient. Go straight to: https://developers.google.com/my-business/content/review-data#get_reviews_from_multiple_locations

// // ReportOnReviews - report on reviews for location
// func ReportOnReviews(client *http.Client, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, reportMonthBack int) (map[string]int, int, time.Month) {
// 	resp, err := client.Get(googleMyBusinessAPIURL + grcfgmbln.GoogleMyBusinessLocationPath + "/reviews")
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve reviews for %s (clientID: %d) for report error: %v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
// 	}
// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalf("Error reading body for %s (clientID: %d) for report err: %v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
// 	}
// 	// fmt.Println(string(body))
// 	firstOfMonth, firstOfNextMonth := reportPeriod(reportMonthBack, grcfgmbln.TimeZone)
// 	// fmt.Printf("firstOfMonth: +%v, firstOfNextMonth: +%v\n", firstOfMonth, firstOfNextMonth)
// 	reviews := reportOnReviewsFromJSON(body, grcfgmbln, firstOfMonth, firstOfNextMonth)
// 	reportFirstOfYear, reportFirstOfMonth, _ := firstOfMonth.Date()
// 	return reviews, reportFirstOfYear, reportFirstOfMonth
// }

// ReportOnReviews - report on reviews for location
func ReportOnReviews(client *http.Client, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, reportMonthBack int) (map[string]int, int, time.Month) {
	var pageToken string
	reviews := make(map[string]int)
	firstOfMonth, firstOfNextMonth := reportPeriod(reportMonthBack, grcfgmbln.TimeZone)
	// fmt.Printf("firstOfMonth: +%v, firstOfNextMonth: +%v\n", firstOfMonth, firstOfNextMonth)
	reportFirstOfYear, reportFirstOfMonth, _ := firstOfMonth.Date()
	for {
		pageReviews, nextPageToken := reportOnReviewsPage(client, grcfgmbln, pageToken, firstOfMonth, firstOfNextMonth)
		// append reviews page
		for k, v := range pageReviews {
			reviews[k] = reviews[k] + v
		}
		if nextPageToken == "" {
			break
		} else {
			pageToken = nextPageToken
		}
	}
	return reviews, reportFirstOfYear, reportFirstOfMonth
}

// reportOnReviewsPage - report on reviews for location per page
func reportOnReviewsPage(client *http.Client, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, pageToken string, startTime time.Time, endTime time.Time) (map[string]int, string) {
	url := googleMyBusinessAPIURL + grcfgmbln.GoogleMyBusinessLocationPath + "/reviews"
	if len(strings.Trim(pageToken, " ")) > 0 {
		url += "?pageToken=" + strings.Trim(pageToken, " ")
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Unable to retrieve reviews for %s (clientID: %d) for report error: %v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading body for %s (clientID: %d) for report err: %v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
	}
	// fmt.Println(string(body))
	reviews, nextPageToken := ReportOnReviewsFromJSON(body, grcfgmbln, startTime, endTime)
	// fmt.Printf("nextPageToken = %s\n", nextPageToken)
	return reviews, nextPageToken
}

// // reportOnReviewsFromJSON - report on reviews from JSON
// func reportOnReviewsFromJSON(in []byte, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, startTime time.Time, endTime time.Time) map[string]int {
// 	// count review star ratings
// 	review_ratings := make(map[string]int)

// 	f := map[string]interface{}{
// 		"key": "value",
// 	}
// 	err := json.Unmarshal([]byte(in), &f)
// 	if err != nil {
// 		log.Fatalf("Error unmarshalling jSON for %s (clientID: %d) err: %v \n", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
// 	}
// 	for key, record := range f {
// 		// fmt.Println(key)
// 		switch key {
// 		case "reviews":
// 			// fmt.Println("interface type", reflect.TypeOf(record))
// 			// test whether a single account, which will be of type map[string]interface{}, if so create a slice and add
// 			b := record
// 			// fmt.Printf("b = %+v\n", b)
// 			if _, ok := b.(map[string]interface{}); ok {
// 				s := make([]interface{}, 1)
// 				b = append(s, record)
// 			}

// 			if r, ok := b.([]interface{}); ok {
// 				for _, a := range r {
// 					if rec, ok := a.(map[string]interface{}); ok {
// 						// fmt.Printf("rec = %+v\n", rec)
// 						// check create time after configured time
// 						createTime := fmt.Sprintf("%s", rec["createTime"])
// 						// fmt.Printf("createTime = %s\n", createTime)
// 						// tm, err := time.Parse("2006-01-02T15:04:05.000Z", createTime)
// 						tm, err := time.Parse(time.RFC3339, createTime)
// 						if err != nil {
// 							log.Printf("Error converting createTime %s string to a time for %s (clientID: %d)\n", createTime, grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID)
// 							continue
// 						}
// 						// fmt.Printf("tm = %+v\n", tm)
// 						if tm.Before(startTime) {
// 							// The reviews are in create time descending order, so safe to stop processing
// 							break
// 						}
// 						if tm.After(endTime) {
// 							continue
// 						}
// 						// star rating can be anyone of these:
// 						// STAR_RATING_UNSPECIFIED - Not specified.
// 						// ONE - One star out of a maximum of five.
// 						// TWO - Two stars out of a maximum of five.
// 						// THREE - Three stars out of a maximum of five.
// 						// FOUR - Four stars out of a maximum of five.
// 						// FIVE - The maximum star rating.
// 						starRating := fmt.Sprintf("%s", rec["starRating"])
// 						// fmt.Printf("starRating = %s\n", starRating)
// 						switch starRating {
// 						case StarRatingUnspecified:
// 							review_ratings[StarRatingUnspecified] += 1
// 						case StarRatingOne:
// 							review_ratings[StarRatingOne] += 1
// 						case StarRatingTwo:
// 							review_ratings[StarRatingTwo] += 1
// 						case StarRatingThree:
// 							review_ratings[StarRatingThree] += 1
// 						case StarRatingFour:
// 							review_ratings[StarRatingFour] += 1
// 						case StarRatingFive:
// 							review_ratings[StarRatingFive] += 1
// 						default:
// 							continue
// 						}
// 					}
// 				}
// 			}
// 			// default:
// 			// fmt.Println(key, ":", record)
// 		}
// 	}
// 	return review_ratings
// }

// ReportOnReviewsFromJSON - report on reviews from JSON
func ReportOnReviewsFromJSON(in []byte, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, startTime time.Time, endTime time.Time) (map[string]int, string) {
	// count review star ratings
	review_ratings := make(map[string]int)
	// nextPageToken default to empty
	var nextPageToken string
	// check complete indicates that all the reviews have been checked for period and do not need to go to next page
	var checkComplete bool

	f := map[string]interface{}{
		"key": "value",
	}
	err := json.Unmarshal([]byte(in), &f)
	if err != nil {
		log.Fatalf("Error unmarshalling jSON for %s (clientID: %d) err: %v \n", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
	}
	for key, record := range f {
		// fmt.Println(key)
		switch key {
		case "reviews":
			// fmt.Println("interface type", reflect.TypeOf(record))
			// test whether a single account, which will be of type map[string]interface{}, if so create a slice and add
			b := record
			// fmt.Printf("b = %+v\n", b)
			if _, ok := b.(map[string]interface{}); ok {
				s := make([]interface{}, 1)
				b = append(s, record)
			}

			if r, ok := b.([]interface{}); ok {
				for _, a := range r {
					if rec, ok := a.(map[string]interface{}); ok {
						// fmt.Printf("rec = %+v\n", rec)
						// check update time if exists, else check create time, after configured time
						checkTime := fmt.Sprintf("%s", rec["createTime"])
						// fmt.Printf("createTime = %s\n", createTime)
						if _, ok := rec["updateTime"]; ok {
							checkTime = fmt.Sprintf("%s", rec["updateTime"])
						}
						// tm, err := time.Parse("2006-01-02T15:04:05.000Z", createTime)
						tm, err := time.Parse(time.RFC3339, checkTime)
						if err != nil {
							log.Printf("Error converting checkTime %s string to a time for %s (clientID: %d)\n", checkTime, grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID)
							continue
						}
						// fmt.Printf("tm = %+v\n", tm)
						if tm.Before(startTime) {
							// The reviews are in update time descending order, so safe to stop processing
							// set the check complete so do not fetch next page of reviews
							checkComplete = true
							break
						}
						if tm.After(endTime) {
							continue
						}
						// star rating can be anyone of these:
						// STAR_RATING_UNSPECIFIED - Not specified.
						// ONE - One star out of a maximum of five.
						// TWO - Two stars out of a maximum of five.
						// THREE - Three stars out of a maximum of five.
						// FOUR - Four stars out of a maximum of five.
						// FIVE - The maximum star rating.
						starRating := fmt.Sprintf("%s", rec["starRating"])
						// fmt.Printf("starRating = %s\n", starRating)
						switch starRating {
						case StarRatingUnspecified:
							review_ratings[StarRatingUnspecified] += 1
						case StarRatingOne:
							review_ratings[StarRatingOne] += 1
						case StarRatingTwo:
							review_ratings[StarRatingTwo] += 1
						case StarRatingThree:
							review_ratings[StarRatingThree] += 1
						case StarRatingFour:
							review_ratings[StarRatingFour] += 1
						case StarRatingFive:
							review_ratings[StarRatingFive] += 1
						default:
							continue
						}
					}
				}
			}
		case "nextPageToken":
			// If the number of reviews exceeded the requested page size, this field is populated
			// with a token to fetch the next page of reviews on a subsequent call to reviews.
			// If there are no more locations, this field is not present in the response.
			nextPageToken = record.(string)
			// default:
			// fmt.Println(key, ":", record)
		}
	}
	// do not need to fetch next page if check is complete (i.e. got all reviews for the specified period)
	if checkComplete {
		nextPageToken = ""
	}
	return review_ratings, nextPageToken
}

// // ReportOnInsights - report on insights for location
// func ReportOnInsights(client *http.Client, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, reportMonthBack int) string {
// 	insights := ""

// 	insightsPath := grcfgmbln.GoogleMyBusinessLocationPath[0:strings.LastIndex(grcfgmbln.GoogleMyBusinessLocationPath, "/locations")] +
// 		"/locations:reportInsights"
// 	// fmt.Printf("insightsPath: %s\n", insightsPath)
// 	firstOfMonth, firstOfNextMonth := reportPeriod(reportMonthBack, grcfgmbln.TimeZone)
// 	// fmt.Printf("firstOfMonth: +%v, firstOfNextMonth: +%v\n", firstOfMonth, firstOfNextMonth)
// 	// insightsBody := fmt.Sprintf("{\"locationNames\":[\"%s\"],"+
// 	// 	"\"basicRequest\":{"+
// 	// 	"\"metricRequests\":["+
// 	// 	"{\"metric\":\"QUERIES_DIRECT\"},"+
// 	// 	"{\"metric\":\"QUERIES_INDIRECT\"},"+
// 	// 	"{\"metric\":\"QUERIES_CHAIN\"},"+
// 	// 	"{\"metric\":\"ACTIONS_WEBSITE\"},"+
// 	// 	"{\"metric\":\"ACTIONS_DRIVING_DIRECTIONS\"},"+
// 	// 	"{\"metric\":\"ACTIONS_PHONE\"}"+
// 	// 	"],"+
// 	// 	"\"timeRange\":{"+
// 	// 	"\"startTime\":\"%s\","+
// 	// 	"\"endTime\":\"%s\""+
// 	// 	"}"+
// 	// 	"}"+
// 	// 	"}",
// 	// 	grcfgmbln.GoogleMyBusinessLocationPath, firstOfMonth.Format("2006-01-02T15:04:05.000000000Z"), firstOfNextMonth.Format("2006-01-02T15:04:05.000000000Z"))
// 	insightsBody := fmt.Sprintf("{\"locationNames\":[\"%s\"],"+
// 		"\"basicRequest\":{"+
// 		"\"metricRequests\":["+
// 		"{\"metric\":\"QUERIES_DIRECT\"},"+
// 		"{\"metric\":\"QUERIES_INDIRECT\"},"+
// 		"{\"metric\":\"QUERIES_CHAIN\"},"+
// 		"{\"metric\":\"ACTIONS_WEBSITE\"},"+
// 		"{\"metric\":\"ACTIONS_DRIVING_DIRECTIONS\"},"+
// 		"{\"metric\":\"ACTIONS_PHONE\"}"+
// 		"],"+
// 		"\"timeRange\":{"+
// 		"\"startTime\":\"%s\","+
// 		"\"endTime\":\"%s\""+
// 		"}"+
// 		"}"+
// 		"}",
// 		grcfgmbln.GoogleMyBusinessLocationPath, firstOfMonth.Format(time.RFC3339), firstOfNextMonth.Format(time.RFC3339))
// 	// fmt.Printf("insightsBody: %s\n", insightsBody)
// 	req, err := http.NewRequest(http.MethodPost, googleMyBusinessAPIURL+insightsPath, strings.NewReader(insightsBody))
// 	// fmt.Printf("req: %+v\n", req)
// 	if err != nil {
// 		log.Printf("Unable to create retrieve insights request for %s (clientID: %d) error: %+v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
// 		return insights
// 	}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Printf("Unable to retrieve insights for %s (clientID: %d) error: %+v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
// 		return insights
// 	}
// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Printf("Error reading body for %s (clientID: %d) for insights err: %+v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
// 		return insights
// 	}
// 	// fmt.Printf("response body: %+v\n", string(body))
// 	insights = reportOnInsightsFromJSON(body, grcfgmbln)
// 	return insights
// }

// ReportOnInsights - report on insights for location
func ReportOnInsights(client *http.Client, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, reportMonthBack int) string {
	// insights := "Business impressions on Google Maps on Desktop devices: " + reportOnDailyMetricsTimeSeriesMetric(client, grcfgmbln, reportMonthBack, "BUSINESS_IMPRESSIONS_DESKTOP_MAPS")
	// insights += "\n"
	// insights += "Business impressions on Google Search on Desktop devices: " + reportOnDailyMetricsTimeSeriesMetric(client, grcfgmbln, reportMonthBack, "BUSINESS_IMPRESSIONS_DESKTOP_SEARCH")
	// insights += "\n"
	// insights += "Business impressions on Google Maps on Mobile devices: " + reportOnDailyMetricsTimeSeriesMetric(client, grcfgmbln, reportMonthBack, "BUSINESS_IMPRESSIONS_MOBILE_MAPS")
	// insights += "\n"
	// insights += "Business impressions on Google Search on Mobile device: " + reportOnDailyMetricsTimeSeriesMetric(client, grcfgmbln, reportMonthBack, "BUSINESS_IMPRESSIONS_MOBILE_SEARCH")
	// insights += "\n"
	// insights += "The number of message conversations received on the business profile: " + reportOnDailyMetricsTimeSeriesMetric(client, grcfgmbln, reportMonthBack, "BUSINESS_CONVERSATIONS")
	// insights += "\n"
	// insights += "The number of times a direction request was requested to the business location: " + reportOnDailyMetricsTimeSeriesMetric(client, grcfgmbln, reportMonthBack, "BUSINESS_DIRECTION_REQUESTS")
	// insights += "\n"
	insights := "The number of times the business profile call button was clicked: " + reportOnDailyMetricsTimeSeriesMetric(client, grcfgmbln, reportMonthBack, "CALL_CLICKS")
	insights += "\n"
	insights += "The number of times the business profile website was clicked: " + reportOnDailyMetricsTimeSeriesMetric(client, grcfgmbln, reportMonthBack, "WEBSITE_CLICKS")
	insights += "\n"
	// insights += "The number of bookings received from the business profile: " + reportOnDailyMetricsTimeSeriesMetric(client, grcfgmbln, reportMonthBack, "BUSINESS_BOOKINGS")
	// insights += "\n"
	// insights += "The number of bookings received from the business profile: " + reportOnDailyMetricsTimeSeriesMetric(client, grcfgmbln, reportMonthBack, "BUSINESS_BOOKINGS")
	// insights += "\n"

	return insights
}

// reportOnDailyMetricsTimeSeries - report on daily metrics time series for a single metric for location
func reportOnDailyMetricsTimeSeriesMetric(client *http.Client, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, reportMonthBack int, metricToRetrieve string) string {
	metric := ""

	// fmt.Printf("grcfgmbln: %+v\n", grcfgmbln)
	firstOfMonth, firstOfNextMonth := reportPeriod(reportMonthBack, grcfgmbln.TimeZone)
	// this reports by day so need to end on last day of month and NOT first of next month
	endOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
	insightsPath := grcfgmbln.GoogleMyBusinessLocationPath[strings.LastIndex(grcfgmbln.GoogleMyBusinessLocationPath, "/locations")+1:] +
		":getDailyMetricsTimeSeries?"
	insightsPath += "dailyMetric=" + metricToRetrieve + "&"
	insightsPath += fmt.Sprintf("dailyRange.start_date.year=%d&dailyRange.start_date.month=%d&dailyRange.start_date.day=%d",
		firstOfMonth.Year(), firstOfMonth.Month(), firstOfMonth.Day())
	insightsPath += fmt.Sprintf("&dailyRange.end_date.year=%d&dailyRange.end_date.month=%d&dailyRange.end_date.day=%d",
		endOfMonth.Year(), endOfMonth.Month(), endOfMonth.Day())
	// fmt.Printf("insightsPath: %s\n", insightsPath)
	req, err := http.NewRequest(http.MethodGet, googleMyBusinessBusinessProfilePerformanceAPIURL+insightsPath, nil)
	// fmt.Printf("req: %+v\n", req)
	if err != nil {
		log.Printf("Unable to create retrieve insights request for %s (clientID: %d) error: %+v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
		return metric
	}
	resp, err := client.Do(req)
	// fmt.Printf("resp: %+v\nerr: %+v\n", resp, err)
	if err != nil {
		log.Printf("Unable to retrieve insights for %s (clientID: %d) error: %+v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
		return metric
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body for %s (clientID: %d) for insights err: %+v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
		return metric
	}
	// fmt.Printf("response body: %s\n", string(body))
	metric = ReportOnDailyMetricFromJSON(body, grcfgmbln)
	return metric
}

// // reportOnInsightsFromJSON - report on insights from JSON
// func reportOnInsightsFromJSON(in []byte, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode) string {
// 	insights := ""

// 	type totalValue struct {
// 		Value string `json:"value"`
// 	}
// 	type metricValue struct {
// 		Metric     string     `json:"metric"`
// 		TotalValue totalValue `json:"totalValue"`
// 	}
// 	type locationMetric struct {
// 		MetricValues []metricValue `json:"metricValues"`
// 	}
// 	type locationMetrics struct {
// 		LocationMetrics []locationMetric `json:"locationMetrics"`
// 	}

// 	var metrics locationMetrics
// 	err := json.Unmarshal([]byte(in), &metrics)
// 	if err != nil {
// 		log.Printf("Error unmarshalling jSON insights for %s (clientID: %d) err: %+v\n", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
// 		return insights
// 	}
// 	// fmt.Printf("metrics: %+v\n", metrics)

// 	m := make(map[string]string)
// 	for _, v := range metrics.LocationMetrics[0].MetricValues {
// 		switch v.Metric {
// 		case "QUERIES_DIRECT":
// 			m["QUERIES_DIRECT"] = v.TotalValue.Value
// 		case "QUERIES_INDIRECT":
// 			m["QUERIES_INDIRECT"] = v.TotalValue.Value
// 		case "QUERIES_CHAIN":
// 			m["QUERIES_CHAIN"] = v.TotalValue.Value
// 		case "ACTIONS_WEBSITE":
// 			m["ACTIONS_WEBSITE"] = v.TotalValue.Value
// 		case "ACTIONS_DRIVING_DIRECTIONS":
// 			m["ACTIONS_DRIVING_DIRECTIONS"] = v.TotalValue.Value
// 		case "ACTIONS_PHONE":
// 			m["ACTIONS_PHONE"] = v.TotalValue.Value
// 		}
// 	}
// 	insights += "\nHow customers search for your business:\n"
// 	insights += fmt.Sprintf("Direct: %s\n", m["QUERIES_DIRECT"])
// 	insights += fmt.Sprintf("Discovery: %s\n", m["QUERIES_INDIRECT"])
// 	insights += fmt.Sprintf("Branded: %s\n", m["QUERIES_CHAIN"])
// 	insights += "\nCustomers actions:\n"
// 	insights += fmt.Sprintf("Visit your website: %s\n", m["ACTIONS_WEBSITE"])
// 	insights += fmt.Sprintf("Request Directions: %s\n", m["ACTIONS_DRIVING_DIRECTIONS"])
// 	insights += fmt.Sprintf("Call you: %s\n", m["ACTIONS_PHONE"])

// 	return insights
// }

// ReportOnDailyMetricFromJSON - report on daily metric from JSON
func ReportOnDailyMetricFromJSON(in []byte, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode) string {
	metric := ""

	type date struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Day   int `json:"day"`
	}

	type datedValue struct {
		Date  date   `json:"date"`
		Value string `json:"value"`
	}

	type timeSeriesObject struct {
		DatedValues []datedValue `json:"datedValues"`
	}

	type timeSeries struct {
		TimeSeries timeSeriesObject `json:"timeSeries"`
	}

	var series timeSeries
	err := json.Unmarshal([]byte(in), &series)
	if err != nil {
		log.Printf("Error unmarshalling jSON insights for %s (clientID: %d) err: %+v\n", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
		return metric
	}
	// fmt.Printf("series: %+v\n", series)

	total := 0
	for _, dv := range series.TimeSeries.DatedValues {
		v, _ := strconv.Atoi(dv.Value)
		total += v
	}
	metric = fmt.Sprintf("%d", total)

	return metric
}

// reportPeriod - get the report period
func reportPeriod(reportMonthBack int, timeZone string) (reportStart, reportEnd time.Time) {
	now := time.Now()
	reportYear, reportMonth, _ := now.Date()
	// set timezone for company (not too important as will be running this monthly)
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		loc = now.Location()
	}
	// fmt.Printf("reportYear: %d, reportMonth: %d\n", reportYear, reportMonth)
	firstOfMonth := time.Date(reportYear, reportMonth, 1, 0, 0, 0, 0, loc)
	firstOfMonth = firstOfMonth.AddDate(0, reportMonthBack*-1, 0)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
	return firstOfMonth, firstOfNextMonth
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// extractReviewText extracts the text from a review
func extractReviewText(review map[string]interface{}) string {
	if comment, ok := review["comment"]; ok {
		// Try the nested map format first
		if commentMap, ok := comment.(map[string]interface{}); ok {
			if text, ok := commentMap["text"]; ok {
				return fmt.Sprintf("%s", text)
			}
		}
		// If not a map, try direct string format
		return fmt.Sprintf("%s", comment)
	}
	return ""
}

// extractReviewerName extracts the reviewer's name
func extractReviewerName(reviewer map[string]interface{}) string {
	if reviewer == nil {
		return ""
	}

	displayName := fmt.Sprintf("%s", reviewer["displayName"])
	names := strings.Split(displayName, " ")
	firstName := strings.Trim(names[0], " ")
	if len(firstName) > 1 {
		return firstName
	}
	return displayName
}

// getTemplateResponse gets a template response based on star rating
func getTemplateResponse(
	starRating string,
	reviewerName string,
	clientConfig database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode,
	replySeparator string,
) (string, bool) {
	// Get template based on star rating
	var reply string
	var enabled bool

	log.Printf("DEBUG: Getting template response for rating %s, client %s (ID: %d)",
		starRating, clientConfig.GoogleMyBusinessLocationName, clientConfig.ClientID)

	switch starRating {
	case StarRatingUnspecified:
		enabled = clientConfig.GoogleMyBusinessReplyToUnspecfifiedStarRating
		reply = clientConfig.GoogleMyBusinessUnspecfifiedStarRatingReply
		log.Printf("DEBUG: Unspecified rating template - enabled: %v, reply: %s", enabled, reply)
	case StarRatingOne:
		enabled = clientConfig.GoogleMyBusinessReplyToOneStarRating
		reply = clientConfig.GoogleMyBusinessOneStarRatingReply
		log.Printf("DEBUG: One star template - enabled: %v, reply: %s", enabled, reply)
	case StarRatingTwo:
		enabled = clientConfig.GoogleMyBusinessReplyToTwoStarRating
		reply = clientConfig.GoogleMyBusinessTwoStarRatingReply
		log.Printf("DEBUG: Two star template - enabled: %v, reply: %s", enabled, reply)
	case StarRatingThree:
		enabled = clientConfig.GoogleMyBusinessReplyToThreeStarRating
		reply = clientConfig.GoogleMyBusinessThreeStarRatingReply
		log.Printf("DEBUG: Three star template - enabled: %v, reply: %s", enabled, reply)
	case StarRatingFour:
		enabled = clientConfig.GoogleMyBusinessReplyToFourStarRating
		reply = clientConfig.GoogleMyBusinessFourStarRatingReply
		log.Printf("DEBUG: Four star template - enabled: %v, reply: %s", enabled, reply)
	case StarRatingFive:
		enabled = clientConfig.GoogleMyBusinessReplyToFiveStarRating
		reply = clientConfig.GoogleMyBusinessFiveStarRatingReply
		log.Printf("DEBUG: Five star template - enabled: %v, reply: %s", enabled, reply)
	default:
		log.Printf("DEBUG: Unknown star rating: %s", starRating)
		return "", false
	}

	if !enabled || reply == "" {
		log.Printf("DEBUG: Template not enabled or empty - enabled: %v, reply length: %d", enabled, len(reply))
		return "", false
	}

	// Handle multiple templates with separator
	if reply != "" && replySeparator != "" {
		rp := strings.Split(reply, replySeparator)
		if len(rp) > 0 {
			r := rand.Intn(len(rp))
			reply = rp[r]
			log.Printf("DEBUG: Selected random template from %d options", len(rp))
		}
	}

	// Replace name placeholder
	if len(reviewerName) > 0 {
		reply = strings.ReplaceAll(reply, "<name>", reviewerName)
		log.Printf("DEBUG: Replaced name placeholder with: %s", reviewerName)
	}

	log.Printf("DEBUG: Final template reply: %s", reply)
	return reply, true
}

// starRatingToInt converts star rating string to int
func starRatingToInt(starRating string) int {
	switch starRating {
	case StarRatingOne:
		return 1
	case StarRatingTwo:
		return 2
	case StarRatingThree:
		return 3
	case StarRatingFour:
		return 4
	case StarRatingFive:
		return 5
	default:
		return 0
	}
}

// getAIResponse generates an AI response for a review
func getAIResponse(
	reviewText string,
	starRating string,
	reviewerName string,
	clientConfig database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode,
	db *sql.DB,
) (string, error) {
	// Get configuration from properties file
	props := config.ReadProperties()

	// Check if OpenAI API key is configured
	if props.OpenAIAPIKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	// Create OpenAI provider using the centralized configuration
	provider, err := ai_service.NewOpenAIProvider(ai_service.GetConfigForUseCase(ai_service.ReviewResponse, props.OpenAIAPIKey))
	if err != nil {
		return "", fmt.Errorf("failed to create OpenAI provider: %w", err)
	}

	// Default contact method if not available in clientConfig
	contactMethod := "Contact our dispatch for assistance"
	if clientConfig.ContactMethod != nil && *clientConfig.ContactMethod != "" {
		contactMethod = *clientConfig.ContactMethod
	}

	// Create review context
	ctx := review.ReviewContext{
		Text:          reviewText,
		Rating:        starRatingToInt(starRating),
		Author:        reviewerName,
		BusinessName:  clientConfig.GoogleMyBusinessLocationName,
		Location:      clientConfig.GoogleMyBusinessPostalCode,
		ContactMethod: contactMethod,
	}

	// Generate response
	generator := review.NewGenerator(provider, review.GeneratorConfig{})
	response, err := generator.Generate(ctx)
	if err != nil {
		return "", fmt.Errorf("AI generation failed: %w", err)
	}

	return response, nil
}

func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength]
}

// FetchReviews fetches reviews for a Google My Business location within a specified time period
//
// Parameters:
// - client: The HTTP client to use for making API requests
// - grcfgmbln: The location configuration containing path and other details
// - reviewsNotBeforeDays: Number of days to look back for reviews
//
// Returns:
// - A slice of review data as maps
// - An error if the operation fails
//
// This function is designed to retrieve reviews without applying responses,
// allowing for analysis and reporting on the review content.
// It handles pagination automatically to fetch all available reviews.
func FetchReviews(client *http.Client, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, reviewsNotBeforeDays int) ([]map[string]interface{}, error) {
	var allReviews []map[string]interface{}
	var pageToken string
	for {
		url := googleMyBusinessAPIURL + grcfgmbln.GoogleMyBusinessLocationPath + "/reviews"
		if len(strings.Trim(pageToken, " ")) > 0 {
			url += "?pageToken=" + strings.Trim(pageToken, " ")
		}
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve reviews for %s (clientID: %d) error: %v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading body for %s (clientID: %d) err: %v", grcfgmbln.GoogleMyBusinessLocationName, grcfgmbln.ClientID, err)
		}
		nextPageToken, reviews := processReviewsFromJSON(client, body, grcfgmbln, reviewsNotBeforeDays, "", false)
		allReviews = append(allReviews, reviews...)
		if nextPageToken == "" {
			break
		} else {
			pageToken = nextPageToken
		}
	}
	return allReviews, nil
}

// FetchReviewsForMonthlyAnalysis fetches reviews for a Google My Business location within a specific month period
//
// Parameters:
// - client: The HTTP client to use for making API requests
// - grcfgmbln: The location configuration containing path and other details
// - monthsBack: Number of months to look back (e.g., 1 for last month)
//
// Returns:
// - A slice of review data as maps
// - The start time of the period
// - The end time of the period
// - An error if the operation fails
//
// This function is designed specifically for the monthly review analysis process,
// filtering reviews to a specific month and handling pagination appropriately.
// It reuses the existing FetchReviews function but adds date-specific filtering.
func FetchReviewsForMonthlyAnalysis(client *http.Client, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, monthsBack int) ([]map[string]interface{}, time.Time, time.Time, error) {
	// Get the start and end dates for the requested month
	startTime, endTime := reportPeriod(monthsBack, grcfgmbln.TimeZone)

	// First fetch all reviews using the existing function
	// Use a large enough reviewsNotBeforeDays to cover the entire period
	// The maximum number of days we might need to look back is approximately 31 * (monthsBack + 1)
	daysToLookBack := (monthsBack + 1) * 31

	allReviews, err := FetchReviews(client, grcfgmbln, daysToLookBack)
	if err != nil {
		return nil, startTime, endTime, fmt.Errorf("failed to fetch reviews for monthly analysis: %w", err)
	}

	// Filter the reviews to only include those within the target month
	var filteredReviews []map[string]interface{}

	for _, review := range allReviews {
		// Extract the create time from the review
		createTimeStr, ok := review["createTime"].(string)
		if !ok {
			log.Printf("Warning: Review missing createTime field for %s", grcfgmbln.GoogleMyBusinessLocationName)
			continue
		}

		// Parse the create time
		createTime, err := time.Parse(time.RFC3339, createTimeStr)
		if err != nil {
			log.Printf("Warning: Failed to parse review createTime '%s': %v", createTimeStr, err)
			continue
		}

		// Check if the review is within our target month
		if (createTime.Equal(startTime) || createTime.After(startTime)) &&
			createTime.Before(endTime) {
			filteredReviews = append(filteredReviews, review)
		}
	}

	log.Printf("Fetched %d reviews for %s (period: %s to %s), filtered to %d reviews",
		len(allReviews),
		grcfgmbln.GoogleMyBusinessLocationName,
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02"),
		len(filteredReviews))

	return filteredReviews, startTime, endTime, nil
}

// FetchReviewsForMonth fetches reviews for a Google My Business location for a specific month
//
// Parameters:
// - client: The HTTP client to use for making API requests
// - grcfgmbln: The location configuration containing path and other details
// - targetMonth: A time.Time value representing the target month (day component is ignored)
//
// Returns:
// - A slice of review data as maps
// - The start time of the period (first day of the month)
// - The end time of the period (first day of the next month)
// - An error if the operation fails
//
// This function is designed specifically for the monthly review analysis process,
// filtering reviews to a specific month and handling pagination appropriately.
// It reuses the existing FetchReviews function but adds date-specific filtering.
func FetchReviewsForMonth(client *http.Client, grcfgmbln database.GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode, targetMonth time.Time) ([]map[string]interface{}, time.Time, time.Time, error) {
	// Load the location's timezone
	loc, err := time.LoadLocation(grcfgmbln.TimeZone)
	if err != nil {
		log.Printf("Warning: Failed to load timezone %s for %s, using local timezone",
			grcfgmbln.TimeZone, grcfgmbln.GoogleMyBusinessLocationName)
		loc = time.Local
	}

	// Calculate the start and end dates for the specified month
	// Start is the first day of the target month
	startTime := time.Date(targetMonth.Year(), targetMonth.Month(), 1, 0, 0, 0, 0, loc)
	// End is the first day of the next month
	endTime := startTime.AddDate(0, 1, 0)

	// Calculate how many days to look back from today to cover the target month
	now := time.Now()
	daysToLookBack := int(now.Sub(startTime).Hours()/24) + 1

	// If the target month is in the future, log a warning and return empty results
	if daysToLookBack < 0 {
		log.Printf("Warning: Target month %s is in the future for %s",
			targetMonth.Format("2006-01"), grcfgmbln.GoogleMyBusinessLocationName)
		return []map[string]interface{}{}, startTime, endTime, nil
	}

	// If the target month is too far in the past, add a reasonable cap to avoid performance issues
	if daysToLookBack > 365*2 { // Cap at 2 years
		log.Printf("Warning: Target month %s is more than 2 years in the past for %s, limiting search depth",
			targetMonth.Format("2006-01"), grcfgmbln.GoogleMyBusinessLocationName)
		daysToLookBack = 365 * 2
	}

	// Fetch all reviews within the calculated lookback period
	allReviews, err := FetchReviews(client, grcfgmbln, daysToLookBack)
	if err != nil {
		return nil, startTime, endTime, fmt.Errorf("failed to fetch reviews for monthly analysis: %w", err)
	}

	// Filter the reviews to only include those within the target month
	var filteredReviews []map[string]interface{}

	for _, review := range allReviews {
		// Extract the create time from the review
		createTimeStr, ok := review["createTime"].(string)
		if !ok {
			log.Printf("Warning: Review missing createTime field for %s", grcfgmbln.GoogleMyBusinessLocationName)
			continue
		}

		// Parse the create time
		createTime, err := time.Parse(time.RFC3339, createTimeStr)
		if err != nil {
			log.Printf("Warning: Failed to parse review createTime '%s': %v", createTimeStr, err)
			continue
		}

		// Check if the review is within our target month
		if (createTime.Equal(startTime) || createTime.After(startTime)) &&
			createTime.Before(endTime) {
			filteredReviews = append(filteredReviews, review)
		}
	}

	log.Printf("Fetched %d reviews for %s (period: %s to %s), filtered to %d reviews",
		len(allReviews),
		grcfgmbln.GoogleMyBusinessLocationName,
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02"),
		len(filteredReviews))

	return filteredReviews, startTime, endTime, nil
}
