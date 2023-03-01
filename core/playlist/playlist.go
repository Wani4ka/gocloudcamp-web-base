package playlist

import (
	"errors"
	"gocloudcamp/core/song"
	"math/rand"
)

type Song struct {
	previous *Song
	data     song.Song
	id       uint32
	next     *Song
}

func (song *Song) define(data song.Song) {
	song.data = data
	song.id = rand.Uint32()
	if song.next == nil {
		song.next = &Song{previous: song}
	}
}

type playlist struct {
	currentSong *Song
	lastSong    *Song
	timer       Timer
	storage     map[uint32]*Song
}

func NewPlaylist() Playlist {
	singleSong := &Song{}
	return &playlist{
		timer:       NewTimer(),
		currentSong: singleSong,
		lastSong:    singleSong,
		storage:     make(map[uint32]*Song),
	}
}

func (playlist *playlist) Play() {
	if playlist.timer.IsPaused() {
		playlist.timer.Resume()
	} else if playlist.currentSong != nil && playlist.currentSong.data.IsValid() {
		playlist.timer.Schedule(playlist.currentSong.data.Length, playlist.Next)
	}
}

func (playlist *playlist) Pause() {
	if !playlist.timer.IsPaused() {
		playlist.timer.Pause()
	}
}

func (playlist *playlist) AddSong(song song.Song) (uint32, error) {
	if !song.IsValid() {
		return 0, errors.New("invalid song")
	}
	playlist.lastSong.define(song)
	id := playlist.lastSong.id
	playlist.storage[id] = playlist.lastSong
	playlist.lastSong = playlist.lastSong.next
	return id, nil
}

func (playlist *playlist) GetSong(id uint32) (song.Song, bool) {
	sng, exists := playlist.storage[id]
	if !exists || sng == nil || !sng.data.IsValid() {
		return song.Song{}, false
	}
	return sng.data, true
}

func (playlist *playlist) ReplaceSong(id uint32, song song.Song) error {
	if !song.IsValid() {
		return errors.New("invalid song")
	}
	sng := playlist.storage[id]
	if sng == nil {
		return NewNoSuchSongError(id)
	}
	sng.data = song
	return nil
}

func (playlist *playlist) RemoveSong(id uint32) (song.Song, error) {
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

func (playlist *playlist) Next() {
	playlist.timer.Stop()
	if playlist.currentSong.next != nil {
		playlist.currentSong = playlist.currentSong.next
	}
	if playlist.currentSong.data.IsValid() {
		playlist.Play()
	}
}

func (playlist *playlist) Prev() {
	playlist.timer.Stop()
	if playlist.currentSong.previous != nil {
		playlist.currentSong = playlist.currentSong.previous
		playlist.Play()
	}
}
