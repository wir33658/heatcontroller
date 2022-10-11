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
	"strconv"

	"github.com/creasty/defaults"
	"github.com/stianeikeland/go-rpio/v4"
)

// check out : https://github.com/creasty/defaults


type Irpio interface {
	Open() (err error)
	Close()
	Output(p rpio.Pin)
	Input(p rpio.Pin)
	Toggle(p rpio.Pin)
}

type RpioStates struct {
	isOpen bool 		`default:"false"`
	pinStates [32]bool 	`default:"[false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false]"`
}



type RealRPIO struct {
	States RpioStates
}

func (r RealRPIO) Open() (err error) {
	fmt.Println("Open")
	r.States.isOpen = true
	return rpio.Open()
}
func (r RealRPIO) Close() {
	fmt.Println("Close")
	r.States.isOpen = false
	rpio.Close()
}
func (r RealRPIO) Output(p rpio.Pin) {
	fmt.Println("Output" + string(p))
	r.States.pinStates[p] = true
	p.Output()
}
func (r RealRPIO) Input(p rpio.Pin) {
	fmt.Println("Input" + string(p))
	r.States.pinStates[p] = false
	p.Input()
}
func (r RealRPIO) Toggle(p rpio.Pin) {
	fmt.Println("Toggle" + string(p))
	r.States.pinStates[p] = !r.States.pinStates[p]
	p.Toggle()
}

type Fuck uint8

type MockRPIO struct {
	States RpioStates	
}

func (m MockRPIO) Open() (err error) {
	fmt.Println("Open")
	m.States.isOpen = true
	m.States.pinStates = [32]bool{false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false}
	return nil
}
func (m MockRPIO) Close() {
	m.States.isOpen = false
	fmt.Println("Close")
}
func (m MockRPIO) Output(p rpio.Pin) {
	m.States.pinStates[p] = true
	fmt.Println("Output" + string(p))
}
func (m MockRPIO) Input(p rpio.Pin) {
	m.States.pinStates[p] = false
	fmt.Println("Input" + string(p))
}
func (m MockRPIO) Toggle(p rpio.Pin) {
	pos := uint64(p)
	var state bool 
	var statenow = m.States.pinStates[pos]
	fmt.Println("Fuck:" + strconv.FormatBool(statenow))
	if(statenow == true){
		state = false 
	} else {
		state = true
	}
	fmt.Println("Fuck2:" + strconv.FormatBool(state))
	*(&m.States.pinStates[pos]) = state
	fmt.Println("Toggle(" + strconv.FormatUint(pos, 10) + "):" + strconv.FormatBool(m.States.pinStates[pos]))
}

var sim = true

func main() {

	fmt.Println("!... Hello GPIO ...!")

	var r Irpio
	if(sim){
		r2 := &MockRPIO{}
		if err := defaults.Set(r2); err != nil {
			panic(err)
		}
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


