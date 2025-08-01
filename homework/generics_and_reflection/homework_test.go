package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	var b strings.Builder

	val := reflect.ValueOf(person)
	typ := reflect.TypeOf(person)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		tag := field.Tag.Get("properties")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		tagName := parts[0]
		omitempty := len(parts) > 1 && parts[1] == "omitempty"

		if omitempty && checkValueIsZero(value) {
			continue
		}

		var strValue string
		switch value.Kind() {
		case reflect.String:
			strValue = value.String()
		case reflect.Int:
			strValue = fmt.Sprintf("%d", value.Int())
		case reflect.Bool:
			strValue = fmt.Sprintf("%t", value.Bool())
		default:
			return ""
		}

		b.WriteString(fmt.Sprintf("%s=%s\n", tagName, strValue))
	}

	result := b.String()

	return strings.TrimSuffix(result, "\n")
}

func checkValueIsZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int:
		return v.Int() == 0
	case reflect.Bool:
		return !v.Bool()
	default:
		return false
	}
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
