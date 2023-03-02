package crud

import (
	"context"
	"gocloudcamp/core/playlist"
	"gocloudcamp/core/song"
	"gocloudcamp/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type ServerImpl struct {
	proto.UnimplementedCRUDServer
	stored playlist.Playlist
}

func NewServer(stored playlist.Playlist) *ServerImpl {
	return &ServerImpl{
		stored: stored,
	}
}

func (server *ServerImpl) AddSong(_ context.Context, data *proto.Song) (*proto.SongLocation, error) {
	added, err := server.stored.AddSong(*song.NewSong(data.GetName(), time.Duration(data.GetSeconds())*time.Second))
	if err != nil {
		return nil, err
	}
	return &proto.SongLocation{Id: uint32(added)}, nil
}

func (server *ServerImpl) UpdateSong(_ context.Context, req *proto.PlaylistEntry) (*emptypb.Empty, error) {
	err := server.stored.ReplaceSong(playlist.SongId(req.GetLocation().GetId()), *song.NewSong(req.GetData().GetName(), time.Duration(req.GetData().GetSeconds())*time.Second))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (server *ServerImpl) GetSong(_ context.Context, req *proto.SongLocation) (*proto.Song, error) {
	result, exists := server.stored.GetSong(playlist.SongId(req.GetId()))
	if !exists {
		return nil, playlist.NewNoSuchSongError(playlist.SongId(req.GetId()))
	}
	return &proto.Song{Name: result.Name, Seconds: uint32(result.Length.Seconds())}, nil
}

func (server *ServerImpl) DeleteSong(_ context.Context, req *proto.SongLocation) (*proto.Song, error) {
	result, err := server.stored.RemoveSong(playlist.SongId(req.GetId()))
	if err != nil {
		return nil, err
	}
	return &proto.Song{Name: result.Name, Seconds: uint32(result.Length.Seconds())}, nil
}
