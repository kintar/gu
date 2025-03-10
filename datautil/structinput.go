package datautil

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// CreateStructByUserInput accepts a struct type parameter and reflects over its fields. For each supported primitive
// type field, the user is prompted to enter a value, read from stdin. The input values are parsed to the appropriate
// type and assigned to a new strut, and the struct is returned to the caller.
// If a parsing error, an invalid field type, or the end of stdin are encountered, an appropriate error will be returned
// along with the struct and any fields which have been successfully set.
func CreateStructByUserInput[T any]() (T, error) {
	stdin := bufio.NewScanner(os.Stdin)

	result := new(T)

	vt := reflect.TypeOf(*result)
	vv := reflect.ValueOf(result).Elem()

	for f := 0; f < vt.NumField(); f++ {
		fieldDef := vt.Field(f)
		fmt.Printf("%s: ", fieldDef.Name)
		if !stdin.Scan() {
			fmt.Print("interrupted\n\n")
			return *result, errors.New("interrupted")
		}

		field := vv.Field(f)
		switch field.Kind() {
		case reflect.String:
			field.SetString(stdin.Text())
		case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Int16:
			if v, e := strconv.ParseInt(stdin.Text(), 10, 64); e != nil {
				return *result, e
			} else {
				field.SetInt(v)
			}
		case reflect.Float64, reflect.Float32:
			if v, e := strconv.ParseFloat(stdin.Text(), 64); e != nil {
				return *result, e
			} else {
				field.SetFloat(v)
			}
		case reflect.Bool:
			if v, e := strconv.ParseBool(stdin.Text()); e != nil {
				return *result, e
			} else {
				field.SetBool(v)
			}
		default:
			fmt.Println("is of an unsupported type, skipping...")
		}
	}

	return *result, nil
}
