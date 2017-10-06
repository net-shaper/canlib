package canlib

import (
    "golang.org/x/sys/unix"
    "net"
    "errors"
)

// SetupCanInterface will set up a CAN file descriptor to be used with sending and recieving CAN message.
// The function takes a string that specifies the interface to open. It returns an integer file descriptor, and an error.
func SetupCanInterface(canInterface string) (int, error) {
    iface, err := net.InterfaceByName(canInterface)
    if err != nil {
        return 0, errors.New("error getting CAN interface by name: " + err.Error())
    }
    var fd int

    fd, err = unix.Socket(unix.AF_CAN, unix.SOCK_RAW, unix.CAN_RAW)
    if err != nil {
        return 0, errors.New("error setting CAN socket: " + err.Error())
    }

    addr := &unix.SockaddrCAN{Ifindex: iface.Index}

    unix.Bind(fd, addr)

    return fd, nil
}
