package canlib

// RawCanFrame represents the data contained in a CAN packet
type RawCanFrame struct {
    OID uint32 // 32-bit CAN_ID before masks applied
    ID uint32 // 32 bit CAN_ID + EFF/RTR/ERR
    Dlc uint8 // Payload length in bytes
    Eff bool // Extended frame flag
    Rtr bool // Remote transmission request flag
    Err bool // Error flag
    Data []byte // Message Payload
    Timestamp int64 // Time message was captured as Unix Timestamp in nanoseconds
}

// ProcessedCanFrame represents a CAN packet and additional data about the packet
type ProcessedCanFrame struct {
    Packet RawCanFrame // CAN packet
    PacketHash string // md5 hash of the Packet's ID and Data fields
    CaptureInterface string // Name of capturing interface
}
