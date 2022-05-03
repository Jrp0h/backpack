package utils

import (
	// "fmt"
	"strings"
	"time"
)

// Ruby to golang format
func FormatDate(format string, t time.Time) string {
	result := strings.ReplaceAll(format, "%a", "Mon")
	result = strings.ReplaceAll(result, "%A", "Monday")
	result = strings.ReplaceAll(result, "%d", "02")

	result = strings.ReplaceAll(result, "%b", "Jan")
	result = strings.ReplaceAll(result, "%B", "January")
	result = strings.ReplaceAll(result, "%m", "01")

	result = strings.ReplaceAll(result, "%Y", "2006")
	result = strings.ReplaceAll(result, "%y", "06")

	result = strings.ReplaceAll(result, "%H", "15")
	result = strings.ReplaceAll(result, "%M", "04")
	result = strings.ReplaceAll(result, "%S", "05")
	result = strings.ReplaceAll(result, "%Z", "MST")

	return t.Format(result)
}

// func pad2(num int) string {
// 	if num < 10 {
// 		return fmt.Sprintf("0%d", num)
// 	}

// 	return fmt.Sprint(num)
// }
