package sub

import (
	"strconv"
)

func GetA(v int) string {
	msg := "getA" + strconv.Itoa(v)
	return msg
}

func GetB(v int) string {
	msg := "getB" + strconv.Itoa(v)
	return msg
}
