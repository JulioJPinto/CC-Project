package error_helpers

import (
	"fmt"
	"strings"
)

type WrapError struct {
	Err error
	Msg string
}

func (err WrapError) Error() string {
	if err.Err != nil {
		return fmt.Sprintf("%s: %v", err.Msg, err.Err)
	}
	return err.Msg
}
func (err WrapError) wrap(inner error) error {
	return WrapError{Msg: err.Msg, Err: inner}
}
func (err WrapError) Unwrap() error {
	return err.Err
}
func (err WrapError) Is(target error) bool {
	ts := target.Error()
	return ts == err.Msg || strings.HasPrefix(ts, err.Msg+": ")
}