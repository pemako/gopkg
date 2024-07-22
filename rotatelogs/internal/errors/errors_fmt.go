package errors

import "fmt"

func New(s string) error {
	return fmt.Errorf(s)
}

func Errorf(s string, args ...any) error {
	return fmt.Errorf(s, args...)
}

func Wrap(err error, s string) error {
	return fmt.Errorf(s+`: %w`, err)
}

// func Wrapf(err error, format string, args ...any) error {
// 	if err == nil {
// 		return nil
// 	}
// 	format = fmt.Sprintf("err: %s", err) + format
// 	return fmt.Errorf(format, args...)
// }
