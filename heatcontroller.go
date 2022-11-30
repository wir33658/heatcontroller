package main

import (
    "fmt"
    "time"
    "strconv"
    "math"
    "github.com/stianeikeland/go-rpio/v4"
)

// check out : https://github.com/creasty/defaults

type HomeStatus struct {
    TRIGGER_STEP float64
    HALF_DEGREE_STEP float64
    MIN_TEMP float64
    MAX_TEMP float64
    RECENT_SET_TEMP float64
    MEAS_DELAY_SECS int64
    LAST_CMD string
}

/*
type Controller interface {
    Setup(r *RealRPIO)
    Calibrate(r *RealRPIO)
    TempDiff(r *RealRPIO, tempdiff float64)
}
*/

func Calibrate(r *RealRPIO) {
    fmt.Println("Calibration started ...")
    r.LeftTurn(r.Home.TRIGGER_STEP)
    time.Sleep(time.Second * 2)
    r.LeftTurn(60.0 * r.Home.HALF_DEGREE_STEP) // should be max(= 30) now
    time.Sleep(time.Second * 10) // Back to 20 degrees
    r.RightTurn(r.Home.TRIGGER_STEP)
    time.Sleep(time.Second * 2)
    r.RightTurn(20 * r.Home.HALF_DEGREE_STEP) // should be 20 now
    r.EngineSet(false, false, false, false)
    r.Home.RECENT_SET_TEMP = 20
    time.Sleep(time.Second * 4)
    fmt.Println("Calibration done.")
    fmt.Println("Set temp should be " + strconv.FormatFloat(r.Home.RECENT_SET_TEMP, 'f', 2, 64))
    r.Home.LAST_CMD = "Calibrate"
}

func TempDiff(rpio *RealRPIO, tempdiff float64) {
    fmt.Println("Tempdiff : " + strconv.FormatFloat(tempdiff, 'f', 3, 64))
    fmt.Println("Recenttemp : " + strconv.FormatFloat(rpio.Home.RECENT_SET_TEMP, 'f', 3, 64))

    // In case the temp diff stays around 0 the recent temp will never drop and
    // might stay on a high heating level which causes unnecessary heating cost.
    // Therefor it will be set to -1 to lower it anyway.
    if(tempdiff < 0.49){
        tempdiff = -1.0
        fmt.Println("Tempdiff readjusted to: " + strconv.FormatFloat(tempdiff, 'f', 3, 64))
    }

    var goal = rpio.Home.RECENT_SET_TEMP + tempdiff
    fmt.Println("goal : " + strconv.FormatFloat(goal, 'f', 3, 64))
    if(goal < rpio.Home.MIN_TEMP){
        goal = rpio.Home.MIN_TEMP
    } else {
        if(goal > rpio.Home.MAX_TEMP){
            goal = rpio.Home.MAX_TEMP
        } else {
            fmt.Println("goal : " + strconv.FormatFloat(goal, 'f', 3, 64))
        }
    }

    var finaltempdiff = roundFloat(goal - rpio.Home.RECENT_SET_TEMP, 0)
    var finaltempdiffabs = finaltempdiff
    if(finaltempdiff < 0){finaltempdiffabs = finaltempdiff * -1}
    fmt.Println("finaltempdiff : "  + strconv.FormatFloat(finaltempdiff, 'f', 3, 64))
    fmt.Println("abs : "   + strconv.FormatFloat(finaltempdiffabs, 'f', 3, 64))

    if(finaltempdiff < 0) {
        rpio.RightTurn(rpio.Home.TRIGGER_STEP)
        time.Sleep(time.Second * 2)
        rpio.RightTurn(finaltempdiffabs * 2 * rpio.Home.HALF_DEGREE_STEP)
        time.Sleep(time.Second * 2)
        rpio.Home.RECENT_SET_TEMP = goal
    } else if(finaltempdiff > 0) {
        rpio.LeftTurn(rpio.Home.TRIGGER_STEP)
        time.Sleep(time.Second * 2)
        rpio.LeftTurn(finaltempdiffabs * 2 * rpio.Home.HALF_DEGREE_STEP)
        time.Sleep(time.Second * 2)
        rpio.Home.RECENT_SET_TEMP = goal
    }
    rpio.EngineSet(false, false, false, false)

    fmt.Println("Adjusted")
    fmt.Println("Set temp should be " + strconv.FormatFloat(rpio.Home.RECENT_SET_TEMP, 'f', 3, 64))
    rpio.Home.LAST_CMD = "TempDiff: " + strconv.FormatFloat(tempdiff, 'f', 3, 64)
}

func Setup(r *RealRPIO) {
    fmt.Println("Setup")

    pin1 := rpio.Pin(17)
    pin2 := rpio.Pin(27)
    pin3 := rpio.Pin(22)
    pin4 := rpio.Pin(4)

    r.Output(pin1)
    r.Output(pin2)
    r.Output(pin3)
    r.Output(pin4)

    r.EngineSet(false, false, false, false)
}

func Wait(d time.Duration){
    time.Sleep(d)
}

func NextCalibration() int64 { // in Secs
    var nc = time.Now().Unix() + (60 * 60)  // = 60 minutes
    return nc
}

func roundFloat(val float64, precision uint) float64 {
    ratio := math.Pow(10, float64(precision))
    return math.Round(val*ratio) / ratio
}

func calcHighestTempDifference(my_home_obj my_home) float64 {
    var highestTempDiff = -100.0
    for key, zone := range my_home_obj.Zones {
        var istr = ""
        if(key == "Toilette oben"){
            istr = " --- is ignored --- "
        }
        fmt.Printf("Zone: %s %s ->\t\t\t", key, istr)
        var recentTemp = zone.ZoneState.SensorDataPoints.InsideTemperature.Celsius
        var goalOverlayTemp = zone.ZoneState.Overlay.Setting.Temperature.Celsius
        var goalSettingTemp = zone.ZoneState.Setting.Temperature.Celsius
        var goalTemp = goalSettingTemp
        if(goalOverlayTemp > 0.0){goalTemp = goalOverlayTemp}
        var tempDiff = float64(goalTemp - recentTemp)
        fmt.Printf("Temps : recent=%f;  ogoal=%f;  sgoal=%f;  diff=%f\n", recentTemp, goalOverlayTemp, goalSettingTemp, tempDiff)
        if(key != "Toilette oben"){
            if(tempDiff > highestTempDiff){
                highestTempDiff = tempDiff
            }
        }
    }
    return highestTempDiff
}