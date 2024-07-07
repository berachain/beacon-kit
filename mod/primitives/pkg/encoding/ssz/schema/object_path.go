package schema

import "strings"

// ObjectPath represents a path to an object
type ObjectPath string

// Split returns the path split by "/"
func (p ObjectPath) Split() []string {
	return strings.Split(string(p), "/")
}
