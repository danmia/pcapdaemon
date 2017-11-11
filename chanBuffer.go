package main

// BUFFER FOR CAPTURE MESSAGES

// RingBufferCapmsg defines a structure of channels to have a non-blocking ring buffer of channels
type RingBufferCapmsg struct {
	inputChannel  <-chan Capmsg
	outputChannel chan Capmsg
}

// NewRingBufferCapmsg instantiates a new RingBufferCapmsg object
func NewRingBufferCapmsg(inputChannel <-chan Capmsg, outputChannel chan Capmsg) *RingBufferCapmsg {
	return &RingBufferCapmsg{inputChannel, outputChannel}
}

// Run starts an infinite loop that handels popping messages off the channel
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
