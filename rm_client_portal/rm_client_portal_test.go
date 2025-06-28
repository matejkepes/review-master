package main

import (
	"fmt"
	"rm_client_portal/google_my_business_api"
	"testing"
)

func TestInit(t *testing.T) {
	fmt.Println("Testing")
	Init()
}

// use this to get a new token
// if the credentials.json change then need to delete the old token.json and run this to create a new token.json
func TestGetClient(t *testing.T) {
	config := GetGoogleCredentials()
	google_my_business_api.SetClient(config)
	fmt.Println(google_my_business_api.Client)
}
