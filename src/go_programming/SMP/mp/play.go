package mp

import "fmt"

type Player interface {
	Play(source string) float64
}

func Play(source, mtype string, ch chan float64) {
	var p Player
	switch mtype {
	case "MP3":
		p = &MP3Player{}
	case "WAV":
		p = &WAVPlayer{}
	default:
		fmt.Println("Unsupported music type", mtype)
		return
	}
	t := p.Play(source)
	fmt.Println("Play times ", t)
	ch <- t
}
