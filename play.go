package main

import (
	"fmt"

	"github.com/heatcontroller/sub"
)

type IGet interface {
	GetA(v int) string
	GetB(v int) string
}

type OneImpl struct{}

func (OneImpl) GetA(v int) string {
	fmt.Println("GetA")
	return "GetA-v-ignored"
}

func (OneImpl) GetB(v int) string {
	fmt.Println("GetB")
	return "GetB-v-ignored"
}

type OneMock struct{}

func (OneMock) GetA(v int) string {
	fmt.Println("GetAMock")
	return "GetAMock-v-ignored"
}

func (OneMock) GetB(v int) string {
	fmt.Println("GetBMock")
	return "GetBMock-v-ignored"
}

func main() {
	fmt.Println("Test")
	v := sub.GetA(4)
	fmt.Println(v)
	sel := false
	var oi IGet
	if sel {
		oi = OneImpl{}
	} else {
		oi = OneMock{}
	}
	oi.GetA(4)
}
