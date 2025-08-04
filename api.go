package main

import (
	"bufio"
	"net/http"
	"strings"
	"math/rand"
	"time"
)

type Playlist struct {
	Name  string
	Songs []Moesic
}

type Moesic struct {
	Name string
	Url  string
}

func getMoesic() []Playlist {
	resp, err := http.Get("https://raw.githubusercontent.com/angga7togk/moesic/main/data/moesic.md")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	playlists := []Playlist{}
	currentPlaylist := ""
	for scanner.Scan() {
		line := scanner.Text()

		// if title
		if strings.HasPrefix(line, "###") {
			playlistName := strings.TrimSpace(strings.TrimPrefix(line, "###"))
			currentPlaylist = playlistName

			// if playlist not exists
			// create new playlist
			if !playlistExists(playlists, playlistName) {
				playlists = append(playlists, Playlist{
					Name:  playlistName,
					Songs: []Moesic{},
				})
			}
			continue
		}

		// If song line
		if strings.HasPrefix(line, "- [") {
			nameStart := strings.Index(line, "[")
			nameEnd := strings.Index(line, "]")
			urlStart := strings.Index(line, "(")
			urlEnd := strings.Index(line, ")")

			if nameStart == -1 || nameEnd == -1 || urlStart == -1 || urlEnd == -1 {
				continue // skip malformed line
			}

			name := line[nameStart+1 : nameEnd]
			url := line[urlStart+1 : urlEnd]

			// Add song to current playlist
			for i, p := range playlists {
				if p.Name == currentPlaylist {
					playlists[i].Songs = append(playlists[i].Songs, Moesic{
						Name: name,
						Url:  url,
					})
					break
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return playlists
}

func playlistExists(playlists []Playlist, name string) bool {
	for _, p := range playlists {
		if p.Name == name {
			return true
		}
	}
	return false
}


func FlatSongs(playlists []Playlist) []Moesic {
	var songs []Moesic
	for _, p := range playlists {
		songs = append(songs, p.Songs...)
	}
	return songs
}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandomSong(songs []Moesic) Moesic {
	index := rnd.Intn(len(songs))
	return songs[index]
}