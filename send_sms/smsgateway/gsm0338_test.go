package smsgateway

import (
	"fmt"
	"log"
	"testing"
)

func TestUTF8ToGsm0338(t *testing.T) {
	s := "Hello World"
	gsm := UTF8ToGsm0338(s)
	log.Printf("gsm: %x\n", gsm)
	h := fmt.Sprintf("%x", gsm)
	log.Printf("h: %s\n", h)
	if h != "48656c6c6f20576f726c64" {
		t.Fatal("Error encoding: " + s)
	}

	s = "£5"
	gsm = UTF8ToGsm0338(s)
	log.Printf("gsm: %x\n", gsm)
	h = fmt.Sprintf("%x", gsm)
	log.Printf("h: %s\n", h)
	if h != "0135" {
		t.Fatal("Error encoding: " + s)
	}
}
