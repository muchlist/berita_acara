package userdao

import (
	"github.com/muchlist/berita_acara/dto"
	"github.com/muchlist/berita_acara/utils/rest_err"
)

type UserDaoAssumer interface {
	UserSaver
	UserReader
}

type UserSaver interface {
	Insert(user dto.User) (*string, rest_err.APIError)
	Edit(userInput dto.User) (*dto.User, rest_err.APIError)
	Delete(userName string) rest_err.APIError
	ChangePassword(input dto.User) (*dto.User, rest_err.APIError)
}

type UserReader interface {
	Get(userName string) (*dto.User, rest_err.APIError)
	Find() ([]dto.User, rest_err.APIError)
}
