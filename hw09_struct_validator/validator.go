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
	var errs []string
	for _, e := range v {
		errs = append(errs, fmt.Sprintf("%s: %v", e.Field, e.Err))
	}
	return strings.Join(errs, ", ")
}

type validateTag struct {
	Key   string
	Value string
}

func Validate(v interface{}) error {
	var errs ValidationErrors

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		rules := parseTag(tag)
		fieldName := field.Name

		for _, rule := range rules {
			err := applyRule(fieldVal, rule)
			if err != nil {
				errs = append(errs, ValidationError{Field: fieldName, Err: err})
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// parseTag parses the validate tag into individual validation rules.
func parseTag(tag string) []validateTag {
	rules := strings.Split(tag, "|")
	var tags []validateTag

	for _, rule := range rules {
		parts := strings.Split(rule, ":")
		if len(parts) == 2 {
			tags = append(tags, validateTag{Key: parts[0], Value: parts[1]})
		}
	}

	return tags
}

// applyRule applies a single validation rule to a field.
func applyRule(fieldVal reflect.Value, rule validateTag) error {
	switch fieldVal.Kind() {
	case reflect.String:
		return validateString(fieldVal.String(), rule)
	case reflect.Int:
		return validateInt(int(fieldVal.Int()), rule)
	case reflect.Slice:
		return validateSlice(fieldVal, rule)
	case reflect.Struct:
		if rule.Key == "nested" {
			return Validate(fieldVal.Interface())
		}
	}
	return nil
}

// validateString validates a string based on a rule.
func validateString(value string, rule validateTag) error {
	switch rule.Key {
	case "len":
		expectedLen, _ := strconv.Atoi(rule.Value)
		if len(value) != expectedLen {
			return fmt.Errorf("must be %d characters long", expectedLen)
		}
	case "regexp":
		re, err := regexp.Compile(rule.Value)
		if err != nil {
			return err
		}
		if !re.MatchString(value) {
			return fmt.Errorf("must match regexp %s", rule.Value)
		}
	case "in":
		validValues := strings.Split(rule.Value, ",")
		for _, v := range validValues {
			if value == v {
				return nil
			}
		}
		return fmt.Errorf("must be one of %v", validValues)
	}
	return nil
}

// validateInt validates an int based on a rule.
func validateInt(value int, rule validateTag) error {
	switch rule.Key {
	case "min":
		minVal, _ := strconv.Atoi(rule.Value)
		if value < minVal {
			return fmt.Errorf("must be at least %d", minVal)
		}
	case "max":
		maxVal, _ := strconv.Atoi(rule.Value)
		if value > maxVal {
			return fmt.Errorf("must be at most %d", maxVal)
		}
	case "in":
		validValues := strings.Split(rule.Value, ",")
		for _, v := range validValues {
			intVal, _ := strconv.Atoi(v)
			if value == intVal {
				return nil
			}
		}
		return fmt.Errorf("must be one of %v", validValues)
	}
	return nil
}

// validateSlice validates each element of a slice based on a rule.
func validateSlice(value reflect.Value, rule validateTag) error {
	for i := 0; i < value.Len(); i++ {
		err := applyRule(value.Index(i), rule)
		if err != nil {
			return fmt.Errorf("element %d: %v", i, err)
		}
	}
	return nil
}

// Example usage
// type User struct {
// 	Name   string `validate:"len:32|regexp:\\d+"`
// 	Age    int    `validate:"min:18|max:65"`
// 	Scores []int  `validate:"min:0|max:100"`
// 	Meta   Meta   `validate:"nested"`
// }

type Meta struct {
	Description string `validate:"len:10|in:foo,bar,baz"`
}
