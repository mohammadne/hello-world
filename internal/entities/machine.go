package entities

import "errors"

type Machine string

const (
	Client Machine = "client"
	// Middleware Machine = "middleware"
	Server Machine = "server"
)

func ValidateMachine(rawMachine string) error {
	if machine := Machine(rawMachine); machine == Client || machine == Server {
		return nil
	}
	return errors.New("Invalid machine type")
}
