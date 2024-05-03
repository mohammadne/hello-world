package validator

import (
	"errors"
	"strconv"
	"strings"
)

const IPv4Length = 4

func ValidateIP(address string) error {
	splits := strings.Split(address, ".")
	if len(splits) != IPv4Length {
		return errors.New("Invalid IP")
	}

	for index := 0; index < IPv4Length; index++ {
		if number, err := strconv.Atoi(splits[index]); err != nil {
			return errors.New("Invalid IP")
		} else if number < 0 || number > 255 {
			return errors.New("Invalid IP")
		} else if (index == 0 || index == 4) && (number == 0 || number == 255) {
			return errors.New("Invalid IP")
		}
	}

	return nil
}
