package main

import (
        "fmt"
	"time"
        "github.com/stianeikeland/go-rpio/v4"
)

func main2() {
	fmt.Println("GPIO Test")

        err := rpio.Open()
        if err != nil {
                panic(fmt.Sprint("unable to open gpio", err.Error()))
        }

        defer rpio.Close()

        pin4 := rpio.Pin(4)
        pin4.Output()

        pin17 := rpio.Pin(17)
        pin17.Output()

        pin27 := rpio.Pin(27)
        pin27.Output()

        pin22 := rpio.Pin(22)
        pin22.Output()

	pin4.High()
	time.Sleep(time.Second * 4)

        pin17.High()
        time.Sleep(time.Second * 4)

        pin27.High()
        time.Sleep(time.Second * 4)

        pin22.High()
        time.Sleep(time.Second * 4)

        pin4.Low()
        time.Sleep(time.Second * 4)
        pin17.Low()
        time.Sleep(time.Second * 4)
        pin27.Low()
        time.Sleep(time.Second * 4)
        pin22.Low()
        time.Sleep(time.Second * 4)
 
	/*
        time.Sleep(time.Second * 2)

        for x := 0; x < 4; x++ {
                pin0.Toggle()
                time.Sleep(time.Second * 2)
        }
	*/
}

