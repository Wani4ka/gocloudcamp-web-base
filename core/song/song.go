package song

import (
	"time"
)

type Song struct {
	Name   string
	Length time.Duration
}

func NewSong(name string, length time.Duration) *Song {
	return &Song{
		Name:   name,
		Length: length,
	}
}

func (song Song) Equal(another Song) bool {
	return song.Name == another.Name && song.Length == another.Length
}

func (song Song) IsValid() bool {
	return song.Length > 0
}
