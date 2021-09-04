package dto

type User struct {
	Username  UppercaseString `json:"username"`
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	Password  string          `json:"-"`
	Role      string          `json:"role"`
	CreatedAt int64           `json:"crated_at"`
	UpdatedAt int64           `json:"updated_at"`
}

type UserRegisterReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UserLoginResponse balikan user ketika sukses login dengan tambahan AccessToken
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserLoginResponse balikan user ketika sukses login dengan tambahan AccessToken
type UserLoginResponse struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Roles        string `json:"role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expired      int64  `json:"expired"`
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
