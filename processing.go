package canlib

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"strconv"
	"strings"
)

// ByteArrayToCanFrame converts a byte array containing a CAN packet and converts it into a RawCanFrame
func ByteArrayToCanFrame(array []byte, canMessage *RawCanFrame, captureTime int64, capIface string) {
	canMessage.OID = binary.LittleEndian.Uint32(array[0:4])
	canMessage.ID = canMessage.OID
	canMessage.CaptureInterface = capIface

	// Check for the RTR Flag
	if canMessage.ID&unix.CAN_RTR_FLAG != 0 {
		canMessage.Rtr = true
	}

	// Check for the error flag
	if canMessage.ID&unix.CAN_ERR_FLAG != 0 {
		canMessage.Err = true
	}

	// Check for extended can and adjust the ID accordingly
	if canMessage.ID&unix.CAN_EFF_FLAG != 0 {
		canMessage.Eff = true
		canMessage.ID = canMessage.ID & unix.CAN_EFF_MASK
	} else {
		canMessage.Eff = false
		canMessage.ID = canMessage.ID & unix.CAN_SFF_MASK
	}

	canMessage.Dlc = array[4]
	canMessage.Data = array[8 : 8+canMessage.Dlc]
	canMessage.Timestamp = captureTime
}

// ProcessRawCan will process a raw can message to add additional contextual information
func ProcessRawCan(processed *ProcessedCanFrame, frame RawCanFrame) {
	processed.Packet = frame
    hash := (fmt.Sprintf("%X",frame.ID)+"#"+fmt.Sprintf("%X",frame.Data))
    processed.PacketHash = hash
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

// ProcessCanalyzeLog will take a canalyze/canalyze-dump/can-dump log and parse it into a raw_can_frame
func ProcessCanalyzeLog(processed *RawCanFrame, frame string) error {

	splitSpaces := strings.Split(frame, " ")
	var cleanSplit []string
	for _, str := range splitSpaces {
		if str != "" {
			cleanSplit = append(cleanSplit, str)
		}
	}
	processed.CaptureInterface = cleanSplit[0]
	time, err := strconv.ParseFloat(cleanSplit[1], 64)
	if err != nil {
		return errors.New("parsing time failed: " + err.Error())
	}
	processed.Timestamp = int64(time * 1000000000)
	oidInt, err := strconv.ParseUint(cleanSplit[2], 16, 32)
	if err != nil {
		return errors.New("parsing OID failed: " + err.Error())
	}
	processed.OID = uint32(oidInt)

	idInt, err := strconv.ParseUint(cleanSplit[6], 16, 32)
	if err != nil {
		return errors.New("parsing ID failed: " + err.Error())
	}
	processed.ID = uint32(idInt)

	dlcInt, err := strconv.ParseUint(cleanSplit[7], 10, 32)
	if err != nil {
		return errors.New("parsing DLC failed: " + err.Error())
	}
	processed.Dlc = uint8(dlcInt)
	dataStr := strings.Join(cleanSplit[8:], "")
	processed.Data, err = hex.DecodeString(dataStr)
	if err != nil {
		return errors.New("parsing data failed: " + err.Error())
	}
	if strings.Contains(cleanSplit[3], "NO") {
		processed.Eff = false
	} else {
		processed.Eff = true
	}
	if strings.Contains(cleanSplit[4], "NO") {
		processed.Rtr = false
	} else {
		processed.Rtr = true
	}
	if strings.Contains(cleanSplit[5], "NO") {
		processed.Err = false
	} else {
		processed.Err = true
	}
	return nil
}
