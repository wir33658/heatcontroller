package main

import (
        "fmt"
	"time"
        "github.com/stianeikeland/go-rpio/v4"
)

var arr1 = [4]int{1,1,0,0}
var arr2 = [4]int{0,1,0,0}

func moveLeft() {
    var arrOUT = append(arr1[1:], arr1[:1]...)               // [1:]+arr1[:1] # rotates array values of 1 digit
    arr1 = arr2
    copy(arr2[0:], arrOUT)
    //GPIO.output(chan_list, arrOUT)
    set(arrOUT)
    time.Sleep(time.Millisecond * 1)
}

func moveRight() {
    var arrOUT = append(arr1[3:], arr1[:3]...) 
    arr1 = arr2
    copy(arr2[0:], arrOUT)
    //GPIO.output(chan_list, arrOUT)
    set(arrOUT)
    time.Sleep(time.Millisecond * 1)
}


func set(arr []int) {
    fmt.Println("arr:", arr)

    pin17 := rpio.Pin(17)
    pin27 := rpio.Pin(27)
    pin22 := rpio.Pin(22)
    pin4 := rpio.Pin(4)

    /*
    pin17.Output()
    pin27.Output()
    pin22.Output()
    pin4.Output()
    */

    if(arr[0] == 0){
        pin17.Low()
    } else {
        pin17.High()
    }
    if(arr[1] == 0){
        pin27.Low()
    } else {
        pin27.High()
    }
    if(arr[2] == 0){
        pin22.Low()
    } else {
        pin22.High()
    }
    if(arr[3] == 0){
        pin4.Low()
    } else {
        pin4.High()
    }
}

func main99() {
	fmt.Println("GPIO Test")

        err := rpio.Open()
        if err != nil {
                panic(fmt.Sprint("unable to open gpio", err.Error()))
        }

        defer rpio.Close()

        pin17 := rpio.Pin(17)
        pin17.Output()
        pin27 := rpio.Pin(27)
        pin27.Output()
        pin22 := rpio.Pin(22)
        pin22.Output()
        pin4 := rpio.Pin(4)
        pin4.Output()

        for true {
            moveRight()
        }

        /*
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
        */
 
	/*
        time.Sleep(time.Second * 2)

        for x := 0; x < 4; x++ {
                pin0.Toggle()
                time.Sleep(time.Second * 2)
        }
	*/
}

/*

func set(arr []int) {

    pin17 := rpio.Pin(17)
    pin17.Output()

    pin27 := rpio.Pin(27)
    pin27.Output()

    pin22 := rpio.Pin(22)
    pin22.Output()

    pin4 := rpio.Pin(4)
    pin4.Output()

    if(arr[0] == 0)pin17.Low() else pin17.High()
    if(arr[1] == 0)pin27.Low() else pin27.High()
    if(arr[2] == 0)pin22.Low() else pin22.High()
    if(arr[3] == 0)pin4.Low() else pin4.High()

}
*/
