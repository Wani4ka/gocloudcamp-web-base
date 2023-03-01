package seek

import (
	"context"
	"gocloudcamp/core/playlist"
	"gocloudcamp/proto"
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

func (ServerImpl) Prev(context.Context, *emptypb.Empty) (*proto.PlaylistEntry, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Prev not implemented")
}
func (ServerImpl) Next(context.Context, *emptypb.Empty) (*proto.PlaylistEntry, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Next not implemented")
}
func (ServerImpl) NowPlaying(context.Context, *emptypb.Empty) (*proto.NowPlayingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NowPlaying not implemented")
}
