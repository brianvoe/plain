package plain

import (
	"fmt"
	"reflect"
	"strings"
)

type Marshaler interface {
	MarshalPlain() ([]byte, error)
}

func Marshal(data any) ([]byte, error) {
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
			tag := field.Tag.Get("plain")
			if tag == "" {
				tag = field.Tag.Get("form")
			}

			if tag == "-" || tag == "" {
				continue
			}

			// If no tag, use field name
			if tag == "" {
				tag = field.Name
			}

			// Check if the field is exported
			if !val.Field(i).CanInterface() {
				continue
			}

			// Check if the field adhears to the Marshaler interface
			if m, ok := val.Field(i).Interface().(Marshaler); ok {
				// If it does then use that to marshal
				marshaled, err := m.MarshalPlain()
				if err != nil {
					return err
				}

				sb.WriteString(rowOutput(tag, string(marshaled)))
				continue
			}

			// Get name and value
			fieldName := tag
			if parent != "" {
				fieldName = parent + "." + tag
			}
			fieldValue := val.Field(i).Interface()

			// If the value is a struct, recurse
			if reflect.ValueOf(fieldValue).Kind() == reflect.Struct {
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
		// if the value is a slice, loop over its elements
		for i := 0; i < val.Len(); i++ {
			fieldValue := val.Index(i).Interface()
			fieldName := parent
			if parent != "" {
				fieldName = parent + "." + fmt.Sprintf("%d", i)
			}
			if reflect.ValueOf(fieldValue).Kind() == reflect.Struct {
				err := plainStruct(sb, reflect.ValueOf(fieldValue), fieldName)
				if err != nil {
					return err
				}
			} else {
				sb.WriteString(rowOutput(fieldName, fieldValue))
			}
		}

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
