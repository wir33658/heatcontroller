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
	"net/http"
//	"math"
	"log"
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
	Low(p rpio.Pin)
	High(p rpio.Pin)

	Setup()
	EngineSet(s7 bool, s0 bool, s2 bool, s3 bool)
	Calibrate()
	TempDiff(tempdiff float64)
}

type RpioStatus struct {
	IsOpen bool 		
	PinStates []rpio.State 	
	PinModes []rpio.Mode 	
}

type HomeStatus struct {
	TRIGGER_STEP float64
	HALF_DEGREE_STEP float64
	MIN_TEMP float64
	MAX_TEMP float64
	RECENT_TEMP float64
	LAST_CMD string
}

type RealRPIO struct {
	Sim bool
	Status RpioStatus
	Home HomeStatus
}

func (r RealRPIO) Open() (err error) {
	fmt.Println("Open")
	if(!r.Sim){
		err = rpio.Open()
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
	if(!r.Sim){p.High()} else {r.Status.PinStates[p] = rpio.High}
	r.printPin("High", p)
}

func (r RealRPIO) Low(p rpio.Pin) {
	if(!r.Sim){p.Low()} else {r.Status.PinStates[p] = rpio.Low}
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
	pin0 := rpio.Pin(0)
	pin2 := rpio.Pin(2)
	pin3 := rpio.Pin(3)
	pin7 := rpio.Pin(7)

	if(s0){r.High(pin0)} else {r.Low(pin0)}
	if(s2){r.High(pin2)} else {r.Low(pin2)}
	if(s3){r.High(pin3)} else {r.Low(pin3)}
	if(s7){r.High(pin7)} else {r.Low(pin7)}

	if(printPinT){
		fmt.Println("-----------------------------")
	}
}

func (r RealRPIO) RightTurn(deg float64){
//	fmt.Println("Right-Turn : " + strconv.FormatFloat(deg, 'e', 3, 64))
	var degree = toDegree(deg)
	r.EngineSet(false, false, false, false)

	for (degree > 0.0) {
//		fmt.Print("d")
//		fmt.Println("degree : " + strconv.FormatFloat(degree, 'e', 3, 64))
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
	fmt.Println("Left-Turn : " + strconv.FormatFloat(deg, 'e', 3, 64))
	var degree = toDegree(deg)
	r.EngineSet(false, false, false, false)

	for (degree > 0.0) {
//		fmt.Print("d")
//		fmt.Println("degree : " + strconv.FormatFloat(degree, 'e', 3, 64))
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

func (r RealRPIO) Calibrate() {
	fmt.Println("Calibration started ...")
	r.LeftTurn(r.Home.TRIGGER_STEP)
	time.Sleep(time.Second * 2)
	r.LeftTurn(40.0 * r.Home.HALF_DEGREE_STEP) // should be max(= 30) now
	time.Sleep(time.Second * 10) // Back to 20 degrees
	r.RightTurn(r.Home.TRIGGER_STEP)
	time.Sleep(time.Second * 2)
	r.RightTurn(20 * r.Home.HALF_DEGREE_STEP) // should be 20 now
	r.EngineSet(false, false, false, false)
	r.Home.RECENT_TEMP = 20
	time.Sleep(time.Second * 2)
	fmt.Println("Calibration done.")
	fmt.Println("Set temp should be " + strconv.FormatFloat(r.Home.RECENT_TEMP, 'e', 3, 64))
	r.Home.LAST_CMD = "Calibrate"
	}

func (r RealRPIO) TempDiff(tempdiff float64) {
	fmt.Println("Tempdiff : " + strconv.FormatFloat(tempdiff, 'e', 3, 64))
	fmt.Println("Recenttemp : " + strconv.FormatFloat(r.Home.RECENT_TEMP, 'e', 3, 64))

	var goal = r.Home.RECENT_TEMP + tempdiff
	fmt.Println("goal : " + strconv.FormatFloat(goal, 'e', 3, 64))
	if(goal < r.Home.MIN_TEMP){
		goal = r.Home.MIN_TEMP
	} else {
		if(goal > r.Home.MAX_TEMP){
			goal = r.Home.MAX_TEMP
		} else {
			fmt.Println("goal : " + strconv.FormatFloat(goal, 'e', 3, 64))
		}
	}

	var finaltempdiff = goal - r.Home.RECENT_TEMP
	var finaltempdiffabs = finaltempdiff
	if(finaltempdiff < 0){finaltempdiffabs = finaltempdiff * -1}
	fmt.Println("finaltempdiff : "  + strconv.FormatFloat(finaltempdiff, 'e', 3, 64))
	fmt.Println("abs : "   + strconv.FormatFloat(finaltempdiffabs, 'e', 3, 64))

	if(finaltempdiff < 0) {
		r.RightTurn(r.Home.TRIGGER_STEP)
		time.Sleep(time.Second * 2)
		r.RightTurn(finaltempdiffabs * 2 * r.Home.HALF_DEGREE_STEP)
		time.Sleep(time.Second * 2)
		r.Home.RECENT_TEMP = goal
	} else if(finaltempdiff > 0) {
		r.LeftTurn(r.Home.TRIGGER_STEP)
		time.Sleep(time.Second * 2)
		r.LeftTurn(finaltempdiffabs * 2 * r.Home.HALF_DEGREE_STEP)
		time.Sleep(time.Second * 2)
		r.Home.RECENT_TEMP = goal
	}
	r.EngineSet(false, false, false, false)

	fmt.Println("Adjusted")
	fmt.Println("Set temp should be " + strconv.FormatFloat(r.Home.RECENT_TEMP, 'e', 3, 64))
	r.Home.LAST_CMD = "TempDiff: " + strconv.FormatFloat(tempdiff, 'e', 3, 64)
}

func (r RealRPIO) Setup() {
	fmt.Println("Setup")
	pin0 := rpio.Pin(0)
	pin2 := rpio.Pin(2)
	pin3 := rpio.Pin(3)
	pin7 := rpio.Pin(7)

	r.Output(pin0)
	r.Output(pin2)
	r.Output(pin3)
	r.Output(pin7)

	r.EngineSet(false, false, false, false)
}

func Wait(d time.Duration){
	time.Sleep(d)
}

func NextCalibration() int64 { // in Secs
	var nc = time.Now().Unix() + (60 * 2)  // = 2 minutes
	return nc
}

var sim = true

func main() {

	fmt.Println("!... Hello GPIO ...!")

	var L = rpio.Low
	var I = rpio.Input
	var r Irpio = RealRPIO{
		Sim : sim,
		Status : RpioStatus {
			IsOpen : false,
			PinStates: []rpio.State{L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L,L},
			PinModes: []rpio.Mode{I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I,I},
		},
		Home : HomeStatus {
			TRIGGER_STEP : 100.0,
			HALF_DEGREE_STEP : 44.0,
			MIN_TEMP : 18.0,
			MAX_TEMP : 23.0,
			RECENT_TEMP : 20.0,
		},
	}
	
	err := r.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	client := http.Client{}
	defer r.Close()
	
	r.Setup()	
	r.Calibrate()

	var nextCalib = NextCalibration() 
	var done = false
	for !done {
		fmt.Print(".")
		/*
		token_obj, _ := getToken(client, 15, 5)
//		fmt.Println("Token :", token_obj)
		me_obj := getMe(client, token_obj)
		printMe(me_obj)
		*/

		token_obj, err := retrier(getTokenI, client, 15, 5)
		if(err != nil){
			log.Println(err)
			panic(err)
		}
	//	printToken(token_obj)

		my_home_obj := getMyHome(client, token_obj)
		
		var highestTempDiff = -100.0
		for key, zone := range my_home_obj.Zones {
			fmt.Printf("Zone: %s ->\t\t\t", key)
			var recentTemp = zone.ZoneState.SensorDataPoints.InsideTemperature.Celsius
			var goalOverlayTemp = zone.ZoneState.Overlay.Setting.Temperature.Celsius
			var goalSettingTemp = zone.ZoneState.Setting.Temperature.Celsius
			var goalTemp = goalSettingTemp
			if(goalOverlayTemp > 0.0){goalTemp = goalOverlayTemp}
			var tempDiff = float64(goalTemp - recentTemp)
			fmt.Printf("Temps : recent=%f;  ogoal=%f;  sgoal=%f;  diff=%f\n", recentTemp, goalOverlayTemp, goalSettingTemp, tempDiff)
			if(tempDiff > highestTempDiff){
				highestTempDiff = tempDiff
			}
		}
		fmt.Printf("Hightest Temp Diff = %f\n\n", highestTempDiff)
		r.TempDiff(highestTempDiff)

		/*
		if(highestTempDiff > 0){ // Rise heat
		} else { // Lower heat
		}
		*/

		Wait(time.Second * 10)
		var now = time.Now().Unix()
		if(now >= nextCalib){
			fmt.Println("Time to calibrate")
			nextCalib = NextCalibration()
			r.Calibrate()
		}
	}



	/*
	for x := 0; x < 4; x++ {
		r.Toggle(pin18)
		time.Sleep(time.Second / 5)
	}
	*/



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


/*
type MockRPIO struct {
	Status RpioStatus	
}

func (m MockRPIO) Open() (err error) {
	fmt.Println("Open")
	m.Status.IsOpen = true
	return nil
}
func (m MockRPIO) Close() {
	m.Status.IsOpen = false
	fmt.Println("Close")
}
func (m MockRPIO) Output(p rpio.Pin) {
	m.Status.PinStates[p] = true
	printPin("Output", p, m.Status.PinStates[p])
}
func (m MockRPIO) Input(p rpio.Pin) {
	m.Status.PinStates[p] = false
	printPin("Input", p, m.Status.PinStates[p])
}
func (m MockRPIO) Toggle(p rpio.Pin) {
	var pin = &(m.Status.PinStates[p])
	if(*pin == true){*pin = false } else { *pin = true }
	printPin("Toggle", p, *pin)
}
*/


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
