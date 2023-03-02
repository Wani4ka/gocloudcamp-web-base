package util

import (
	"gocloudcamp/core/song"
	"gocloudcamp/proto"
)

func SongToProto(song song.Song) *proto.Song {
	return &proto.Song{
		Name:    song.Name,
		Seconds: uint32(song.Length.Seconds()),
	}
}
