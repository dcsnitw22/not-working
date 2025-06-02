package nas

import "strconv"

func DecodeSchemaOutput(byteArray []byte) string {
	res := ""
	for _, num := range byteArray {
		first4bits := (num >> 4) & 0b00001111
		last4bits := num & 0b00001111
		if int(first4bits) == 15 {
			res += strconv.Itoa(int(last4bits))
		} else {
			res += strconv.Itoa(int(last4bits))
			res += strconv.Itoa(int(first4bits))
		}

	}
	return res
}
