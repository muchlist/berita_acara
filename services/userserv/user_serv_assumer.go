package userserv

import (
	"context"
	"github.com/muchlist/berita_acara/dto"
	"github.com/muchlist/berita_acara/utils/rest_err"
)

type UserServiceAssumer interface {
	UserServiceAccess
	UserServiceModifier
	UserServiceReader
}

type UserServiceReader interface {
	GetUser(ctx context.Context, userID int) (*dto.User, rest_err.APIError)
	FindUsers(ctx context.Context, search string, limit int, cursor int) ([]dto.User, rest_err.APIError)
}

type UserServiceAccess interface {
	Login(ctx context.Context, login dto.UserLoginRequest) (*dto.UserLoginResponse, rest_err.APIError)
	Refresh(ctx context.Context, payload dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, rest_err.APIError)
}

type UserServiceModifier interface {
	InsertUser(ctx context.Context, user dto.User) (int, rest_err.APIError)
	EditUser(ctx context.Context, request dto.User) (*dto.User, rest_err.APIError)
	DeleteUser(ctx context.Context, userID int) rest_err.APIError
}
