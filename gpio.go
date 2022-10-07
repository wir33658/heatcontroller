package main

import (
	"fmt"

	"github.com/stianeikeland/go-rpio/v4"
)

func main3() {

	fmt.Println("!... Hello GPIO ...!")

	err := rpio.Open()
}
