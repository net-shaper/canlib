package canlib

import (
    "testing"
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
        t.Errorf("%s != %s", expected, result)
    }
}

// TestTimestampToSeconds makes sure that the function works
func TestTimestampToSeconds(t *testing.T) {
    fakeTime := int64(1000000000)
    expected := float64(1)
    result := TimestampToSeconds(fakeTime)
    if expected != result {
        t.Errorf("%s != %s", expected, result)
    }
}

// TestProcessedCanFrameToString makes sure that ProcessedCanFrameToString works
func TestProcessedCanFrameToString(t *testing.T) {
    testRawFrame := RawCanFrame{
                            OID: 1,
                            ID: 1,
                            Dlc: 1,
                            Eff: false,
                            Rtr: false,
                            Err: false,
                            Data: []byte{1},
                            Timestamp: 1000000000,
    }
    testProcessedFrame := ProcessedCanFrame{Packet: testRawFrame,
                                            CaptureInterface: "test",
                                            PacketHash: "testHash"}
    expected := "test,1,1,NOEFF,NORTR,NOERR,1,1,01,testHash"
    result := ProcessedCanFrameToString(testProcessedFrame, ",")
    if expected != result {
        t.Errorf("%s != %s", expected, result)
    }
}
