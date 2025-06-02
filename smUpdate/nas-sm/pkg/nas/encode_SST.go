package nas

import (
	"fmt"
	"strconv"
)

func EncodeSST(sst int) (string, error) {
	binaryFormat := fmt.Sprintf("%08s", strconv.FormatInt(int64(sst), 2))
	return binaryFormat, nil
}
