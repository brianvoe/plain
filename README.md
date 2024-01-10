# Plain

`plain` is a Go package designed for marshaling and unmarshaling structured data to and from plain text. It supports basic types, nested structs, and custom logic for types that fulfill the Marshaler and Unmarshaler interfaces. This package is ideal for dealing with simple, structured plain text data.

## Features
- Marshaling Go structs to plain text.
- Unmarshaling plain text into Go structs.
- Support for basic data types (string, int, float64, bool).
- Handling nested structs with dot-separated keys.
- Custom marshaling/unmarshaling for types implementing the Marshaler/Unmarshaler interfaces.
- Ignoring fields tagged with -.
- Lenient handling of unknown fields during unmarshaling.
- [Warnings](#warnings)

## Installation
```bash
go get github.com/brianvoe/plain
```

## Basic Usage
### Marshaling
```go
import "github.com/brianvoe/plain"

type Employee struct {
    Name string `plain:"name"`
    Age  int   `plain:"age"`
}

emp := Employee{
    Name: "John Doe",
    Age:  42,
}

data, _ := plain.Marshal(emp)

fmt.Println(string(data))
```

#### Output
```text
name: John Doe
age: 42
```

### Unmarshaling
```go
import "github.com/brianvoe/plain"

type Employee struct {
    Name string `plain:"name"`
    Age  int   `plain:"age"`
}

data := []byte(`name: John Doe
age: 42`)

emp := Employee{}
plain.Unmarshal(data, &emp)

fmt.Printf("%+v", emp)
```

### Output
```text
{Name:John Doe Age:42}
```

## Nested Structs
```go
import "github.com/brianvoe/plain"

type Address struct {
    City  string `plain:"city"`
    State string `plain:"state"`
}

type Person struct {
    Name    string  `plain:"name"`
    Age     int     `plain:"age"`
    Address Address `plain:"address"`
}

// Marshal
data, _err_ := plain.Marshal(Person{
    Name: "John Doe",
    Age:  42,
    Address: Address{
        City:  "New York",
        State: "NY",
    },
})

fmt.Println(string(data))

// Unmarshal
person := Person{}
plain.Unmarshal(data, &person)

fmt.Printf("%+v", person)
```

#### Output
```text
// Marshal
name: John Doe
age: 42
address.city: New York
address.state: NY

// Unmarshal
{Name:John Doe Age:42 Address:{City:New York State:NY}}
```

## Custom Marshaling/Unmarshaling
```go
import "github.com/brianvoe/plain"

type Employee struct {
    Name string `plain:"name"`
    Age  int   `plain:"age"`
}

func (e Employee) MarshalPlain() ([]byte, error) {
    return []byte(fmt.Sprintf("%s: %d", e.Name, e.Age)), nil
}

func (e *Employee) UnmarshalPlain(data []byte) error {
    parts := strings.Split(string(data), ":")
    e.Name = strings.TrimSpace(parts[0])
    e.Age, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
    return nil
}

// Marshal
data, _ := plain.Marshal(Employee{
    Name: "John Doe",
    Age:  42,
})

fmt.Println(string(data))

// Unmarshal
emp := Employee{}
plain.Unmarshal([]byte("John Doe: 42"), &emp)

fmt.Printf("%+v", emp)
```

#### Output
```text
// Marshal
John Doe: 42

// Unmarshal
{Name:John Doe Age:42}
```

## Warnings
### Handling of Newlines in Data

Please be aware that the `plain` package uses newline characters (`\n`) and double newlines (`\n\n`) to interpret the structure of the input data. If your input data includes newline characters within the values themselves, it could lead to unexpected behavior or unmarshaling issues. Ensure that your data is formatted correctly, avoiding newlines within individual values, to prevent parsing errors or data loss during the unmarshaling process.