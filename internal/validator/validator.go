package validator

import (
	"regexp"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	v.Errors[key] = message
}

func (v *Validator) Check(ok bool, key, message string) {
	if ok {
		v.AddError(key, message)
	}
}

func Out(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return false
		}
	}

	return true
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unmatches(value string, rx *regexp.Regexp) bool {
	return !rx.MatchString(value)
}

func HasDuplicate(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) != len(uniqueValues)
}
