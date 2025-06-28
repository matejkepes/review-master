package utils

import (
	"strings"

	"github.com/dongri/phonenumber"
)

// TelephoneParse - parse telephone only returning country specific mobile numbers
func TelephoneParse(telephone string, country string) string {
	// clean up
	tel := strings.TrimSpace(telephone)
	// 00 does not work with Dutch numbers (I do not know why)
	tel = strings.TrimPrefix(tel, "00")
	// a correctly formatted foreign telephone number with a + prefix will return the foreign number
	tel = strings.TrimPrefix(tel, "+")

	// problems with Dutch numbers
	if country == "NL" {
		tel = strings.TrimPrefix(tel, "31")
		tel = "0" + tel
	}

	t := phonenumber.Parse(tel, country)

	return t
}
