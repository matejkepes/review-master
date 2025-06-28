package utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dongri/phonenumber"
)

func TestPhonenumberMobile(t *testing.T) {
	test := phonenumber.Parse("07123456789", "GB")
	fmt.Println(test)
	if test != "447123456789" {
		t.Error("Error formatting mobile phone number")
	}
}

func TestPhonenumberLandline(t *testing.T) {
	test := phonenumber.Parse("01132345678", "GB")
	fmt.Println(test)
	if test != "" {
		t.Error("Error landline phone numbers should return an empty string")
	}
}

func TestPhonenumberForeign(t *testing.T) {
	test := phonenumber.Parse("17752259889", "GB")
	fmt.Println(test)
	if test != "" {
		t.Error("Error foreign phone numbers should return an empty string")
	}
}

// This is the strange one that needs coding for. If the format is as expected with a + return the number rather than an empty string
func TestPhonenumberForeignWithPlus(t *testing.T) {
	expectedCountryCode := "GB"
	// test := phonenumber.Parse("+17752259889", "GB")
	test := phonenumber.Parse("+17752259889", expectedCountryCode)
	fmt.Println(test)
	// fmt.Printf("%+v\n", phonenumber.GetISO3166())
	// fmt.Printf("%+v\n", phonenumber.GetISO3166ByNumber("+17752259889", false))
	// fmt.Printf("%+v\n", phonenumber.GetISO3166ByNumber("17752259889", false))
	// fmt.Printf("%+v\n", phonenumber.GetISO3166ByNumber(test, false))
	// country := phonenumber.GetISO3166ByNumber("17752259889", false)
	country := phonenumber.GetISO3166ByNumber(test, false)
	fmt.Println(country.Alpha2)
	// if test != "" {
	// 	t.Error("Error phone numbers with a plus sign should return an empty string")
	// }
	// if country.Alpha2 != "GB" {
	if country.Alpha2 == expectedCountryCode {
		t.Error("Error phone numbers with a plus sign should return an empty string")
	}
}

// This is a quick fix by making 2 requests to the phonenumber library since the first one removes the +
func TestPhonenumberForeignWithPlus2(t *testing.T) {
	expectedCountryCode := "GB"
	// test := phonenumber.Parse("+17752259889", "GB")
	test := phonenumber.Parse("+17752259889", expectedCountryCode)
	fmt.Println(test)
	test1 := phonenumber.Parse(test, expectedCountryCode)
	fmt.Println(test1)
	if test1 != "" {
		t.Error("Error phone numbers with a plus sign should return an empty string")
	}
}

// This is a quick fix by making 2 requests to the phonenumber library since the first one removes the +
func TestPhonenumberMobileWithPlus2(t *testing.T) {
	expectedCountryCode := "GB"
	test := phonenumber.Parse("+447123456789", expectedCountryCode)
	fmt.Println(test)
	test1 := phonenumber.Parse(test, expectedCountryCode)
	fmt.Println(test1)
	if test1 != "447123456789" {
		t.Error("Error formatting mobile phone number")
	}
}

// This is a quick fix by removing any + sign at the beginning
func TestPhonenumberMobileWithPlus3(t *testing.T) {
	expectedCountryCode := "GB"
	tel := "+447123456789"
	tel1 := strings.TrimPrefix(tel, "+")
	fmt.Println(tel1)
	test := phonenumber.Parse(tel1, expectedCountryCode)
	fmt.Println(test)
	if test != "447123456789" {
		t.Error("Error formatting mobile phone number")
	}
}

func TestPhonenumberDutch(t *testing.T) {
	expectedCountryCode := "GB"
	test := phonenumber.Parse("+31615363999", expectedCountryCode)
	fmt.Println(test)
	// fmt.Printf("%+v\n", phonenumber.GetISO3166())
	country := phonenumber.GetISO3166ByNumber(test, false)
	fmt.Println(country.Alpha2)
	if country.Alpha2 == expectedCountryCode {
		t.Error("Error phone numbers with a plus sign should return an empty string")
	}

	expectedCountryCode = "NL"
	test = phonenumber.Parse("+31615363999", expectedCountryCode)
	fmt.Println(test)
	test = phonenumber.Parse("31615363999", expectedCountryCode)
	fmt.Println(test)
	test = phonenumber.Parse("0031615363999", expectedCountryCode)
	fmt.Println(test)
	test = phonenumber.Parse("0615363999", expectedCountryCode)
	fmt.Println(test)
}

func TestTelephoneParse1(t *testing.T) {
	test := TelephoneParse("07123456789", "GB")
	fmt.Println(test)
	if test != "447123456789" {
		t.Error("Error formatting mobile phone number")
	}
}

func TestTelephoneParse2(t *testing.T) {
	test := TelephoneParse("01132345678", "GB")
	fmt.Println(test)
	if test != "" {
		t.Error("Error landline phone numbers should return an empty string")
	}
}

func TestTelephoneParse3(t *testing.T) {
	test := TelephoneParse("17752259889", "GB")
	fmt.Println(test)
	if test != "" {
		t.Error("Error foreign phone numbers should return an empty string")
	}
}

func TestTelephoneParse4(t *testing.T) {
	test := TelephoneParse("+17752259889", "GB")
	fmt.Println(test)
	if test != "" {
		t.Error("Error foreign phone numbers should return an empty string")
	}
}

func TestTelephoneParse5(t *testing.T) {
	test := TelephoneParse("0017752259889", "GB")
	fmt.Println(test)
	if test != "" {
		t.Error("Error foreign phone numbers should return an empty string")
	}
}

func TestTelephoneParse6(t *testing.T) {
	// Dutch
	test := TelephoneParse("0031615363999", "NL")
	fmt.Println(test)
	if test != "31615363999" {
		t.Error("Error formatting Dutch phone numbers")
	}

	test = TelephoneParse("+31615363999", "NL")
	fmt.Println(test)
	if test != "31615363999" {
		t.Error("Error formatting Dutch phone numbers")
	}

	test = TelephoneParse("0615363999", "NL")
	fmt.Println(test)
	if test != "31615363999" {
		t.Error("Error formatting Dutch phone numbers")
	}

	test = phonenumber.Parse("00447123456789", "NL")
	fmt.Println(test)
	if test != "" {
		t.Error("Error foreign phone numbers should return an empty string")
	}
}

// Testing US and Canadian numbers
func TestTelephoneParse7(t *testing.T) {
	// US mobile in Canada
	test := TelephoneParse("04242098463", "CA")
	fmt.Println(test)
	if test != "" {
		t.Error("Error foreign phone numbers should return an empty string")
	}

	// US mobile in US
	test = TelephoneParse("04242098463", "US")
	fmt.Println(test)
	if test != "14242098463" {
		t.Error("Error formatting mobile phone number")
	}

	// Canadian mobile in Canada
	test = TelephoneParse("05142601521", "CA")
	fmt.Println(test)
	if test != "15142601521" {
		t.Error("Error formatting mobile phone number")
	}

	// Canadian mobile in US
	test = TelephoneParse("05142601521", "US")
	fmt.Println(test)
	if test != "" {
		t.Error("Error foreign phone numbers should return an empty string")
	}
}
