package canlib

import (
	"encoding/binary"
	"errors"

	"golang.org/x/sys/unix"
	"reflect"
)

// SendCan will send the provided CAN message on the given CAN interface
func SendCan(canInterface string, message RawCanFrame) error {

	if (message.Dlc > 8) || (len(message.Data) != int(message.Dlc)) {

		return errors.New("CAN message to send is invalid")
	}

	canFD, err := SetupCanInterface(canInterface)
	if err != nil {
		return errors.New("error setting up CAN interface: " + err.Error())
	}

	frame := make([]byte, 16)
	binary.LittleEndian.PutUint32(frame[0:4], message.OID)
	frame[4] = byte(message.Dlc)
	copy(frame[8:], message.Data)
	unix.Write(canFD, frame)

	return nil
}

// SendCanConcurrent will utilize a channel to send CAN messages on the given CAN interface
func SendCanConcurrent(canInterface string, canChannel <-chan RawCanFrame, errorChannel chan<- error) {
	canFD, err := SetupCanInterface(canInterface)
	if err != nil {
		errorChannel <- errors.New("error setting up CAN interface: " + err.Error())
		return
	}

	for message := range canChannel {
		frame := make([]byte, 16)
		binary.LittleEndian.PutUint32(frame[0:4], message.OID)
		frame[4] = byte(message.Dlc)
		copy(frame[8:], message.Data)
		unix.Write(canFD, frame)
	}

	errorChannel <- nil
}
