package main

import (
    "github.com/buffersandbeer/canlib"
    "flag"
    "fmt"
)

func main() {
    caniface := flag.String("canif", "vcan0", "The CAN interface to capture on")
    flag.Parse()
    c := make(chan canlib.RawCanFrame, 100)
    err := make(chan error)
    go canlib.CaptureCan(*caniface, c, err)
    go printCan(c)
    <-err
}

func printCan(ch <-chan canlib.RawCanFrame) {
    for n:= range ch {
        fmt.Println(canlib.RawCanFrameToString(n, " \t"))
    }
}
