package server

import (
	"context"
	"gocloudcamp/core/playlist"
	"gocloudcamp/core/song"
	"gocloudcamp/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type PlaylistServerImpl struct {
	proto.UnimplementedPlaylistServer
	stored playlist.Playlist
}

func NewPlaylistServer() *PlaylistServerImpl {
	list := playlist.NewPlaylist()
	return &PlaylistServerImpl{
		stored: list,
	}
}

func (server PlaylistServerImpl) AddSong(_ context.Context, data *proto.Song) (*proto.SongLocation, error) {
	added, err := server.stored.AddSong(*song.NewSong(data.GetName(), time.Duration(data.GetSeconds())*time.Second))
	if err != nil {
		return nil, err
	}
	return &proto.SongLocation{Id: added}, nil
}

func (server PlaylistServerImpl) UpdateSong(_ context.Context, req *proto.UpdateSongRequest) (*emptypb.Empty, error) {
	err := server.stored.ReplaceSong(req.GetLocation().GetId(), *song.NewSong(req.GetData().GetName(), time.Duration(req.GetData().GetSeconds())*time.Second))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (server PlaylistServerImpl) GetSong(_ context.Context, req *proto.SongLocation) (*proto.Song, error) {
	result, exists := server.stored.GetSong(req.GetId())
	if !exists {
		return nil, playlist.NewNoSuchSongError(req.GetId())
	}
	return &proto.Song{Name: result.Name, Seconds: uint32(result.Length.Seconds())}, nil
}

func (server PlaylistServerImpl) DeleteSong(_ context.Context, req *proto.SongLocation) (*proto.Song, error) {
	result, err := server.stored.RemoveSong(req.GetId())
	if err != nil {
		return nil, err
	}
	return &proto.Song{Name: result.Name, Seconds: uint32(result.Length.Seconds())}, nil
}
