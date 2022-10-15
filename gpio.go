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

//	"github.com/creasty/defaults"
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
	IsOpen bool 		// `default:"false"`
	PinStates []bool 	// `default:"[false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false]"`
}



type RealRPIO struct {
	States RpioStates
}

func (r RealRPIO) Open() (err error) {
	fmt.Println("Open")
	r.States.IsOpen = true
	return rpio.Open()
}
func (r RealRPIO) Close() {
	fmt.Println("Close")
	r.States.IsOpen = false
	rpio.Close()
}
func (r RealRPIO) Output(p rpio.Pin) {
	fmt.Println("Output" + string(p))
	r.States.PinStates[p] = true
	p.Output()
}
func (r RealRPIO) Input(p rpio.Pin) {
	fmt.Println("Input" + string(p))
	r.States.PinStates[p] = false
	p.Input()
}
func (r RealRPIO) Toggle(p rpio.Pin) {
	fmt.Println("Toggle" + string(p))
	r.States.PinStates[p] = !r.States.PinStates[p]
	p.Toggle()
}

type Fuck uint8

type MockRPIO struct {
	States RpioStates	
}

func (m MockRPIO) Open() (err error) {
	fmt.Println("Open")
	m.States.IsOpen = true
	return nil
}
func (m MockRPIO) Close() {
	m.States.IsOpen = false
	fmt.Println("Close")
}
func (m MockRPIO) Output(p rpio.Pin) {
	m.States.PinStates[p] = true
	printPin("Output", p, m.States.PinStates[p])
}
func (m MockRPIO) Input(p rpio.Pin) {
	m.States.PinStates[p] = false
	printPin("Input", p, m.States.PinStates[p])
}
func (m MockRPIO) Toggle(p rpio.Pin) {
	var pin = &(m.States.PinStates[p])
	if(*pin == true){*pin = false } else { *pin = true }
	printPin("Toggle", p, *pin)
}

func printPin(action string, p rpio.Pin, state bool) {
	fmt.Print(action + " : ")
	fmt.Print("Pin:" + strconv.FormatUint(uint64(p), 10))
	fmt.Println("  " + strconv.FormatBool(state))
}


var sim = true

func main() {

	fmt.Println("!... Hello GPIO ...!")

	var r Irpio
	if(sim){
		r = MockRPIO{
			States : RpioStates {
				IsOpen : false,
				PinStates: []bool{false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false,false},
			},
		}
	} else {
		r = RealRPIO{}
	}
	
	err := r.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	defer r.Close()
	
	pin18 := rpio.Pin(18)
	r.Output(pin18)

	for x := 0; x < 4; x++ {
		r.Toggle(pin18)
		time.Sleep(time.Second / 5)
	}



	/*
	//    ar1 := []bool{true,true,true}
	var as = ha{}
	as.as.ar1 = []bool{true,true,true}

	fmt.Println("Fuck1:" + strconv.FormatBool(as.as.ar1[1]))
	as.as.ar1[1] = false
	fmt.Println("Fuck2:" + strconv.FormatBool(as.as.ar1[1]))
	as.as.ar1[1] = true
	fmt.Println("Fuck3:" + strconv.FormatBool(as.as.ar1[1]))
	par11 := &as.as.ar1[1]
	*par11 = false
	fmt.Println("Fuck4:" + strconv.FormatBool(as.as.ar1[1]))
	*par11 = true
	fmt.Println("Fuck5:" + strconv.FormatBool(as.as.ar1[1]))

	var a ai
	a = as
	a.t()

	fmt.Println("Fuck9:" + strconv.FormatBool(as.as.ar1[1]))
	fmt.Println("Fuck10:" + strconv.FormatBool(as.as.ar1[1]))

	*/
}


type ai interface {
	t()
}

type arstr struct {
	ar1 []bool 
}

type ha struct {
	as arstr
}

func (as ha) t() {

	as.as.ar1[1] = false
	fmt.Println("Fuck6:" + strconv.FormatBool(as.as.ar1[1]))
	as.as.ar1[1] = true
	fmt.Println("Fuck7:" + strconv.FormatBool(as.as.ar1[1]))
	as.as.ar1[1] = false
	fmt.Println("Fuck8:" + strconv.FormatBool(as.as.ar1[1]))
	/*
	par11 := &as.as.ar1[1]
	*par11 = false
	fmt.Println("Fuck8:" + strconv.FormatBool(as.as.ar1[1]))
	*/
}
