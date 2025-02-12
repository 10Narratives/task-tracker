package validation

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/10Narratives/task-tracker/internal/lib"
	"github.com/go-playground/validator/v10"
)

// IsDateValid checks if date field is in YYYYMMDD format
func IsDateValid(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	_, err := time.Parse(lib.DateFormat, date)
	return err == nil
}

// IsTitleValid checks if title is non-empty
func IsTitleValid(fl validator.FieldLevel) bool {
	title := fl.Field().String()
	return len(title) > 0
}

var (
	DailyRepeatPattern   string = "^d (?:[1-9]|[1-9][0-9]|[1-3][0-9]{2}|400)$"
	WeeklyRepeatPattern  string = "^w ([1-7](,[1-7])*)$"
	MonthlyRepeatPattern string = "^m (-?[1-9]|-1|-2|[12][0-9]|3[01])(,(-?[1-9]|-1|-2|[12][0-9]|3[01]))*( (1[0-2]|[1-9])(,(1[0-2]|[1-9]))*)?$"
	YearlyRepeatPattern  string = "^y$"
)

// IsRepeatValid checks if the repeat string is either empty or matches the expected patterns.
func IsRepeatValid(fl validator.FieldLevel) bool {
	repeat := fl.Field().String()
	if repeat == "" {
		return true
	}
	dailyMatch, _ := regexp.MatchString(DailyRepeatPattern, repeat)
	weeklyMatch, _ := regexp.MatchString(WeeklyRepeatPattern, repeat)
	monthlyMatch, _ := regexp.MatchString(MonthlyRepeatPattern, repeat)
	yearlyMatch, _ := regexp.MatchString(YearlyRepeatPattern, repeat)
	return dailyMatch || weeklyMatch || monthlyMatch || yearlyMatch
}

func ValidationErrorMsg(errs validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		case "dateformat":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be in YYYYMMDD date format", err.Field()))
		case "title":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be non-empty", err.Field()))
		case "repeat":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must satisfy expected patterns", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return strings.Join(errMsgs, ", ")
}
