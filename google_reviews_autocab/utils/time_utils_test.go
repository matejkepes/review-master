package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckTime(t *testing.T) {
	test := CheckTime("10:00", "14:00", "Europe/London")
	fmt.Println(test)
}

func TestDiffTimeRFC3339_1(t *testing.T) {
	dt, cf := DiffTimeRFC3339("2019-11-04T12:35:02+00:00", "2019-11-04T12:37:29+00:00")
	fmt.Println(dt, cf)
}

func TestDiffTimeRFC3339_2(t *testing.T) {
	dt, cf := DiffTimeRFC3339("2019-11-04T12:35:02+00:00", "2019-11-04T12:35:02+00:00")
	fmt.Println(dt, cf)
}

func TestCheckDiffTimeRFC3339_1(t *testing.T) {
	ct := CheckDiffTimeRFC3339("2019-11-04T12:35:02+00:00", "2019-11-04T12:37:29+00:00", "3")
	fmt.Println(ct)
	if !ct {
		t.Fatal("Error should be true")
	}
}

func TestCheckDiffTimeRFC3339_2(t *testing.T) {
	ct := CheckDiffTimeRFC3339("2019-11-04T12:35:02+00:00", "2019-11-04T12:39:29+00:00", "3")
	fmt.Println(ct)
	if ct {
		t.Fatal("Error should be false")
	}
}

func TestCheckWeekday_1(t *testing.T) {
	wd := CheckWeekday(true, true, true, true, true, true, true, "Europe/London")
	fmt.Println(wd)
	if !wd {
		t.Fatal("Error should be true")
	}
}

func TestCheckWeekday_2(t *testing.T) {
	wd := CheckWeekday(false, false, false, false, false, false, false, "Europe/London")
	fmt.Println(wd)
	if wd {
		t.Fatal("Error should be false")
	}
}

func TestConvertToTimeZone_1(t *testing.T) {
	c1 := ConvertToTimeZone(time.Now(), "UTC")
	fmt.Println(c1)
	c2 := ConvertToTimeZone(time.Now(), "Europe/London")
	fmt.Println(c2)
	c3 := ConvertToTimeZone(time.Now(), "America/Caracas")
	fmt.Println(c3)
}
