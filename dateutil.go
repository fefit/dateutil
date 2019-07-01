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
	dateFormats = FormatList{
		"dd": "([0-2]?[0-9]|3[0-1])(?:st|nd|rd|th)?",
		"DD": "(0[0-9]|[1-2][0-9]|3[0-1])",
		"m":  "(january|february|march|april|may|june|july|august|september|october|november|december|jan|feb|mar|apr|may|jun|jul|aug|sep|sept|oct|nov|dec|I|II|III|IV|V|VI|VII|VIII|IX|X|XI|XII)",
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
		"tzcorrection": "(?:GMT)?[+-](0?[0-9]|1[0-2]):?(?:[0-5][0-9])?",
	}
	allFormats = map[string]*FormatList{
		"date": &dateFormats,
		"time": &timeFormats,
	}
	allPatternInfo = map[string]*PatternInfo{}
)

// Pattern struct
type Pattern struct {
	Rule *regexp.Regexp
	Keys []string
	Type string
}

// Match method
func (pattern *Pattern) Match(target string) (FormatResult, []int, bool) {
	rule, keys := pattern.Rule, pattern.Keys
	if loc := rule.FindStringIndex(target); loc != nil {
		result := FormatResult{}
		replaceWith(rule, target, func(args ...string) string {
			for index, value := range args[1:] {
				key := keys[index]
				result[key] = value
			}
			return ""
		})
		return result, loc, true
	}
	return nil, nil, false
}

// StrToTime func
func StrToTime(target interface{}) (time.Time, error) {
	switch t := target.(type) {
	case int:
	case float64:
	case string:
		var lasts FormatResult
		timeFormat := t
		startIndex := 0
		total := len(t)
		if result, loc, ok := matchDateFormat(t); ok {
			if loc[0] != startIndex {
				return time.Time{}, fmt.Errorf("wrong start of date")
			}
			lasts = result
			if loc[1] < total {
				timeFormat = t[loc[1]:]
				startIndex = 1
			}
		}
		if result, loc, ok := matchTimeFormat(timeFormat); ok {
			if loc[0] != startIndex || loc[1] != len(timeFormat) {
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
		if lasts != nil {
			return makeFormatDateTime(lasts)
		}
	}
	return time.Time{}, fmt.Errorf("wrong ")
}
func noEmptyField(target FormatResult, args ...string) string {
	for _, field := range args {
		if cur, ok := target[field]; ok && cur != "" {
			return cur
		}
	}
	return ""
}

func makeFormatDateTime(result FormatResult) (time.Time, error) {
	now := time.Now()
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
	fmt.Println("year", year)
	return time.Time{}, nil
}

/*
* https://gist.github.com/elliotchance/d419395aa776d632d897
 */
func replaceWith(re *regexp.Regexp, str string, repl func(args ...string) string) string {
	result := ""
	lastIndex := 0
	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}
		result += str[lastIndex:v[0]] + repl(groups...)
		lastIndex = v[1]
	}
	return result + str[lastIndex:]
}

func makePatterns(t string, rules ...string) (*PatternInfo, error) {
	if formatList, ok := allFormats[t]; ok {
		regRule, _ := regexp.Compile("\\$\\{[A-Za-z]+}")
		ptns := []*Pattern{}
		for _, rule := range rules {
			pattern := new(Pattern)
			pattern.Type = t
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
			"${dd}[.\t-]${mm}[.-]${YY}",
			"${dd}[.\t]${mm}\\.${yy}",
			"(?i)${dd}[ \t.-]*${m}[ \t.-]*${y}",
			"(?i)${m}[ \t.-]*${YY}",
			"(?i)${YY}[ \t.-]*${m}",
			"(?i)${m}[ .\t-]*${dd}[,.stndrh\t ]+${y}",
			"(?i)${m}[ .\t-]*${dd}[,.stndrh\t ]*",
			"(?i)${d}[ .\t-]*${m}",
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
func matchTimeFormat(target string) (FormatResult, []int, bool) {
	var patterns []*Pattern
	if info, ok := allPatternInfo["time"]; ok {
		patterns = info.Patterns
	} else {
		info, _ := makePatterns("time",
			"(?i)t?${HH}[.:]${MN}[.:]${II}",
			"(?i)${hh}[.:]${MN}[.:]${II}${space}?${meridian}",
			"(?i)${hh}${space}?${meridian}",
			"(?i)${hh}[.:]${MN}${space}?${meridian}",
			"(?i)${hh}:${MN}:${II}[.:][0-9]+${meridian}",
			"(?i)t?${HH}[.:]${MN}",
			"(?i)t?${HH}${MN}",
			"(?i)t?${HH}${MN}${II}",
			"(?i)t?${HH}[.:]${MN}[.:]${II}${space}?(?:${tzcorrection}|${tz})",
			"(?i)t?${HH}[.:]${MN}[.:]${II}${frac}",
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
		if cur, err := StrToTime(target); err == nil {
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
