package viper

import (
	"path/filepath"
	"reflect"

	"github.com/spf13/viper"
)

// WriteStructToFile takes any struct labelled with 'mapstructure' tags
// and writes it to a file using the viper package.
func WriteStructToFile(path, name string, data any) error {
	v := viper.New()
	v.SetConfigFile(filepath.Join(path, name))

	// Use reflection to iterate over the fields of the SpecData struct
	val := reflect.ValueOf(data)
	typeOfSpecData := val.Type()
	// Check if the data is a pointer
	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		// Get the underlying value of the pointer
		data = reflect.ValueOf(data).Elem().Interface()
		val = reflect.ValueOf(data)
	}

	for i := 0; i < val.NumField(); i++ {
		// Get the field tag value
		tag := typeOfSpecData.Field(i).Tag.Get("mapstructure")

		// If the tag is not an empty string, set it in the viper instance
		if tag != "" {
			field := val.Field(i)
			if field.Kind() == reflect.Ptr {
				field = field.Elem()
			}
			v.Set(tag, field.Interface())
		}
	}

	return v.WriteConfig()
}
