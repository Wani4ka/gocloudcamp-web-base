package playlist

import (
	"errors"
	"gocloudcamp/core/song"
	"math/rand"
	"time"
)

type Song struct {
	previous *Song
	data     song.Song
	id       SongId
	next     *Song
}

func (song *Song) define(data song.Song) {
	song.data = data
	song.id = SongId(rand.Uint32())
	if song.next == nil {
		song.next = &Song{previous: song}
	}
}

type playlist struct {
	currentSong *Song
	lastSong    *Song
	timer       Timer
	storage     map[SongId]*Song
}

func NewPlaylist() Playlist {
	singleSong := &Song{}
	return &playlist{
		timer:       NewTimer(),
		currentSong: singleSong,
		lastSong:    singleSong,
		storage:     make(map[SongId]*Song),
	}
}

func (playlist *playlist) IsPlaying() bool {
	return playlist.timer.IsScheduled() && !playlist.timer.IsPaused()
}

func (playlist *playlist) GetNowPlaying() (SongId, song.Song, time.Duration, bool) {
	if playlist.currentSong == nil || !playlist.currentSong.data.IsValid() {
		return 0, song.Song{}, 0, false
	}
	return playlist.currentSong.id, playlist.currentSong.data, playlist.timer.ElapsedTime(), playlist.IsPlaying()
}

func (playlist *playlist) Play() {
	if playlist.timer.IsPaused() {
		playlist.timer.Resume()
	} else if playlist.currentSong != nil && playlist.currentSong.data.IsValid() {
		playlist.timer.Schedule(playlist.currentSong.data.Length, playlist.next)
	}
}

func (playlist *playlist) Pause() {
	if !playlist.timer.IsPaused() {
		playlist.timer.Pause()
	}
}

func (playlist *playlist) AddSong(song song.Song) (SongId, error) {
	if !song.IsValid() {
		return 0, errors.New("invalid song")
	}
	playlist.lastSong.define(song)
	id := playlist.lastSong.id
	playlist.storage[id] = playlist.lastSong
	playlist.lastSong = playlist.lastSong.next
	return id, nil
}

func (playlist *playlist) GetSong(id SongId) (song.Song, bool) {
	sng, exists := playlist.storage[id]
	if !exists || sng == nil || !sng.data.IsValid() {
		return song.Song{}, false
	}
	return sng.data, true
}

func (playlist *playlist) ReplaceSong(id SongId, song song.Song) error {
	if !song.IsValid() {
		return errors.New("invalid song")
	}
	sng := playlist.storage[id]
	if sng == nil {
		return NewNoSuchSongError(id)
	}
	if playlist.currentSong == sng {
		return NewSongIsCurrentlyPlayingError(id)
	}
	sng.data = song
	return nil
}

func (playlist *playlist) RemoveSong(id SongId) (song.Song, error) {
	sng, exists := playlist.storage[id]
	if !exists {
		return song.Song{}, NewNoSuchSongError(id)
	}
	if playlist.currentSong == sng {
		return song.Song{}, NewSongIsCurrentlyPlayingError(id)
	}
	if sng.previous != nil && sng.next != nil {
		sng.previous.next = sng.next
		sng.next.previous = sng.previous
	}
	playlist.storage[id] = nil
	return sng.data, nil
}

func (playlist *playlist) Next() (SongId, error) {
	playlist.timer.Stop()
	if playlist.currentSong.next != nil {
		playlist.currentSong = playlist.currentSong.next
	}
	if playlist.currentSong.data.IsValid() {
		playlist.Play()
		return playlist.currentSong.id, nil
	}
	return 0, errors.New("playlist is empty")
}

func (playlist *playlist) next() {
	_, _ = playlist.Next()
}

func (playlist *playlist) Prev() (SongId, error) {
	playlist.timer.Stop()
	if playlist.currentSong.previous != nil {
		playlist.currentSong = playlist.currentSong.previous
	}
	if playlist.currentSong.data.IsValid() {
		playlist.Play()
		return playlist.currentSong.id, nil
	}
	return 0, errors.New("playlist is empty")
}
