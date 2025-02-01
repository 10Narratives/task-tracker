package nextdate

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
	NextDate(startDate string) (string, error)
}

func Validate(repeat string, rule string) error {
	re := regexp.MustCompile(rule)
	if !re.MatchString(repeat) {
		return ErrRepeatRuleIsNotValid
	}
	return nil
}

func StringToTime(str string) (time.Time, error) {
	parsedStr, err := time.Parse(DateLayout, str)
	if err != nil {
		return time.Time{}, errors.Join(ErrStartDateCanNotBeParsed, err)
	}
	return parsedStr, nil
}
