package dateutil

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// FormatList for all formats
type FormatList map[string]string

// FormatResult for result
type FormatResult map[string]string

// ReplaceFn for format
type ReplaceFn func(map[string]string) time.Time

// PatternInfo for all types
type PatternInfo struct {
	Patterns  []*Pattern
	ReplaceFn ReplaceFn
}

const (
	RFCSYMBOL = "rfc"
)

var (
	// months names
	allMonthExp = []string{"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december", "jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "sept", "oct", "nov", "dec", "viii", "iii", "vii", "xii", "iv", "vi", "ix", "xi", "ii", "i", "v", "x"}
	// weekday names
	weekdayFullNames = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	// weekday short
	weekdayShortNames = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	// roman number to arabic number
	romanNumHash = map[rune]int{
		'i': 1,
		'v': 5,
		'x': 10,
	}
	// the below formats can be seen in:
	// https://www.php.net/manual/en/datetime.formats.time.php
	dateFormats = FormatList{
		// year/month/day
		"dd": "(3[0-1]|[0-2]?[0-9])(?:st|nd|rd|th)?",
		"DD": "(0[0-9]|[1-2][0-9]|3[0-1])",
		"m":  "(" + strings.Join(allMonthExp, "|") + ")",
		"M":  "(jan|feb|mar|apr|may|jun|jul|aug|sep|sept|oct|nov|dec)",
		"mm": "(1[0-2]|0?[0-9])",
		"MM": "(0[0-9]|1[0-2])",
		"y":  "([0-9]{1,4})",
		"yy": "([0-9]{2})",
		"YY": "([0-9]{4})",
		// weekday
		"l": "(" + strings.Join(weekdayFullNames, "|") + ")",
		"D": "(" + strings.Join(weekdayShortNames, "|") + ")",
	}
	dateRules = []string{
		// keep the orders, make sure match as much as more characters
		"^[+-]?${YY}-${MM}-${DD}",                    // "-0002-07-26", "+1978-04-17", "1814-05-17"
		"^${dd}[.\\t-]${mm}[.-]${YY}",                // "30-6-2008", "22.12.1978"
		"^${YY}-${mm}-${dd}",                         // "2008-6-30", "1978-12-22"
		"^${yy}-${MM}-${DD}",                         // "08-06-30", "78-12-22"
		"^${y}-${mm}-${dd}",                          // "2008-6-30", "78-12-22", "8-6-21"
		"^(?i)${M}-${DD}-${y}",                       // "May-09-78", "Apr-17-1790"
		"^(?i)${y}-${M}-${DD}",                       // "78-Dec-22", "1814-MAY-17"
		"^${YY}\\/${MM}\\/${DD}",                     // "2008/06/30", "1978/12/22"
		"^${YY}\\/${mm}\\/${dd}",                     // "2008/6/30", "1978/12/22"
		"^${mm}\\/${dd}\\/${y}",                      // "12/22/78", "1/17/2006", "1/17/6"
		"^${dd}[.\\t]${mm}\\.${yy}",                  // "30.6.08", "22\t12.78"
		"^${YY}${MM}${DD}",                           // "15810726", "19780417", "18140517"
		"^${mm}\\/${dd}",                             // "5/12", "10/27"
		"^${YY}-${mm}",                               // "2008-6", "2008-06", "1978-12"
		"(?i)^${dd}[ \\t.-]*${m}[ \\t.-]*${y}",       // "30-June 2008", "22DEC78", "14 III 1879"
		"(?i)^${m}[ \\t.-]*${YY}",                    // "June 2008", "DEC1978", "March 1879"
		"(?i)^${YY}[ \\t.-]*${m}",                    // "2008 June", "1978-XII", "1879.MArCH"
		"(?i)^${m}[ .\\t-]*${dd}[,.stndrh\\t ]+${y}", // "July 1st, 2008", "April 17, 1790", "May.9,78"
		"(?i)^${m}[ .\\t-]*${dd}[,.stndrh\\t ]*",     // "July 1st,", "Apr 17", "May.9"
		"(?i)^${dd}[ .\\t-]*${m}",                    // "1 July", "17 Apr", "9.May"
		"^${YY}",                                     // "1978", "2008"
		"(?i)^${m}",                                  // "March", "jun", "DEC"
	}
	timeFormats = FormatList{
		"frac":               "([0-9]{1,9})",
		"hh":                 "(1[0-2]|0?[0-9])",
		"HH":                 "(1[0-9]|2[0-4]|0?[0-9])",
		"meridian":           "([AaPp]\\.?[Mm](?:\\.?|\\b|$))",
		"MN":                 "([1-5][0-9]|0?[0-9])",
		"MNA":                "([0-5][0-9])",
		"II":                 "([1-5][0-9]|0?[0-9])",
		"IIA":                "([0-5][0-9])",
		"tz":                 "([A-Z][a-z]+(?:[_/][A-Z][a-z]+)+|\\([A-Za-z]{1,6}\\)|[A-Za-z]{1,6})",
		"tz_plain":           "([A-Z]{1,6})",
		"tzcorrection":       "((?:GMT)?[+-](?:1[0-2]|0?[0-9]):?(?:[0-5][0-9])?)",
		"tzcorrection_plain": "([+-](?:1[0-2]|0?[0-9]):?(?:[0-5][0-9])?)",
	}
	timeRules = []string{
		"(?i)^${hh}:${MN}:${II}[.:]${frac}${meridian}$",                       // "4:08:39:12313am"
		"(?i)^${hh}[.:]${MN}[.:]${II}[ \\t]?${meridian}$",                     // "4:08:37 am", "7:19:19P.M."
		"(?i)^t?${HH}[.:]?${MNA}[.:]?${IIA}[ \\t]?(?:${tzcorrection}|${tz})$", // "040837CEST", "T191919-0700"
		"(?i)^${hh}[.:]${MN}[ \\t]?${meridian}$",                              // "4:08 am", "7:19P.M."
		"(?i)^t?${HH}[.:]${MN}[.:]${II}\\.${frac}$",                           // "04.08.37.81412", "19:19:19.532453"
		"(?i)^t?${HH}[.:]${MN}[.:]${II}$",                                     // "04.08.37", "t19:19:19"
		"(?i)^t?${HH}[.:]${MN}$",                                              // "04:08", "19.19", "T23:43"
		"(?i)^${hh}[ \\t]?${meridian}$",                                       // "4 am", "5PM"
		"(?i)^t?${HH}${MNA}${IIA}$",                                           // "040837", "T191919"
		"(?i)^t?${HH}${MNA}$",                                                 // "0408", "t1919", "T2343"
		"(?i)^(?:${tzcorrection}|${tz})$",                                     // "CEST", "Europe/Amsterdam", "+0430", "GMT-06:00"
	}
	// will fill next
	rfcFormats = FormatList{}
	rfcRules   = []string{
		// golang format time.String(): "2006-01-02 15:04:05.999999999 -0700 MST"
		"^${YY}-${MM}-${DD}[ \\t]+${HH}:${MN}:${II}(?:\\.${frac})?[ \\t]+${tzcorrection_plain}[ \\t]+${tz_plain}",
		// golang default laytout: "01/02 03:04:05PM '06 -0700"
		"^${MM}\\/${DD}[ \\t]+${hh}-${MN}-${II}${meridian}[ \\t]+'${yy}[ \\t]+${tzcorrection_plain}",
		// ANSIC: "Mon Jan _2 15:04:05 2006"
		// UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
		// RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
		"^(?i)${D}[ \\t]+${M}[ \\t]+${DD}[ \\t]+${HH}:${MN}:${II}(?:[ \\t]+(?:${tz_plain}|${tzcorrection_plain}))?[ \\t]+${YY}",
		// RFC822      = "02 Jan 06 15:04 MST"
		// RFC822Z     = "02 Jan 06 15:04 -0700"
		"^(?i)${DD}[ \\t]+${M}[ \\t]+${yy}[ \\t]+${HH}:${MN}[ \\t]+(?:${tz_plain}|${tzcorrection_plain})",
		// RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
		"^(?i)${l},[ \\t]+${DD}-${M}-${yy}[ \\t]+${HH}:${MN}:${II}[ \\t]+${tz_plain}",
		// RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
		// RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700"
		"^(?i)${D},[ \\t]+${DD}[ \\t]+${M}[ \\t]+${YY}[ \\t]+${HH}:${MN}:${II}[ \\t]+(?:${tz_plain}|${tzcorrection_plain})",
	}
	allFormats = map[string]*FormatList{
		"date": &dateFormats,
		"time": &timeFormats,
	}
	allPatternInfo = map[string]*PatternInfo{}
)

// get weekday number
func getWeekdayNum(weekday string) int {
	for num, name := range weekdayShortNames {
		if strings.EqualFold(name, weekday[0:3]) {
			return num
		}
	}
	return -1
}

// check if is not a timezone but an en month
func transTimezoneResult(result FormatResult) (FormatResult, bool) {
	if isResultTimezone(result) {
		if tz, ok := result["tz"]; ok && tz != "" {
			month := strings.ToLower(tz)
			for _, name := range allMonthExp {
				if name == month {
					if len(name) == 3 {
						return FormatResult{
							"M": month,
						}, true
					} else {
						return FormatResult{
							"m": month,
						}, true
					}
				}
			}
		}
	}
	return FormatResult{}, false
}

func isResultTimezone(result FormatResult) bool {
	if len(result) > 2 {
		return false
	}
	timezoneFields := map[string]bool{
		"tz":                 true,
		"tz_plain":           true,
		"tzcorrection":       true,
		"tzcorrection_plain": true,
	}
	for key := range result {
		if _, ok := timezoneFields[key]; !ok {
			return false
		}
	}
	return true
}

// Pattern struct
type Pattern struct {
	Rule     *regexp.Regexp
	Keys     []string
	Type     string
	Original string
}

// Match method
func (pattern *Pattern) Match(target string) (FormatResult, []int, bool) {
	rule, keys := pattern.Rule, pattern.Keys
	if loc := rule.FindStringIndex(target); loc != nil {
		result := FormatResult{}
		matchs := rule.FindStringSubmatch(target)
		if len(matchs) == len(keys)+1 {
			for index, value := range matchs[1:] {
				result[keys[index]] = value
			}
		}
		return result, loc, true
	}
	return nil, nil, false
}

// StrToTime to timestamp
func StrToTime(target interface{}) (int64, error) {
	nowTime, err := DateTime(target)
	if err == nil {
		return nowTime.Unix(), nil
	}
	return 0, err
}

// DateTime func
func DateTime(target interface{}) (time.Time, error) {
	var timestamp int64
	switch t := target.(type) {
	case time.Time:
		return t.In(time.Local), nil
	case int64:
		timestamp = t
	case int:
		timestamp = int64(t)
	case int32:
		timestamp = int64(t)
	case float64:
		timestamp = int64(t)
	case string:
		var (
			lasts     FormatResult
			isRFCTime bool
		)
		t = strings.TrimSpace(t)
		// match golang and rfc format first
		if result, loc, ok := matchRFCFormat(t); ok {
			if len(t) == loc[1] {
				lasts = result
				isRFCTime = true
			}
		}
		// not rfc time
		if !isRFCTime {
			// match time first
			if result, _, ok := matchTimeFormat(t); ok {
				// if timezone, but match a en month
				if timeResult, ok := transTimezoneResult(result); ok {
					lasts = timeResult
				} else {
					lasts = result
				}
			} else {
				// not a time format, so maybe a date or a datetime
				// match the date format first
				if result, loc, ok := matchDateFormat(t); ok {
					// set lasts
					lasts = result
					// set next index
					nextIndex := loc[1]
					// get the left characters after date string
					suffix := t[nextIndex:]
					// no more characters
					timeFormat := strings.TrimSpace(suffix)
					// special date
					if timeFormat != "" {
						if result, _, ok := matchTimeFormat(timeFormat); ok {
							for key, value := range result {
								lasts[key] = value
							}
						} else {
							return time.Time{}, fmt.Errorf("wrong time format:'%s'", t)
						}
					}
				} else {
					return time.Time{}, fmt.Errorf("wrong date or datetime:'%s'", t)
				}
			}
		}
		if lasts != nil {
			return makeFormatDateTime(lasts)
		}
		return time.Time{}, fmt.Errorf("wrong datetime string:%s", t)
	default:
		// other conditions
		return time.Time{}, fmt.Errorf("can't parse the datetime: %#v", target)
	}
	return time.Unix(timestamp, 0).In(time.Local), nil
}

// get any of the argument fields in the target format result
// if no field found, return an empty string
func noEmptyField(target FormatResult, args ...string) string {
	for _, field := range args {
		if cur, ok := target[field]; ok && cur != "" {
			return cur
		}
	}
	return ""
}

// translate result information to a time struct
func makeFormatDateTime(result FormatResult) (time.Time, error) {
	// tz, tzcorrection
	var lastTime time.Time
	// set default timezone as 'Local'
	timezone := "Local"
	// get matched tz
	hasTimezone := false
	tz := noEmptyField(result, "tz", "tz_plain")
	tzcorrection := noEmptyField(result, "tzcorrection", "tzcorrection_plain")
	needCorrection := false
	// set the timezone
	if tz != "" {
		timezone = tz
		hasTimezone = true
	}
	// correction the time to +0000 of timezone
	if tzcorrection != "" {
		// set timezone to GMT
		if !hasTimezone {
			timezone = "GMT"
		}
		needCorrection = true
	}
	// load location
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	// plain timezone, set hour/minute/second/nanoseconds to now time
	if (hasTimezone || needCorrection) && isResultTimezone(result) {
		// use UTC time
		now := time.Now().UTC()
		lastTime = time.Date(now.Year(), time.Month(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), location)
	} else {
		// current time
		now := time.Now()
		// get full year of current
		year := now.Year()
		strYear := strconv.Itoa(year)
		rnYear := []rune(strYear)
		curYear := noEmptyField(result, "YY", "yy", "y")
		if curYear != "" {
			rns := []rune(curYear)
			total := len(curYear)
			switch total {
			case 4:
			case 3:
				// e.g. '032' -> '2032'
				if rns[0] == '0' {
					rns = append(rnYear[0:1], rns...)
				}
			case 2:
				rns = append(rnYear[0:2], rns...)
			case 1:
				rns = append(rnYear[0:2], []rune{'0', rns[0]}...)
			}
			year, _ = strconv.Atoi(string(rns))
		}
		curMonth := noEmptyField(result, "MM", "mm")
		// month
		var month int
		if curMonth != "" {
			month, _ = strconv.Atoi(curMonth)
		} else {
			curMonth = noEmptyField(result, "M", "m")
			if curMonth != "" {
				curMonth = strings.ToLower(curMonth)
				for index, name := range allMonthExp {
					if name == curMonth {
						if index >= 24 {
							// roman
							month = 0
							prev := 0
							rns := []rune(name)
							for idx, num := range rns {
								value := romanNumHash[num]
								if idx > 0 && value > prev {
									month = value - prev
								} else {
									month += value
								}
								prev = value
							}
						} else {
							month = index%12 + 1
						}
						break
					}
				}
			} else {
				month = int(now.Month())
			}
		}
		// day
		var day int
		curDay := noEmptyField(result, "DD", "dd")
		if curDay != "" {
			day, _ = strconv.Atoi(curDay)
		} else {
			day = now.Day()
		}
		// hour
		var hour int
		curHour := noEmptyField(result, "HH", "hh")
		if curHour != "" {
			hour, _ = strconv.Atoi(curHour)
			if meridian, ok := result["meridian"]; ok && (meridian[0] == 'P' || meridian[0] == 'p') {
				hour += 12
			}
		} else {
			hour = 0
		}
		// minute
		var minute int
		curMinute := noEmptyField(result, "MN", "MNA")
		if curMinute != "" {
			minute, _ = strconv.Atoi(curMinute)
		} else {
			minute = 0
		}
		// second
		var second int
		curSecond := noEmptyField(result, "II", "IIA")
		if curSecond != "" {
			second, _ = strconv.Atoi(curSecond)
		} else {
			second = 0
		}
		// nanoSeconds
		var nanoSeconds int
		fracSeconds := noEmptyField(result, "frac")
		if fracSeconds != "" {
			exp := 9 - len(fracSeconds)
			nanoSeconds, _ = strconv.Atoi(fracSeconds)
			if exp > 1 {
				nanoSeconds = nanoSeconds * int(math.Pow10(exp))
			}
		} else {
			nanoSeconds = 0
		}
		// make date time
		lastTime = time.Date(year, time.Month(month), day, hour, minute, second, nanoSeconds, location)
		// weekday
		weekday := noEmptyField(result, "l", "D")
		if weekday != "" {
			// if the weekday is not current day
			// fix the day to that weekday
			curWeekday := int(lastTime.Weekday())
			relWeekday := getWeekdayNum(weekday)
			forwardDays := relWeekday - curWeekday
			if forwardDays != 0 {
				// make sure the days is increased
				if forwardDays < 0 {
					forwardDays += 7
				}
				// add date
				lastTime = lastTime.AddDate(0, 0, forwardDays)
			}
		}
	}
	// fix time to GMT/UTC+0000 time
	if needCorrection {
		tzcorrection = strings.TrimPrefix(tzcorrection, "GMT")
		rns := []rune(tzcorrection)
		count := len(rns)
		var (
			multi, addMinute, addHour int64
		)
		multi = 1
		if rns[0] == '+' {
			multi = -1
		}
		addHour, _ = strconv.ParseInt(string(rns[1:3]), 10, 64)
		if count >= 5 {
			addMinute, _ = strconv.ParseInt(string(rns[count-2:]), 10, 64)
		}
		correct := (time.Duration(addHour)*time.Hour + time.Duration(addMinute)*time.Minute) * time.Duration(multi)
		lastTime = lastTime.Add(correct)
	}
	// Change time to local time
	if timezone != "Local" {
		lastTime = lastTime.In(time.Local)
	}
	return lastTime, nil
}

// make patterns
// save the patterns into global variable 'allPatternInfo'
func makePatterns(t string, rules ...string) (*PatternInfo, error) {
	var (
		ptns       = []*Pattern{}
		ok         bool
		formatList *FormatList
	)
	if formatList, ok = allFormats[t]; ok {
		// finded in all formats
	} else if t == RFCSYMBOL {
		// append to rfc formats
		for _, list := range allFormats {
			for k, v := range *list {
				rfcFormats[k] = v
			}
		}
		// dynamic add rfc to all formats
		allFormats[RFCSYMBOL] = &rfcFormats
		// set the variables
		formatList = &rfcFormats
		ok = true
	}
	if ok {
		regRule, _ := regexp.Compile(`\$\{[A-Za-z_]+}`)
		for _, rule := range rules {
			pattern := new(Pattern)
			pattern.Type = t
			pattern.Original = rule
			keys := []string{}
			context := regRule.ReplaceAllStringFunc(rule, func(all string) string {
				rns := []rune(all)
				key := string(rns[2 : len(rns)-1])
				if seg, ok := (*formatList)[key]; ok {
					keys = append(keys, key)
					return seg
				}
				return all
			})
			curRule := regexp.MustCompile(context)
			pattern.Rule = curRule
			pattern.Keys = keys
			ptns = append(ptns, pattern)
		}
		info := new(PatternInfo)
		info.Patterns = ptns
		allPatternInfo[t] = info
		return info, nil
	}
	return nil, fmt.Errorf("the format type '%s' doesn't exist", t)
}

// factory match format
func factoryMatchFormat(key string, rules []string, target string) (FormatResult, []int, bool) {
	var patterns []*Pattern
	if info, ok := allPatternInfo[key]; ok {
		// get patterns from cache
		patterns = info.Patterns
	} else {
		// save patterns to cache
		info, _ := makePatterns(key, rules...)
		patterns = info.Patterns
	}
	for _, pattern := range patterns {
		if result, loc, ok := pattern.Match(target); ok {
			return result, loc, true
		}
	}
	return nil, nil, false
}

// golang/RFC formats
func matchRFCFormat(target string) (FormatResult, []int, bool) {
	return factoryMatchFormat(RFCSYMBOL, rfcRules, target)
}

// date fomrats
func matchDateFormat(target string) (FormatResult, []int, bool) {
	return factoryMatchFormat("date", dateRules, target)
}

// time formats
func matchTimeFormat(target string) (FormatResult, []int, bool) {
	return factoryMatchFormat("time", timeRules, target)
}

// DateFormat func
func DateFormat(target interface{}, format string) (string, error) {
	// the golang layout, use the format of the golang born time
	// this is used for strings.NewReplacer
	// so they are appeared pairs with a string array
	layouts := []string{
		// year
		"Y", "2006",
		"y", "06",
		// month
		"m", "01",
		"n", "1",
		// date
		"d", "02",
		"j", "2",
		// hours
		"h", "03",
		"g", "3",
		"G", "15",
		// minutes
		"i", "04",
		// seconds
		"s", "05",
	}
	formats := map[string]string{
		// am, pm
		"a": "pm",
		"A": "PM",
		// month
		"F": "January",
		"M": "Jan",
		// week
		"D": "Mon",
		"l": "Monday",
	}
	// weekday, from monday to sunday
	// monday return 1, sunday return 7
	// it's the same as golang's weekday of time.Time.
	N := func(t time.Time) string {
		weekday := t.Weekday()
		return fmt.Sprintf("%d", int(weekday))
	}
	// weekday, from sunday to saturday
	// sunday return 0, saturday return 6
	w := func(t time.Time) string {
		weekday := t.Weekday()
		dayNum := (int(weekday) + 1) % 7
		return fmt.Sprintf("%d", dayNum)
	}
	// the day of the year, from 0 to 365
	// the golang's yearday is from 1 to 366
	// so here need reduce by one day
	z := func(t time.Time) string {
		yearday := t.YearDay()
		return fmt.Sprintf("%d", yearday-1)
	}
	// the nth week of a year
	W := func(t time.Time) string {
		_, week := t.ISOWeek()
		return fmt.Sprintf("%d", week)
	}
	// check if is leap year
	L := func(t time.Time) string {
		yearday := t.YearDay()
		if yearday > 365 {
			return "1"
		}
		return "0"
	}
	// get how many days of the month
	t := func(t time.Time) string {
		nums := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
		// if is leap Year, the second month has 29 days.
		monthIndex := int(t.Month()) - 1
		if monthIndex == 1 && t.YearDay() > 365 {
			return "29"
		}
		return fmt.Sprintf("%d", nums[monthIndex])
	}
	// get the hour of the time, 2 digits with leading zero
	H := func(t time.Time) string {
		hour := t.Hour()
		return fmt.Sprintf("%02d", hour)
	}
	// get the microseconds of the time
	u := func(t time.Time) string {
		nano := t.Nanosecond()
		return fmt.Sprintf("%06d", nano/1e3)
	}
	// get the milliseconds of the time
	v := func(t time.Time) string {
		nano := t.Nanosecond()
		return fmt.Sprintf("%03d", nano/1e6)
	}
	fns := map[string]func(t time.Time) string{
		"N": N,
		"w": w,
		"z": z,
		"W": W,
		"L": L,
		"t": t,
		"H": H,
		"u": u,
		"v": v,
	}
	// numbers in date can format first by the golang layout
	// "Y-m-da" will translate into -> '2006-01-02a'
	repRule := strings.NewReplacer(layouts...)
	layout := repRule.Replace(format)
	// The other format keyword such as before 'a'
	// Will use golang layout format or functions
	regRule := func() *regexp.Regexp {
		var str strings.Builder
		str.WriteString("[")
		for key := range formats {
			str.WriteString(key)
		}
		for key := range fns {
			str.WriteString(key)
		}
		str.WriteString("]")
		rule, _ := regexp.Compile(str.String())
		return rule
	}()
	// Change target to time struct
	var timeTarget time.Time
	if cur, ok := target.(time.Time); ok {
		timeTarget = cur
	} else {
		if cur, err := DateTime(target); err == nil {
			timeTarget = cur
		} else {
			return "", err
		}
	}
	// Format date numbers first, no letter characters will making
	result := timeTarget.Format(layout)
	// Replace the keyword letter character into real value
	result = regRule.ReplaceAllStringFunc(result, func(name string) string {
		if layout, ok := formats[name]; ok {
			return timeTarget.Format(layout)
		} else if fn, ok := fns[name]; ok {
			return fn(timeTarget)
		}
		return ""
	})
	return result, nil
}
