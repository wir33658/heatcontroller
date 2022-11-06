package main

import (
    "fmt"
    "time"
)

var arr1 = [4]int{1,1,0,0}
var arr2 = [4]int{0,1,0,0}

func main() {
    fmt.Println("GPIO Test2")
    for true {
        move()
    }
}

func move99() {
    var arrOUT = append(arr1[1:], arr1[:1]...)               // [1:]+arr1[:1] # rotates array values of 1 digit
    arr1 = arr2
    copy(arr2[0:], arrOUT)
    //GPIO.output(chan_list, arrOUT)
    set(arrOUT)
    time.Sleep(time.Second * 1)
}

func set(arr []int) {
    fmt.Println("arr:", arr)
}