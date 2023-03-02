package playlist

import "fmt"

type NoSuchSongError struct {
	SongId SongId
}

func NewNoSuchSongError(songId SongId) NoSuchSongError {
	return NoSuchSongError{SongId: songId}
}
func (err NoSuchSongError) Error() string {
	return fmt.Sprintf("Song with ID %v does not exist", err.SongId)
}

type SongIsCurrentlyPlayingError struct {
	SongId SongId
}

func NewSongIsCurrentlyPlayingError(songId SongId) SongIsCurrentlyPlayingError {
	return SongIsCurrentlyPlayingError{SongId: songId}
}
func (err SongIsCurrentlyPlayingError) Error() string {
	return fmt.Sprintf("Can't operate with song %v while it's playing", err.SongId)
}
