package plain

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestPlain_Marshal(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		type TestSingle struct {
			Name string `form:"name"`
		}

		resp, err := Marshal(TestSingle{Name: "test"})
		if err != nil {
			t.Fatalf("Was not expecting an error got %s", err)
		}

		expected := "name: test"
		if string(resp) != expected {
			t.Fatalf("Single was expecting %s\n got %s", expected, string(resp))
		}
	})

	t.Run("Int", func(t *testing.T) {
		type TestMultiple struct {
			Name string `form:"name"`
			Age  int    `form:"age"`
		}

		resp, err := Marshal(TestMultiple{Name: "test", Age: 35})
		if err != nil {
			t.Fatalf("Was not expecting an error got %s", err)
		}

		expected := "name: test\nage: 35"
		if string(resp) != expected {
			t.Fatalf("Multiple was expecting %s\n got %s", expected, string(resp))
		}
	})

	t.Run("Float", func(t *testing.T) {
		type TestFloat struct {
			Age float64 `form:"age"`
		}

		resp, err := Marshal(TestFloat{Age: 35.5})
		if err != nil {
			t.Fatalf("Was not expecting an error got %s", err)
		}

		expected := "age: 35.5"
		if string(resp) != expected {
			t.Fatalf("Float was expecting %s\n got %s", expected, string(resp))
		}
	})

	t.Run("Empty", func(t *testing.T) {
		type TestEmpty struct {
			Empty string `form:"-"`
		}

		resp, err := Marshal(TestEmpty{})
		if err != nil {
			t.Fatalf("Was not expecting an error got %s", err)
		}

		expected := ""
		if string(resp) != expected {
			t.Fatalf("Empty was expecting %s\n got\n %s", expected, string(resp))
		}
	})

	t.Run("Bool", func(t *testing.T) {
		type TestSliceSingle struct {
			Name string `form:"name"`
		}

		resp, err := Marshal([]TestSliceSingle{{Name: "test1"}, {Name: "test2"}})
		if err != nil {
			t.Fatalf("Was not expecting an error got %s", err)
		}

		expected := "name: test1\n\nname: test2"
		if string(resp) != expected {
			t.Fatalf("Slice Single was expecting %s\n got %s", expected, string(resp))
		}
	})

	t.Run("Time", func(t *testing.T) {
		type TestTime struct {
			TimeUTC  time.Time `form:"time_utc"`
			TimeUnix int64     `form:"time_unix"`
		}

		resp, err := Marshal(TestTime{
			TimeUTC:  time.Date(2024, 5, 3, 1, 4, 0, 0, time.UTC),
			TimeUnix: 1672530248,
		})
		if err != nil {
			t.Fatalf("Was not expecting an error got %s", err)
		}

		expected := "time_utc: 2024-05-03 01:04:00 +0000 UTC\ntime_unix: 1672530248"
		if string(resp) != expected {
			t.Fatalf("Time was expecting %s\n got %s", expected, string(resp))
		}

	})

	t.Run("Struct", func(t *testing.T) {
		type TestSliceMultiple struct {
			Name string `form:"name"`
			Age  int    `form:"age"`
		}

		resp, err := Marshal([]TestSliceMultiple{{Name: "test1", Age: 35}, {Name: "test2", Age: 36}})
		if err != nil {
			t.Fatalf("Was not expecting an error got %s", err)
		}

		expected := "name: test1\nage: 35\n\nname: test2\nage: 36"
		if string(resp) != expected {
			t.Fatalf("Slice Multiple was expecting %s\n got %s", expected, string(resp))
		}
	})

	t.Run("Slice in Struct", func(t *testing.T) {
		type TestSliceStruct struct {
			Names []string `form:"names"`
		}

		resp, err := Marshal(TestSliceStruct{Names: []string{"test1", "test2"}})
		if err != nil {
			t.Fatalf("Was not expecting an error got %s", err)
		}

		expected := "names: [test1, test2]"
		if string(resp) != expected {
			t.Fatalf("Slice in Struct was expecting %s\n got %s", expected, string(resp))
		}
	})

	type TestSubSubStruct struct {
		Name string `form:"name"`
		Age  int    `form:"age"`
	}
	type TestSubStruct struct {
		Name string `form:"name"`
		Age  int    `form:"age"`

		Sub TestSubSubStruct `form:"sub"`
	}

	t.Run("Slice in Struct Slice", func(t *testing.T) {
		resp, err := Marshal(TestSubStruct{Name: "test", Age: 35, Sub: TestSubSubStruct{Name: "test2", Age: 36}})
		if err != nil {
			t.Fatalf("Was not expecting an error got %s", err)
		}

		expected := "name: test\nage: 35\nsub.name: test2\nsub.age: 36"
		if string(resp) != expected {
			t.Fatalf("Sub Struct was expecting %s\n got %s", expected, string(resp))
		}
	})

	t.Run("Slice in Struct Slice", func(t *testing.T) {
		resp, err := Marshal([]TestSubStruct{{Name: "test", Age: 35, Sub: TestSubSubStruct{Name: "test2", Age: 36}}, {Name: "test3", Age: 37, Sub: TestSubSubStruct{Name: "test4", Age: 38}}})
		if err != nil {
			t.Fatalf("Was not expecting an error got %s", err)
		}

		expected := "name: test\nage: 35\nsub.name: test2\nsub.age: 36\n\nname: test3\nage: 37\nsub.name: test4\nsub.age: 38"
		if string(resp) != expected {
			t.Fatalf("Sub Struct Slice was expecting %s\n got %s", expected, string(resp))
		}
	})
}

type TestMarshalerString string

func (t TestMarshalerString) MarshalPlain() ([]byte, error) {
	return []byte(t), nil
}

type TestMarshalerInt int

func (t TestMarshalerInt) MarshalPlain() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", t)), nil
}

type TestMarshalerFloat float64

func (t TestMarshalerFloat) MarshalPlain() ([]byte, error) {
	return []byte(fmt.Sprintf("%f", t)), nil
}

type TestMarshalerBool bool

func (t TestMarshalerBool) MarshalPlain() ([]byte, error) {
	return []byte(fmt.Sprintf("%t", t)), nil
}

type TestMarshalerStruct struct {
	Name string `form:"name"`
}

func (t TestMarshalerStruct) MarshalPlain() ([]byte, error) {
	return []byte(fmt.Sprintf("name: %s", t.Name)), nil
}

type TestMarshalerStructMultiple struct {
	Name string `form:"name"`
	Age  int    `form:"age"`
}

func (t TestMarshalerStructMultiple) MarshalPlain() ([]byte, error) {
	return []byte(fmt.Sprintf("name: %s - age: %d", t.Name, t.Age)), nil
}

type TestMarshalerStructSlice []TestMarshalerStructMultiple

func (t TestMarshalerStructSlice) MarshalPlain() ([]byte, error) {
	var sb strings.Builder
	for _, v := range t {
		sb.WriteString(fmt.Sprintf("name: %s - age: %d\n", v.Name, v.Age))
	}

	return []byte(sb.String()), nil
}

func TestPlain_Marshaler(t *testing.T) {
	// String
	resp, err := Marshal(TestMarshalerString("test"))
	if err != nil {
		t.Fatalf("Was not expecting an error got %s", err)
	}

	expected := "test"
	if string(resp) != expected {
		t.Fatalf("Marshaler was expecting %s got %s", expected, string(resp))
	}

	// Int
	resp, err = Marshal(TestMarshalerInt(35))
	if err != nil {
		t.Fatalf("Was not expecting an error got %s", err)
	}

	expected = "35"
	if string(resp) != expected {
		t.Fatalf("Marshaler was expecting %s got %s", expected, string(resp))
	}

	// Float
	resp, err = Marshal(TestMarshalerFloat(35.5))
	if err != nil {
		t.Fatalf("Was not expecting an error got %s", err)
	}

	expected = "35.5"
	if string(resp) != expected {
		t.Fatalf("Marshaler was expecting %s got %s", expected, string(resp))
	}

	// Bool
	resp, err = Marshal(TestMarshalerBool(true))
	if err != nil {
		t.Fatalf("Was not expecting an error got %s", err)
	}

	expected = "true"
	if string(resp) != expected {
		t.Fatalf("Marshaler was expecting %s got %s", expected, string(resp))
	}

	// Struct
	resp, err = Marshal(TestMarshalerStruct{Name: "test"})
	if err != nil {
		t.Fatalf("Was not expecting an error got %s", err)
	}

	expected = "name: test"
	if string(resp) != expected {
		t.Fatalf("Marshaler was expecting %s got %s", expected, string(resp))
	}

	// Struct Multiple
	resp, err = Marshal(TestMarshalerStructMultiple{Name: "test", Age: 35})
	if err != nil {
		t.Fatalf("Was not expecting an error got %s", err)
	}

	expected = "name: test - age: 35"
	if string(resp) != expected {
		t.Fatalf("Marshaler was expecting %s got %s", expected, string(resp))
	}

	// Struct Slice
	resp, err = Marshal(TestMarshalerStructSlice{{Name: "test", Age: 35}, {Name: "test2", Age: 36}})
	if err != nil {
		t.Fatalf("Was not expecting an error got %s", err)
	}

	expected = "name: test - age: 35\n\nname: test2 - age: 36"
	if string(resp) != expected {
		t.Fatalf("Marshaler was expecting %s got %s", expected, string(resp))
	}
}
