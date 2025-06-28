package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"send_sms_ui/config"
)

func TestSendSmsHandler1(t *testing.T) {
	// read config file
	config := config.ReadProperties()
	router := SetupRouter(config, "../templates/**/*")

	w := httptest.NewRecorder()
	data := url.Values{}
	data.Set("tels", "07123456789\n07987654321")
	data.Set("msg", "testing")
	req, _ := http.NewRequest("POST", "/sendsms", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.Username, config.UserPassword)
	router.ServeHTTP(w, req)
	// fmt.Printf("code: %d\n", w.Code)
	if w.Code != 200 {
		t.Fatalf("received incorrect HTTP code, received: %d\n", w.Code)
	}
	fmt.Printf("response: %v\n", w)
	if !strings.Contains(w.Body.String(), "447123456789") {
		t.Fatal("does not contain telephone")
	}
	if !strings.Contains(w.Body.String(), "447987654321") {
		t.Fatal("does not contain telephone")
	}
}

func TestSendSmsHandler2(t *testing.T) {
	// read config file
	config := config.ReadProperties()
	router := SetupRouter(config, "../templates/**/*")

	w := httptest.NewRecorder()
	data := url.Values{}
	data.Set("tels", "07123456789\n07123456789")
	data.Set("msg", "testing")
	req, _ := http.NewRequest("POST", "/sendsms", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.Username, config.UserPassword)
	router.ServeHTTP(w, req)
	// fmt.Printf("code: %d\n", w.Code)
	if w.Code != 200 {
		t.Fatalf("received incorrect HTTP code, received: %d\n", w.Code)
	}
	fmt.Printf("response: %v\n", w)
	if !strings.Contains(w.Body.String(), "447123456789") {
		t.Fatal("does not contain telephone")
	}
}

// NOTE: set config.properties value maxTelephones=5 or less than 6 to check the number sent
func TestSendSmsHandler3(t *testing.T) {
	// read config file
	config := config.ReadProperties()
	router := SetupRouter(config, "../templates/**/*")

	w := httptest.NewRecorder()
	data := url.Values{}
	data.Set("tels", "07987654321\n07987654322\n07987654323\n07987654324\n07987654325\n07987654326\n")
	data.Set("msg", "testing")
	req, _ := http.NewRequest("POST", "/sendsms", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.Username, config.UserPassword)
	router.ServeHTTP(w, req)
	// fmt.Printf("code: %d\n", w.Code)
	if w.Code != 200 {
		t.Fatalf("received incorrect HTTP code, received: %d\n", w.Code)
	}
	fmt.Printf("response: %v\n", w)
	if !strings.Contains(w.Body.String(), "447987654321") {
		t.Fatal("does not contain telephone")
	}
	if strings.Contains(w.Body.String(), "447987654326") {
		t.Fatal("should not contain telephone")
	}
}

// no correctly configured telephone numbers
func TestSendSmsHandler4(t *testing.T) {
	// read config file
	config := config.ReadProperties()
	router := SetupRouter(config, "../templates/**/*")

	w := httptest.NewRecorder()
	data := url.Values{}
	data.Set("tels", "0712345678\n0798765432")
	data.Set("msg", "testing")
	req, _ := http.NewRequest("POST", "/sendsms", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.Username, config.UserPassword)
	router.ServeHTTP(w, req)
	// fmt.Printf("code: %d\n", w.Code)
	if w.Code != 200 {
		t.Fatalf("received incorrect HTTP code, received: %d\n", w.Code)
	}
	fmt.Printf("response: %v\n", w)
	if !strings.Contains(w.Body.String(), "No correctly configured telephone numbers") {
		t.Fatal("should not contain any correctly configured telephone numbers")
	}
}

// no message to send
func TestSendSmsHandler5(t *testing.T) {
	// read config file
	config := config.ReadProperties()
	router := SetupRouter(config, "../templates/**/*")

	w := httptest.NewRecorder()
	data := url.Values{}
	data.Set("tels", "07123456789\n07987654321")
	data.Set("msg", "   ")
	req, _ := http.NewRequest("POST", "/sendsms", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.Username, config.UserPassword)
	router.ServeHTTP(w, req)
	// fmt.Printf("code: %d\n", w.Code)
	if w.Code != 200 {
		t.Fatalf("received incorrect HTTP code, received: %d\n", w.Code)
	}
	fmt.Printf("response: %v\n", w)
	if !strings.Contains(w.Body.String(), "No message to send") {
		t.Fatal("should not send because no message")
	}
}

// are not mobile numbers
func TestSendSmsHandler6(t *testing.T) {
	// read config file
	config := config.ReadProperties()
	router := SetupRouter(config, "../templates/**/*")

	w := httptest.NewRecorder()
	data := url.Values{}
	data.Set("tels", "01123456789\n01987654321")
	data.Set("msg", "testing")
	req, _ := http.NewRequest("POST", "/sendsms", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.Username, config.UserPassword)
	router.ServeHTTP(w, req)
	// fmt.Printf("code: %d\n", w.Code)
	if w.Code != 200 {
		t.Fatalf("received incorrect HTTP code, received: %d\n", w.Code)
	}
	fmt.Printf("response: %v\n", w)
	if strings.Contains(w.Body.String(), "441123456789") {
		t.Fatal("should not contain telephone")
	}
	if strings.Contains(w.Body.String(), "441987654321") {
		t.Fatal("should not contain telephone")
	}
	if !strings.Contains(w.Body.String(), "No correctly configured telephone numbers to send") {
		t.Fatal("should not contain telephone")
	}
}
