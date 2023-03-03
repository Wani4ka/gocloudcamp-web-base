package playlist

import "fmt"

type NoSuchSongError struct {
	SongId SongId
}

func NewNoSuchSongError(songId SongId) error {
	return NoSuchSongError{SongId: songId}
}
func (err NoSuchSongError) Error() string {
	return fmt.Sprintf("song with ID %v does not exist", err.SongId)
}

type SongIsCurrentlyPlayingError struct {
	SongId SongId
}

func NewSongIsCurrentlyPlayingError(songId SongId) error {
	return SongIsCurrentlyPlayingError{SongId: songId}
}
func (err SongIsCurrentlyPlayingError) Error() string {
	return fmt.Sprintf("can't operate with song %v while it's playing", err.SongId)
}

type InvalidSongError struct{}

func NewInvalidSongError() error {
	return InvalidSongError{}
}
func (err InvalidSongError) Error() string {
	return "invalid song"
}

type EmptyPlaylistError struct{}

func NewEmptyPlaylistError() error {
	return EmptyPlaylistError{}
}
func (err EmptyPlaylistError) Error() string {
	return "playlist is empty or its end reached"
}
