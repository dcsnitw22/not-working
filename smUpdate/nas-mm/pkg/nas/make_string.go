package nas

import "strconv"

func MakeString(numbers []int) string {
	var result string
	for _, num := range numbers {
		result += strconv.Itoa(num)
	}
	return result

}
