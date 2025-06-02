package nas

import (
	"strconv"
	"strings"
)

func EncodePduAddress(ipAddr string) ([]byte, error) {
	// Create an empty byte array
	byteArray := make([]byte, 0)

	numArr := strings.Split(ipAddr, ".")

	for _, num := range numArr {
		number, err := strconv.Atoi(num)
		if err != nil {
			return nil, err
		}
		byteArray = append(byteArray, byte(number))
	}

	return byteArray, nil
}
