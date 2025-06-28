package barred

import (
	"log"
	"testing"
)

func TestReadBarredFile(t *testing.T) {
	bars, err := ReadBarredFile("../config/barred_telephone_prefixes.txt")
	if err != nil {
		log.Printf("err: %s", err)
		t.Fatal("Error reading barred")
	}
	log.Printf("barred: %s\n", bars)
}

func TestCheckBarred1(t *testing.T) {
	bars, _ := ReadBarredFile("../config/barred_telephone_prefixes.txt")
	b := CheckBarred("447123456789", bars)
	if b {
		t.Fatal("Error telephone should not be barred")
	}
}

func TestCheckBarred2(t *testing.T) {
	bars, _ := ReadBarredFile("../config/barred_telephone_prefixes.txt")
	// NOTE: telephone prefix must be in the file list of barred telephone numbers
	b := CheckBarred("447418456789", bars)
	if !b {
		t.Fatal("Error telephone should be barred")
	}
}
