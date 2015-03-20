package ui

import (
	"fmt"

	"code.google.com/p/portaudio-go/portaudio"
)

type Audio struct {
	stream  *portaudio.Stream
	channel chan byte
}

func NewAudio() *Audio {
	a := Audio{}
	a.channel = make(chan byte, 44100)
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
	if err := stream.Start(); err != nil {
		return err
	}
	a.stream = stream
	return nil
}

func (a *Audio) Stop() error {
	return a.stream.Close()
}

func (a *Audio) Callback(out []byte) {
	count := 0
	for i := range out {
		select {
		case sample := <-a.channel:
			out[i] = sample
		default:
			out[i] = 0
			count++
		}
	}
	if count > 0 {
		fmt.Println(count)
	}
}
