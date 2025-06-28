package config

import (
	"fmt"
	"testing"
)

func TestReadProperties(t *testing.T) {
	config := ReadProperties()
	fmt.Printf("%+v\n", config)
}
