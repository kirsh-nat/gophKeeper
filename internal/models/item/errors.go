package item

import (
	"fmt"
	"strings"
)

type UndefinedDataTypeError struct {
	level string
	Err   error
}

type DataStrutureError struct {
	level string
	Err   error
}

func (le *UndefinedDataTypeError) Error() string {
	return fmt.Sprintf("[%s] %v", le.level, le.Err)
}

func NewUndefinedDataTypeError(label string, err error) error {
	return &UndefinedDataTypeError{
		level: strings.ToUpper(label),
		Err:   err,
	}
}

func (le *DataStrutureError) Error() string {
	return fmt.Sprintf("[%s] %v", le.level, le.Err)
}

func NewDataStrutureError(label string, err error) error {
	return &DataStrutureError{
		level: strings.ToUpper(label),
		Err:   err,
	}
}
