package canlib

// RawCanChannelMultiplex will take a RawCanFrame sent into the input channel and relay it to all output channels
func RawCanChannelMultiplex(input <-chan RawCanFrame, output ...chan<- RawCanFrame) {

    for message := range input {
        for _, out := range output {
            out <- message
        }
    }
}
