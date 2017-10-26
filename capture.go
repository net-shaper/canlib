package canlib

import (
    "golang.org/x/sys/unix"
    "time"
    "errors"
)

// CaptureCan will listen to the provided SocketCAN interface and add any messages seen to the provided channel
func CaptureCan(canInterface string, canChannel chan<- RawCanFrame, errorChannel chan<- error) {
    canFD, err := SetupCanInterface(canInterface)
    if err != nil {
        errorChannel <- errors.New("error setting up CAN interface: " + err.Error())
        return
    }

    frame := make([]byte, 16)
    canmsg := new(RawCanFrame)
    for {
        unix.Read(canFD, frame)
        captime := time.Now().UnixNano()
        ByteArrayToCanFrame(frame, canmsg, captime, canInterface)
        canChannel <- *canmsg
    }

    errorChannel <- nil
    close(errorChannel)
}
