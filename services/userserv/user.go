package userserv

import (
	"github.com/muchlist/berita_acara/dao/userdao"
	"github.com/muchlist/berita_acara/dto"
	"github.com/muchlist/berita_acara/utils/mcrypt"
	"github.com/muchlist/berita_acara/utils/mjwt"
	"github.com/muchlist/berita_acara/utils/rest_err"
	"net/http"
	"time"
)

func NewUserService(dao userdao.UserDaoAssumer, crypto mcrypt.BcryptAssumer, jwt mjwt.JWTAssumer) UserServiceAssumer {
	return &userService{
		dao:    dao,
		crypto: crypto,
		jwt:    jwt,
	}
}

type userService struct {
	dao    userdao.UserDaoAssumer
	crypto mcrypt.BcryptAssumer
	jwt    mjwt.JWTAssumer
}

type UserServiceAssumer interface {
	Login(login dto.UserLoginRequest) (*dto.UserLoginResponse, rest_err.APIError)
	InsertUser(user dto.User) (int, rest_err.APIError)
	EditUser(request dto.User) (*dto.User, rest_err.APIError)
	Refresh(payload dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, rest_err.APIError)
	DeleteUser(userID int) rest_err.APIError
	GetUser(userID int) (*dto.User, rest_err.APIError)
	FindUsers() ([]dto.User, rest_err.APIError)
}

// Login
func (u *userService) Login(login dto.UserLoginRequest) (*dto.UserLoginResponse, rest_err.APIError) {
	user, err := u.dao.Get(login.UserID)
	if err != nil {
		return nil, rest_err.NewBadRequestError("Username atau password tidak valid")
	}

	if !u.crypto.IsPWAndHashPWMatch(login.Password, user.Password) {
		return nil, rest_err.NewUnauthorizedError("Username atau password tidak valid")
	}

	AccessClaims := mjwt.CustomClaim{
		Identity:    user.ID,
		Name:        string(user.Name),
		Roles:       user.Roles,
		ExtraMinute: 60 * 24 * 1, // 1 Hour
		Type:        mjwt.Access,
		Fresh:       true,
	}

	RefreshClaims := mjwt.CustomClaim{
		Identity:    user.ID,
		Name:        string(user.Name),
		Roles:       user.Roles,
		ExtraMinute: 60 * 24 * 10, // 60 days
		Type:        mjwt.Refresh,
	}

	accessToken, err := u.jwt.GenerateToken(AccessClaims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := u.jwt.GenerateToken(RefreshClaims)
	if err != nil {
		return nil, err
	}

	userResponse := dto.UserLoginResponse{
		ID:           user.ID,
		Email:        user.Email,
		Name:         string(user.Name),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expired:      time.Now().Add(time.Minute * time.Duration(60*24*1)).Unix(),
	}

	return &userResponse, nil
}

// InsertUser melakukan register user
func (u *userService) InsertUser(user dto.User) (int, rest_err.APIError) {
	hashPassword, err := u.crypto.GenerateHash(user.Password)
	if err != nil {
		return 0, err
	}

	user.Password = hashPassword
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	insertedUserID, err := u.dao.Insert(user)
	if err != nil {
		return 0, err
	}
	return insertedUserID, nil
}

// EditUser
func (u *userService) EditUser(request dto.User) (*dto.User, rest_err.APIError) {
	request.UpdatedAt = time.Now().Unix()
	result, err := u.dao.Edit(request)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Refresh token
func (u *userService) Refresh(payload dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, rest_err.APIError) {
	token, apiErr := u.jwt.ValidateToken(payload.RefreshToken)
	if apiErr != nil {
		return nil, apiErr
	}
	claims, apiErr := u.jwt.ReadToken(token)
	if apiErr != nil {
		return nil, apiErr
	}

	// cek apakah tipe claims token yang dikirim adalah tipe refresh (1)
	if claims.Type != mjwt.Refresh {
		return nil, rest_err.NewAPIError("Token tidak valid", http.StatusUnprocessableEntity, "jwt_error", []interface{}{"not a refresh token"})
	}

	// mendapatkan data terbaru dari user
	user, apiErr := u.dao.Get(claims.Identity)
	if apiErr != nil {
		return nil, apiErr
	}

	AccessClaims := mjwt.CustomClaim{
		Identity:    user.ID,
		Name:        string(user.Name),
		Roles:       user.Roles,
		ExtraMinute: time.Duration(60 * 60 * 1),
		Type:        mjwt.Access,
		Fresh:       false,
	}

	accessToken, err := u.jwt.GenerateToken(AccessClaims)
	if err != nil {
		return nil, err
	}

	userRefreshTokenResponse := dto.UserRefreshTokenResponse{
		AccessToken: accessToken,
		Expired:     time.Now().Add(time.Minute * time.Duration(60*60*1)).Unix(),
	}

	return &userRefreshTokenResponse, nil
}

// DeleteUser
func (u *userService) DeleteUser(userID int) rest_err.APIError {
	err := u.dao.Delete(userID)
	if err != nil {
		return err
	}
	return nil
}

// GetUser mendapatkan user dari database
func (u *userService) GetUser(userID int) (*dto.User, rest_err.APIError) {
	user, err := u.dao.Get(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindUsers
func (u *userService) FindUsers() ([]dto.User, rest_err.APIError) {
	userList, err := u.dao.Find()
	if err != nil {
		return nil, err
	}
	return userList, nil
}
