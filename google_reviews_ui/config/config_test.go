package config

import (
	"fmt"
	"testing"
)

func TestReadProperties(t *testing.T) {
	ReadProperties()
	if Conf.Users[0].Username != "sal" {
		t.Fatal("Error reading Users[0].Username")
	}
	fmt.Println(Conf.Users)
	fmt.Println(Conf.LogServers)
	fmt.Printf("config: %+v\n", Conf)
}
