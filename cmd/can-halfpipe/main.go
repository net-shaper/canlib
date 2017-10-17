package main

import (
    "github.com/buffersandbeer/canlib"
    "flag"
    "fmt"
    "sync"
)

func main() {
    canifaceIn := flag.String("global", "vcan0", "The CAN interface to capture on")
    canifaceOut := flag.String("target", "vcan1", "The CAN interface to pipe to")
    flag.Parse()

    canGlobalChan := make(chan canlib.RawCanFrame, 100)
    canTargetChan := make(chan canlib.RawCanFrame, 100)
    canMultiplexOne := make(chan canlib.RawCanFrame, 100)
    canMultiplexTwo := make(chan canlib.RawCanFrame, 100)
    output := make(chan canlib.RawCanFrame, 100)
    err := make(chan error)

    go canlib.CaptureCan(*canifaceIn, canGlobalChan, err)
    go canlib.CaptureCan(*canifaceOut, canTargetChan, err)
    go canlib.SendCanConcurrent(*canifaceOut, canMultiplexOne, err)
    go globalMultiplex(canGlobalChan, canMultiplexOne, canMultiplexTwo)
    go processing(canTargetChan, canMultiplexTwo, output)

    for message := range output {
        fmt.Println(message)
    }

}

// globalMultiplex will read a value from globalChan and sent that value to both mplexOne and mplexTwo
func globalMultiplex(globalChan <-chan canlib.RawCanFrame, mplexOne chan<- canlib.RawCanFrame,
                     mplexTwo chan<- canlib.RawCanFrame) {

    for message := range globalChan {
        mplexOne <- message
        mplexTwo <- message

    }
}

// processing will start another process to load an array of known messages and then diff that with the target captures
func processing(targetChan <-chan canlib.RawCanFrame, globalChan <-chan canlib.RawCanFrame, output chan<- canlib.RawCanFrame) {
    var seenMessages = []canlib.RawCanFrame{}
    var mutex = &sync.Mutex{}

    go func() {
        for globalMessage := range globalChan {
            mutex.Lock()
            if (canlib.RawFrameInSlice(globalMessage, seenMessages) == false) {
                seenMessages = append(seenMessages, globalMessage)
            }
            mutex.Unlock()
        }
    }()

    for newMessage := range targetChan {
        mutex.Lock()
        if (canlib.RawFrameInSlice(newMessage, seenMessages) == false) {
            canlib.RawCanFrameToString(newMessage, "\t")
        }
        mutex.Unlock()
    }
}
