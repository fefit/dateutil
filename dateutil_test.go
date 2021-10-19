package dateutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func makeTestTime() time.Time {
	return time.Date(2021, time.September, 9, 18, 7, 6, 12345678, time.UTC)
}

type EnumDate = int32

const (
	YEAR EnumDate = 1 << iota
	MONTH
	DAY
	HOUR
	MINUTE
	SECOND
	MILLISECOND
	MICROSECOND
	YMD = YEAR | MONTH | DAY
)

func isSameDate(date *time.Time, mode EnumDate) bool {
	testDate := makeTestTime()
	flag := true
	// test if year is equal
	if flag && (mode&YEAR > 0) {
		flag = date.Year() == testDate.Year()
	}
	// test if month is equal
	if flag && (mode&MONTH > 0) {
		flag = date.Month() == testDate.Month()
	}
	// test if day is equal
	if flag && (mode&DAY > 0) {
		flag = date.Day() == testDate.Day()
	}
	return flag
}
func TestDateFormat(t *testing.T) {
	curTime := makeTestTime()
	/*
	* Test the 'Year'
	 */
	// Full year
	if Y, err := DateFormat(curTime, "Y"); err == nil {
		assert.Equal(t, Y, "2021")
	} else {
		assert.Fail(t, "Format full year 'Y' is not ok.")
	}
	// Short year
	if y, err := DateFormat(curTime, "y"); err == nil {
		assert.Equal(t, y, "21")
	} else {
		assert.Fail(t, "Format short year 'y' is not ok.")
	}
	// Leap year
	if L, err := DateFormat(curTime, "L"); err == nil {
		assert.Equal(t, L, "0")
	} else {
		assert.Fail(t, "Format leap year 'L' is not ok.")
	}
	/*
	* Test the 'Month'
	 */
	// A full textual representation of a month
	if F, err := DateFormat(curTime, "F"); err == nil {
		assert.Equal(t, F, "September")
	} else {
		assert.Fail(t, "Format month full name 'F' is not ok.")
	}
	// A short textual representation of a month
	if M, err := DateFormat(curTime, "M"); err == nil {
		assert.Equal(t, M, "Sep")
	} else {
		assert.Fail(t, "Format month short name 'M' is not ok.")
	}
	// Numeric representation of a month, with leading zeros
	if m, err := DateFormat(curTime, "m"); err == nil {
		assert.Equal(t, m, "09")
	} else {
		assert.Fail(t, "Format month numberic with leading zeros 'm' is not ok.")
	}
	// Numeric representation of a month, without leading zeros
	if n, err := DateFormat(curTime, "n"); err == nil {
		assert.Equal(t, n, "9")
	} else {
		assert.Fail(t, "Format month numberic without leading zeros 'm' is not ok.")
	}
	// Number of days in the given month
	if _t, err := DateFormat(curTime, "t"); err == nil {
		assert.Equal(t, _t, "30")
	} else {
		assert.Fail(t, "Format number of days in the given month 't' is not ok.")
	}
	/*
	* Test the 'Day'
	 */
	// Day of the month, 2 digits with leading zeros
	if d, err := DateFormat(curTime, "d"); err == nil {
		assert.Equal(t, d, "09")
	} else {
		assert.Fail(t, "Format day of the month 'd' is not ok.")
	}
	// A textual representation of a day, three letters
	if D, err := DateFormat(curTime, "D"); err == nil {
		assert.Equal(t, D, "Thu")
	} else {
		assert.Fail(t, "Format a textual representation of a day 'D' is not ok.")
	}
	// Day of the month without leading zeros
	if j, err := DateFormat(curTime, "j"); err == nil {
		assert.Equal(t, j, "9")
	} else {
		assert.Fail(t, "Format day of the month without leading zeros 'j' is not ok.")
	}
	// A full textual representation of the day of the week
	if l, err := DateFormat(curTime, "l"); err == nil {
		assert.Equal(t, l, "Thursday")
	} else {
		assert.Fail(t, "Format a full textual representation of the day of the week 'l' is not ok.")
	}
	// ISO-8601 numeric representation of the day of the week
	if N, err := DateFormat(curTime, "N"); err == nil {
		assert.Equal(t, N, "4")
	} else {
		assert.Fail(t, "Format ISO-8601 numeric representation of the day of the week 'N' is not ok.")
	}
	// Numeric representation of the day of the week
	if w, err := DateFormat(curTime, "w"); err == nil {
		assert.Equal(t, w, "5")
	} else {
		assert.Fail(t, "Format numeric representation of the day of the week 'w' is not ok.")
	}
	// The day of the year
	if z, err := DateFormat(curTime, "z"); err == nil {
		assert.Equal(t, z, "251")
	} else {
		assert.Fail(t, "Format numeric representation of the day of the week 'z' is not ok.")
	}
	/*
	* Test the 'Week'
	 */
	// ISO-8601 week number of year, weeks starting on Monday
	if W, err := DateFormat(curTime, "W"); err == nil {
		assert.Equal(t, W, "36")
	} else {
		assert.Fail(t, "Format ISO-8601 week number of year 'W' is not ok.")
	}
	/*
	* Test the 'Time'
	 */
	// Lowercase Ante meridiem and Post meridiem
	if a, err := DateFormat(curTime, "a"); err == nil {
		assert.Equal(t, a, "pm")
	} else {
		assert.Fail(t, "Format lowercase Ante meridiem and Post meridiem 'a' is not ok.")
	}
	// Uppercase Ante meridiem and Post meridiem
	if A, err := DateFormat(curTime, "A"); err == nil {
		assert.Equal(t, A, "PM")
	} else {
		assert.Fail(t, "Format uppercase Ante meridiem and Post meridiem 'A' is not ok.")
	}
	// 12-hour format of an hour without leading zeros
	if g, err := DateFormat(curTime, "g"); err == nil {
		assert.Equal(t, g, "6")
	} else {
		assert.Fail(t, "Format 12-hour format of an hour without leading zeros 'g' is not ok.")
	}
	// 24-hour format of an hour without leading zeros
	if G, err := DateFormat(curTime, "G"); err == nil {
		assert.Equal(t, G, "18")
	} else {
		assert.Fail(t, "Format 24-hour format of an hour without leading zeros 'G' is not ok.")
	}
	// 12-hour format of an hour with leading zeros
	if h, err := DateFormat(curTime, "h"); err == nil {
		assert.Equal(t, h, "06")
	} else {
		assert.Fail(t, "Format 12-hour format of an hour with leading zeros 'h' is not ok.")
	}
	// 24-hour format of an hour with leading zeros
	if H, err := DateFormat(curTime, "H"); err == nil {
		assert.Equal(t, H, "18")
	} else {
		assert.Fail(t, "Format 24-hour format of an hour with leading zeros 'H' is not ok.")
	}
	// Minutes with leading zeros
	if i, err := DateFormat(curTime, "i"); err == nil {
		assert.Equal(t, i, "07")
	} else {
		assert.Fail(t, "Format minutes with leading zeros 'i' is not ok.")
	}
	// Seconds with leading zeros
	if s, err := DateFormat(curTime, "s"); err == nil {
		assert.Equal(t, s, "06")
	} else {
		assert.Fail(t, "Format seconds with leading zeros 's' is not ok.")
	}
	// Microseconds
	if u, err := DateFormat(curTime, "u"); err == nil {
		assert.Equal(t, u, "012345")
	} else {
		assert.Fail(t, "Format microseconds 'u' is not ok.")
	}
	// Milliseconds
	if v, err := DateFormat(curTime, "v"); err == nil {
		assert.Equal(t, v, "012")
	} else {
		assert.Fail(t, "Format milliseconds 'v' is not ok.")
	}
}

func TestStrToTime(t *testing.T) {
	// 2021-09-09
	if date, err := DateTime("2021-09-09"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021-09-09 fail")
	}
	// 2021-9-9
	if date, err := DateTime("2021-9-9"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021-9-9 fail")
	}
	// 2021-09-9
	if date, err := DateTime("2021-09-9"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021-09-9 fail")
	}
	// 2021-9-09
	if date, err := DateTime("2021-9-09"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021-9-09 fail")
	}
	// 21-9-09
	if date, err := DateTime("21-9-09"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 21-9-09 fail")
	}
	// 21-9-9
	if date, err := DateTime("21-9-9"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 21-9-9 fail")
	}
	// 21-09-9
	if date, err := DateTime("21-09-9"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 21-09-9 fail")
	}
	// Sep-09-21
	if date, err := DateTime("Sep-09-21"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime Sep-09-21 fail")
	}
	// Sep-09-2021
	if date, err := DateTime("Sep-09-2021"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime Sep-09-2021 fail")
	}
	// Sep-09-021
	if date, err := DateTime("Sep-09-021"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime Sep-09-021 fail")
	}
	// 1-09-9
	if date, err := DateTime("1-09-9"); err == nil {
		assert.True(t, isSameDate(&date, MONTH|DAY))
	} else {
		assert.Fail(t, "StrToTime 1-09-9 fail")
	}
	// 20210909
	if date, err := DateTime("20210909"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 20210909 fail")
	}
	// 202199, wrong date
	if _, err := DateTime("202199"); err == nil {
		assert.Fail(t, "StrToTime wrong date 202199 format ok")
	}
	// 2021/09/09
	if date, err := DateTime("2021/09/09"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021/09/09 fail")
	}
	// 2021/9/9
	if date, err := DateTime("2021/9/9"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021/9/9 fail")
	}
	// 2021/09/9
	if date, err := DateTime("2021/09/9"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021/09/9 fail")
	}
	// 2021/9/09
	if date, err := DateTime("2021/9/09"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021/9/09 fail")
	}
	// 9/9
	if date, err := DateTime("9/9"); err == nil {
		assert.True(t, isSameDate(&date, MONTH|DAY))
	} else {
		assert.Fail(t, "StrToTime 9/9 fail")
	}
	// four number, take as time first
	if date, err := DateTime("2021"); err == nil {
		assert.Equal(t, date.Hour(), 20)
		assert.Equal(t, date.Minute(), 21)
	} else {
		assert.Fail(t, "StrToTime 2021 fail")
	}
	// four number, take as time first, but not a correct time
	if date, err := DateTime("2061"); err == nil {
		assert.Equal(t, date.Year(), 2061)
	} else {
		assert.Fail(t, "StrToTime 2021 fail")
	}
	// golang datetime
	// if date, err := DateTime("2021-09-09 18:07:06 +0000 UTC"); err == nil {
	// 	assert.True(t, isSameDate(&date, YMD))
	// } else {
	// 	assert.Fail(t, "StrToTime 2021-09-09 18:07:06 +0000 UTC fail")
	// }
	// seconds with fraction
	if date, err := DateTime("2021-09-09 06:07:06.123456789PM"); err == nil {
		assert.True(t, isSameDate(&date, YMD|HOUR|MINUTE|SECOND))
	} else {
		assert.Fail(t, "StrToTime 2021 fail")
	}
}
