package nas

import "errors"

func BoolToByte(b bool) (byte, error) {
	switch b {
	case false:
		return 0b0, nil
	case true:
		return 0b1, nil
	}

	return 0, errors.New("invalid Input")

}
