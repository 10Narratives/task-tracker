package nextdate

// TODO: Write docs

import (
	"errors"
	"regexp"
	"time"
)

var (
	DateLayout                 string = "20060102"
	ErrRepeatRuleIsNotValid    error  = errors.New("repeat rule is not valid")
	ErrStartDateCanNotBeParsed error  = errors.New("can not parse start date")
)

type DateIterator interface {
	// TODO: Change on Next(base time.Time) time.Time
	NextDate(startDate string) (string, error)
}

func Validate(repeat string, rule string) error {
	re := regexp.MustCompile(rule)
	if !re.MatchString(repeat) {
		return ErrRepeatRuleIsNotValid
	}
	return nil
}

// TODO: Refactor names in this function
func StringToTime(str string) (time.Time, error) {
	parsedStr, err := time.Parse(DateLayout, str)
	if err != nil {
		return time.Time{}, errors.Join(ErrStartDateCanNotBeParsed, err)
	}
	return parsedStr, nil
}
