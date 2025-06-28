package config

import (
	"fmt"
	"testing"
)

func TestReadProperties(t *testing.T) {
	ReadProperties()
	if Conf.DbName != "google_reviews" {
		t.Fatal("Error reading DbName")
	}

	fmt.Printf("%+v\n", Conf)
}
