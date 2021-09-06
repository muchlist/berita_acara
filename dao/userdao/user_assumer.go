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
	Insert(user dto.User) (int, rest_err.APIError)
	Edit(userInput dto.User) (*dto.User, rest_err.APIError)
	Delete(id int) rest_err.APIError
	ChangePassword(input dto.User) rest_err.APIError
}

type UserReader interface {
	Get(id int) (*dto.User, rest_err.APIError)
	Find() ([]dto.User, rest_err.APIError)
}
