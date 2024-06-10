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

type SoftwareError struct {
	Message string
}

func (e SoftwareError) Error() string {
	return e.Message
}

func Validate(v interface{}) error {
	var validationErrs ValidationErrors

	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	valueType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)
		fieldVal := value.Field(i)

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		rules, err := parseTag(tag)
		if err != nil {
			return fmt.Errorf("error parsing tag on field %s: %w", field.Name, err)
		}
		fieldName := field.Name

		for _, rule := range rules {
			err := applyRule(fieldVal, rule)
			if err != nil {
				if _, ok := err.(SoftwareError); ok {
					return err
				}
				validationErrs = append(validationErrs, ValidationError{Field: fieldName, Err: err})
			}
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

// parseTag parses the validate tag into individual validation rules.
func parseTag(tag string) ([]validateTag, error) {
	rules := strings.Split(tag, "|")
	var tags []validateTag

	for _, rule := range rules {
		parts := strings.Split(rule, ":")
		if len(parts) == 0 || (len(parts) != 2 && parts[0] != "nested") {
			return nil, fmt.Errorf("invalid tag format: %s", rule)
		}
		key := parts[0]
		var value string
		if parts[0] == "nested" {
			value = parts[0]
		} else {
			value = parts[1]
		}
		tags = append(tags, validateTag{Key: key, Value: value})
	}

	return tags, nil
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
		expectedLen, err := strconv.Atoi(rule.Value)
		if err != nil {
			return SoftwareError{Message: err.Error()}
		}
		if len(value) != expectedLen {
			return fmt.Errorf("must be %d characters long", expectedLen)
		}
	case "regexp":
		re, err := regexp.Compile(rule.Value)
		if err != nil {
			return SoftwareError{Message: err.Error()}
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
		minVal, err := strconv.Atoi(rule.Value)
		if err != nil {
			return SoftwareError{Message: err.Error()}
		}
		if value < minVal {
			return fmt.Errorf("must be at least %d", minVal)
		}
	case "max":
		maxVal, err := strconv.Atoi(rule.Value)
		if err != nil {
			return SoftwareError{Message: err.Error()}
		}
		if value > maxVal {
			return fmt.Errorf("must be at most %d", maxVal)
		}
	case "in":
		validValues := strings.Split(rule.Value, ",")
		for _, v := range validValues {
			intVal, err := strconv.Atoi(v)
			if err != nil {
				return SoftwareError{Message: err.Error()}
			}
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
			if _, ok := err.(SoftwareError); ok {
				return err
			}
			return fmt.Errorf("element %d: %v", i, err)
		}
	}
	return nil
}
