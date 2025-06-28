package google_my_business_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"

	"rm_client_portal/database"
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

// token from Google OAuth
var googleToken *oauth2.Token

// client for calling Google APIs
var Client *http.Client

// Retrieve a token, saves the token, then sets the client.
func SetClient(config *oauth2.Config) {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		getTokenFromWeb(config)
		if googleToken != nil {
			saveToken(tokFile, googleToken)
		}
	}
	Client = config.Client(context.Background(), tok)
}

// Request a token from the web, the token is set in the global variable googleToken.
func getTokenFromWeb(config *oauth2.Config) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	log.Printf("getTokenFromWeb: starting HTTP server\n")
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	startHttpServerForTokenUri(httpServerExitDone, config)
	// wait for goroutine started in startHttpServerForTokenUri() to stop
	httpServerExitDone.Wait()
	log.Printf("getTokenFromWeb: done. exiting\n")
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// start HTTP server for receiving the OAuth token from Google.
// The credentials.json file contain the "redirect_uris":["http://localhost"] to reference this server.
// request example: http://localhost/?state=state-token&code=4/0AX4XfWjODE_aomPgdCJRX5AWAzfzIWUbRxrZJ_M4Z2NY8VvdyeH6XbsM1Sm3sYLcIwgFag&scope=https://www.googleapis.com/auth/business.manage
func startHttpServerForTokenUri(wg *sync.WaitGroup, config *oauth2.Config) {
	srv := &http.Server{
		Addr: ":80",
	}
	idleConnsClosed := make(chan struct{})

	// sends an HTTP response to the browser
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("startHttpServerForTokenUri: request raw query: %s\n", r.URL.RawQuery)
		authCode := r.URL.Query().Get("code")
		if authCode != "" {
			log.Println("startHttpServerForTokenUri: Found authorization code:", authCode)
			var err error
			googleToken, err = config.Exchange(context.Background(), authCode)
			if err != nil {
				log.Fatalf("startHttpServerForTokenUri: Unable to retrieve token from web: %v", err)
			}
			log.Println("startHttpServerForTokenUri: Shutting down ...")
			// graceful-shutdown
			err = srv.Shutdown(context.Background())
			if err != nil {
				log.Printf("startHttpServerForTokenUri: server.Shutdown: %v\n", err)
			}
			close(idleConnsClosed)
		} else {
			fmt.Fprintln(w, "No token found") // Just display a default page if the query string does not have a token (code) in it return to server at /
		}
	})

	go func() {
		defer wg.Done() // let calling function know done cleaning up
		// always return error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Printf("startHttpServerForTokenUri: Error ListenAndServer(): %v\n", err)
		}
		log.Println("startHttpServerForTokenUri: Bye.")
	}()
}

// GetAccounts - get accounts
func GetAccounts() []string {
	resp, err := Client.Get(googleMyBusinessAccountManagementAPIURL + "accounts")
	if err != nil {
		log.Fatalf("Unable to retrieve accounts error: %v", err)
	}
	defer resp.Body.Close()
	// fmt.Println(resp)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading body err: %v", err)
	}
	// fmt.Println(body)
	// fmt.Println(string(body))
	accounts := GetAccountsFromJSON(body)
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

// GetLocationsCheckClientID - get locations for account checking the client ID
// report parameter indicates if this is for reporting rather than replying
func GetLocationsCheckClientID(account string, clientIDs []uint64) []database.GoogleReviewsConfigAndGoogleMyBusinessLocation {
	var pageToken string
	var grcagmblns []database.GoogleReviewsConfigAndGoogleMyBusinessLocation
	for {
		grcagmblnsp, nextPageToken := getLocationsCheckClientIDPage(account, clientIDs, pageToken)
		grcagmblns = append(grcagmblns, grcagmblnsp...)
		if nextPageToken == "" {
			break
		} else {
			pageToken = nextPageToken
		}
	}
	return grcagmblns
}

// getLocationsCheckClientIDPage - get locations for account page at a time (default page size is 100)
// check the client ID
// pageToken parameter indicates when the locations has already been called and there are more results
func getLocationsCheckClientIDPage(account string, clientIDs []uint64, pageToken string) ([]database.GoogleReviewsConfigAndGoogleMyBusinessLocation, string) {
	locationsParams := ""
	if pageToken != "" {
		locationsParams += "?pageToken=" + pageToken
	}
	readMaskParams := "readMask=name,title,storefrontAddress"
	if locationsParams == "" {
		readMaskParams = "?" + readMaskParams
	} else {
		readMaskParams = "&" + readMaskParams
	}
	resp, err := Client.Get(googleMyBusinessBusinessInformationAPIURL + account + "/locations" + locationsParams + readMaskParams)
	if err != nil {
		log.Fatalf("Unable to retrieve locations error: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading body err: %v", err)
	}
	// fmt.Println(string(body))
	grcagmblns, nextPageToken := GetLocationsCheckClientIDFromJSON(body, account, clientIDs)
	// fmt.Printf("nextPageToken = %s\n", nextPageToken)
	return grcagmblns, nextPageToken
}

// GetLocationsCheckClientIDFromJSON - get locations from JSON check client ID
func GetLocationsCheckClientIDFromJSON(in []byte, account string, clientIDs []uint64) ([]database.GoogleReviewsConfigAndGoogleMyBusinessLocation, string) {
	f := map[string]interface{}{
		"key": "value",
	}
	err := json.Unmarshal([]byte(in), &f)
	if err != nil {
		log.Fatal(err)
	}
	var grcagmblns []database.GoogleReviewsConfigAndGoogleMyBusinessLocation
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
						locationName := fmt.Sprintf("%s", rec["title"])
						// fmt.Printf("locationName = %s\n", locationName)
						address := rec["storefrontAddress"]
						// fmt.Printf("reviewer = %+v\n", reviewer)
						postalCode := ""
						// address for display
						displayAddress := ""
						locality := ""
						if v, ok := address.(map[string]interface{}); ok {
							postalCode = fmt.Sprintf("%s", v["postalCode"])
							// fmt.Printf("postalCode = %s\n", postalCode)
							addressLines := v["addressLines"]
							if als, ok := addressLines.([]interface{}); ok {
								for _, al := range als {
									if a, ok := al.(string); ok {
										if len(displayAddress) > 0 {
											displayAddress += ", "
										}
										displayAddress += a
									}
								}
							}
							locality = fmt.Sprintf("%s", v["locality"])
						}
						// require a postal code
						if len(postalCode) == 0 {
							continue
						}
						grcagmbln := database.ConfigFromGoogleMyBusinessLocationNameAndPostalCode(locationName, postalCode)
						// fmt.Printf("google reviews config fields from google my business location name: %v\n", grcagmbln)
						// the name holds the path to be appended to the google api URL
						name := fmt.Sprintf("%s/%s", account, rec["name"])
						grcagmbln.GoogleMyBusinessLocationPath = name
						grcagmbln.GoogleMyBusinessLocationAddress = displayAddress
						grcagmbln.GoogleMyBusinessLocality = locality
						if grcagmbln.ClientID != 0 {
							// check client ID
							found := false
							for _, id := range clientIDs {
								if grcagmbln.ClientID == id {
									found = true
									break
								}
							}
							if found {
								grcagmblns = append(grcagmblns, grcagmbln)
							}
						}
					}
				}
			}
		case "nextPageToken":
			// If the number of locations exceeded the requested page size, this field is populated
			// with a token to fetch the next page of locations on a subsequent call to locations.list.
			// If there are no more locations, this field is not present in the response.
			nextPageToken = record.(string)
			// default:
			// fmt.Println(key, ":", record)
		}
	}
	// fmt.Printf("grcagmblns = %+v", grcagmblns)
	return grcagmblns, nextPageToken
}

// ReportOnReviewsWeb - report on reviews for location
func ReportOnReviewsWeb(grcagmbln database.GoogleReviewsConfigAndGoogleMyBusinessLocation, reportStart, reportEnd string) database.GoogleReviewRatingsFromGoogleMyBusiness {
	var pageToken string
	var reviews database.GoogleReviewRatingsFromGoogleMyBusiness
	reportStartTime := ReportTime(reportStart, grcagmbln.TimeZone)
	reportEndTime := ReportTime(reportEnd, grcagmbln.TimeZone)
	for {
		pageReviews, nextPageToken := reportOnReviewsWebPage(grcagmbln, pageToken, reportStartTime, reportEndTime)
		// append reviews page
		reviews.StarRatingUnspecified += pageReviews.StarRatingUnspecified
		reviews.StarRatingOne += pageReviews.StarRatingOne
		reviews.StarRatingTwo += pageReviews.StarRatingTwo
		reviews.StarRatingThree += pageReviews.StarRatingThree
		reviews.StarRatingFour += pageReviews.StarRatingFour
		reviews.StarRatingFive += pageReviews.StarRatingFive
		if nextPageToken == "" {
			break
		} else {
			pageToken = nextPageToken
		}
	}
	return reviews
}

// reportOnReviewsWebPage - report on reviews for location per page
func reportOnReviewsWebPage(grcagmbln database.GoogleReviewsConfigAndGoogleMyBusinessLocation, pageToken string, startTime time.Time, endTime time.Time) (database.GoogleReviewRatingsFromGoogleMyBusiness, string) {
	url := googleMyBusinessAPIURL + grcagmbln.GoogleMyBusinessLocationPath + "/reviews"
	if len(strings.Trim(pageToken, " ")) > 0 {
		url += "?pageToken=" + strings.Trim(pageToken, " ")
	}
	resp, err := Client.Get(url)
	if err != nil {
		log.Fatalf("Unable to retrieve reviews for %s (clientID: %d) for report error: %v", grcagmbln.GoogleMyBusinessLocationName, grcagmbln.ClientID, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading body for %s (clientID: %d) for report err: %v", grcagmbln.GoogleMyBusinessLocationName, grcagmbln.ClientID, err)
	}
	// fmt.Println(string(body))
	reviews, nextPageToken := ReportOnReviewsWebFromJSON(body, grcagmbln, startTime, endTime)
	// fmt.Printf("nextPageToken = %s\n", nextPageToken)
	return reviews, nextPageToken
}

// ReportOnReviewsWebFromJSON - report on reviews from JSON
func ReportOnReviewsWebFromJSON(in []byte, grcagmbln database.GoogleReviewsConfigAndGoogleMyBusinessLocation, startTime time.Time, endTime time.Time) (database.GoogleReviewRatingsFromGoogleMyBusiness, string) {
	// count review star ratings
	var grrfgmb database.GoogleReviewRatingsFromGoogleMyBusiness
	// nextPageToken default to empty
	var nextPageToken string
	// check complete indicates that all the reviews have been checked for period and do not need to go to next page
	var checkComplete bool

	f := map[string]interface{}{
		"key": "value",
	}
	err := json.Unmarshal([]byte(in), &f)
	if err != nil {
		log.Fatalf("Error unmarshalling jSON for %s (clientID: %d) err: %v \n", grcagmbln.GoogleMyBusinessLocationName, grcagmbln.ClientID, err)
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
							log.Printf("Error converting checkTime %s string to a time for %s (clientID: %d)\n", checkTime, grcagmbln.GoogleMyBusinessLocationName, grcagmbln.ClientID)
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
							grrfgmb.StarRatingUnspecified += 1
						case StarRatingOne:
							grrfgmb.StarRatingOne += 1
						case StarRatingTwo:
							grrfgmb.StarRatingTwo += 1
						case StarRatingThree:
							grrfgmb.StarRatingThree += 1
						case StarRatingFour:
							grrfgmb.StarRatingFour += 1
						case StarRatingFive:
							grrfgmb.StarRatingFive += 1
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
	return grrfgmb, nextPageToken
}

// ReportOnInsightsWeb - report on insights for location
func ReportOnInsightsWeb(grcagmbln database.GoogleReviewsConfigAndGoogleMyBusinessLocation, reportStart, reportEnd string) database.GoogleInsightsFromGoogleMyBusiness {
	// results := make(map[string]int)
	// // results["Business impressions on Google Maps on Desktop devices"] = reportOnDailyMetricsTimeSeriesMetric(grcagmbln, reportStart, reportEnd, "BUSINESS_IMPRESSIONS_DESKTOP_MAPS")
	// // results["Business impressions on Google Search on Desktop devices"] = reportOnDailyMetricsTimeSeriesMetric(grcagmbln, reportStart, reportEnd, "BUSINESS_IMPRESSIONS_DESKTOP_SEARCH")
	// // results["Business impressions on Google Maps on Mobile devices"] = reportOnDailyMetricsTimeSeriesMetric(grcagmbln, reportStart, reportEnd, "BUSINESS_IMPRESSIONS_MOBILE_MAPS")
	// // results["Business impressions on Google Search on Mobile device"] = reportOnDailyMetricsTimeSeriesMetric(grcagmbln, reportStart, reportEnd, "BUSINESS_IMPRESSIONS_MOBILE_SEARCH")
	// // results["The number of message conversations received on the business profile"] = reportOnDailyMetricsTimeSeriesMetric(grcagmbln, reportStart, reportEnd, "BUSINESS_CONVERSATIONS")
	// // results["The number of times a direction request was requested to the business location"] = reportOnDailyMetricsTimeSeriesMetric(grcagmbln, reportStart, reportEnd, "BUSINESS_DIRECTION_REQUESTS")
	// results["The number of times the business profile call button was clicked"] = reportOnDailyMetricsTimeSeriesMetric(grcagmbln, reportStart, reportEnd, "CALL_CLICKS")
	// results["The number of times the business profile website was clicked"] = reportOnDailyMetricsTimeSeriesMetric(grcagmbln, reportStart, reportEnd, "WEBSITE_CLICKS")
	// // results["The number of bookings received from the business profile"] = reportOnDailyMetricsTimeSeriesMetric(grcagmbln, reportStart, reportEnd, "BUSINESS_BOOKINGS")
	// // results["The number of bookings received from the business profile"] = reportOnDailyMetricsTimeSeriesMetric(grcagmbln, reportStart, reportEnd, "BUSINESS_BOOKINGS")

	var results database.GoogleInsightsFromGoogleMyBusiness
	results.NumberOfBusinessProfileCallButtonClicked = reportOnDailyMetricsTimeSeriesMetric(grcagmbln, reportStart, reportEnd, "CALL_CLICKS")
	results.NumberOfBusinessProfileWebsiteClicked = reportOnDailyMetricsTimeSeriesMetric(grcagmbln, reportStart, reportEnd, "WEBSITE_CLICKS")

	return results
}

// reportOnDailyMetricsTimeSeries - report on daily metrics time series for a single metric for location
func reportOnDailyMetricsTimeSeriesMetric(grcagmbln database.GoogleReviewsConfigAndGoogleMyBusinessLocation, reportStart, reportEnd, metricToRetrieve string) int {
	metric := 0

	// fmt.Printf("grcagmbln: %+v\n", grcagmbln)
	reportStartTime := ReportTime(reportStart, grcagmbln.TimeZone)
	reportEndTime := ReportTime(reportEnd, grcagmbln.TimeZone)
	insightsPath := grcagmbln.GoogleMyBusinessLocationPath[strings.LastIndex(grcagmbln.GoogleMyBusinessLocationPath, "/locations")+1:] +
		":getDailyMetricsTimeSeries?"
	insightsPath += "dailyMetric=" + metricToRetrieve + "&"
	insightsPath += fmt.Sprintf("dailyRange.start_date.year=%d&dailyRange.start_date.month=%d&dailyRange.start_date.day=%d",
		reportStartTime.Year(), reportStartTime.Month(), reportStartTime.Day())
	insightsPath += fmt.Sprintf("&dailyRange.end_date.year=%d&dailyRange.end_date.month=%d&dailyRange.end_date.day=%d",
		reportEndTime.Year(), reportEndTime.Month(), reportEndTime.Day())
	// fmt.Printf("insightsPath: %s\n", insightsPath)
	req, err := http.NewRequest(http.MethodGet, googleMyBusinessBusinessProfilePerformanceAPIURL+insightsPath, nil)
	// fmt.Printf("req: %+v\n", req)
	if err != nil {
		log.Printf("Unable to create retrieve insights request for %s (clientID: %d) error: %+v", grcagmbln.GoogleMyBusinessLocationName, grcagmbln.ClientID, err)
		return metric
	}
	resp, err := Client.Do(req)
	// fmt.Printf("resp: %+v\nerr: %+v\n", resp, err)
	if err != nil {
		log.Printf("Unable to retrieve insights for %s (clientID: %d) error: %+v", grcagmbln.GoogleMyBusinessLocationName, grcagmbln.ClientID, err)
		return metric
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body for %s (clientID: %d) for insights err: %+v", grcagmbln.GoogleMyBusinessLocationName, grcagmbln.ClientID, err)
		return metric
	}
	// fmt.Printf("response body: %s\n", string(body))
	metric = ReportOnDailyMetricFromJSON(body, grcagmbln)
	return metric
}

// ReportOnDailyMetricFromJSON - report on daily metric from JSON
func ReportOnDailyMetricFromJSON(in []byte, grcagmbln database.GoogleReviewsConfigAndGoogleMyBusinessLocation) int {
	metric := 0

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
		log.Printf("Error unmarshalling jSON insights for %s (clientID: %d) err: %+v\n", grcagmbln.GoogleMyBusinessLocationName, grcagmbln.ClientID, err)
		return metric
	}
	// fmt.Printf("series: %+v\n", series)

	total := 0
	for _, dv := range series.TimeSeries.DatedValues {
		v, _ := strconv.Atoi(dv.Value)
		total += v
	}
	metric = total

	return metric
}

// ReportTime - convert the report time sent as a string to a time in the correct time zone
func ReportTime(reportTimeStr string, timeZone string) time.Time {
	now := time.Now()
	r, err := time.Parse(time.RFC3339, reportTimeStr)
	if err != nil {
		log.Printf("Error converting report time string: %s, to time", reportTimeStr)
		// set time to now
		r = now
	}
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		loc = now.Location()
	}
	rt := r.In(loc)
	return rt
}
