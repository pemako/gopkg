package errors

import "fmt"

func New(s string) error {
	return fmt.Errorf(s)
}

func Errorf(s string, args ...interface{}) error {
	return fmt.Errorf(s, args...)
}

func Wrap(err error, s string) error {
	return fmt.Errorf(s+`: %w`, err)
}
