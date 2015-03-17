package ui

import "code.google.com/p/portaudio-go/portaudio"

type Audio struct {
	stream  *portaudio.Stream
	channel chan byte
}

func NewAudio() *Audio {
	a := Audio{}
	a.channel = make(chan byte)
	return &a
}

func (a *Audio) Start() error {
	host, err := portaudio.DefaultHostApi()
	if err != nil {
		return err
	}
	parameters := portaudio.HighLatencyParameters(nil, host.DefaultOutputDevice)
	stream, err := portaudio.OpenStream(parameters, a.Callback)
	if err != nil {
		return err
	}
	if err = stream.Start(); err != nil {
		return err
	}
	a.stream = stream
	return nil
}

func (a *Audio) Stop() error {
	return a.stream.Close()
}

func (a *Audio) Enqueue(sample byte) {
	a.channel <- sample
}

func (a *Audio) Callback(out []byte) {
	for i := range out {
		select {
		case sample := <-a.channel:
			out[i] = sample
		default:
			out[i] = 0
		}
	}
}
