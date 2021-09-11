package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	ID        int             `json:"id" example:"1"`
	Email     string          `json:"email" example:"example@example.com"`
	Name      UppercaseString `json:"name" example:"muchlis"`
	Password  string          `json:"-"`
	Roles     []string        `json:"roles" example:"ADMIN,NORMAL"`
	CreatedAt int64           `json:"created_at" example:"1631341964"`
	UpdatedAt int64           `json:"updated_at" example:"1631341964"`
}

func (u *User) Prepare() {
	if u.Roles == nil {
		u.Roles = make([]string, 0)
	}
}

type UserRegisterReq struct {
	ID       int      `json:"id" example:"1"`
	Email    string   `json:"email" example:"example@example.com"`
	Name     string   `json:"name" example:"muchlis"`
	Password string   `json:"password" example:"password123"`
	Roles    []string `json:"roles" example:"ADMIN,NORMAL,BASIC"`
}

func (u UserRegisterReq) Validate() error {
	if err := validation.ValidateStruct(&u,
		validation.Field(&u.ID, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Roles, validation.NotNil),
		validation.Field(&u.Password, validation.Required, validation.Length(3, 20)),
	); err != nil {
		return err
	}
	return nil
}

type UserEditRequest struct {
	Email string   `json:"email" example:"example@example.com"`
	Name  string   `json:"name" example:"muchlis"`
	Roles []string `json:"roles" example:"ADMIN,NORMAL"`
}

func (u UserEditRequest) Validate() error {
	if err := validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Roles, validation.NotNil),
	); err != nil {
		return err
	}
	return nil
}

// UserLoginResponse balikan user ketika sukses login dengan tambahan AccessToken
type UserLoginRequest struct {
	UserID   int    `json:"user_id"`
	Password string `json:"password"`
}

// UserLoginResponse balikan user ketika sukses login dengan tambahan AccessToken
type UserLoginResponse struct {
	ID           int      `json:"id" example:"1"`
	Email        string   `json:"email" example:"example@example.com"`
	Name         string   `json:"name" example:"muchlis"`
	Roles        []string `json:"roles" example:"ADMIN,NORMAL,BASIC"`
	AccessToken  string   `json:"access_token" example:"eyJhbGciOiJIUzI1N.ywibmFtZSI6IkR5cGUiOjB9.aFjz4esDQ4-_K3dMUmo"`
	RefreshToken string   `json:"refresh_token" example:"eyJhbGciOiJIUzI1N.ywibmFtZSI6IkR5cGUiOjB9.aFjz4esDQ4-_K3dMUmo"`
	Expired      int64    `json:"expired" example:"1631341964"`
}

type UserRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1N.ywibmFtZSI6IkR5cGUiOjB9.aFjz4esDQ4-_K3dMUmo"`
}

// UserRefreshTokenResponse mengembalikan token dengan claims yang
// sama dengan token sebelumnya dengan expired yang baru
type UserRefreshTokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1N.ywibmFtZSI6IkR5cGUiOjB9.aFjz4esDQ4-_K3dMUmo"`
	Expired     int64  `json:"expired" example:"1631341964"`
}
