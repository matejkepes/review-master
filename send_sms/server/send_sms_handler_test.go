package server

import (
	"fmt"
	"testing"

	"send_sms/config"
	"send_sms/shared"
)

func TestSendTo(t *testing.T) {

	config := config.ReadProperties()
	fmt.Printf("%+v\n", config.Gateways)
	// fmt.Printf("%s\n", config.Gateways[0].NumberOfSims)
	shared.Initialise(len(config.Gateways))
	fmt.Printf("%d \n", SendTo(config, -1))
	fmt.Printf("%d \n", SendTo(config, -1))
	fmt.Printf("%d \n", SendTo(config, -1))
	fmt.Printf("%d \n", SendTo(config, -1))
	fmt.Printf("%d \n", SendTo(config, -1))
	fmt.Printf("%d \n", SendTo(config, -1))
	fmt.Printf("%d \n", SendTo(config, 0))
	fmt.Printf("%d \n", SendTo(config, 0))
	fmt.Printf("%d \n", SendTo(config, 0))
	fmt.Printf("%d \n", SendTo(config, 0))
	fmt.Printf("%d \n", SendTo(config, 0))
	fmt.Printf("%d \n", SendTo(config, 0))
	fmt.Printf("%d \n", SendTo(config, 1))
	fmt.Printf("%d \n", SendTo(config, 1))
	fmt.Printf("%d \n", SendTo(config, 1))
	fmt.Printf("%d \n", SendTo(config, 1))
	fmt.Printf("%d \n", SendTo(config, 1))
	fmt.Printf("%d \n", SendTo(config, 1))
}
