package main

import (
	"log"

	vlc "github.com/adrg/libvlc-go/v3"
)

func playVLC(url string, onEnd func()) *vlc.Player {
	if err := vlc.Init("--no-xlib", "--quiet"); err != nil {
		log.Fatalf("libVLC init error: %v", err)
	}

	// Player
	p, err := vlc.NewPlayer()
	if err != nil {
		log.Fatalf("new player error: %v", err)
	}

	// Load media
	media, err := p.LoadMediaFromURL(url)
	if err != nil {
		log.Fatalf("load media error: %v", err)
	}
	defer media.Release()

	// Attach event selesai lagu
	em, err := p.EventManager()
	if err != nil {
		log.Fatalf("event manager error: %v", err)
	}

	_, err = em.Attach(vlc.MediaPlayerEndReached, func(e vlc.Event, i interface{}) {
		onEnd()
	}, em)
	if err != nil {
		log.Fatalf("attach event error: %v", err)
	}

	// Play
	if err := p.Play(); err != nil {
		log.Fatalf("play error: %v", err)
	}

	return p
}
