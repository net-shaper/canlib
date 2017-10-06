package canlib

import (
    "testing"
    "bytes"
)

//TestByteArrayToCanFrame checks that ByteArrayToCanFrame accurately converts an Extended CAN frame into a RawCanFrame
func TestByteArrayToCanFrame(t *testing.T) {
    frame := []byte{109, 237, 19, 137, 8, 0, 0, 0, 15, 234, 197, 79, 101, 147, 251, 118}
    expected := RawCanFrame{
        OID: 2299784557,
        ID: 152300909,
        Rtr: false,
        Err: false,
        Eff: true,
        Dlc: 8,
        Data: []byte{15, 234, 197, 79, 101, 147, 251, 118},
	}
    var result = new(RawCanFrame)
	ByteArrayToCanFrame(frame, result, 0)
	if (result.OID != expected.OID) {
		t.Error("OID mismatch")
	} else if result.ID != expected.ID {
		t.Error("ID mismatch")
	} else if result.Rtr != expected.Rtr {
		t.Error("RTR mismatch")
	} else if result.Err != expected.Err {
        t.Error("ERR mismatch")
    } else if result.Eff != expected.Eff {
        t.Error("EFF mismatch")
    } else if result.Dlc != expected.Dlc {
        t.Error("data length mismatch")
    } else if bytes.Equal(result.Data, expected.Data) != true {
        t.Error("data value mismatch")
    }
}
