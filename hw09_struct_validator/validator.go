package hw09structvalidator

import (
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

func (v ValidationErrors) Error() string {
	panic("implement me")
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

func validateValue(fName string, validateTag string, rv reflect.Value) {
	rules := strings.Split(validateTag, "|")

	for _, rule := range rules {
		r := strings.Split(rule, ":")
		if len(r) != 2 {
			continue
		}

		rType := r[0]
		rValue := r[1]

		switch rType {
		case "len":
			isValid := checkLen(rv, rValue)
			if !isValid {
				fmt.Println("not valid", fName, rValue, rv.Interface())
			}
		case "regexp":
			isValid := checkRegex(rv, rValue)
			if !isValid {
				fmt.Println("not valid", fName, rValue, rv.Interface())
			}
		case "min":
			isValid := checkMin(rv, rValue)
			if !isValid {
				fmt.Println("not valid", fName, rValue, rv.Interface())
			}
		case "max":
			isValid := checkMax(rv, rValue)
			if !isValid {
				fmt.Println("not valid", fName, rValue, rv.Interface())
			}
		case "in":
			isValid := checkIn(rv, rValue)
			if !isValid {
				fmt.Println("not valid", fName, rValue, rv.Interface())
			}
		}
	}
}

func checkValue(fName string, validateTag string, rv reflect.Value) {
	switch rv.Kind() {
	case reflect.String:
		validateValue(fName, validateTag, rv)
	case reflect.Int:
		validateValue(fName, validateTag, rv)
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			checkValue(fName, validateTag, rv.Index(i))
		}
	}
}

func Validate(v interface{}) error {
	// Place your code here.
	iv := reflect.ValueOf(v)
	if iv.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, but received %T", v)
	}

	t := iv.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i) // reflect.StructField
		fv := iv.Field(i)   // reflect.Value

		validateTag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		checkValue(field.Name, validateTag, fv)
	}

	return nil
}
