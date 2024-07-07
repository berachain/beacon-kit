package tree

import (
	"strings"
)

type ObjectPath string

func (p ObjectPath) Split() []string {
	return strings.Split(string(p), "/")
}
