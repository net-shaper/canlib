package canlib

import (
    "testing"
    "fmt"
)

// TestRawCanFrameToString will verify that a CAN message is formatted properly
func TestRawCanFrameToString(t *testing.T) {
    testFrame := RawCanFrame{
                            OID: 1,
                            ID: 1,
                            Dlc: 1,
                            Eff: false,
                            Rtr: false,
                            Err: false,
                            Data: []byte{1},
                            Timestamp: 1000000000,
    }
    expected := "1,1,NOEFF,NORTR,NOERR,1,1,01"
    result := RawCanFrameToString(testFrame, ",")
    if expected != result {
        t.Errorf(fmt.Sprintf("%s != %s", expected, result))
    }
}
