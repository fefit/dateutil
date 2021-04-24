package dateutil

import (
	"fmt"
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
		"frac":         "(\\.[0-9]+)",
		"hh":           "(0?[0-9]|1[0-2])",
		"HH":           "([01][0-9]|2[0-4])",
		"meridian":     "([AaPp]\\.?[Mm]\\.?[\\0\\t])",
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
		total := len(t)
		if result, loc, ok := matchDateFormat(t); ok {
			if loc[0] != 0 {
				// ignore, try time format
			} else {
				suffix := t[loc[1]:]
				plainSuffix := strings.TrimSpace(suffix)
				isJustSpaces := suffix == "" || plainSuffix == ""
				if _, ok := result["YY"]; ok && len(result) == 1 && isJustSpaces {
					// ignore, use time format first
				} else {
					lasts = result
					if loc[1] < total && !isJustSpaces {
						if strings.HasPrefix(suffix, " ") {
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
					return time.Time{}, fmt.Errorf("wrong start of time")
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
func noEmptyField(target FormatResult, args ...string) string {
	for _, field := range args {
		if cur, ok := target[field]; ok && cur != "" {
			return cur
		}
	}
	return ""
}

// translate result info to time
func makeFormatDateTime(result FormatResult) (time.Time, error) {
	now := time.Now()
	// year
	year := now.Year()
	strYear := strconv.Itoa(year)
	rnYear := []rune(strYear)
	curYear := noEmptyField(result, "YY", "yy", "y")
	if curYear != "" {
		rns := []rune(curYear)
		total := len(curYear)
		switch total {
		case 3, 4:
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
	// millisecond
	var milliSecond int
	curMillisecond := noEmptyField(result, "frac")
	if curMillisecond != "" {
		milliSecond, _ = strconv.Atoi(curMillisecond[1:])
	} else {
		milliSecond = 0
	}
	// tz, tzcorrection
	var lastTime time.Time
	timezone := "Local"
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
	lastTime = time.Date(year, time.Month(month), day, hour, minute, second, milliSecond*1e6, location)
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
		info, _ := makePatterns("date",
			"[+-]?${YY}-${MM}-${DD}",
			"${YY}-${mm}-${dd}",
			"${YY}\\/${MM}\\/${DD}",
			"${yy}-${MM}-${DD}",
			"${YY}${MM}${DD}",
			"${mm}\\/${dd}",         // 5/12
			"${mm}\\/${dd}\\/${y}",  // 5/12/06
			"${YY}\\/${mm}\\/${dd}", // 2006/5/12
			"${YY}-${mm}",
			"${y}-${mm}-${dd}",
			"${dd}[.\\t-]${mm}[.-]${YY}",
			"${dd}[.\\t]${mm}\\.${yy}",
			"(?i)${dd}[ \\t.-]*${m}[ \\t.-]*${y}",
			"(?i)${m}[ \\t.-]*${YY}",
			"(?i)${YY}[ \\t.-]*${m}",
			"(?i)${m}[ .\\t-]*${dd}[,.stndrh\\t ]+${y}",
			"(?i)${m}[ .\\t-]*${dd}[,.stndrh\\t ]*",
			"(?i)${d}[ .\\t-]*${m}",
			"(?i)${M}-${DD}-${y}",
			"(?i)${y}-${M}-${DD}",
			"${YY}",
			"(?i)${m}",
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
		info, _ := makePatterns("time",
			"(?i)t?${HH}[.:]${MN}[.:]${II}${space}?(?:${tzcorrection}|${tz})",
			"(?i)t?${HH}[.:]${MN}[.:]${II}",
			"(?i)${hh}[.:]${MN}[.:]${II}${space}?${meridian}",
			"(?i)${hh}${space}?${meridian}",
			"(?i)${hh}[.:]${MN}${space}?${meridian}",
			"(?i)${hh}:${MN}:${II}[.:][0-9]+${meridian}",
			"(?i)t?${HH}[.:]${MN}",
			"(?i)t?${HH}${MN}",
			"(?i)t?${HH}[.:]${MN}[.:]${II}${frac}",
			"(?i)t?${HH}${MN}${II}",
			"(?i)(?:${tzcorrection}|${tz})",
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
	N := func(t time.Time) string {
		weekday := t.Weekday()
		return fmt.Sprintf("%d", int(weekday))
	}
	w := func(t time.Time) string {
		weekday := t.Weekday()
		dayNum := int(weekday) % 7
		return fmt.Sprintf("%d", dayNum)
	}
	z := func(t time.Time) string {
		yearday := t.YearDay()
		return fmt.Sprintf("%d", yearday-1)
	}
	W := func(t time.Time) string {
		_, week := t.ISOWeek()
		return fmt.Sprintf("%d", week)
	}
	L := func(t time.Time) string {
		yearday := t.YearDay()
		if yearday > 365 {
			return "1"
		}
		return "0"
	}
	t := func(t time.Time) string {
		nums := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
		monthIndex := int(t.Month()) - 1
		if monthIndex == 1 && t.YearDay() > 365 {
			return "29"
		}
		return fmt.Sprintf("%d", nums[monthIndex])
	}
	H := func(t time.Time) string {
		hour := t.Hour()
		return fmt.Sprintf("%02d", hour)
	}
	fns := map[string]func(t time.Time) string{
		"N": N,
		"w": w,
		"z": z,
		"W": W,
		"L": L,
		"t": t,
		"H": H,
	}
	repRule := strings.NewReplacer(layouts...)
	layout := repRule.Replace(format)
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
	result := timeTarget.Format(layout)
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
