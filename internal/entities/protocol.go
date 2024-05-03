package entities

import "errors"

type Protocol string

const (
	Reality Protocol = "reality"
)

func ValidateProtocol(rawProtocol string) error {
	if protocol := Protocol(rawProtocol); protocol == Reality {
		return nil
	}
	return errors.New("Invalid protocol")
}
