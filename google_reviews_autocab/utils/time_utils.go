package utils

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

// CheckTime - check time is between start and end
func CheckTime(start string, end string, timeZone string) bool {
	// check start format
	m, err := regexp.MatchString("^\\d{2}:\\d{2}$", strings.Trim(start, " "))
	if err != nil {
		log.Println(err)
		return false
	}
	if !m {
		return false
	}
	// fmt.Println(m)

	// check end format
	m, err = regexp.MatchString("^\\d{2}:\\d{2}$", strings.Trim(end, " "))
	if err != nil {
		log.Println(err)
		return false
	}
	if !m {
		return false
	}
	// fmt.Println(m)

	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Println(err)
		return false
	}

	now := time.Now().In(loc)
	// fmt.Printf("now: %v\n", now)
	nowDate := fmt.Sprintf("%d/%02d/%02d", now.Year(), now.Month(), now.Day())
	// fmt.Printf("nowDate: %s\n", nowDate)

	startTime, err := time.ParseInLocation("2006/01/02 15:04", nowDate+" "+start, loc)
	if err != nil {
		log.Println(err)
		return false
	}
	// fmt.Printf("startTime: %v\n", startTime)

	endTime, err := time.ParseInLocation("2006/01/02 15:04", nowDate+" "+end, loc)
	if err != nil {
		log.Println(err)
		return false
	}
	// fmt.Printf("endTime: %v\n", endTime)

	if now.Before(startTime) || now.After(endTime) {
		return false
	}
	return true
}

// DiffTimeRFC3339 - get difference between from and to and whether correct format
// from and to are in RFC3339 (ISO 8601) format
func DiffTimeRFC3339(from string, to string) (time.Duration, bool) {
	start, err := time.Parse(time.RFC3339, from)
	if err != nil {
		log.Println(err)
		return 0, false
	}
	// fmt.Println(start)

	end, err := time.Parse(time.RFC3339, to)
	if err != nil {
		log.Println(err)
		return 0, false
	}
	// fmt.Println(end)

	return end.Sub(start), true
}

// CheckDiffTimeRFC3339 - check difference between from and to is less than minutesDiff in minutes
// from and to are in RFC3339 (ISO 8601) format
// minutesDiff is a number as a string
func CheckDiffTimeRFC3339(from string, to string, minutesDiff string) bool {
	diff, correctFormat := DiffTimeRFC3339(from, to)
	if !correctFormat {
		return false
	}
	m, err := time.ParseDuration(minutesDiff + "m")
	if err != nil {
		log.Println(err)
		return false
	}
	return diff.Seconds() < m.Seconds()
}

// CheckWeekday - check that the current weekday is enabled
func CheckWeekday(sunday bool, monday bool, tuesday bool, wednesday bool, thursday bool, friday bool, saturday bool, timeZone string) bool {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Println(err)
		return false
	}

	now := time.Now().In(loc)
	wd := now.Weekday()
	switch wd {
	case 0:
		return sunday
	case 1:
		return monday
	case 2:
		return tuesday
	case 3:
		return wednesday
	case 4:
		return thursday
	case 5:
		return friday
	case 6:
		return saturday
	}
	return false
}

// ConvertToTimeZone - convert time to timezone
func ConvertToTimeZone(tm time.Time, timeZone string) time.Time {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Println(err)
		return tm
	}
	tm1 := tm.In(loc)

	return tm1
}
