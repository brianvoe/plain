package plain

import (
	"reflect"
	"testing"
)

type TestData struct {
	Name    string  `plain:"name"`
	Age     int     `plain:"age"`
	Active  bool    `plain:"active"`
	Balance float64 `plain:"balance"`
}

type TestDataSub struct {
	Name string   `plain:"name"`
	Age  int      `plain:"age"`
	Sub  TestData `plain:"sub"`
}

func TestUnmarshal(t *testing.T) {
	// Test case for unmarshaling into struct
	t.Run("Unmarshal into struct", func(t *testing.T) {
		data := []byte("name: John Doe\nage: 30\nactive: true\nbalance: 123.45")
		var result TestData
		err := Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal into struct: %v", err)
		}
		expected := TestData{"John Doe", 30, true, 123.45}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Unmarshal into struct: got %v, want %v", result, expected)
		}
	})

	// Test case for unmarshaling into struct with sub-struct
	t.Run("Unmarshal into struct with sub-struct", func(t *testing.T) {
		data := []byte("name: John Doe\nage: 30\nsub.name: Jane Doe\nsub.age: 25\nsub.active: true\nsub.balance: 123.45")
		var result TestDataSub
		err := Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal into struct with sub-struct: %v", err)
		}
		expected := TestDataSub{"John Doe", 30, TestData{"Jane Doe", 25, true, 123.45}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Unmarshal into struct with sub-struct: got %v, want %v", result, expected)
		}
	})

	// Test case for unmarshaling into an array of structs
	t.Run("Unmarshal into an array of structs", func(t *testing.T) {
		data := []byte("name: John Doe\nage: 30\nactive: true\nbalance: 123.45\n\nname: Jane Doe\nage: 25\nactive: false\nbalance: 543.21")
		var result []TestData
		err := Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal into an array of structs: %v", err)
		}
		expected := []TestData{
			{"John Doe", 30, true, 123.45},
			{"Jane Doe", 25, false, 543.21},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Unmarshal into an array of structs: got %v, want %v", result, expected)
		}
	})

	// Test case for unmarshaling into string
	t.Run("Unmarshal into string", func(t *testing.T) {
		data := "sample string"
		var result string
		err := Unmarshal([]byte(data), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal into string: %v", err)
		}
		if result != data {
			t.Errorf("Unmarshal into string: got %v, want %v", result, data)
		}
	})

	// Add more test cases for int, float64, and bool...
	// Test case for unmarshaling into int
	t.Run("Unmarshal into int", func(t *testing.T) {
		data := "42"
		var result int
		err := Unmarshal([]byte(data), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal into int: %v", err)
		}
		if result != 42 {
			t.Errorf("Unmarshal into int: got %v, want 42", result)
		}
	})

	// Test case for unmarshaling into float64
	t.Run("Unmarshal into float64", func(t *testing.T) {
		data := "123.45"
		var result float64
		err := Unmarshal([]byte(data), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal into float64: %v", err)
		}
		if result != 123.45 {
			t.Errorf("Unmarshal into float64: got %v, want 123.45", result)
		}
	})

	// Test case for unmarshaling into bool
	t.Run("Unmarshal into bool", func(t *testing.T) {
		data := "true"
		var result bool
		err := Unmarshal([]byte(data), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal into bool: %v", err)
		}
		if result != true {
			t.Errorf("Unmarshal into bool: got %v, want true", result)
		}
	})
}

type TestUnmarshalerString string

func (t *TestUnmarshalerString) UnmarshalPlain(data []byte) error {
	*t = "test"
	return nil
}

type TestUnmarshalerInt int

func (t *TestUnmarshalerInt) UnmarshalPlain(data []byte) error {
	*t = 42
	return nil
}

type TestUnmarshalerFloat float64

func (t *TestUnmarshalerFloat) UnmarshalPlain(data []byte) error {
	*t = 123.45
	return nil
}

type TestUnmarshalerBool bool

func (t *TestUnmarshalerBool) UnmarshalPlain(data []byte) error {
	*t = true
	return nil
}

type TestUnmarshalerStruct struct {
	Name     string `plain:"name"`
	Age      int    `plain:"age"`
	IsActive bool   `plain:"is_active"`
}

func (t *TestUnmarshalerStruct) UnmarshalPlain(data []byte) error {
	t.Name = "John Doe"
	t.Age = 30
	t.IsActive = true
	return nil
}

func TestUnmarshaler(t *testing.T) {
	// String
	t.Run("Unmarshaler string", func(t *testing.T) {
		var result TestUnmarshalerString
		err := Unmarshal([]byte(""), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal string: %v", err)
		}
		if result != "test" {
			t.Errorf("Unmarshal string: got %v, want test", result)
		}
	})

	// Add more test cases for int, float64, and bool...
	// Int
	t.Run("Unmarshaler int", func(t *testing.T) {
		var result TestUnmarshalerInt
		err := Unmarshal([]byte(""), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal int: %v", err)
		}
		if result != 42 {
			t.Errorf("Unmarshal int: got %v, want 42", result)
		}
	})

	// Float64
	t.Run("Unmarshaler float64", func(t *testing.T) {
		var result TestUnmarshalerFloat
		err := Unmarshal([]byte(""), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal float64: %v", err)
		}
		if result != 123.45 {
			t.Errorf("Unmarshal float64: got %v, want 123.45", result)
		}
	})

	// Bool
	t.Run("Unmarshaler bool", func(t *testing.T) {
		var result TestUnmarshalerBool
		err := Unmarshal([]byte(""), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal bool: %v", err)
		}
		if result != true {
			t.Errorf("Unmarshal bool: got %v, want true", result)
		}
	})

	// Struct
	t.Run("Unmarshaler struct", func(t *testing.T) {
		var result TestUnmarshalerStruct
		err := Unmarshal([]byte(""), &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal struct: %v", err)
		}
		expected := TestUnmarshalerStruct{"John Doe", 30, true}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Unmarshal struct: got %v, want %v", result, expected)
		}
	})
}
