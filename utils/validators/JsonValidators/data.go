package jsonV

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)


// Input represents form input values and validations
type JInput struct {
	Values  JsonInput
	VErrors ValidationErrors

}

type JsonInput map[string]interface{}

func (ji JsonInput) Get(field string) interface{} {
	ves,ok := ji[field]
	if !ok {
		return ""
	}
	return ves
}

// MinLength checks if a given minium length is satisfied
func (inVal *JInput) MinLength(field string, d int) {
	value := inVal.Values.Get(field)

	if value == "" {
		return
	}

	if utf8.RuneCountInString(value.(string)) < d {
		inVal.VErrors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

// Required checks if list of provided form input fields have values
func (inVal *JInput) Required(fields ...string) {
	for _, f := range fields {
		value := inVal.Values.Get(f)
		if value == "" {
			inVal.VErrors.Add(f, "This field is required field")
		}
	}
}

// MatchesPattern checks if a given input form field matchs a given pattern
func (inVal *JInput) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := inVal.Values.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value.(string)) {
		inVal.VErrors.Add(field, "The value entered is invalid")
	}
}

// PasswordMatches checks if Password and Confirm Password fields match
func (inVal *JInput) PasswordMatches(password string, confPassword string) {
	pwd := inVal.Values.Get(password)
	confPwd := inVal.Values.Get(confPassword)

	if pwd == "" || confPwd == "" {
		return
	}

	if pwd != confPwd {
		inVal.VErrors.Add(password, "The Password and Confim Password values did not match")
		inVal.VErrors.Add(confPassword, "The Password and Confim Password values did not match")
	}
}

// Valid checks if any form input validation has failed or not
func (inVal *JInput) Valid() bool {
	return len(inVal.VErrors) == 0
}
