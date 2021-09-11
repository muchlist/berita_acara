package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	ID        int             `json:"id"`
	Email     string          `json:"email"`
	Name      UppercaseString `json:"name"`
	Password  string          `json:"-"`
	Roles     []string        `json:"roles"`
	CreatedAt int64           `json:"crated_at"`
	UpdatedAt int64           `json:"updated_at"`
}

func (u *User) Prepare() {
	if u.Roles == nil {
		u.Roles = make([]string, 0)
	}
}

type UserRegisterReq struct {
	ID       int      `json:"id"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
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

// UserLoginResponse balikan user ketika sukses login dengan tambahan AccessToken
type UserLoginRequest struct {
	UserID   int    `json:"user_id"`
	Password string `json:"password"`
}

// UserLoginResponse balikan user ketika sukses login dengan tambahan AccessToken
type UserLoginResponse struct {
	ID           int      `json:"id"`
	Email        string   `json:"email"`
	Name         string   `json:"name"`
	Roles        []string `json:"roles"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Expired      int64    `json:"expired"`
}

type UserRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// UserRefreshTokenResponse mengembalikan token dengan claims yang
// sama dengan token sebelumnya dengan expired yang baru
type UserRefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	Expired     int64  `json:"expired"`
}
