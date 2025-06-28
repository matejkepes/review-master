package config

import (
	"fmt"
	"testing"
)

func TestReadProperties(t *testing.T) {
	ReadProperties()
	fmt.Printf("config: %+v\n", Conf)
}
