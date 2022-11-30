package main

// Using : https://github.com/stianeikeland/go-rpio

import (
	"fmt"
	"time"
	"strconv"
    "encoding/json"
    "strings"
	"github.com/stianeikeland/go-rpio/v4"
)

// check out : https://github.com/creasty/defaults


type Irpio interface {
	Open() (err error)
	Close()
	Output(p rpio.Pin)
	Input(p rpio.Pin)
	Toggle(p rpio.Pin)
	Low(p rpio.Pin)
	High(p rpio.Pin)

	EngineSet(s7 bool, s0 bool, s2 bool, s3 bool)
}

type RpioStatus struct {
	IsOpen bool 		
	PinStates []rpio.State 	
	PinModes []rpio.Mode 	
}

type RealRPIO struct {
	Sim bool
	Status RpioStatus
	Home HomeStatus
}

func (h *HomeStatus) toString() (string) {
    js, _ := json.MarshalIndent(m, "", "   ")
    fmt.Println(string(js))
}

func (r RealRPIO) Open() (err error) {
	fmt.Println("Open")
	if(!r.Sim){
		err = rpio.Open()
	} else {
		fmt.Println("-- Simulation --")
	}
	if err == nil {
		r.Status.IsOpen = true
	}
	return err
}

func (r RealRPIO) Close() {
	fmt.Println("Close")
	if(!r.Sim){rpio.Close()}
	r.Status.IsOpen = false
}

func (r RealRPIO) Input(p rpio.Pin) { // Input Mode : into the raspi
	if(!r.Sim){p.Input()} else {r.Status.PinModes[p] = rpio.Input}
	r.printPin("Input", p)
}
func (r RealRPIO) Output(p rpio.Pin) { // Output Mode : out from the raspi
	if(!r.Sim){p.Output()} else {r.Status.PinModes[p] = rpio.Output}
	r.printPin("Output", p)
}

func (r RealRPIO) Toggle(p rpio.Pin) {
	if(!r.Sim){p.Toggle()} else {
		if(r.Status.PinStates[p] == rpio.High){r.Status.PinStates[p] = rpio.Low} else {r.Status.PinStates[p] = rpio.High}
	}
	r.printPin("Toggle", p)
}

func (r RealRPIO) High(p rpio.Pin) {
	if(!r.Sim){
		p.High()
//		var state = rpio.ReadPin(p)
//		fmt.Printf("High : Read-State=%d\n",state)
	} else {r.Status.PinStates[p] = rpio.High}
	r.printPin("High", p)
}

func (r RealRPIO) Low(p rpio.Pin) {
	if(!r.Sim){
		p.Low()
//		var state = rpio.ReadPin(p)
//		fmt.Printf("Low : Read-State=%d\n",state)
	} else {r.Status.PinStates[p] = rpio.Low}
	r.printPin("Low", p)
}

var printPinT = false
func (r RealRPIO) printPin(action string, p rpio.Pin) {
	if(printPinT){
		var pinState rpio.State
		if(!r.Sim){
			pinState = p.Read()		
		} else {
			pinState = r.Status.PinStates[p]
		}
		r.Status.PinStates[p] = pinState 
		pinMode := r.Status.PinModes[p]

		fmt.Print(action + " : ")
		fmt.Print("Pin:" + strconv.FormatUint(uint64(p), 10))
		fmt.Print("  State=" + strconv.FormatUint(uint64(pinState), 10))
		fmt.Println("  Mode=" + strconv.FormatUint(uint64(pinMode), 10))
	}
}

var FULL_CIRCLE float64 = 510.0

func toDegree(deg float64) float64 {
	return FULL_CIRCLE / 360 * deg
}

func (r RealRPIO) EngineSet(s7 bool, s0 bool, s2 bool, s3 bool){
// 0 , 2, 3, 7

	pin1 := rpio.Pin(17)
	pin2 := rpio.Pin(27)
	pin3 := rpio.Pin(22)
        pin4 := rpio.Pin(4)

/*
	pin0.Output()
	pin2.Output()
	pin3.Output()
	pin7.Output()
*/
	if(s0){r.High(pin1)} else {r.Low(pin1)}
	if(s2){r.High(pin2)} else {r.Low(pin2)}
	if(s3){r.High(pin3)} else {r.Low(pin3)}
	if(s7){r.High(pin4)} else {r.Low(pin4)}

        time.Sleep(time.Millisecond * 1)

	if(printPinT){
		fmt.Println("-----------------------------")
	}
}

func (r RealRPIO) RightTurn(deg float64){
	fmt.Println("Right-Turn : " + strconv.FormatFloat(deg, 'f', 3, 64))
	var degree = toDegree(deg)
	r.EngineSet(false, false, false, false)

	for (degree > 0.0) {
//		fmt.Print("d")
//		fmt.Println("degree : " + strconv.FormatFloat(degree, 'f', 3, 64))
		r.EngineSet(true, false, false, false)
		r.EngineSet(true, true, false, false)
		r.EngineSet(false, true, false, false)
		r.EngineSet(false, true, true, false)
		r.EngineSet(false, false, true, false)
		r.EngineSet(false, false, true, true)
		r.EngineSet(false, false, false, true)
		r.EngineSet(true, false, false, true)
		degree -= 1
	}
//	fmt.Println()
}

func (r RealRPIO) LeftTurn(deg float64){
	fmt.Println("Left-Turn : " + strconv.FormatFloat(deg, 'f', 3, 64))
	var degree = toDegree(deg)
	r.EngineSet(false, false, false, false)

	for (degree > 0.0) {
//		fmt.Print("d")
//		fmt.Println("degree : " + strconv.FormatFloat(degree, 'f', 3, 64))
		r.EngineSet(true, false, false, true)
		r.EngineSet(false, false, false, true)
		r.EngineSet(false, false, true, true)
		r.EngineSet(false, false, true, false)
		r.EngineSet(false, true, true, false)
		r.EngineSet(false, true, false, false)
		r.EngineSet(true, true, false, false)
		r.EngineSet(true, false, false, false)
		degree -= 1
	}
//	fmt.Println()
}
