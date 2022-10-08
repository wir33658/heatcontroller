package main

// Using : https://github.com/stianeikeland/go-rpio
/*
err := rpio.Open()

pin := rpio.Pin(10)

pin.Output()       // Output mode
pin.High()         // Set pin High
pin.Low()          // Set pin Low
pin.Toggle()       // Toggle pin (Low -> High -> Low)

pin.Input()        // Input mode
res := pin.Read()  // Read state from pin (High / Low)

pin.Mode(rpio.Output)   // Alternative syntax
pin.Write(rpio.High)


pin.PullUp()
pin.PullDown()
pin.PullOff()

pin.Pull(rpio.PullUp)

rpio.Close()

Example (LED blink on Pin 18):

	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	defer rpio.Close()

	pin := rpio.Pin(18)
	pin.Output()

	for x := 0; x < 20; x++ {
		pin.Toggle()
		time.Sleep(time.Second / 5)
	}
*/

import (
	"fmt"
//	"time"

//	"github.com/stianeikeland/go-rpio/v4"
)

type GPIO interface {
	open() string
	close()
	//pin(nr int) 
}

func open() string {
	return "open"
}

func main() {

	fmt.Println("!... Hello GPIO ...!")

	/*
	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	defer rpio.Close()

	pin := rpio.Pin(18)
	pin.Output()

	for x := 0; x < 20; x++ {
		pin.Toggle()
		time.Sleep(time.Second / 5)
	}

	gpio := GPIO

	rv := gpio.open()
	fmt.Println(rv)
	*/
}


