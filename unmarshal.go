package plain

import (
	"bytes"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Unmarshaler interface {
	UnmarshalPlain([]byte) error
}

// Unmarshal parses the plain text data and fills the provided target variable.
func Unmarshal(data []byte, v any) error {
	// Ensure v is a pointer
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("v must be a non-nil pointer")
	}

	// Dereference the pointer
	rv = rv.Elem()

	// Check if the target type implements Unmarshaler interface
	if unmarshaler, ok := v.(Unmarshaler); ok {
		return unmarshaler.UnmarshalPlain(data)
	}

	// Handle basic types (string, int, float64, bool)
	if isBasicType(rv.Kind()) {
		return unmarshalBasicType(data, rv)
	}

	// Handle slice of any type
	if rv.Kind() == reflect.Slice {
		return unmarshalSlice(data, rv)
	}

	// Handle struct
	if rv.Kind() == reflect.Struct {
		return unmarshalStruct(data, rv)
	}

	return errors.New("unsupported type for unmarshaling")
}

// unmarshalBasicType handles unmarshaling of basic data types.
func unmarshalBasicType(data []byte, v reflect.Value) error {
	trimmedData := strings.TrimSpace(string(data))

	switch v.Kind() {
	case reflect.String:
		v.SetString(trimmedData)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(trimmedData, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(intValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(trimmedData, 64)
		if err != nil {
			return err
		}
		v.SetFloat(floatValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(trimmedData)
		if err != nil {
			return err
		}
		v.SetBool(boolValue)
	default:
		return errors.New("unsupported basic type")
	}

	return nil
}

// unmarshalSlice handles unmarshaling of slice types.
func unmarshalSlice(data []byte, v reflect.Value) error {
	elementType := v.Type().Elem()

	// Split the data into separate elements by newline
	elementsData := bytes.Split(data, []byte("\n\n"))
	for _, elementData := range elementsData {
		trimmedData := strings.TrimSpace(string(elementData))

		// Check if the element is in array format
		if strings.HasPrefix(trimmedData, "[") && strings.HasSuffix(trimmedData, "]") {
			// Process as an array formatted string
			arrayContent := trimmedData[1 : len(trimmedData)-1]
			arrayElements := strings.Split(arrayContent, ",")
			for _, arrayElement := range arrayElements {
				trimmedElement := strings.TrimSpace(arrayElement)
				if err := processElement(trimmedElement, elementType, v); err != nil {
					return err
				}
			}
		} else {
			// Process as a single element
			if err := processElement(trimmedData, elementType, v); err != nil {
				return err
			}
		}
	}

	return nil
}

// processElement handles the creation and setting of a new element in the slice.
func processElement(elementData string, elementType reflect.Type, v reflect.Value) error {
	newElement := reflect.New(elementType).Elem()

	// Check if element is a struct or a basic type
	if elementType.Kind() == reflect.Struct {
		err := unmarshalStruct([]byte(elementData), newElement)
		if err != nil {
			return err
		}
	} else if isBasicType(elementType.Kind()) {
		err := unmarshalBasicType([]byte(elementData), newElement)
		if err != nil {
			return err
		}
	} else {
		return errors.New("unsupported slice element type")
	}

	v.Set(reflect.Append(v, newElement))
	return nil
}

// unmarshalStruct handles unmarshaling of struct types
func unmarshalStruct(data []byte, v reflect.Value) error {
	if v.Kind() != reflect.Struct {
		return errors.New("expected a struct type")
	}

	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		pair := bytes.SplitN(line, []byte(":"), 2)
		if len(pair) != 2 {
			continue // skip invalid lines
		}

		key := strings.TrimSpace(string(pair[0]))
		value := strings.TrimSpace(string(pair[1]))

		// Handle nested fields indicated by a dot separator
		if err := setFieldValue(v, key, value); err != nil {
			return err
		}
	}

	return nil
}

// setFieldValue sets the value of a field, handling nested structs.
func setFieldValue(v reflect.Value, key, value string) error {
	keys := strings.Split(key, ".")

	for i, k := range keys {
		if v.Kind() == reflect.Struct {
			found := false
			for j := 0; j < v.NumField(); j++ {
				field := v.Field(j)
				fieldType := v.Type().Field(j)
				tagValue := fieldType.Tag.Get("plain")
				if tagValue == "" {
					tagValue = fieldType.Tag.Get("form")
				}

				// Check if tag is "-" or empty
				if tagValue == "-" || tagValue == "" {
					continue
				}

				if strings.EqualFold(tagValue, k) {
					if i == len(keys)-1 {
						// Last key, set the value
						return setValue(field, value)
					} else if field.Kind() == reflect.Struct {
						// Nested struct, proceed to the next level
						v = field
						found = true
						break
					} else {
						return errors.New("non-struct field found in nested path: " + k)
					}
				}
			}
			if !found {
				// If no matching field is found, ignore and continue
				return nil
			}
		} else {
			return errors.New("attempted to navigate into non-struct field")
		}
	}

	return nil
}

// setValue sets the field with the provided value, handling type conversion.
func setValue(field reflect.Value, value string) error {
	if !field.CanSet() {
		return errors.New("cannot set field")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			field.SetInt(intValue)
		} else {
			return err
		}
	case reflect.Float32, reflect.Float64:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			field.SetFloat(floatValue)
		} else {
			return err
		}
	case reflect.Bool:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			field.SetBool(boolValue)
		} else {
			return err
		}
	case reflect.Slice:
		return unmarshalSlice([]byte(value), field)
	default:
		return errors.New("unsupported field type")
	}

	return nil
}

// isBasicType checks if the provided kind is a basic type.
func isBasicType(kind reflect.Kind) bool {
	switch kind {
	case reflect.String, reflect.Int, reflect.Float64, reflect.Bool:
		return true
	default:
		return false
	}
}
