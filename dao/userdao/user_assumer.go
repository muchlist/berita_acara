package userdao

import (
	"context"
	"github.com/muchlist/berita_acara/dto"
	"github.com/muchlist/berita_acara/utils/rest_err"
)

type UserDaoAssumer interface {
	UserSaver
	UserReader
}

type UserSaver interface {
	Insert(ctx context.Context, user dto.User) (int, rest_err.APIError)
	Edit(ctx context.Context, userInput dto.User) (*dto.User, rest_err.APIError)
	Delete(ctx context.Context, id int) rest_err.APIError
	ChangePassword(ctx context.Context, input dto.User) rest_err.APIError
}

type UserReader interface {
	Get(ctx context.Context, id int) (*dto.User, rest_err.APIError)
	FindWithCursor(ctx context.Context, search string, limit uint64, cursor int) ([]dto.User, rest_err.APIError)
}
