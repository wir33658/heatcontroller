package main

import (
//    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

var tlsCertFile = "./certbundle.pem"
var tlsKeyFile = "./server.key"

func apiHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World")
}

// https://www.golinuxcloud.com/golang-http/
func main(){
    fmt.Println("!... Heatcontroller Server ...!")

    mux := http.NewServeMux()
    mux.HandleFunc("/api", apiHandler)

    server := http.ListenAndServe(":8080", mux)
    log.Fatal(server)

//    server2 := http.ListenAndServeTLS(":8443", tlsCertFile, tlsKeyFile, mux)
//    log.Fatal(server2)
}
