package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrLen            = errors.New("length")
	ErrRegex          = errors.New("regex")
	ErrMin            = errors.New("greater")
	ErrMax            = errors.New("less")
	ErrIn             = errors.New("lots of")
	ErrInvalidValues  = errors.New("invalid values")
	ErrExpectedStruct = errors.New("expected a struct")
)

func (v ValidationErrors) Error() string {
	errStrs := make([]string, 0)

	for _, err := range v {
		errStrs = append(errStrs, fmt.Sprintf("%s %s", err.Field, err.Err))
	}

	return strings.Join(errStrs, ", ")
}

func checkLen(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.String {
		intValue, err := strconv.Atoi(ruleValue)
		if err != nil {
			return false
		}

		return rv.Len() == intValue
	}

	return false
}

func checkRegex(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.String {
		rx, err := regexp.Compile(ruleValue)
		if err != nil {
			return false
		}

		return rx.Match([]byte(rv.String()))
	}

	return false
}

func checkMin(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.Int {
		intValue := int(rv.Int())
		min, err := strconv.Atoi(ruleValue)
		if err != nil {
			return false
		}

		return intValue > min
	}

	return false
}

func checkMax(rv reflect.Value, ruleValue string) bool {
	if rv.Kind() == reflect.Int {
		intValue := int(rv.Int())
		max, err := strconv.Atoi(ruleValue)
		if err != nil {
			return false
		}

		return intValue < max
	}

	return false
}

func checkIn(rv reflect.Value, ruleValue string) bool {
	ins := strings.Split(ruleValue, ",")
	isValid := false

	switch rv.Kind() {
	case reflect.Int:
		intValue := int(rv.Int())

		for _, in := range ins {
			in, err := strconv.Atoi(in)
			if err != nil {
				continue
			}

			if in == intValue {
				isValid = true
			}
		}
	case reflect.String:
		strValue := rv.String()

		for _, in := range ins {
			if in == strValue {
				isValid = true
			}
		}
	}

	return isValid
}

func validateValue(validateTag string, rv reflect.Value) []error {
	rules := strings.Split(validateTag, "|")
	errs := make([]error, 0)

	for _, rule := range rules {
		r := strings.Split(rule, ":")
		if len(r) != 2 {
			continue
		}

		rType, rValue := r[0], r[1]

		var err error

		switch rType {
		case "len":
			if !checkLen(rv, rValue) {
				err = fmt.Errorf("value %v %w must be %s", rv.Interface(), ErrLen, rValue)
			}
		case "regexp":
			if !checkRegex(rv, rValue) {
				err = fmt.Errorf("value %v must match %w %s", rv.Interface(), ErrRegex, rValue)
			}
		case "min":
			if !checkMin(rv, rValue) {
				err = fmt.Errorf("value %v must be %w than %s", rv.Interface(), ErrMin, rValue)
			}
		case "max":
			if !checkMax(rv, rValue) {
				err = fmt.Errorf("value %v must be %w than %s", rv.Interface(), ErrMax, rValue)
			}
		case "in":
			if !checkIn(rv, rValue) {
				err = fmt.Errorf("value %v must be %w %s", rv.Interface(), ErrIn, rValue)
			}
		default:
			continue
		}

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func checkValue(valErrs *ValidationErrors, fName string, validateTag string, rv reflect.Value) {
	var errs []error

	switch rv.Kind() {
	case reflect.String:
		errs = validateValue(validateTag, rv)
	case reflect.Int:
		errs = validateValue(validateTag, rv)
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			checkValue(valErrs, fName, validateTag, rv.Index(i))
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			*valErrs = append(*valErrs, ValidationError{fName, err})
		}
	}
}

func Validate(v interface{}) (ValidationErrors, error) {
	errs := make(ValidationErrors, 0)

	iv := reflect.ValueOf(v)
	if iv.Kind() != reflect.Struct {
		return nil, ErrExpectedStruct
	}

	t := iv.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i) // reflect.StructField
		fv := iv.Field(i)   // reflect.Value

		validateTag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		checkValue(&errs, field.Name, validateTag, fv)
	}

	if len(errs) > 0 {
		return errs, ErrInvalidValues
	}

	return nil, nil
}
