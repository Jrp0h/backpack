package utils

import (
	"fmt"
	"strings"
	"time"
)

func FormatDate(format string, t time.Time) string {
	result := strings.ReplaceAll(format, "%Y", fmt.Sprint(t.Year()))
	result = strings.ReplaceAll(result, "%m", fmt.Sprintf("%d", t.Month()))
	result = strings.ReplaceAll(result, "%d", fmt.Sprint(t.Day()))

	result = strings.ReplaceAll(result, "%H", fmt.Sprint(t.Hour()))
	result = strings.ReplaceAll(result, "%M", fmt.Sprint(t.Minute()))
	result = strings.ReplaceAll(result, "%S", fmt.Sprint(t.Second()))


	return result
}