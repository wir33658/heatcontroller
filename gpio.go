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
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

type Irpio interface {
	Open() (err error)
	Close()
	Output(p rpio.Pin)
	Input(p rpio.Pin)
	Toggle(p rpio.Pin)
}

type RealRPIO struct {}

func (RealRPIO) Open() (err error) {
	fmt.Println("Open")
	return rpio.Open()
}
func (RealRPIO) Close() {
	fmt.Println("Close")
	rpio.Close()
}
func (RealRPIO) Output(p rpio.Pin) {
	fmt.Println("Output" + string(p))
	p.Output()
}
func (RealRPIO) Input(p rpio.Pin) {
	fmt.Println("Input" + string(p))
	p.Input()
}
func (RealRPIO) Toggle(p rpio.Pin) {
	fmt.Println("Toggle" + string(p))
	p.Toggle()
}

type MockRPIO struct {
	State bool `default:false`
	
}

func (m MockRPIO) Open() (err error) {
	fmt.Println("Open")
	
	return nil
}
func (MockRPIO) Close() {
	fmt.Println("Close")
}
func (MockRPIO) Output(p rpio.Pin) {
	fmt.Println("Output" + string(p))
}
func (MockRPIO) Input(p rpio.Pin) {
	fmt.Println("Input" + string(p))
}
func (MockRPIO) Toggle(p rpio.Pin) {
	fmt.Println("Toggle" + string(p))
}

var sim = true

func main() {

	fmt.Println("!... Hello GPIO ...!")

	var r Irpio
	if(sim){
		r2 := &MockRPIO{}
		Set(r2, "default")
		r = r2
	} else {
		r = &RealRPIO{}
	}
	
	err := r.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	defer r.Close()
	
	pin18 := rpio.Pin(18)
	r.Output(pin18)

	for x := 0; x < 20; x++ {
		r.Toggle(pin18)
		time.Sleep(time.Second / 5)
	}
}


