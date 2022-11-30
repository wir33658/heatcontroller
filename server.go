package main

import (
//    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

var tlsCertFile = "./certbundle.pem"
var tlsKeyFile = "./server.key"

type IHeatControllerServer interface {
    Open() (err error)
    Close()
}

type HeatControllerServer struct {
    rpio *RealRPIO
}

func Open(rpio *RealRPIO) {
    fmt.Println("Open server")

    mux := http.NewServeMux()
    mux.HandleFunc("/api", apiHandler)

    server := http.ListenAndServe(":8080", mux)
    log.Fatal(server)


}

func (s *HeatControllerServer) getStatus(w http.ResponseWriter, r *http.Request) {
    response := s.rpio.
    fmt.Fprintf(w, "Hello World")
}

// https://www.golinuxcloud.com/golang-http/
func main66(){
    fmt.Println("!... Heatcontroller Server ...!")

    mux := http.NewServeMux()
    mux.HandleFunc("/api", apiHandler)

    server := http.ListenAndServe(":8080", mux)
    log.Fatal(server)

//    server2 := http.ListenAndServeTLS(":8443", tlsCertFile, tlsKeyFile, mux)
//    log.Fatal(server2)
}
