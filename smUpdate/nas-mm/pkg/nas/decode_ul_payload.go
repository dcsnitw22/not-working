package nas

import "errors"

func DecodeULPayload(payloadconType string, byteArray []byte) (interface{}, error) {
	if payloadconType == "N1 SM information" {
		epd, mt, err := Classify(byteArray)
		if err != nil {
			return nil, err
		}
		result, err := ReRouteDecode(epd, mt, byteArray)
		if err != nil {
			return nil, err
		}
		return result, nil

	}

	return nil, errors.New("currently not handling this Container Type")

}
