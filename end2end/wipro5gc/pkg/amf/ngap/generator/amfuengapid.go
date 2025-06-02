package generator

import "math"

var AmfUeNgapId uint64 = 0

func GenerateAmfUeNgapId() uint64 {
	if AmfUeNgapId <= uint64(math.Pow(2, 40)-1) {
		AmfUeNgapId += 1
	}
	return AmfUeNgapId
}
