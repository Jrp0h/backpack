package utils

import (
	"fmt"
	"strings"
	"time"
)

func FormatDate(format string, t time.Time) string {
	result := strings.ReplaceAll(format, "%Y", pad2(t.Year()))
	result = strings.ReplaceAll(result, "%m", pad2(int(t.Month())))
	result = strings.ReplaceAll(result, "%d", pad2(t.Day()))

	result = strings.ReplaceAll(result, "%H", pad2(t.Hour()))
	result = strings.ReplaceAll(result, "%M", pad2(t.Minute()))
	result = strings.ReplaceAll(result, "%S", pad2(t.Second()))


	return result
}

func pad2(num int) string {
	if num < 10 {
		return fmt.Sprintf("0%d", num)
	}

	return fmt.Sprint(num)
}