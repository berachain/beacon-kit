package viper

import (
	"net/url"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"
)

// StringToSliceHookFunc returns a DecodeHookFunc that converts
// string to []string by splitting on the given sep.
func StringToExecutionAddressFunc() mapstructure.DecodeHookFunc {
	return StringTo(
		func(s string) (common.Address, error) {
			return common.HexToAddress(s), nil
		},
	)
}

// StringToDialURLFunc returns a DecodeHookFunc that converts
// string to *url.URL by parsing the string.
func StringToDialURLFunc() mapstructure.DecodeHookFunc {
	return StringTo(
		func(s string) (*url.URL, error) {
			url, err := url.Parse(s)
			if err != nil {
				return nil, err
			}
			return url, nil
		},
	)
}

// // FilePathToJWTSecretFunc returns a DecodeHookFunc that converts
// func FilePathToJWTSecretFunc() mapstructure.DecodeHookFunc {
// 	return StringTo(
// 		func(s string) (*jwt.Secret, error) {
// 			return jwt.NewFromFile(s)
// 		},
// 	)
// }

// string to *jwt.Secret by reading the file at the given path.
func StringTo[T any](constructor func(string) (T, error)) mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		var retType T
		if t != reflect.TypeOf(retType) {
			return data, nil
		}

		// Convert it by parsing
		return constructor(data.(string))
	}
}
