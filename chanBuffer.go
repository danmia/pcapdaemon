package main

// BUFFER FOR CAPTURE MESSAGES

type RingBufferCapmsg struct {
    inputChannel  <-chan Capmsg
    outputChannel chan Capmsg
}


func NewRingBufferCapmsg(inputChannel <-chan Capmsg, outputChannel chan Capmsg) *RingBufferCapmsg {
    return &RingBufferCapmsg{inputChannel, outputChannel}
}

func (r *RingBufferCapmsg) Run() {
    for v := range r.inputChannel {
        select {
        case r.outputChannel <- v:
        default:
            <-r.outputChannel
            r.outputChannel <- v
        }
    }
    close(r.outputChannel)
}
