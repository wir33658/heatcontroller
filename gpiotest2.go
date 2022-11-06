package main

import (
    "fmt"
    "time"
)

var arr3 = [4]int{1,1,0,0}
var arr4 = [4]int{0,1,0,0}

func main() {
    fmt.Println("GPIO Test2")
    for true {
        move2()
    }
}

func move2() {
    var arrOUT = append(arr3[1:], arr3[:1]...)               // [1:]+arr1[:1] # rotates array values of 1 digit
    arr3 = arr4
    copy(arr4[0:], arrOUT)
    //GPIO.output(chan_list, arrOUT)
    set2(arrOUT)
    time.Sleep(time.Second * 1)
}

func set2(arr []int) {
    fmt.Println("arr:", arr)

    if(arr[0] == 0){
        fmt.Println("17 : Low")
    } else {
        fmt.Println("17 : High")
    }
}