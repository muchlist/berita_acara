package userserv

import (
	"github.com/muchlist/berita_acara/dto"
	"github.com/muchlist/berita_acara/utils/rest_err"
)

type UserServiceAssumer interface {
	UserServiceAccess
	UserServiceModifier
	UserServiceReader
}

type UserServiceReader interface {
	GetUser(userID int) (*dto.User, rest_err.APIError)
	FindUsers(search string, limit int, cursor int) ([]dto.User, rest_err.APIError)
}

type UserServiceAccess interface {
	Login(login dto.UserLoginRequest) (*dto.UserLoginResponse, rest_err.APIError)
	Refresh(payload dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, rest_err.APIError)
}

type UserServiceModifier interface {
	InsertUser(user dto.User) (int, rest_err.APIError)
	EditUser(request dto.User) (*dto.User, rest_err.APIError)
	DeleteUser(userID int) rest_err.APIError
}
