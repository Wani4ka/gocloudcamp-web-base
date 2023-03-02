package playlist

import (
	"gocloudcamp/core/song"
	"time"
)

type SongId uint32

type Playlist interface {
	Play()
	Pause()
	AddSong(song song.Song) (SongId, error)
	Next() (SongId, error)
	Prev() (SongId, error)
	GetSong(id SongId) (song.Song, bool)
	ReplaceSong(id SongId, song song.Song) error
	RemoveSong(id SongId) (song.Song, error)
	IsPlaying() bool
	GetNowPlaying() (SongId, song.Song, time.Duration, bool)
}
