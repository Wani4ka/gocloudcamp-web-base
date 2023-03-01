package playlist

import (
	"gocloudcamp/core/song"
	"time"
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
	IsPlaying() bool
	GetNowPlaying() (song.Song, time.Duration, bool)
}
