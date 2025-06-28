package config

import (
	"fmt"
	"testing"
)

func TestReadProperties(t *testing.T) {
	config := ReadProperties()
	// fmt.Print(config.GatewayAddress)
	// if config.GatewayAddress != "localhost" {
	// 	t.Fatal("Error reading GatewayAddress, got: " + config.GatewayAddress)
	// }
	if config.Gateways[0].NumberOfSims != "8" {
		t.Fatal("Error reading Gateways[0].NumberOfSims")
	}
	fmt.Println(config.Gateways)
	fmt.Println(config.Tokens)

	fmt.Println(config.EmailTo)
	fmt.Printf("%+v\n", config)
}
