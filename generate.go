package canlib

import (
    "errors"
    "golang.org/x/sys/unix"
)

// CreateRawFrame will take an ID, Data, and Flags to generate a valid RawCanFrame
func CreateRawFrame(targetFrame *RawCanFrame, id uint32, data []byte, eff bool, rtr bool, err bool) error {
    targetFrame.ID = id
    targetFrame.Eff = eff
    targetFrame.Rtr = rtr
    targetFrame.Err = err
    dataLength := uint8(len(data))
    if dataLength > 8 {
        return errors.New("data too long. Data must be < 8 bytes")
    }
    targetFrame.Dlc = dataLength
    targetFrame.Data = data

    oid := id
    if eff {
        oid = id & unix.CAN_EFF_FLAG
    }
    if err {
        oid = oid & unix.CAN_ERR_FLAG
    }
    if rtr {
        oid = oid & unix.CAN_RTR_FLAG
    }
    targetFrame.OID = oid
    return nil
}
