package util

import (
	"gocloudcamp/core/playlist"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getCode(err error) codes.Code {
	switch err.(type) {
	case playlist.NoSuchSongError:
		return codes.NotFound
	case playlist.SongIsCurrentlyPlayingError:
		return codes.Unavailable
	case playlist.InvalidSongError:
		return codes.InvalidArgument
	case playlist.EmptyPlaylistError:
		return codes.OutOfRange
	default:
		return codes.Internal
	}
}

func WrapErrorToGRPC(err error) error {
	return status.Error(getCode(err), err.Error())
}
