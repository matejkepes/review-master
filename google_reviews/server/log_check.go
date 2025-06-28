package server

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ErrorResult - represents an error result.
type ErrorResult struct {
	ClientID  int       `json:"client_id"`  // client id
	Frequency int       `json:"frequency"`  // frequency
	LastError time.Time `json:"last_error"` // last error
}

var failedLogResponse = []byte(``)

// CheckLogHandler - Check log file handler
// For security require a token to protect route that should be sent and checked
func CheckLogHandler(logFileName, logTk string) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		logToken, ok := req.URL.Query()["log_token"]
		if !ok || len(logToken[0]) < 1 || strings.Trim(logToken[0], " ") != logTk {
			log.Printf("Error, log_token %s is incorrect\n", logToken[0])
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(failedLogResponse)
			return
		}
		hb := 24
		hoursBack, ok := req.URL.Query()["hours_back"]
		if ok && len(hoursBack[0]) > 0 {
			h, err := strconv.Atoi(hoursBack[0])
			if err == nil {
				hb = h
			}
		}
		// fmt.Printf("hb = %d\n", hb)
		// w.Write([]byte(CheckLog(logFileName, hb)))
		cl := CheckLog(logFileName, hb)
		c, _ := json.Marshal(cl)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(c))
	}

	return http.HandlerFunc(fn)
}

// CheckLog - check log file
func CheckLog(logFileName string, hoursBack int) []ErrorResult {
	var errorResult []ErrorResult
	type errFreqValue struct {
		Freq     int
		LastTime time.Time
	}
	var sendErrFreq = make(map[string]errFreqValue)
	// fmt.Fprintln(w, "Logs")
	// fmt.Fprintf(w, "logFileName: %s\n", logFileName)
	f, err := os.Open(logFileName) // For read access.
	if err != nil {
		log.Printf("Error, opening log file: %s to view logs with error: %v\n", logFileName, err)
		return errorResult
	}
	defer f.Close()

	// check time after
	checkAfter := time.Now().Add(-time.Duration(hoursBack) * time.Hour)
	scanner := bufio.NewScanner(f)
	re, _ := regexp.Compile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} Error sending message for clientID: .*`)
	for scanner.Scan() {
		txt := scanner.Text()
		// log format starts with time
		// get required message relating to failure to send message
		// m, err := regexp.MatchString("^\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} Error sending message for clientID: .*", txt)
		// if err != nil {
		// 	// log.Println(err)
		// 	continue
		// }
		m := re.MatchString(txt)
		if m {
			// check date
			ta := strings.Split(txt, " ")
			t, err := time.Parse("2006/01/02 15:04:05", ta[0]+" "+ta[1])
			if err != nil {
				continue
			}
			// fmt.Println(t)
			if t.After(checkAfter) {
				// fmt.Println(ta[7])
				c := strings.Split(txt, "clientID:")
				// fmt.Println(c[1])
				ci := strings.Split(strings.Trim(c[1], " "), " ")
				// fmt.Println(ci[0])
				clientID := ci[0]
				// v, _ := sendErrFreq[clientID]
				v := sendErrFreq[clientID]
				sendErrFreq[clientID] = errFreqValue{Freq: v.Freq + 1, LastTime: t}
			}
		}
	}

	// for cid, v := range sendErrFreq {
	// 	l := fmt.Sprintf("{\"client_id\":\"%s\",\"frequency\":\"%d\",\"last_error\":\"%s\"},", cid, v.Freq, v.LastTime.Format("2006/01/02 15:04:05"))
	// 	returnTxt = returnTxt + l
	// }
	// if len(returnTxt) > 0 {
	// 	returnTxt = "[" + strings.TrimSuffix(returnTxt, ",") + "]"
	// }
	for cid, v := range sendErrFreq {
		var e ErrorResult
		i, err := strconv.Atoi(cid)
		if err != nil {
			continue
		}
		e.ClientID = i
		e.Frequency = v.Freq
		e.LastError = v.LastTime
		errorResult = append(errorResult, e)
	}

	if err = scanner.Err(); err != nil {
		log.Printf("Error, reading log file: %s to check logs with error: %v\n", logFileName, err)
	}
	return errorResult
}
