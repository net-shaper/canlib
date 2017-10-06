package canlib

import (
    "golang.org/x/sys/unix"
    "encoding/binary"
)

// ByteArrayToCanFrame converts a byte array containing a CAN packet and converts it into a RawCanFrame
func ByteArrayToCanFrame(array []byte, canMessage *RawCanFrame, captureTime int64) {
    canMessage.OID = binary.LittleEndian.Uint32(array[0:4])
    canMessage.ID = canMessage.OID

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
