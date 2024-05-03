package validator

import (
	"errors"
	"strings"
)

const DomainMinLength = 2

func ValidateDomain(address string) error {
	splits := strings.Split(address, ".")
	if len(splits) < DomainMinLength {
		return errors.New("Invalid Domain")
	}

	return nil
}
