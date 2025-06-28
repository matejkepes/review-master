package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

	"google_reviews_ui/config"
	"google_reviews_ui/database"
	"google_reviews_ui/shared"
)

const statsLogPrefix = "stats: "
const statsFilename = "stats.json"

// ClientsHandler - retrieve a list of clients
func ClientsHandler(c *gin.Context) {
	clients, err := database.ListAllClients(getPartnerID(c))
	if err != nil {
		log.Printf("error retrieving client list, err: %+v\n", err)
	}
	c.JSON(200, gin.H{
		"clients": clients,
	})
}

// SimpleGetClientHandler - retrieve a config for a client assuming only 1 config and 1 time
func SimpleGetClientHandler(c *gin.Context) {
	var simpleConfig database.SimpleConfig
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		log.Printf("error (client id error) retrieving simple config for a client id: %s, err: %+v\n", c.Query("id"), err)
	} else {
		simpleConfig, err = database.GetSimpleClient(id, getPartnerID(c))
		if err != nil {
			log.Printf("error retrieving simple config for a client id: %s, err: %+v\n", c.Query("id"), err)
		}
	}
	c.JSON(200, gin.H{
		"simpleConfig": simpleConfig,
	})
}

// SimpleUpdateClientHandler - update a config for a client assuming only 1 config and 1 time
func SimpleUpdateClientHandler(c *gin.Context) {
	success := true
	var errStr string
	var simpleConfig database.SimpleConfig
	if err := c.ShouldBind(&simpleConfig); err != nil {
		log.Printf("err: %+v\n", err)
		errStr = fmt.Sprintf("error: %+v", err)
		success = false
	} else {
		// log.Printf("simpleConfig: %+v\n", simpleConfig)
		err := database.UpdateSimpleClient(simpleConfig)
		if err != nil {
			log.Printf("error updating simple config for a client, err: %+v\n", err)
			errStr = fmt.Sprintf("error updating simple config for a client, error: %+v", err)
			success = false
		}
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
	})
}

// SimpleCreateClientHandler - create a config for a client assuming only 1 config and 1 time
func SimpleCreateClientHandler(c *gin.Context) {
	success := true
	var errStr string
	var simpleConfig database.SimpleConfig
	if err := c.ShouldBind(&simpleConfig); err != nil {
		log.Printf("err: %+v\n", err)
		errStr = fmt.Sprintf("error: %+v", err)
		success = false
	} else {
		// log.Printf("simpleConfig: %+v\n", simpleConfig)
		err := database.CreateSimpleClient(simpleConfig, getPartnerID(c))
		if err != nil {
			log.Printf("error creating simple config for a client, err: %+v\n", err)
			errStr = fmt.Sprintf("error creating simple config for a client, error: %+v", err)
			success = false
		}
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
	})
}

// GetClientHandler - retrieve configs for a client with many configs with many times
func GetClientHandler(c *gin.Context) {
	var config database.ClientConfig
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		log.Printf("error (client id error) retrieving configs for a client id: %s, err: %+v\n", c.Query("id"), err)
	} else {
		config, err = database.GetClient(id, getPartnerID(c))
		if err != nil {
			log.Printf("error retrieving configs for a client id: %s, err: %+v\n", c.Query("id"), err)
		}
	}
	c.JSON(200, gin.H{
		"config": config,
	})
}

// UpdateClientHandler - update a config for a client with many configs with many times
func UpdateClientHandler(c *gin.Context) {
	success := true
	var errStr string
	var clientConfig database.ClientConfig

	// Copy the request body
	var buf bytes.Buffer
	tee := io.TeeReader(c.Request.Body, &buf)
	body, _ := io.ReadAll(tee)
	c.Request.Body = io.NopCloser(&buf)

	// Log the raw request
	log.Printf("Raw request data: %s\n", string(body))

	if err := c.ShouldBind(&clientConfig); err != nil {
		log.Printf("Binding error: %+v\n", err)
		errStr = fmt.Sprintf("error: %+v", err)
		success = false
	} else {
		// Add debug logging for parsed config
		log.Printf("Parsed clientConfig: %+v\n", clientConfig)
		err := database.UpdateClient(clientConfig, getPartnerID(c))
		if err != nil {
			log.Printf("error updating config for a client, err: %+v\n", err)
			errStr = fmt.Sprintf("error updating config for a client, error: %+v", err)
			success = false
		}
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
	})
}

// CreateClientHandler - create a config for a client with many configs with many times
func CreateClientHandler(c *gin.Context) {
	success := true
	var errStr string
	var clientConfig database.ClientConfig
	if err := c.ShouldBind(&clientConfig); err != nil {
		log.Printf("err: %+v\n", err)
		errStr = fmt.Sprintf("error: %+v", err)
		success = false
	} else {
		// log.Printf("clientConfig: %+v\n", clientConfig)
		err := database.CreateClient(clientConfig, getPartnerID(c))
		if err != nil {
			log.Printf("error creating config for a client, err: %+v\n", err)
			errStr = fmt.Sprintf("error creating config for a client, error: %+v", err)
			success = false
		}
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
	})
}

// CreateClientConfigHandler - create a config for a client
func CreateClientConfigHandler(c *gin.Context) {
	success := true
	var errStr string
	var googleReviewsConfig database.GoogleReviewsConfig
	if err := c.ShouldBind(&googleReviewsConfig); err != nil {
		log.Printf("err: %+v\n", err)
		errStr = fmt.Sprintf("error: %+v", err)
		success = false
	} else {
		// log.Printf("googleReviewsConfig: %+v\n", googleReviewsConfig)
		err := database.CreateGRConfig(googleReviewsConfig)
		if err != nil {
			log.Printf("error creating config for a client, err: %+v\n", err)
			errStr = fmt.Sprintf("error creating config for a client, error: %+v", err)
			success = false
		}
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
	})
}

// CreateClientConfigTimeHandler - create a config time for a client
func CreateClientConfigTimeHandler(c *gin.Context) {
	success := true
	var errStr string
	var googleReviewsConfigTime database.GoogleReviewsConfigTime
	if err := c.ShouldBind(&googleReviewsConfigTime); err != nil {
		log.Printf("err: %+v\n", err)
		errStr = fmt.Sprintf("error: %+v", err)
		success = false
	} else {
		// log.Printf("googleReviewsConfigTime: %+v\n", googleReviewsConfigTime)
		err := database.CreateGRCTime(googleReviewsConfigTime)
		if err != nil {
			log.Printf("error creating config time for a client, err: %+v\n", err)
			errStr = fmt.Sprintf("error creating config time for a client, error: %+v", err)
			success = false
		}
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
	})
}

// StatsHandler - retrieve a list of stats
func StatsHandler(c *gin.Context) {
	success := true
	var errStr string
	var stats []database.StatsResult
	var statsRequest database.StatsRequest
	startDay := c.Query("start_day")
	endDay := c.Query("end_day")
	timeGrouping := c.Query("time_grouping")
	// log.Printf("start_day: %s\n", startDay)
	statsRequest.StartDay = startDay
	statsRequest.EndDay = endDay
	statsRequest.TimeGrouping = timeGrouping
	// log.Printf("statsRequest: %+v\n", statsRequest)
	stats, err := database.Stats(statsRequest, getPartnerID(c))
	if err != nil {
		log.Printf("error retrieving stats list, err: %+v\n", err)
		errStr = fmt.Sprintf("error retrieving stats list, error: %+v", err)
		success = false
	}
	// Need record stats to file so can retrieve it later.
	// This was originally done in the log.
	// This is done due to network failure when retrieving stats (I was unable to solve this).
	// Need to marshall to JSON before putting in log
	m, err := json.Marshal(stats)
	if err != nil {
		log.Printf("Error marshalling stats: %+v to a string, err: %v\n", stats, err)
	} else {
		// log.Printf("%s%s\n", statsLogPrefix, m)
		createStatsFile(string(m))
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
		"stats":   stats,
	})
}

// StatsFromLogHandler - retrieve a list of stats from the log file. The log is created in the StatsHandler.
// This is done because the frontend is failing to retrieve the stats sometimes with a network error.
func StatsFromLogHandler(c *gin.Context) {
	lf := c.MustGet(shared.LoggerFilename)
	logFilename := fmt.Sprintf("%s", lf)
	success := true
	var errStr string
	var stats []database.StatsResult
	statsString := ""
	f, err := os.Open(logFilename) // For read access.
	if err != nil {
		errStr = fmt.Sprintf("Error, opening log file: %s to view logs with error: %v\n", logFilename, err)
		log.Print(errStr)
		success = false
	}
	defer f.Close()
	if success {
		scanner := bufio.NewScanner(f)
		// log entry time format
		re := regexp.MustCompile(`^\d\d\d\d\/\d\d\/\d\d \d\d:\d\d:\d\d `)
		for scanner.Scan() {
			t := scanner.Text()
			// remove time (contains date and time at beginning of each log entry it contains a date and time separated by a space)
			s := re.ReplaceAllString(t, ``)
			if strings.HasPrefix(s, statsLogPrefix) {
				s = strings.Replace(s, statsLogPrefix, "", 1)
				statsString = s
			}
		}
		if err = scanner.Err(); err != nil {
			errStr = fmt.Sprintf("Error, reading log file: %s to view logs with error: %v\n", logFilename, err)
			log.Print(errStr)
			success = false
		}
		if statsString != "" {
			err = json.Unmarshal([]byte(statsString), &stats)
			if err != nil {
				errStr = fmt.Sprintf("Error, unmarshalling: %s from log to JSON, error: %v\n", statsString, err)
				log.Print(errStr)
				success = false
			}
		} else {
			success = false
		}
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
		"stats":   stats,
	})
}

// StatsFromFileHandler - retrieve a list of stats from the stats file. The stats file is created in the StatsHandler.
// This is done because the frontend is failing to retrieve the stats sometimes with a network error.
// NOTE: The stats file is created everytime new stats is requested so there is only one list of stats in the file.
func StatsFromFileHandler(c *gin.Context) {
	success := true
	var errStr string
	var stats []database.StatsResult
	statsString := ""
	// create file in same directory as executable
	fp := "./"
	ex, err := os.Executable()
	if err != nil {
		log.Printf("Unable to determine executable location for reading stats file, err: %+v\n", err)
	} else {
		fp = filepath.Dir(ex) + "/"
	}
	statsFile := fp + statsFilename
	// open file
	f, err := os.Open(statsFile)
	if err != nil {
		log.Printf("Error opening stats file for reading, error: %+v\n", err)
		errStr = fmt.Sprintf("Error, opening stats file: %s to read stats with error: %v\n", statsFile, err)
		log.Print(errStr)
		success = false
	}
	defer f.Close()
	if success {
		// Need to use reader rather than scanner because scanner cannot handle really long lines
		r := bufio.NewReader(f)
		for {
			s, err := r.ReadString(10) // 0x0A separator = newline
			if err == io.EOF {
				log.Println("Reached end of file")
				break
			}
			statsString += s
		}
		if statsString != "" {
			// statsString = strings.Trim(statsString, "\x00")
			err = json.Unmarshal([]byte(statsString), &stats)
			if err != nil {
				errStr = fmt.Sprintf("Error, unmarshalling: %s from stats file to JSON, error: %v\n", statsString, err)
				log.Print(errStr)
				success = false
			}
		} else {
			success = false
		}
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
		"stats":   stats,
	})
}

// StatsNewHandler - retrieve a list of stats using stats table
func StatsNewHandler(c *gin.Context) {
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
	stats, err := database.StatsNew(statsRequest, getPartnerID(c))
	if err != nil {
		log.Printf("error retrieving stats list, err: %+v\n", err)
		errStr = fmt.Sprintf("error retrieving stats list, error: %+v", err)
		success = false
	}
	// Need record stats to file so can retrieve it later.
	// This was originally done in the log.
	// This is done due to network failure when retrieving stats (I was unable to solve this).
	// Need to marshall to JSON before putting in log
	m, err := json.Marshal(stats)
	if err != nil {
		log.Printf("Error marshalling stats: %+v to a string, err: %v\n", stats, err)
	} else {
		// log.Printf("%s%s\n", statsLogPrefix, m)
		createStatsFile(string(m))
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
		"stats":   stats,
	})
}

// // CheckNothingSentHandler - check no messages sent for a company for a specific period
// func CheckNothingSentHandler(c *gin.Context) {
// 	db := c.MustGet(shared.DatabaseConn).(*sql.DB)
// 	success := true
// 	var errStr string
// 	var nothingSentResults []database.NothingSentResult
// 	hoursBack := c.Query("hours_back")
// 	// log.Printf("hours_back: %s\n", hoursBack)
// 	hb, err := strconv.Atoi(hoursBack)
// 	if err != nil {
// 		log.Printf("error retrieving no message sent for companies list, error with parameter hours back: %s , err: %+v\n", hoursBack, err)
// 	}
// 	nothingSentResults, err = database.CheckNothingSent(hb, getPartnerID(c))
// 	if err != nil {
// 		log.Printf("error retrieving no message sent for companies list, err: %+v\n", err)
// 		errStr = fmt.Sprintf("error retrieving no message sent for companies list, error: %+v", err)
// 		success = false
// 	}
// 	// log.Printf("nothingSentResults: %+v\n", nothingSentResults)
// 	c.JSON(200, gin.H{
// 		"success":            success,
// 		"err":                errStr,
// 		"nothingSentResults": nothingSentResults,
// 	})
// }

// CheckNothingSentHandler - check no messages sent for a company for a specific period
func CheckNothingSentHandler(c *gin.Context) {
	success := true
	var errStr string
	var nothingSentResults []database.NothingSentResult
	daysBack := c.Query("days_back")
	// log.Printf("days_back: %s\n", daysBack)
	daysB, err := strconv.Atoi(daysBack)
	if err != nil {
		log.Printf("error retrieving no message sent for companies list, error with parameter days back: %s , err: %+v\n", daysBack, err)
	}
	nothingSentResults, err = database.CheckNothingSent(daysB, getPartnerID(c))
	if err != nil {
		log.Printf("error retrieving no message sent for companies list, err: %+v\n", err)
		errStr = fmt.Sprintf("error retrieving no message sent for companies list, error: %+v", err)
		success = false
	}
	// log.Printf("nothingSentResults: %+v\n", nothingSentResults)
	c.JSON(200, gin.H{
		"success":            success,
		"err":                errStr,
		"nothingSentResults": nothingSentResults,
	})
}

// GoogleMyBusinessReportHandler - run Google My Business report (adhoc)
// NOTE: This needs the Google My Business code running on the same machine.
func GoogleMyBusinessReportHandler(c *gin.Context) {
	success := true
	var errStr string
	monthsBack := c.Query("months_back")
	csvReportOnly := c.Query("cvs_report_only")
	// log.Printf("monthsBack: %s, csvReportOnly: %s\n", monthsBack, csvReportOnly)
	// check value is a positive number
	monthsB, err := strconv.Atoi(monthsBack)
	if err != nil || monthsB < 0 {
		log.Printf("error running Google My Business Report, error with parameter months back: %s, must be 0 or a positive number, setting to 0, err: %+v\n", monthsBack, err)
		monthsBack = "0"
	}
	csvRO, err := strconv.ParseBool(csvReportOnly)
	if err != nil {
		log.Printf("error running Google My Business Report, error with parameter CSV report only: %s, should be a boolean, setting to false, err: %v\n", csvReportOnly, err)
		csvRO = false
	}
	// var cmd *exec.Cmd
	// if csvRO {
	// 	cmd = exec.Command("./my_business", "-reportmonthback", monthsBack, "-csvreportonly")
	// } else {
	// 	cmd = exec.Command("./my_business", "-reportmonthback", monthsBack)
	// }
	cmd := exec.Command("./my_business", "-reportmonthback", monthsBack, "-csvreportonly")
	if !csvRO {
		cmd = exec.Command("./my_business", "-reportmonthback", monthsBack)
	}
	cmd.Dir = config.Conf.GoogleMyBusinessDirectory
	// cmd.Path = reportPath
	err = cmd.Run()
	if err != nil {
		log.Printf("error running Google My Business Report, error: %+v\n", err)
		errStr = fmt.Sprintf("error running Google My Business Report, error: %+v", err)
		success = false
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
	})
}

// Get the partner id from the JWT claims
func getPartnerID(c *gin.Context) int {
	// log.Printf("claims: %s\n", jwt.ExtractClaims(c))
	claims := jwt.ExtractClaims(c)
	partnerID, err := strconv.Atoi(claims["partnerID"].(string))
	if err != nil {
		log.Printf("error retrieving partner ID from claims, err: %+v\n", err)
		return -1
	}
	return partnerID
}

// Create stats file - everytime new stats are created
func createStatsFile(stats string) {
	// create file in same directory as executable
	fp := "./"
	ex, err := os.Executable()
	if err != nil {
		log.Printf("Unable to determine executable location for creating stats file, err: %+v\n", err)
	} else {
		fp = filepath.Dir(ex) + "/"
	}
	statsFile := fp + statsFilename
	// create and open file
	f, err := os.Create(statsFile)
	if err != nil {
		log.Printf("Error creating stats file, error: %+v\n", err)
		return
	}
	defer f.Close()
	// write stats to file
	_, err = f.WriteString(stats + "\n")
	if err != nil {
		log.Printf("Error writing stats to file, error: %+v\n", err)
		return
	}
	log.Printf("Stats written successfully to file: %s\n", statsFile)
}

// UsersHandler - retrieve a list of users
func UsersHandler(c *gin.Context) {
	users, err := database.ListAllUsers(getPartnerID(c))
	if err != nil {
		log.Printf("error retrieving user list, err: %+v\n", err)
	}
	c.JSON(200, gin.H{
		"users": users,
	})
}

// GetUserHandler - retrieve a user
func GetUserHandler(c *gin.Context) {
	var userClients database.UserClients
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		log.Printf("error (user id error) retrieving user and clients for user id: %s, err: %+v\n", c.Query("id"), err)
	} else {
		userClients, err = database.GetUser(id, getPartnerID(c))
		if err != nil {
			log.Printf("error retrieving user for a user id: %s, err: %+v\n", c.Query("id"), err)
		}
	}
	c.JSON(200, gin.H{
		"user": userClients,
	})
}

// UpdateUserHandler - update a user and associated clients
func UpdateUserHandler(c *gin.Context) {
	success := true
	var errStr string
	var userClients database.UserClients
	if err := c.ShouldBind(&userClients); err != nil {
		log.Printf("err: %+v\n", err)
		errStr = fmt.Sprintf("error: %+v", err)
		success = false
	} else {
		// log.Printf("clientConfig: %+v\n", clientConfig)
		err := database.UpdateUser(userClients, getPartnerID(c))
		if err != nil {
			log.Printf("error updating user, err: %+v\n", err)
			errStr = fmt.Sprintf("error updating user, error: %+v", err)
			success = false
		}
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
	})
}

// CreateUserHandler - create a user and associated clients
func CreateUserHandler(c *gin.Context) {
	success := true
	var errStr string
	var userClients database.UserClients
	if err := c.ShouldBind(&userClients); err != nil {
		log.Printf("err: %+v\n", err)
		errStr = fmt.Sprintf("error: %+v", err)
		success = false
	} else {
		// log.Printf("userClientConfig: %+v\n", userClients)
		err := database.CreateUser(userClients, getPartnerID(c))
		if err != nil {
			log.Printf("error user, err: %+v\n", err)
			errStr = fmt.Sprintf("error creating user, error: %+v", err)
			success = false
		}
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
	})
}

// DeleteUserHandler - delete a user assuming
func DeleteUserHandler(c *gin.Context) {
	success := true
	var errStr string
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		log.Printf("error (user id error) deleting user and clients for user id: %s, err: %+v\n", c.Query("id"), err)
		errStr = fmt.Sprintf("error (user id error) deleting user and clients for user id: %s, err: %+v", c.Query("id"), err)
		success = false
	} else {
		err = database.DeleteUser(id, getPartnerID(c))
		if err != nil {
			log.Printf("error deleting user for a user id: %s, err: %+v\n", c.Query("id"), err)
			errStr = fmt.Sprintf("error deleting user for a user id: %s, err: %+v", c.Query("id"), err)
			success = false
		}
	}
	c.JSON(200, gin.H{
		"success": success,
		"err":     errStr,
	})
}
