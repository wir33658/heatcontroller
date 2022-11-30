package main

import (
    "fmt"
    "time"
    "log"
    "net/http"
    "github.com/stianeikeland/go-rpio/v4"
)

func main(){
    fmt.Println("!... Hello GPIO ...!")

    var L = rpio.Low
    var I = rpio.Input
    var rpioImpl = RealRPIO{
        Sim : true,
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
            MEAS_DELAY_SECS : 15,
            RECENT_SET_TEMP : 20.0,
        },
    }

    err := rpioImpl.Open()
    if err != nil {
        panic(fmt.Sprint("unable to open gpio", err.Error()))
    }

    client := http.Client{}
    defer rpioImpl.Close()

    Setup(&rpioImpl)
//            Calibrate(&rpioImpl)

    var nextCalib = NextCalibration()
    var done = false
    for !done {
        fmt.Print(".")

        token_obj, err := retrier(getTokenI, client, 15, 5)
        if(err != nil){
            log.Println(err)
            panic(err)
        }

        my_home_obj := getMyHome(client, token_obj)

        var highestTempDiff = calcHighestTempDifference(my_home_obj)

        fmt.Printf("Hightest Temp Diff = %f\n\n", highestTempDiff)
        TempDiff(&rpioImpl, highestTempDiff)

        Wait(time.Second * time.Duration(rpioImpl.Home.MEAS_DELAY_SECS))

        var now = time.Now().Unix()
        if(now >= nextCalib){
            fmt.Println("Time to calibrate")
            nextCalib = NextCalibration()
            Calibrate(&rpioImpl)
        }
    }
}