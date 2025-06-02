package nas

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func EncodeQoSIdentifier(qri string) (string, error) {
	parts := strings.Fields(qri)
	if len(parts) == 2 {
		num, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}
		binaryFormat := fmt.Sprintf("%08s", strconv.FormatInt(int64(num), 2))
		return binaryFormat, nil
	}

	return "", errors.New("invalid input")
}
