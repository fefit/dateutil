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

var (
	allMonthExp  = []string{"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december", "jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "sept", "oct", "nov", "dec", "viii", "iii", "vii", "xii", "iv", "vi", "ix", "xi", "ii", "i", "v", "x"}
	romanNumHash = map[rune]int{
		'i': 1,
		'v': 5,
		'x': 10,
	}
	// the below formats can be seen in:
	// https://www.php.net/manual/en/datetime.formats.time.php
	dateFormats = FormatList{
		"dd": "([0-2]?[0-9]|3[0-1])(?:st|nd|rd|th)?",
		"DD": "(0[0-9]|[1-2][0-9]|3[0-1])",
		"m":  "(" + strings.Join(allMonthExp, "|") + ")",
		"M":  "(jan|feb|mar|apr|may|jun|jul|aug|sep|sept|oct|nov|dec)",
		"mm": "(0?[0-9]|1[0-2])",
		"MM": "(0[0-9]|1[0-2])",
		"y":  "([0-9]{1,4})",
		"yy": "([0-9]{2})",
		"YY": "([0-9]{4})",
	}
	timeFormats = FormatList{
		"frac":         "([0-9]{1,9})",
		"hh":           "(0?[0-9]|1[0-2])",
		"HH":           "([01][0-9]|2[0-4])",
		"meridian":     "([AaPp]\\.?[Mm]\\.?)\\b",
		"MN":           "([0-5][0-9])",
		"II":           "([0-5][0-9])",
		"space":        "([ \\t])",
		"tz":           "(\\(?[A-Za-z]{1,6}\\)?|[A-Z][a-z]+(?:[_/][A-Z][a-z]+)+)",
		"tzcorrection": "((?:GMT)?[+-](?:0?[0-9]|1[0-2]):?(?:[0-5][0-9])?)",
	}
	allFormats = map[string]*FormatList{
		"date": &dateFormats,
		"time": &timeFormats,
	}
	allPatternInfo = map[string]*PatternInfo{}
)

func isTimeFormat(format string) bool {
	numbers := strings.Split(format, "")
	if len(numbers) == 4 {
		maxNums := []int{2, 3, 5, 9}
		for index, strNum := range numbers {
			num, err := strconv.Atoi(strNum)
			if err != nil {
				return false
			}
			if num > maxNums[index] {
				return false
			}
		}
		return true
	}
	return false
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
	switch t := target.(type) {
	case int64:
		return t, nil
	case int:
		return int64(t), nil
	case int32:
		return int64(t), nil
	case float64:
		return int64(t), nil
	case string:
		nowTime, err := DateTime(t)
		if err == nil {
			return nowTime.UTC().Unix(), nil
		}
		return 0, err
	default:
		return 0, fmt.Errorf("unsupport type %T for StrToTime func", t)
	}
}

// DateTime func
func DateTime(target interface{}) (time.Time, error) {
	var timestamp int64
	switch t := target.(type) {
	case int64:
		timestamp = t
	case int:
		timestamp = int64(t)
	case int32:
		timestamp = int64(t)
	case float64:
		timestamp = int64(t)
	case string:
		var lasts FormatResult
		t = strings.TrimSpace(t)
		timeFormat := t
		// match the date format first
		if result, loc, ok := matchDateFormat(t); ok {
			if loc[0] != 0 {
				// ignore, try time format
			} else {
				nextIndex := loc[1]
				// get the left characters after date string
				suffix := t[nextIndex:]
				// no more characters
				plainSuffix := strings.TrimSpace(suffix)
				isJustSpaces := suffix == "" || plainSuffix == ""
				// check if just have four numbers, if true, and is a time format, take it as time format
				if year, ok := result["YY"]; ok && len(result) == 1 && isJustSpaces && isTimeFormat(year) {
					// ignore, use time format first
				} else {
					lasts = result
					// if have more characters
					if !isJustSpaces {
						// the next characters begin with a whitespace or 't/T'
						if strings.HasPrefix(suffix, " ") || plainSuffix[0] == 't' || plainSuffix[0] == 'T' {
							timeFormat = plainSuffix
						} else {
							return time.Time{}, fmt.Errorf("wrong datetime %s", t)
						}
					} else {
						// no need for time format
						timeFormat = ""
					}
				}
			}
		}
		if timeFormat != "" {
			if result, loc, ok := matchTimeFormat(timeFormat); ok {
				if loc[0] != 0 || loc[1] != len(timeFormat) {
					return time.Time{}, fmt.Errorf("can't format the time string: \"%s\"", timeFormat)
				}
				if lasts == nil {
					lasts = result
				} else {
					for key, value := range result {
						lasts[key] = value
					}
				}
			}
		}
		if lasts != nil {
			return makeFormatDateTime(lasts)
		}
		return time.Time{}, fmt.Errorf("wrong datetime string:%s", t)
	default:
		// other conditions
		return time.Time{}, fmt.Errorf("wrong datatime %v", target)
	}
	location, _ := time.LoadLocation("Local")
	return time.Unix(timestamp, 0).UTC().In(location), nil
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
			month = 1
		}
	}
	// day
	var day int
	curDay := noEmptyField(result, "DD", "dd")
	if curDay != "" {
		day, _ = strconv.Atoi(curDay)
	} else {
		day = 1
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
	curMinute := noEmptyField(result, "MN")
	if curMinute != "" {
		minute, _ = strconv.Atoi(curMinute)
	} else {
		minute = 0
	}
	// second
	var second int
	curSecond := noEmptyField(result, "II")
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
	// tz, tzcorrection
	var lastTime time.Time
	// set default timezone as 'Local'
	timezone := "Local"
	// get matched tz
	tz := noEmptyField(result, "tz")
	tzcorrection := noEmptyField(result, "tzcorrection")
	needCorrection := false
	if tz != "" {
		timezone = tz
	} else {
		if tzcorrection != "" {
			timezone = "UTC"
			needCorrection = true
		}
	}
	location, _ := time.LoadLocation(timezone)
	lastTime = time.Date(year, time.Month(month), day, hour, minute, second, nanoSeconds, location)
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
	if timezone != "Local" {
		location, _ = time.LoadLocation("Local")
		lastTime = lastTime.In(location)
	}
	return lastTime, nil
}

// make patterns
// save the patterns into global variable 'allPatternInfo'
func makePatterns(t string, rules ...string) (*PatternInfo, error) {
	if formatList, ok := allFormats[t]; ok {
		regRule, _ := regexp.Compile(`\$\{[A-Za-z]+}`)
		ptns := []*Pattern{}
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

// date fomrats
func matchDateFormat(target string) (FormatResult, []int, bool) {
	var patterns []*Pattern
	if info, ok := allPatternInfo["date"]; ok {
		patterns = info.Patterns
	} else {
		// keep the orders
		// make sure the match is greedy, as much characters as it can
		info, _ := makePatterns("date",
			"[+-]?${YY}-${MM}-${DD}",                    // "-0002-07-26", "+1978-04-17", "1814-05-17"
			"${YY}-${mm}-${dd}",                         // "2008-6-30", "1978-12-22"
			"${yy}-${MM}-${DD}",                         // "08-06-30", "78-12-22"
			"${y}-${mm}-${dd}",                          // "2008-6-30", "78-12-22", "8-6-21"
			"(?i)${M}-${DD}-${y}",                       // "May-09-78", "Apr-17-1790"
			"(?i)${y}-${M}-${DD}",                       // "78-Dec-22", "1814-MAY-17"
			"${YY}\\/${MM}\\/${DD}",                     // "2008/06/30", "1978/12/22"
			"${YY}\\/${mm}\\/${dd}",                     // "2008/6/30", "1978/12/22"
			"${mm}\\/${dd}\\/${y}",                      // "12/22/78", "1/17/2006", "1/17/6"
			"${dd}[.\\t-]${mm}[.-]${YY}",                // "30-6-2008", "22.12.1978"
			"${dd}[.\\t]${mm}\\.${yy}",                  // "30.6.08", "22\t12.78"
			"${YY}${MM}${DD}",                           // "15810726", "19780417", "18140517"
			"${mm}\\/${dd}",                             // "5/12", "10/27"
			"${YY}-${mm}",                               // "2008-6", "2008-06", "1978-12"
			"(?i)${dd}[ \\t.-]*${m}[ \\t.-]*${y}",       // "30-June 2008", "22DEC78", "14 III 1879"
			"(?i)${m}[ \\t.-]*${YY}",                    // "June 2008", "DEC1978", "March 1879"
			"(?i)${YY}[ \\t.-]*${m}",                    // "2008 June", "1978-XII", "1879.MArCH"
			"(?i)${m}[ .\\t-]*${dd}[,.stndrh\\t ]+${y}", // "July 1st, 2008", "April 17, 1790", "May.9,78"
			"(?i)${m}[ .\\t-]*${dd}[,.stndrh\\t ]*",     // "July 1st,", "Apr 17", "May.9"
			"(?i)${d}[ .\\t-]*${m}",                     // "1 July", "17 Apr", "9.May"
			"${YY}",                                     // "1978", "2008"
			"(?i)${m}",                                  // "March", "jun", "DEC"
		)
		patterns = info.Patterns
	}
	for _, pattern := range patterns {
		if result, loc, ok := pattern.Match(target); ok {
			return result, loc, true
		}
	}
	return nil, nil, false
}

// time formats
func matchTimeFormat(target string) (FormatResult, []int, bool) {
	var patterns []*Pattern
	if info, ok := allPatternInfo["time"]; ok {
		patterns = info.Patterns
	} else {
		// keep the orders
		info, _ := makePatterns("time",
			"(?i)${hh}:${MN}:${II}[.:]${frac}${meridian}",                     // "4:08:39:12313am"
			"(?i)${hh}[.:]${MN}[.:]${II}${space}?${meridian}",                 // "4:08:37 am", "7:19:19P.M."
			"(?i)t?${HH}[.:]${MN}[.:]${II}${space}?(?:${tzcorrection}|${tz})", // "040837CEST", "T191919-0700"
			"(?i)${hh}[.:]${MN}${space}?${meridian}",                          // "4:08:37 am", "7:19:19P.M."
			"(?i)t?${HH}[.:]${MN}[.:]${II}\\.${frac}",                         // "04.08.37.81412", "19:19:19.532453"
			"(?i)t?${HH}[.:]${MN}[.:]${II}",                                   // "04.08.37", "t19:19:19"
			"(?i)t?${HH}[.:]${MN}",                                            //	"04:08", "19.19", "T23:43"
			"(?i)${hh}${space}?${meridian}",                                   // "4 am", "5PM"
			"(?i)t?${HH}${MN}${II}",                                           // "04.08.37", "t19:19:19"
			"(?i)t?${HH}${MN}",                                                // "0408", "t1919", "T2343"
			"(?i)(?:${tzcorrection}|${tz})",                                   // "CEST", "Europe/Amsterdam", "+0430", "GMT-06:00"
		)
		patterns = info.Patterns
	}
	for _, pattern := range patterns {
		if result, loc, ok := pattern.Match(target); ok {
			return result, loc, true
		}
	}
	return nil, nil, false
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
