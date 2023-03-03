package playlist

import (
	"bytes"
	"encoding/gob"
	songmodule "gocloudcamp/core/song"
	"os"
	"sync"
)

type Storage interface {
	Save(pl Playlist) error
	Load() (Playlist, error)
}

type storage struct {
	dir  string
	path string
	lock sync.Mutex
}

func NewStorage(dir string, file string) Storage {
	return &storage{
		dir:  dir,
		path: dir + "/" + file,
	}
}

func (st *storage) createEmptyPlaylist() *playlist {
	pl := NewPlaylist().(*playlist)
	pl.storage = st
	return pl
}

func (st *storage) mkdirs() error {
	return os.MkdirAll(st.dir, os.ModePerm)
}

func (st *storage) toBytes(pl *playlist) (buf bytes.Buffer, err error) {
	enc := gob.NewEncoder(&buf)
	song := pl.currentSong
	for song.previous != nil {
		song = song.previous
	}
	for song.data.IsValid() {
		err := enc.Encode(song.data)
		if err != nil {
			return bytes.Buffer{}, err
		}
		err = enc.Encode(song == pl.currentSong)
		if err != nil {
			return bytes.Buffer{}, err
		}
		song = song.next
	}
	return
}

func (st *storage) save(pl Playlist) error {
	err := st.mkdirs()
	if err != nil {
		return err
	}
	file, err := os.Create(st.path)
	if err != nil {
		return err
	}
	buf, err := st.toBytes(pl.(*playlist))
	if err != nil {
		return err
	}
	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (st *storage) Save(pl Playlist) error {
	st.lock.Lock()
	err := st.save(pl)
	st.lock.Unlock()
	return err
}

func (st *storage) fromBytes(buf *bytes.Buffer) (Playlist, error) {
	dec := gob.NewDecoder(buf)
	pl := NewPlaylist().(*playlist)
	var song songmodule.Song
	var current bool
	var currentSong *Song
	for dec.Decode(&song) == nil {
		_, err := pl.AddSong(song)
		if err != nil {
			return st.createEmptyPlaylist(), err
		}
		err = dec.Decode(&current)
		if err != nil {
			return st.createEmptyPlaylist(), err
		}
		if current {
			currentSong = pl.lastSong.previous
		}
	}
	pl.currentSong = currentSong
	pl.storage = st
	return pl, nil
}

func (st *storage) load() (Playlist, error) {
	err := st.mkdirs()
	if err != nil {
		return st.createEmptyPlaylist(), err
	}
	stream, err := os.ReadFile(st.path)
	if err != nil {
		return st.createEmptyPlaylist(), err
	}
	return st.fromBytes(bytes.NewBuffer(stream))
}

func (st *storage) Load() (Playlist, error) {
	st.lock.Lock()
	pl, err := st.load()
	st.lock.Unlock()
	return pl, err
}
