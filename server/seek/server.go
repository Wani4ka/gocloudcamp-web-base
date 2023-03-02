package seek

import (
	"context"
	"gocloudcamp/core/playlist"
	"gocloudcamp/proto"
	"gocloudcamp/server/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServerImpl struct {
	proto.UnimplementedSeekServer
	stored playlist.Playlist
}

func NewServer(stored playlist.Playlist) *ServerImpl {
	return &ServerImpl{
		stored: stored,
	}
}

func (server *ServerImpl) Prev(context.Context, *emptypb.Empty) (*proto.PlaylistEntry, error) {
	id, err := server.stored.Prev()
	if err != nil {
		return nil, err
	}
	song, _ := server.stored.GetSong(id)
	return &proto.PlaylistEntry{
		Location: &proto.SongLocation{Id: uint32(id)},
		Data:     util.SongToProto(song),
	}, nil
}

func (server *ServerImpl) Next(context.Context, *emptypb.Empty) (*proto.PlaylistEntry, error) {
	id, err := server.stored.Next()
	if err != nil {
		return nil, err
	}
	song, _ := server.stored.GetSong(id)
	return &proto.PlaylistEntry{
		Location: &proto.SongLocation{Id: uint32(id)},
		Data:     util.SongToProto(song),
	}, nil
}
func (server *ServerImpl) NowPlaying(context.Context, *emptypb.Empty) (*proto.NowPlayingResponse, error) {
	id, song, elapsed, playing := server.stored.GetNowPlaying()
	if !song.IsValid() {
		return nil, status.Errorf(codes.OutOfRange, "playlist is empty")
	}
	return &proto.NowPlayingResponse{
		Entry: &proto.PlaylistEntry{
			Location: &proto.SongLocation{Id: uint32(id)},
			Data:     util.SongToProto(song),
		},
		Elapsed: uint32(elapsed.Seconds()),
		Playing: playing,
	}, nil
}

func (server *ServerImpl) Play(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	server.stored.Play()
	return &emptypb.Empty{}, nil
}
func (server *ServerImpl) Pause(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	server.stored.Pause()
	return &emptypb.Empty{}, nil
}
