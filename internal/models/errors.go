package models

import "fmt"

func errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
