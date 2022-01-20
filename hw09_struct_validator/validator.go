package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrRegexp      = errors.New("regexp")
	ErrLen         = errors.New("len")
	ErrIn          = errors.New("in")
	ErrMin         = errors.New("min")
	ErrMax         = errors.New("max")
	ErrNotStruct   = errors.New("not a struct")
	ErrUnknownTag  = errors.New("unknown tag")
	ErrUnknownType = errors.New("unknown type")
)

func (v ValidationErrors) Error() string {
	var res strings.Builder
	for _, val := range v {
		res.WriteString(fmt.Sprintf("%v: %v\n", val.Field, val.Err))
	}
	return res.String()
}

func ValidateIntMin(ve ValidationErrors, cur int, iMin string, fieldName string) (ValidationErrors, error) {
	min, err := strconv.Atoi(iMin)
	if err != nil {
		return ve, err
	}
	if cur < min {
		ve = append(ve, ValidationError{
			Field: fieldName,
			Err:   ErrMin,
		})
	}
	return ve, nil
}

func ValidateIntMax(ve ValidationErrors, cur int, iMax string, fieldName string) (ValidationErrors, error) {
	max, err := strconv.Atoi(iMax)
	if err != nil {
		return ve, err
	}
	if cur > max {
		ve = append(ve, ValidationError{
			Field: fieldName,
			Err:   ErrMax,
		})
	}
	return ve, nil
}

func ValidateIntIn(ve ValidationErrors, cur int, in string, fieldName string) (ValidationErrors, error) {
	for _, iVal := range strings.Split(in, ",") {
		val, err := strconv.Atoi(iVal)
		if err != nil {
			return ve, err
		}
		if val == cur {
			return ve, nil
		}
	}
	ve = append(ve, ValidationError{
		Field: fieldName,
		Err:   ErrIn,
	})
	return ve, nil
}

func ValidateStringRegexp(ve ValidationErrors, cur string, re string, fieldName string) (ValidationErrors, error) {
	ok, err := regexp.Match(re, []byte(cur))
	if err != nil {
		return ve, err
	}
	if !ok {
		ve = append(ve, ValidationError{
			Field: fieldName,
			Err:   ErrRegexp,
		})
	}
	return ve, nil
}

func ValidateStringLen(ve ValidationErrors, cur string, sLen string, fieldName string) (ValidationErrors, error) {
	counter := utf8.RuneCount([]byte(cur))
	iLen, err := strconv.Atoi(sLen)
	if err != nil {
		return ve, err
	}
	if counter != iLen {
		ve = append(ve, ValidationError{
			Field: fieldName,
			Err:   ErrLen,
		})
	}
	return ve, nil
}

func ValidateStringIn(ve ValidationErrors, cur string, in string, fieldName string) (ValidationErrors, error) {
	for _, val := range strings.Split(in, ",") {
		if val == cur {
			return ve, nil
		}
	}
	ve = append(ve, ValidationError{
		Field: fieldName,
		Err:   ErrIn,
	})
	return ve, nil
}

func ValidateInt(ve ValidationErrors, fieldName string, tags string, value reflect.Value) (ValidationErrors, error) {
	var err error
	for _, tag := range strings.Split(tags, "|") {
		tagValue := strings.Split(tag, ":")
		curValue := int(value.Int())
		switch tagValue[0] {
		case "min":
			ve, err = ValidateIntMin(ve, curValue, tagValue[1], fieldName)
		case "max":
			ve, err = ValidateIntMax(ve, curValue, tagValue[1], fieldName)
		case "in":
			ve, err = ValidateIntIn(ve, curValue, tagValue[1], fieldName)
		default:
			err = ErrUnknownTag
		}
		if err != nil {
			return ve, err
		}
	}
	return ve, nil
}

func ValidateString(ve ValidationErrors, fieldName string, tags string, value reflect.Value) (ValidationErrors, error) {
	var err error
	for _, tag := range strings.Split(tags, "|") {
		tagValue := strings.Split(tag, ":")
		curValue := value.String()
		switch tagValue[0] {
		case "regexp":
			ve, err = ValidateStringRegexp(ve, curValue, tagValue[1], fieldName)
		case "len":
			ve, err = ValidateStringLen(ve, curValue, tagValue[1], fieldName)
		case "in":
			ve, err = ValidateStringIn(ve, curValue, tagValue[1], fieldName)
		default:
			err = ErrUnknownTag
		}
		if err != nil {
			return ve, err
		}
	}
	return ve, nil
}

func ValidatorSwitch(ve ValidationErrors, structField reflect.StructField,
	value reflect.Value) (ValidationErrors, error) {
	var err error
	fieldName := structField.Name
	tag := structField.Tag.Get("validate")
	if tag == "" && value.Kind() != reflect.Struct {
		return ve, nil
	}
	switch value.Kind() { //nolint:exhaustive
	case reflect.Struct:
		ve, err = Validate(ve, value.Interface())
	case reflect.Slice:
		sl := value.Slice(0, value.Len())
		switch value.Type().String() {
		case "[]string":
			for i := 0; i < value.Len(); i++ {
				ve, err = ValidateString(ve, fieldName+fmt.Sprintf("[%d]", i), tag, sl.Index(i))
				if err != nil {
					return ve, err
				}
			}
		case "[]int":
			for i := 0; i < value.Len(); i++ {
				ve, err = ValidateInt(ve, fieldName+fmt.Sprintf("[%d]", i), tag, sl.Index(i))
				if err != nil {
					return ve, err
				}
			}
		default:
			err = ErrUnknownType
		}
	case reflect.Int:
		ve, err = ValidateInt(ve, fieldName, tag, value)
	case reflect.String:
		ve, err = ValidateString(ve, fieldName, tag, value)
	}
	return ve, err
}

func Validate(ve ValidationErrors, iv interface{}) (ValidationErrors, error) {
	var err error
	v := reflect.ValueOf(iv)
	if v.Kind() != reflect.Struct {
		return ve, ErrNotStruct
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i)
		value := v.Field(i)
		ve, err = ValidatorSwitch(ve, structField, value)
		if err != nil {
			break
		}
	}
	fmt.Println(ve.Error())
	return ve, err
}
