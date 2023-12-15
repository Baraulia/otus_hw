package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrorNotInSlice                    = errors.New("value is not in slice")
	ErrorDoesNotMatchSpecifiedLength   = errors.New("value does not match the specified length")
	ErrorDoesNotMatchRegularExpression = errors.New("value does not match the given regular expression")
	ErrorLessThanMin                   = errors.New("value is less than minimum")
	ErrorMoreThanMax                   = errors.New("value is more than maximum")
)

const tagName = "validate"

type Element interface {
	~string | ~int
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	for _, err := range v {
		builder.WriteString(fmt.Sprintf("FieldName: %s, error: %s;", err.Field, err.Err))
	}

	return builder.String()
}

func ValidateStruct(str interface{}) error {
	var validationErrors ValidationErrors
	value := reflect.ValueOf(str)

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		fieldValue := value.Field(i)

		if field.Type.Kind() == reflect.Struct {
			err := ValidateStruct(fieldValue.Interface())
			if err != nil {
				var validationErr ValidationErrors
				if errors.As(err, &validationErr) {
					validationErrors = append(validationErrors, validationErr...)
				} else {
					fmt.Printf("invalid error: %s\n", err)
				}
			}
		}

		tag := value.Type().Field(i).Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}

		splitTags := strings.Split(tag, "|")
		for _, t := range splitTags {
			err := ValidateField(fieldValue, t)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: field.Name,
					Err:   err,
				})
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func ValidateField(fieldValue reflect.Value, tag string) error {
	args := strings.Split(tag, ":")
	switch args[0] {
	case "len":
		specifiedLen, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("unsupported type for len tag value")
		}
		return ValidateLen(fieldValue, specifiedLen)
	case "regexp":
		return ValidateRegexp(fieldValue, args[1])
	case "in":
		sliceElements := strings.Split(args[1], ",")
		return ValidateIncluding(fieldValue, sliceElements)
	case "min":
		specifiedMin, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("unsupported type for min value: %w", err)
		}
		return ValidateMin(fieldValue, specifiedMin)
	case "max":
		specifiedMax, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("unsupported type for max value: %w", err)
		}
		return ValidateMax(fieldValue, specifiedMax)
	default:
		return fmt.Errorf("unsupported tag: %s", args[0])
	}
}

func ValidateLen(val reflect.Value, specLen int) error {
	//nolint:exhaustive
	switch val.Kind() {
	case reflect.String:
		if val.Len() != specLen {
			return ErrorDoesNotMatchSpecifiedLength
		}
		return nil
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			if elem.Kind() != reflect.String {
				return fmt.Errorf("unsupported type of elements in the slice: %s", elem.Kind().String())
			}
			if elem.Len() != specLen {
				return ErrorDoesNotMatchSpecifiedLength
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported type of field for validation: %s", val.Kind().String())
	}
}

func ValidateRegexp(val reflect.Value, reg string) error {
	//nolint:exhaustive
	switch val.Kind() {
	case reflect.String:
		r, err := regexp.Compile(reg)
		if err != nil {
			return fmt.Errorf("it is not a regular expresion in tag: %s", reg)
		}
		if !r.MatchString(val.String()) {
			return ErrorDoesNotMatchRegularExpression
		}
		return nil
	default:
		return fmt.Errorf("unsupported type of field for validation: %s", val.Kind().String())
	}
}

func ValidateIncluding(val reflect.Value, slice []string) error {
	var valueString string

	//nolint:exhaustive
	switch val.Kind() {
	case reflect.Int:
		valueString = strconv.Itoa(int(val.Int()))
	case reflect.String:
		valueString = val.String()
	default:
		return fmt.Errorf("unsupported type of value: %s", val.Kind().String())
	}

	for _, s := range slice {
		if s == valueString {
			return nil
		}
	}

	return ErrorNotInSlice
}

func ValidateMin(val reflect.Value, min int) error {
	//nolint:exhaustive
	switch val.Kind() {
	case reflect.Int:
		if val.Int() < int64(min) {
			return ErrorLessThanMin
		}
		return nil
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			if elem.Kind() != reflect.Int {
				return fmt.Errorf("unsupported type of elements in the slice: %s", elem.Kind().String())
			}
			if elem.Int() < int64(min) {
				return ErrorLessThanMin
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported type of field for validation: %s", val.Kind().String())
	}
}

func ValidateMax(val reflect.Value, max int) error {
	//nolint:exhaustive
	switch val.Kind() {
	case reflect.Int:
		if val.Int() > int64(max) {
			return ErrorMoreThanMax
		}
		return nil
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			if elem.Kind() == reflect.Int && elem.Int() > int64(max) {
				return ErrorMoreThanMax
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported type of field for validation: %s", val.Kind().String())
	}
}
