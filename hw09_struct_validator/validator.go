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
			validationErr, err := applyRule(fieldVal, rule)
			if err != nil {
				return err
			}
			if validationErr != nil {
				if _, ok := err.(ValidationErrors); ok {
					validationErrs = append(validationErrs, ValidationError{Field: fieldName, Err: err})
				}
			}
			// if validationErr != nil {
			// 	validationErrs = append(validationErrs, ValidationError{Field: fieldName, Err: validationErr})
			// }
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
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid tag format: %s", rule)
		}
		tags = append(tags, validateTag{Key: parts[0], Value: parts[1]})
	}

	return tags, nil
}

// applyRule applies a single validation rule to a field.
func applyRule(fieldVal reflect.Value, rule validateTag) (error, error) {
	switch fieldVal.Kind() {
	case reflect.String:
		return validateString(fieldVal.String(), rule)
	case reflect.Int:
		return validateInt(int(fieldVal.Int()), rule)
	case reflect.Slice:
		return validateSlice(fieldVal, rule)
	case reflect.Struct:
		if rule.Key == "nested" {
			return Validate(fieldVal.Interface()), nil
		}
	}
	return nil, nil
}

// validateString validates a string based on a rule.
func validateString(value string, rule validateTag) (error, error) {
	switch rule.Key {
	case "len":
		expectedLen, err := strconv.Atoi(rule.Value)
		if err != nil {
			return nil, err
		}
		if len(value) != expectedLen {
			return fmt.Errorf("must be %d characters long", expectedLen), nil
		}
	case "regexp":
		re, err := regexp.Compile(rule.Value)
		if err != nil {
			return nil, err
		}
		if !re.MatchString(value) {
			return fmt.Errorf("must match regexp %s", rule.Value), nil
		}
	case "in":
		validValues := strings.Split(rule.Value, ",")
		for _, v := range validValues {
			if value == v {
				return nil, nil
			}
		}
		return fmt.Errorf("must be one of %v", validValues), nil
	}
	return nil, nil
}

// validateInt validates an int based on a rule.
func validateInt(value int, rule validateTag) (error, error) {
	switch rule.Key {
	case "min":
		minVal, err := strconv.Atoi(rule.Value)
		if err != nil {
			return nil, err
		}
		if value < minVal {
			return fmt.Errorf("must be at least %d", minVal), nil
		}
	case "max":
		maxVal, err := strconv.Atoi(rule.Value)
		if err != nil {
			return nil, err
		}
		if value > maxVal {
			return fmt.Errorf("must be at most %d", maxVal), nil
		}
	case "in":
		validValues := strings.Split(rule.Value, ",")
		for _, v := range validValues {
			intVal, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			if value == intVal {
				return nil, nil
			}
		}
		return fmt.Errorf("must be one of %v", validValues), nil
	}
	return nil, nil
}

// validateSlice validates each element of a slice based on a rule.
func validateSlice(value reflect.Value, rule validateTag) (error, error) {
	for i := 0; i < value.Len(); i++ {
		validationErr, err := applyRule(value.Index(i), rule)
		if err != nil {
			return nil, err
		}
		if validationErr != nil {
			return fmt.Errorf("element %d: %v", i, err), nil
		}
	}
	return nil, nil
}
