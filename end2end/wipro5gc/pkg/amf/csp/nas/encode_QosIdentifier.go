package nas

import (
	"errors"
	"strconv"
	"strings"
)

func EncodeQoSIdentifier(qri string) (byte, error) {
	if qri == "no QoS rule identifier assigned" {
		return 0b00000000, nil
	}
	parts := strings.Fields(qri)
	if len(parts) == 2 {
		num, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0b0, err
		}
		// binaryFormat := fmt.Sprintf("%08s", strconv.FormatInt(int64(num), 2))
		return byte(num), nil
	}

	return 0b0, errors.New("invalid input")
}
