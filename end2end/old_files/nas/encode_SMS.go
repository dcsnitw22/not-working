package nas

import "errors"

func EncodeSMS(sms string) (byte, error) {
	switch sms {
	case "SMS over NAS not allowed":
		return NotAllowed, nil
	case "SMS over NAS allowed":
		return Allowed, nil
	default:
		return 0, errors.New("invalid SMS")
	}
}
