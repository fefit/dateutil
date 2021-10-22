package dateutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
	YMD    = YEAR | MONTH | DAY
	His    = HOUR | MINUTE | SECOND
	YMDHis = YMD | His
)

var (
	localLocation, _ = time.LoadLocation("Asia/Shanghai")
)

func makeTestTime() time.Time {
	return time.Date(2021, time.September, 5, 18, 7, 6, 12345678, time.UTC)
}

func isSameDate(date *time.Time, mode EnumDate, args ...bool) bool {
	var testDate time.Time
	if len(args) == 1 && args[0] == true {
		testDate = time.Now()
	} else {
		testDate = makeTestTime()
	}
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
	// test if Hour is equal
	if flag && (mode&HOUR > 0) {
		flag = date.Hour() == testDate.Hour()
	}
	// test if Minute is equal
	if flag && (mode&MINUTE > 0) {
		flag = date.Minute() == testDate.Minute()
	}
	// test if Second is equal
	if flag && (mode&SECOND > 0) {
		flag = date.Second() == testDate.Second()
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
		assert.Equal(t, d, "05")
	} else {
		assert.Fail(t, "Format day of the month 'd' is not ok.")
	}
	// A textual representation of a day, three letters
	if D, err := DateFormat(curTime, "D"); err == nil {
		assert.Equal(t, D, "Sun")
	} else {
		assert.Fail(t, "Format a textual representation of a day 'D' is not ok.")
	}
	// Day of the month without leading zeros
	if j, err := DateFormat(curTime, "j"); err == nil {
		assert.Equal(t, j, "5")
	} else {
		assert.Fail(t, "Format day of the month without leading zeros 'j' is not ok.")
	}
	// A full textual representation of the day of the week
	if l, err := DateFormat(curTime, "l"); err == nil {
		assert.Equal(t, l, "Sunday")
	} else {
		assert.Fail(t, "Format a full textual representation of the day of the week 'l' is not ok.")
	}
	// ISO-8601 numeric representation of the day of the week
	if N, err := DateFormat(curTime, "N"); err == nil {
		assert.Equal(t, N, "0")
	} else {
		assert.Fail(t, "Format ISO-8601 numeric representation of the day of the week 'N' is not ok.")
	}
	// Numeric representation of the day of the week
	if w, err := DateFormat(curTime, "w"); err == nil {
		assert.Equal(t, w, "1")
	} else {
		assert.Fail(t, "Format numeric representation of the day of the week 'w' is not ok.")
	}
	// The day of the year
	if z, err := DateFormat(curTime, "z"); err == nil {
		assert.Equal(t, z, "247")
	} else {
		assert.Fail(t, "Format numeric representation of the day of the week 'z' is not ok.")
	}
	/*
	* Test the 'Week'
	 */
	// ISO-8601 week number of year, weeks starting on Monday
	if W, err := DateFormat(curTime, "W"); err == nil {
		assert.Equal(t, W, "35")
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

func TestNumberDate(t *testing.T) {
	/**
	 * Test timestamp date
	 */
	// 2021-09-05 18:07:06 since 1970
	var timestamp int = 1630836426
	if date, err := DateTime(timestamp); err == nil {
		assert.True(t, isSameDate(&date, YMDHis))
	} else {
		assert.Fail(t, "DateTime int fail")
	}
	if date, err := DateTime(int64(timestamp)); err == nil {
		assert.True(t, isSameDate(&date, YMDHis))
	} else {
		assert.Fail(t, "DateTime int64 fail")
	}
	if date, err := DateTime(int32(timestamp)); err == nil {
		assert.True(t, isSameDate(&date, YMDHis))
	} else {
		assert.Fail(t, "DateTime int32 fail")
	}
	if date, err := DateTime(float64(timestamp)); err == nil {
		assert.True(t, isSameDate(&date, YMDHis))
	} else {
		assert.Fail(t, "DateTime float64 fail")
	}
}

func TestStringDate(t *testing.T) {
	/**
	 * Test dates
	 */
	// 2021-09-05
	if date, err := DateTime("2021-09-05"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021-09-05 fail")
	}
	// 2021-9-5
	if date, err := DateTime("2021-9-5"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021-9-5 fail")
	}
	// 2021-09-5
	if date, err := DateTime("2021-09-5"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021-09-5 fail")
	}
	// 2021-9-05
	if date, err := DateTime("2021-9-05"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021-9-05 fail")
	}
	// 21-9-05
	if date, err := DateTime("21-9-05"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 21-9-05 fail")
	}
	// 21-9-5
	if date, err := DateTime("21-9-5"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 21-9-5 fail")
	}
	// 21-09-5
	if date, err := DateTime("21-09-5"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 21-09-5 fail")
	}
	// Sep-05-21
	if date, err := DateTime("Sep-05-21"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime Sep-05-21 fail")
	}
	// Sep-05-2021
	if date, err := DateTime("Sep-05-2021"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime Sep-05-2021 fail")
	}
	// Sep-05-021
	if date, err := DateTime("Sep-05-021"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime Sep-05-021 fail")
	}
	// 1-09-5
	if date, err := DateTime("1-09-5"); err == nil {
		assert.True(t, isSameDate(&date, MONTH|DAY))
	} else {
		assert.Fail(t, "StrToTime 1-09-5 fail")
	}
	// 20210905
	if date, err := DateTime("20210905"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 20210905 fail")
	}
	// 202195, wrong date
	if _, err := DateTime("202195"); err == nil {
		assert.Fail(t, "StrToTime wrong date 202195 format ok")
	}
	// 2021/09/05
	if date, err := DateTime("2021/09/05"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021/09/05 fail")
	}
	// 2021/9/5
	if date, err := DateTime("2021/9/5"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021/9/5 fail")
	}
	// 2021/09/5
	if date, err := DateTime("2021/09/5"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021/09/5 fail")
	}
	// 2021/9/05
	if date, err := DateTime("2021/9/05"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021/9/05 fail")
	}
	// 9/5
	if date, err := DateTime("9/5"); err == nil {
		assert.True(t, isSameDate(&date, MONTH|DAY))
	} else {
		assert.Fail(t, "StrToTime 9/5 fail")
	}
	// Sep-05-21
	if date, err := DateTime("Sep-05-21"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime Sep-05-21 fail")
	}
	// Sep-05-2021
	if date, err := DateTime("Sep-05-2021"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime Sep-05-2021 fail")
	}
	// 2021-Sep-05
	if date, err := DateTime("2021-Sep-05"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 2021-Sep-05 fail")
	}
	// 21-Sep-05
	if date, err := DateTime("21-Sep-05"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 21-Sep-05 fail")
	}
	// 9/5/21
	if date, err := DateTime("9/5/21"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 9/5/21 fail")
	}
	// 9/5/2021
	if date, err := DateTime("9/5/2021"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 9/5/2021 fail")
	}
	// 5-9-2021
	if date, err := DateTime("5-9-2021"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 5-9-2021 fail")
	}
	// 5.9.2021
	if date, err := DateTime("5.9.2021"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 5.9.2021 fail")
	}
	// 5.9.21 :time first
	if date, err := DateTime("5.9.21"); err == nil {
		assert.Equal(t, date.Hour(), 5)
		assert.Equal(t, date.Minute(), 9)
		assert.Equal(t, date.Second(), 21)
	} else {
		assert.Fail(t, "StrToTime 5.9.21 fail")
	}
	// 5.9.2021T
	if _, err := DateTime("5.9.2021T"); err == nil {
		assert.Fail(t, "StrToTime 5.9.2021T fail")
	}
	// 20210905
	if date, err := DateTime("20210905"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 20210905 fail")
	}
	// 2021-9
	if date, err := DateTime("2021-9"); err == nil {
		assert.True(t, isSameDate(&date, YEAR|MONTH))
	} else {
		assert.Fail(t, "StrToTime 2021-9 fail")
	}
	// 2021-09
	if date, err := DateTime("2021-09"); err == nil {
		assert.True(t, isSameDate(&date, YEAR|MONTH))
	} else {
		assert.Fail(t, "StrToTime 2021-09 fail")
	}
	// 05-September 2021
	if date, err := DateTime("05-September 2021"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 05-September 2021 fail")
	}
	// 05Sep2021
	if date, err := DateTime("05Sep2021"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 05Sep2021 fail")
	}
	// 05 ix 2021
	if date, err := DateTime("05 ix 2021"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime 05 ix 2021 fail")
	}
	// September 2021
	if date, err := DateTime("September 2021"); err == nil {
		assert.True(t, isSameDate(&date, YEAR|MONTH))
	} else {
		assert.Fail(t, "StrToTime September 2021 fail")
	}
	// Sep2021
	if date, err := DateTime("Sep2021"); err == nil {
		assert.True(t, isSameDate(&date, YEAR|MONTH))
	} else {
		assert.Fail(t, "StrToTime Sep2021 fail")
	}
	// 2021 September
	if date, err := DateTime("2021 September"); err == nil {
		assert.True(t, isSameDate(&date, YEAR|MONTH))
	} else {
		assert.Fail(t, "StrToTime 2021 September fail")
	}
	// 2021.September
	if date, err := DateTime("2021.September"); err == nil {
		assert.True(t, isSameDate(&date, YEAR|MONTH))
	} else {
		assert.Fail(t, "StrToTime 2021.September fail")
	}
	// 2021-ix
	if date, err := DateTime("2021-ix"); err == nil {
		assert.True(t, isSameDate(&date, YEAR|MONTH))
	} else {
		assert.Fail(t, "StrToTime 2021-ix fail")
	}
	// September 5th, 2021
	if date, err := DateTime("September 5th, 2021"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime September 5th, 2021 fail")
	}
	// September 5, 2021
	if date, err := DateTime("September 5, 2021"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime September 5, 2021 fail")
	}
	// September.5.21
	if date, err := DateTime("September.5.21"); err == nil {
		assert.True(t, isSameDate(&date, YMD))
	} else {
		assert.Fail(t, "StrToTime September.5.21 fail")
	}
	// September 5th
	if date, err := DateTime("September 5th"); err == nil {
		assert.True(t, isSameDate(&date, MONTH|DAY))
	} else {
		assert.Fail(t, "StrToTime September 5th fail")
	}
	// Sep 5
	if date, err := DateTime("Sep 5"); err == nil {
		assert.True(t, isSameDate(&date, MONTH|DAY))
	} else {
		assert.Fail(t, "StrToTime Sep 5 fail")
	}
	// Sep.5
	if date, err := DateTime("Sep.5"); err == nil {
		assert.True(t, isSameDate(&date, MONTH|DAY))
	} else {
		assert.Fail(t, "StrToTime Sep.5 fail")
	}
	// 5 Sep
	if date, err := DateTime("5 Sep"); err == nil {
		assert.True(t, isSameDate(&date, MONTH|DAY))
	} else {
		assert.Fail(t, "StrToTime 5 Sep fail")
	}
	// 5.Sep
	if date, err := DateTime("5.Sep"); err == nil {
		assert.True(t, isSameDate(&date, MONTH|DAY))
	} else {
		assert.Fail(t, "StrToTime 5.Sep fail")
	}
	// Sep
	if date, err := DateTime("Sep"); err == nil {
		assert.True(t, isSameDate(&date, MONTH))
	} else {
		assert.Fail(t, "StrToTime Sep fail")
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
}

func TestLayoutDateTime(t *testing.T) {
	/**
	 * Test golang layouts
	 */
	// golang datetime
	if date, err := DateTime("2021-09-05 18:07:06 +0000 UTC"); err == nil {
		assert.True(t, isSameDate(&date, YEAR|MONTH))
	} else {
		assert.Fail(t, "StrToTime 2021-09-05 18:07:06 +0000 UTC fail")
	}
	// golang datetime
	if date, err := DateTime("2021-09-05 18:07:06.123456789 +0000 UTC"); err == nil {
		assert.True(t, isSameDate(&date, YEAR|MONTH))
	} else {
		assert.Fail(t, "StrToTime 2021-09-05 18:07:06.123456789 +0000 UTC fail")
	}
	// ANSIC
	if date, err := DateTime("Mon Jan 02 15:04:05 2006"); err == nil {
		assert.Equal(t, int(date.Weekday()), 1)
		assert.Equal(t, date.Day(), 2)
	} else {
		assert.Fail(t, "StrToTime Mon Jan 02 15:04:05 2006 fail")
	}
	// ANSIC, change the weekday
	if date, err := DateTime("Fri Jan 02 15:04:05 2006"); err == nil {
		assert.Equal(t, int(date.Weekday()), 5)
		assert.Equal(t, date.Day(), 6)
	} else {
		assert.Fail(t, "StrToTime Fri Jan 02 15:04:05 2006 fail")
	}
	// UnixDate
	if date, err := DateTime("Mon Jan 02 15:04:05 MST 2006"); err == nil {
		// UTC-0700 => UTC+0800 -> +15h
		date = date.In(localLocation)
		assert.Equal(t, date.Day(), 3)
		assert.Equal(t, date.Hour(), 6)
	} else {
		assert.Fail(t, "StrToTime Mon Jan 02 15:04:05 MST 2006 fail")
	}
	// RubyDate
	if date, err := DateTime("Mon Jan 02 15:04:05 -0700 2006"); err == nil {
		// UTC-0700 => UTC+0800 -> +15h
		date = date.In(localLocation)
		assert.Equal(t, date.Day(), 3)
		assert.Equal(t, date.Hour(), 6)
	} else {
		assert.Fail(t, "StrToTime Mon Jan 02 15:04:05 -0700 2006 fail")
	}
	// RubyDate
	if date, err := DateTime("Fri Jan 02 15:04:05 -0700 2006"); err == nil {
		// UTC-0700 => UTC+0800 -> +15h
		// Mon->Fri => +4d
		date = date.In(localLocation)
		assert.Equal(t, date.Day(), 7)
		assert.Equal(t, date.Hour(), 6)
	} else {
		assert.Fail(t, "StrToTime Fri Jan 02 15:04:05 -0700 2006 fail")
	}
	// RFC850
	if date, err := DateTime("Monday, 02-Jan-06 15:04:05 MST"); err == nil {
		// UTC-0700 => UTC+0800 -> +15h
		date = date.In(localLocation)
		assert.Equal(t, date.Day(), 3)
		assert.Equal(t, date.Hour(), 6)
	} else {
		assert.Fail(t, "StrToTime Monday, 02-Jan-06 15:04:05 MST fail")
	}
	// RFC850
	if date, err := DateTime("Friday, 02-Jan-06 15:04:05 MST"); err == nil {
		// UTC-0700 => UTC+0800 -> +15h
		// Mon->Fri => +4d
		date = date.In(localLocation)
		assert.Equal(t, date.Day(), 7)
		assert.Equal(t, date.Hour(), 6)
	} else {
		assert.Fail(t, "StrToTime Friday, 02-Jan-06 15:04:05 MST fail")
	}
	// RFC822
	if date, err := DateTime("02 Jan 06 15:04 MST"); err == nil {
		// UTC-0700 => UTC+0800 -> +15h
		date = date.In(localLocation)
		assert.Equal(t, date.Day(), 3)
		assert.Equal(t, date.Hour(), 6)
	} else {
		assert.Fail(t, "StrToTime 02 Jan 06 15:04 MST fail")
	}
	// RFC822Z
	if date, err := DateTime("02 Jan 06 15:04 -0700"); err == nil {
		// UTC-0700 => UTC+0800 -> +15h
		date = date.In(localLocation)
		assert.Equal(t, date.Day(), 3)
		assert.Equal(t, date.Hour(), 6)
	} else {
		assert.Fail(t, "StrToTime 02 Jan 06 15:04 MST fail")
	}
	// RFC1123
	if date, err := DateTime("Mon, 02 Jan 2006 15:04:05 MST"); err == nil {
		// UTC-0700 => UTC+0800 -> +15h
		date = date.In(localLocation)
		assert.Equal(t, date.Day(), 3)
		assert.Equal(t, date.Hour(), 6)
	} else {
		assert.Fail(t, "StrToTime Mon, 02 Jan 2006 15:04:05 MST fail")
	}
	// RFC1123Z
	if date, err := DateTime("Fri, 02 Jan 2006 15:04:05 -0700"); err == nil {
		// UTC-0700 => UTC+0800 -> +15h
		// Mon->Fri => +4d
		date = date.In(localLocation)
		assert.Equal(t, date.Day(), 7)
		assert.Equal(t, date.Hour(), 6)
	} else {
		assert.Fail(t, "StrToTime Fri, 02 Jan 2006 15:04:05 -0700 fail")
	}
}

func TestStringTime(t *testing.T) {
	// 6:07:06:12313pm
	if date, err := DateTime("6:07:06:12313pm"); err == nil {
		assert.True(t, isSameDate(&date, His))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime 6:07:06:12313pm fail")
	}
	// 6:07:06 pm
	if date, err := DateTime("6:07:06 pm"); err == nil {
		assert.True(t, isSameDate(&date, His))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime 6:07:06 pm fail")
	}
	// 6.07.06P.M.
	if date, err := DateTime("6.07.06P.M."); err == nil {
		assert.True(t, isSameDate(&date, His))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime 6.07.06P.M. fail")
	}
	// 6.07.06P.M
	if date, err := DateTime("6.07.06P.M"); err == nil {
		assert.True(t, isSameDate(&date, His))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime 6.07.06P.M fail")
	}
	// 180706CET
	if date, err := DateTime("180706CET"); err == nil {
		assert.True(t, isSameDate(&date, MINUTE|SECOND))
	} else {
		assert.Fail(t, "StrToTime 180706CET fail")
	}
	// T180706+0800
	if date, err := DateTime("T180706+0800"); err == nil {
		date = date.In(localLocation)
		assert.True(t, isSameDate(&date, His))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime T180706+0800 fail")
	}
	// 6:07 pm
	if date, err := DateTime("6:07 pm"); err == nil {
		assert.True(t, isSameDate(&date, HOUR|MINUTE))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime 6:07 pm fail")
	}
	// 6:07P.M.
	if date, err := DateTime("6:07P.M."); err == nil {
		assert.True(t, isSameDate(&date, HOUR|MINUTE))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime 6:07P.M. fail")
	}
	// 18:07:06.123456
	if date, err := DateTime("18:07:06.123456"); err == nil {
		assert.True(t, isSameDate(&date, His))
		assert.True(t, isSameDate(&date, YMD, true))
		assert.Equal(t, date.Nanosecond(), 123456000)
	} else {
		assert.Fail(t, "StrToTime 18:07:06.123456 fail")
	}
	// 18.07.06.123
	if date, err := DateTime("18.07.06.123"); err == nil {
		assert.True(t, isSameDate(&date, His))
		assert.True(t, isSameDate(&date, YMD, true))
		assert.Equal(t, date.Nanosecond(), 123000000)
	} else {
		assert.Fail(t, "StrToTime 18.07.06.123 fail")
	}
	// 5.9.21T18.07.06.123
	if date, err := DateTime("5.9.21T18.07.06.123"); err == nil {
		assert.True(t, isSameDate(&date, YMDHis))
		assert.Equal(t, date.Nanosecond(), 123000000)
	} else {
		assert.Fail(t, "StrToTime 5.9.21T18.07.06.123 fail")
	}
	// 18.07.06
	if date, err := DateTime("18.07.06"); err == nil {
		assert.True(t, isSameDate(&date, His))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime 18.07.06 fail")
	}
	// 18.07
	if date, err := DateTime("18.07"); err == nil {
		assert.True(t, isSameDate(&date, HOUR|MINUTE))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime 18.07 fail")
	}
	// 18:07
	if date, err := DateTime("18:07"); err == nil {
		assert.True(t, isSameDate(&date, HOUR|MINUTE))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime 18:07 fail")
	}
	// T18:07
	if date, err := DateTime("T18:07"); err == nil {
		assert.True(t, isSameDate(&date, HOUR|MINUTE))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime T18:07 fail")
	}
	// 6 pm
	if date, err := DateTime("6 pm"); err == nil {
		assert.True(t, isSameDate(&date, HOUR))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime 6 pm fail")
	}
	// 6am
	if date, err := DateTime("6am"); err == nil {
		// not equal
		assert.False(t, isSameDate(&date, HOUR))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime 6am fail")
	}
	// 180706
	if date, err := DateTime("180706"); err == nil {
		assert.True(t, isSameDate(&date, His))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime 180706 fail")
	}
	// t180706
	if date, err := DateTime("t180706"); err == nil {
		assert.True(t, isSameDate(&date, His))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime t180706 fail")
	}
	// 1807
	if date, err := DateTime("1807"); err == nil {
		assert.True(t, isSameDate(&date, HOUR|MINUTE))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime 1807 fail")
	}
	// t1807
	if date, err := DateTime("t1807"); err == nil {
		assert.True(t, isSameDate(&date, HOUR|MINUTE))
		assert.True(t, isSameDate(&date, YMD, true))
	} else {
		assert.Fail(t, "StrToTime t1807 fail")
	}
}

func TestStringDateTime(t *testing.T) {
	// seconds with fraction
	if date, err := DateTime("2021-09-05 06:07:06.123456789PM"); err == nil {
		assert.True(t, isSameDate(&date, YMDHis))
	} else {
		assert.Fail(t, "StrToTime 2021-09-05 06:07:06.123456789PM fail")
	}
	// seconds with fraction
	if date, err := DateTime("2021/09/05 06:07:06.123456789PM"); err == nil {
		assert.True(t, isSameDate(&date, YMDHis))
	} else {
		assert.Fail(t, "StrToTime 2021/09/05 06:07:06.123456789PM fail")
	}
	// seconds with fraction
	if date, err := DateTime("2021/09/05T18:07:06.123456789"); err == nil {
		assert.True(t, isSameDate(&date, YMDHis))
	} else {
		assert.Fail(t, "StrToTime 2021/09/05T18:07:06.123456789 fail")
	}
}

func TestStrToTime(t *testing.T) {
	// number date
	var timestamp int64 = 1630865226
	if date, err := StrToTime(int(timestamp)); err == nil {
		assert.Equal(t, date, timestamp)
	} else {
		assert.Fail(t, "StrToTime int fail")
	}
	if date, err := StrToTime(timestamp); err == nil {
		assert.Equal(t, date, timestamp)
	} else {
		assert.Fail(t, "StrToTime int64 fail")
	}
	if date, err := StrToTime(int32(timestamp)); err == nil {
		assert.Equal(t, date, timestamp)
	} else {
		assert.Fail(t, "StrToTime int32 fail")
	}
	if date, err := StrToTime(float64(timestamp)); err == nil {
		assert.Equal(t, date, timestamp)
	} else {
		assert.Fail(t, "StrToTime float64 fail")
	}
}
