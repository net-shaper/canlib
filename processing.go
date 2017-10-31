package canlib

import (
    "golang.org/x/sys/unix"
    "fmt"
	"strings"
	"strconv"
	"encoding/hex"
	"encoding/binary"
    "crypto/md5"
	"errors"
)

// ByteArrayToCanFrame converts a byte array containing a CAN packet and converts it into a RawCanFrame
func ByteArrayToCanFrame(array []byte, canMessage *RawCanFrame, captureTime int64, capIface string) {
    canMessage.OID = binary.LittleEndian.Uint32(array[0:4])
    canMessage.ID = canMessage.OID
    canMessage.CaptureInterface = capIface

    // Check for the RTR Flag
    if canMessage.ID & unix.CAN_RTR_FLAG != 0 {
        canMessage.Rtr = true
    }

    // Check for the error flag
    if canMessage.ID & unix.CAN_ERR_FLAG != 0 {
        canMessage.Err = true
    }

    // Check for extended can and adjust the ID accordingly
    if canMessage.ID & unix.CAN_EFF_FLAG != 0 {
        canMessage.Eff = true
        canMessage.ID = canMessage.ID & unix.CAN_EFF_MASK
    } else {
        canMessage.Eff = false
        canMessage.ID = canMessage.ID & unix.CAN_SFF_MASK
    }

    canMessage.Dlc = array[4]
    canMessage.Data = array[8:8+canMessage.Dlc]
    canMessage.Timestamp = captureTime
}

// ProcessRawCan will process a raw can message to add additional contextual information
func ProcessRawCan(processed *ProcessedCanFrame, frame RawCanFrame) {
    processed.Packet = frame
    toHash := append(frame.Data, byte(frame.ID))
    processed.PacketHash = fmt.Sprintf("%x", md5.Sum(toHash))
}

// ProcessCandump will take a Socketcan/candump log and parse it into a raw_can_frame
func ProcessCandump(processed *RawCanFrame, frame string) error {

	// Setting to default values since not all values are added when converting from candump
	processed.OID = 0
	processed.Rtr = false
	processed.Eff = false
	processed.Err = false

	splitSpaces := strings.Split(frame, " ")
	time := strings.Split(splitSpaces[0], "(")[1]
	time = strings.Split(time, ")")[0]
	timeFloat, err := strconv.ParseFloat(time, 64)
	if err != nil {
		return errors.New("parsing time failed: " + err.Error())
	}
	processed.Timestamp = int64(timeFloat * 1000000000)
	processed.CaptureInterface = splitSpaces[1]
	splitPacket := strings.Split(splitSpaces[2], "#")
	idInt, err := strconv.ParseUint(splitPacket[0], 16, 32)
	if err != nil {
		return errors.New("parsing id failed: " + err.Error())
	}
	processed.ID = uint32(idInt)
	processed.Data, err = hex.DecodeString(splitPacket[1])
	if err != nil {
		return errors.New("parsing data failed: " + err.Error())
	}
	processed.Dlc = uint8(binary.Size(processed.Data))
	return nil
}
