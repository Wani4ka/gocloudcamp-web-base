package playlist

import (
	"gocloudcamp/core/song"
)

type Playlist interface {
	Play()
	Pause()
	AddSong(song song.Song) (uint32, error)
	Next()
	Prev()
	GetSong(id uint32) (song.Song, bool)
	ReplaceSong(id uint32, song song.Song) error
	RemoveSong(id uint32) (song.Song, error)
}
