package server

import (
	"encoding/json"
	"google_reviews/database"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// pairingResult - represents pairing result with queue id (used for sending messages for each client)
// for convenience it will be the same as the client id (0 if failed) unless it is set to the master queue id.
type pairingResult struct {
	QueueID string `json:"queue_id"` // queue id
}

var failedPairingResponse = []byte(``)

// ReviewMasterSMSGatewayPairingHandler - Review Master SMS Gateway pairing handler
// For security require a token to protect route that should be sent and checked
func ReviewMasterSMSGatewayPairingHandler(reviewMasterSMSGatewayPairingTk string) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		pairingToken := req.Header.Get("api-token")
		if len(pairingToken) < 1 || pairingToken != reviewMasterSMSGatewayPairingTk {
			log.Printf("Error, api-token (pairing token) %s is incorrect\n", pairingToken)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(failedPairingResponse)
			return
		}
		pairingCode, ok := req.URL.Query()["pairing_code"]
		if !ok || len(pairingCode[0]) < 1 {
			log.Printf("Error, pairing code is missing\n")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(failedPairingResponse)
			return
		}
		// queueID will be set to 0 if not found
		queueID := database.QueueIDFromReviewMasterPairCode(strings.Trim(pairingCode[0], " "))
		w.Header().Set("Content-Type", "application/json")
		w.Write(resultPairing(queueID))
	}

	return http.HandlerFunc(fn)
}

// resultPairing - result for pairing
// returns the pairing result as JSON to return to caller (queue id set to 0 for not found)
func resultPairing(queueID int) []byte {
	var pr pairingResult
	pr.QueueID = strconv.Itoa(queueID)
	r, _ := json.Marshal(pr)
	return r
}
