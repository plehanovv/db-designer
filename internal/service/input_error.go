package service

import "errors"

type InputError struct {
	Message string
}

func (err InputError) Error() string {
	return err.Message
}

func badInput(message string) error {
	return InputError{Message: message}
}

func IsInputError(err error) bool {
	var inputErr InputError
	return errors.As(err, &inputErr)
}
