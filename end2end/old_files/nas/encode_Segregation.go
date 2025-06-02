package nas

import "errors"

func EncodeSegregation(seg string) (byte, error) {
	switch seg {
	case "Segregation not requested":
		return NotRequested, nil
	case "Segregation requested":
		return Requested, nil
	}

	return 0, errors.New("invalid Input")

}
