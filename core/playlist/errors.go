package playlist

import "fmt"

type NoSuchSongError struct {
	SongId uint32
}

func NewNoSuchSongError(songId uint32) NoSuchSongError {
	return NoSuchSongError{SongId: songId}
}
func (err NoSuchSongError) Error() string {
	return fmt.Sprintf("Song with ID %v does not exist", err.SongId)
}

type SongIsCurrentlyPlayingError struct {
	SongId uint32
}

func NewSongIsCurrentlyPlayingError(songId uint32) SongIsCurrentlyPlayingError {
	return SongIsCurrentlyPlayingError{SongId: songId}
}
func (err SongIsCurrentlyPlayingError) Error() string {
	return fmt.Sprintf("Can't operate with song %v while it's playing", err.SongId)
}
