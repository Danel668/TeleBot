package utils

import(
	"strconv"
	"fmt"
	"strings"
	"errors"
	"time"
)

func StringToInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

func ToString(n int64) string {
	return fmt.Sprintf("%d", n)
}

func ToGoName(jsonName string) string {
	var result string
	parts := strings.Split(jsonName, "_")

	for _, part := range parts {
		if len(part) > 0 {
			result += strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return result
}

func ConvertTimeZone(s string) (string, error) {
	switch s {
	case "UTC-11":
		return "Pacific/Pago_Pago", nil
	case "UTC-10":
		return "Pacific/Honolulu", nil
	case "UTC-9":
		return "America/Anchorage", nil
	case "UTC-8":
		return "America/Los_Angeles", nil
	case "UTC-7":
		return "America/Denver", nil
	case "UTC-6":
		return "America/Chicago", nil
	case "UTC-5":
		return "America/New_York", nil
	case "UTC-4":
		return "America/Barbados", nil
	case "UTC-3":
		return "America/Argentina/Buenos_Aires", nil
	case "UTC-2":
		return "Atlantic/South_Georgia", nil
	case "UTC-1":
		return "Atlantic/Azores", nil
	case "UTC":
		return "Europe/London", nil
	case "UTC+1":
		return "Europe/Berlin", nil
	case "UTC+2":
		return "Europe/Kiev", nil
	case "UTC+3":
		return "Europe/Moscow", nil
	case "UTC+4":
		return "Asia/Dubai", nil
	case "UTC+5":
		return "Asia/Karachi", nil
	case "UTC+6":
		return "Asia/Dhaka", nil
	case "UTC+7":
		return "Asia/Bangkok", nil
	case "UTC+8":
		return "Asia/Shanghai", nil
	case "UTC+9":
		return "Asia/Tokyo", nil
	case "UTC+10":
		return "Australia/Sydney", nil
	case "UTC+11":
		return "Pacific/Guadalcanal", nil
	case "UTC+12":
		return "Pacific/Fiji", nil
	case "UTC+13":
		return "Pacific/Tongatapu", nil
	case "UTC+14":
		return "Pacific/Kiritimati", nil
	default:
		return "", errors.New("not valid timezone")
	}
}

func ParseHoursMinsString(s string, timezone string) (time.Time, error) {
	now := time.Now()

	parsedTime, err := time.Parse("15:04", s)
	if err != nil {
		return time.Time{}, err
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	resultTime := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		parsedTime.Hour(),
		parsedTime.Minute(),
		0,
		0,
		location,
	)

	return resultTime, nil
}

func ParseTimeInCommonFormat(s string, timezone string) (time.Time, error) {
	parsedTime, err := time.Parse("02.01.2006 15:04", s)
	if err != nil {
		return time.Time{}, err
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	resultTime := time.Date(
		parsedTime.Year(),
		parsedTime.Month(),
		parsedTime.Day(),
		parsedTime.Hour(),
		parsedTime.Minute(),
		0,
		0,
		location,
	)

	return resultTime, nil
}

func ParseDateInCommonFormat(s string) (time.Time, error) {
	parsedDate, err := time.Parse("02.01.2006", s)
	if err != nil {
		return time.Time{}, err
	}

	now := time.Now()
	resultTime := time.Date(
		parsedDate.Year(),
		parsedDate.Month(),
		parsedDate.Day(),
		now.Hour(),
		now.Minute(),
		0,
		0,
		now.Location(),
	)
	return resultTime, nil
}
