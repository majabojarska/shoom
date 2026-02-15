package main

import (
	"fmt"
	"time"

	"github.com/ebitengine/oto/v3"
)

type Noise struct {
	x1 byte
	x2 byte
}

func (noise *Noise) Init() {
	noise.x1 = 0x67
	noise.x2 = 0xef
}

func (noise *Noise) Read(p []byte) (n int, err error) {
	// fmt.Printf("x1=%d, x2=%d\n", noise.x1, noise.x2)

	count := 1024
	for i := range count {
		noise.x1 ^= noise.x2
		p[i] = noise.x2
		// fmt.Printf("sample %d\n", noise.x2)
		noise.x2 += noise.x1
	}

	return 1024, nil
}

func main() {
	// Read the mp3 file into memory
	// fileBytes, err := os.ReadFile("./song.mp3")
	// if err != nil {
	// 	panic("reading song.mp3 failed: " + err.Error())
	// }

	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
	// fileBytesReader := bytes.NewReader(fileBytes)

	// Decode file
	// decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	// if err != nil {
	// 	panic("mp3.NewDecoder failed: " + err.Error())
	// }

	// Prepare an Oto context (this will use your default audio device) that will
	// play all our sounds. Its configuration can't be changed later.

	op := &oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	// Create a new 'player' that will handle our sound. Paused by default.
	// player := otoCtx.NewPlayer(decodedMp3)
	noise := Noise{}
	noise.Init()
	fmt.Printf("%v", noise)

	player := otoCtx.NewPlayer(&Noise{})
	// fmt.Printf("%d", player.BufferedSize())
	// buf := make([]byte, 4)
	// for range 10000000 {
	// 	_, err := noise.Read(buf)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// Play starts playing the sound and returns without waiting for it (Play() is async).
	player.Play()
	//
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
}
