package main

import (
    "github.com/buffersandbeer/canlib"
    "flag"
    "fmt"
)

type CanInstance struct {
    frame canlib.RawCanFrame
    src int
}

func main() {
    canGlobal := flag.String("global", "", "The CAN interface for the greater CAN network")
    canTarget := flag.String("target", "", "The CAN interface for the targeted CAN device")
    flag.Parse()

    canGlobalIn := make(chan canlib.RawCanFrame,100)
    canGlobalOut := make(chan canlib.RawCanFrame, 1)
    canTargetIn := make(chan canlib.RawCanFrame, 100)
    canTargetOut := make(chan canlib.RawCanFrame, 1)
    ctcInput := make(chan CanInstance)
    err := make(chan error)

    go canlib.CaptureCan(*canGlobal, canGlobalIn, err)
    go canlib.SendCanConcurrent(*canGlobal, canGlobalOut, err)
    go processFrames(canGlobalIn, ctcInput, 1)

    go canlib.CaptureCan(*canTarget, canTargetIn, err)
    go canlib.SendCanConcurrent(*canTarget, canTargetOut, err)
    go processFrames(canTargetIn, ctcInput, 0)

    canTrafficControl(ctcInput, canTargetOut, canGlobalOut)
}

func canTrafficControl(input <-chan CanInstance, targetOut chan<- canlib.RawCanFrame, globalOut chan<- canlib.RawCanFrame) {
    history := []CanInstance{}
    printTemplate := "%s:\t%s\n"
    for update := range input {
        known := false
        var lastSeen CanInstance

        for _, entry := range history {
            known = canlib.CompareRawFrames(entry.frame, update.frame)
            if known != false {
                lastSeen = entry
                break
            }
        }

        if known == false {
            history = append(history, update)
            lastSeen = update
        }

        if update.src == 1 {
            fmt.Printf(printTemplate, "target", canlib.RawCanFrameToString(lastSeen.frame, " "))
            globalOut <- lastSeen.frame
        } else if update.src == 0 {
            fmt.Printf(printTemplate, "global", canlib.RawCanFrameToString(lastSeen.frame, " "))
            targetOut <- lastSeen.frame
        }
    }
}

func processFrames(captureChan <-chan canlib.RawCanFrame, ctcChan chan<- CanInstance, id int) {
    for newMessage := range captureChan {
        newInstance := CanInstance{
            frame: newMessage,
            src: id,
        }
        ctcChan <- newInstance
    }
}
