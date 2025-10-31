package plain

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Marshaler interface {
	MarshalPlain() ([]byte, error)
}

func Marshal(data any) ([]byte, error) {
	// Handle []interface{} case
	if slice, ok := data.([]interface{}); ok && len(slice) > 0 {
		// Try to convert to typed slice using reflection
		// Get the type of the first element
		firstElem := slice[0]
		elemType := reflect.TypeOf(firstElem)

		// Create a new slice of the correct type
		newSlice := reflect.MakeSlice(reflect.SliceOf(elemType), len(slice), len(slice))
		for i, item := range slice {
			newSlice.Index(i).Set(reflect.ValueOf(item))
		}

		return Marshal(newSlice.Interface())
	}

	var sb strings.Builder
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			err := plainStruct(&sb, val.Index(i), "")
			if err != nil {
				return nil, err
			}

			sb.WriteString("\n")
		}
	} else {
		err := plainStruct(&sb, val, "")
		if err != nil {
			return nil, err
		}
	}

	buildStr := sb.String()

	// remove any trailing newlines
	for strings.HasSuffix(buildStr, "\n") {
		buildStr = strings.TrimSuffix(buildStr, "\n")
	}

	return []byte(buildStr), nil
}

func plainStruct(sb *strings.Builder, val reflect.Value, parent string) error {
	typ := val.Type()

	// switch on the type of the value
	switch val.Kind() {
	case reflect.Ptr:
		// if the value is a pointer, dereference it
		val = val.Elem()
		err := plainStruct(sb, val, parent)
		if err != nil {
			return err
		}

		return nil
	case reflect.Struct:
		// Check if struct has Marshaler interface
		if m, ok := val.Interface().(Marshaler); ok {
			// If it does then use that to marshal
			marshaled, err := m.MarshalPlain()
			if err != nil {
				return err
			}

			sb.WriteString(rowOutput(parent, string(marshaled)))
			return nil
		}

		// if the value is a struct, loop over its fields
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			fieldVal := val.Field(i)

			tag := field.Tag.Get("plain")
			if tag == "" {
				tag = field.Tag.Get("form")
			}

			// Parse the tag to check for omitempty
			tagName, omitEmpty := parseTag(tag)

			if tagName == "-" || tagName == "" {
				continue
			}

			// If no tag name, use field name
			if tagName == "" {
				tagName = field.Name
			}

			// Check if the field is exported
			if !fieldVal.CanInterface() {
				continue
			}

			// Check for omitempty and skip if empty
			if omitEmpty && isEmpty(fieldVal) {
				continue
			}

			// Check if the field adhears to the Marshaler interface
			if m, ok := val.Field(i).Interface().(Marshaler); ok {
				// If it does then use that to marshal
				marshaled, err := m.MarshalPlain()
				if err != nil {
					return err
				}

				sb.WriteString(rowOutput(tagName, string(marshaled)))
				continue
			}

			// Get name and value
			fieldName := tagName
			if parent != "" {
				fieldName = parent + "." + tagName
			}
			fieldValue := val.Field(i).Interface()

			// If the value is a struct check for time.Time, otherwise recurse
			if reflect.ValueOf(fieldValue).Kind() == reflect.Struct {
				// Check if the struct is a time.Time
				if t, ok := fieldValue.(time.Time); ok {
					sb.WriteString(rowOutput(fieldName, t.Format(time.RFC3339)))
					continue
				}

				err := plainStruct(sb, reflect.ValueOf(fieldValue), fieldName)
				if err != nil {
					return err
				}
				continue
			} else if reflect.ValueOf(fieldValue).Kind() == reflect.Slice {
				err := plainStruct(sb, reflect.ValueOf(fieldValue), fieldName)
				if err != nil {
					return err
				}
				continue
			}

			sb.WriteString(rowOutput(fieldName, fieldValue))
		}

		return nil
	case reflect.Slice:
		var sliceValues []string
		for i := 0; i < val.Len(); i++ {
			fieldValue := val.Index(i).Interface()

			// Check if the slice element is a struct or another slice, and process it accordingly
			if reflect.ValueOf(fieldValue).Kind() == reflect.Struct || reflect.ValueOf(fieldValue).Kind() == reflect.Slice {
				var nestedSB strings.Builder
				err := plainStruct(&nestedSB, reflect.ValueOf(fieldValue), "")
				if err != nil {
					return err
				}
				sliceValues = append(sliceValues, nestedSB.String())
			} else {
				// For simple types, just convert to string and append
				sliceValues = append(sliceValues, fmt.Sprintf("%v", fieldValue))
			}
		}

		// Join all slice values into a single string with the required format
		sliceStr := fmt.Sprintf("[%s]", strings.Join(sliceValues, ", "))
		if parent != "" {
			parent = parent + ": "
		}
		sb.WriteString(parent + sliceStr + "\n")

		return nil
	}

	// if the value is anything else, print it
	sb.WriteString(rowOutput(parent, val.Interface()))

	return nil
}

func rowOutput(field string, value any) string {
	if field == "" {
		return fmt.Sprintf("%v\n", value)
	}

	return fmt.Sprintf("%s: %v\n", field, value)
}

// parseTag parses the tag value and returns the name and whether omitempty is set.
func parseTag(tag string) (string, bool) {
	if tag == "" {
		return "", false
	}

	// Check for omitempty
	if strings.Contains(tag, ",omitempty") {
		return strings.TrimSuffix(tag, ",omitempty"), true
	}

	return tag, false
}

// isEmpty checks if a value is considered empty (zero value).
func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Struct:
		// For structs, check if it's time.Time and if it's zero
		if v.Type() == reflect.TypeOf(time.Time{}) {
			return v.Interface().(time.Time).IsZero()
		}
		// For other structs, check if all fields are zero
		for i := 0; i < v.NumField(); i++ {
			if !isEmpty(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
