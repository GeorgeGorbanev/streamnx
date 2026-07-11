package duration

import (
	"math"
	"regexp"
	"strconv"
)

var iso8601Re = regexp.MustCompile(`^P(?:(\d+)D)?T?(?:(\d+)H)?(?:(\d+)M)?(?:(\d+(?:\.\d+)?)S)?$`)

func ISO8601ToSeconds(raw string) int {
	matches := iso8601Re.FindStringSubmatch(raw)
	if len(matches) == 0 {
		return 0
	}

	days := atoi(matches[1])
	hours := atoi(matches[2])
	minutes := atoi(matches[3])
	seconds := atof(matches[4])

	return int(math.Round(float64(days*24*60*60+hours*60*60+minutes*60) + seconds))
}

func MsToSeconds(milliseconds int) int {
	if milliseconds <= 0 {
		return 0
	}
	return (milliseconds + 500) / 1000
}

func atoi(raw string) int {
	if raw == "" {
		return 0
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return value
}

func atof(raw string) float64 {
	if raw == "" {
		return 0
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0
	}
	return value
}
