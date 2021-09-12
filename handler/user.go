package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/berita_acara/dto"
	"github.com/muchlist/berita_acara/payload"
	"github.com/muchlist/berita_acara/services/userserv"
	"github.com/muchlist/berita_acara/utils/mjwt"
	"github.com/muchlist/berita_acara/utils/rest_err"
	"github.com/muchlist/berita_acara/utils/sfunc"
	"strconv"
	"time"
)

func NewUserHandler(userService userserv.UserServiceAssumer) *UserHandler {
	return &UserHandler{
		service: userService,
	}
}

type UserHandler struct {
	service userserv.UserServiceAssumer
}

// Login login
// @Summary login
// @Description login menggunakan userID dan password untuk mendapatkan JWT Token
// @ID user-login
// @Accept json
// @Produce json
// @Tags Access
// @Param ReqBody body dto.UserLoginRequest true "Body raw JSON"
// @Success 200 {object} payload.RespWrap{data=dto.UserLoginResponse}
// @Failure 400 {object} payload.RespWrap{error=payload.ErrorExample400}
// @Failure 500 {object} payload.RespWrap{error=payload.ErrorExample500}
// @Router /api/v1/login [post]
func (u *UserHandler) Login(c *fiber.Ctx) error {
	var login dto.UserLoginRequest
	if err := c.BodyParser(&login); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if login.UserID == 0 || login.Password == "" {
		apiErr := rest_err.NewBadRequestError("user_id atau password tidak boleh kosong")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	response, apiErr := u.service.Login(c.Context(), login)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

// Register menambahkan user
// @Summary register user
// @Description added user to repository
// @ID user-register
// @Accept json
// @Produce json
// @Tags Access
// @Security BearerAuth
// @Param ReqBody body dto.UserRegisterReq true "Body raw JSON"
// @Success 200 {object} payload.RespMsgExample
// @Failure 400 {object} payload.RespWrap{error=payload.ErrorExample400}
// @Failure 500 {object} payload.RespWrap{error=payload.ErrorExample500}
// @Router /api/v1/users [post]
func (u *UserHandler) Register(c *fiber.Ctx) error {
	var user dto.UserRegisterReq
	if err := c.BodyParser(&user); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := user.Validate(); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	insertUserID, apiErr := u.service.InsertUser(c.Context(), dto.User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      dto.UppercaseString(user.Name),
		Password:  user.Password,
		Roles:     user.Roles,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Register berhasil, ID: %d", insertUserID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

// Edit
// @Summary edit user
// @Description melakukan perubahan data pada user
// @ID user-edit
// @Accept json
// @Produce json
// @Tags Access
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param ReqBody body dto.UserEditRequest true "Body raw JSON"
// @Success 200 {object} payload.RespWrap{data=dto.User}
// @Failure 400 {object} payload.RespWrap{error=payload.ErrorExample400}
// @Failure 500 {object} payload.RespWrap{error=payload.ErrorExample500}
// @Router /api/v1/users/{id} [put]
func (u *UserHandler) Edit(c *fiber.Ctx) error {
	userID := c.Params("id")

	var req dto.UserEditRequest
	var err error
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	userEdited, apiErr := u.service.EditUser(c.Context(), dto.User{
		ID:    userIDInt,
		Email: req.Email,
		Name:  dto.UppercaseString(req.Name),
		Roles: req.Roles,
	})
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": userEdited})
}

// RefreshToken
// @Summary refresh token
// @Description mendapatkan token dengan tambahan waktu expired menggunakan refresh token
// @ID user-refresh
// @Accept json
// @Produce json
// @Tags Access
// @Param ReqBody body dto.UserRefreshTokenRequest true "Body raw JSON"
// @Success 200 {object} payload.RespWrap{data=dto.UserRefreshTokenResponse}
// @Failure 400 {object} payload.RespWrap{error=payload.ErrorExample400}
// @Failure 500 {object} payload.RespWrap{error=payload.ErrorExample500}
// @Router /api/v1/refresh [post]
func (u *UserHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.UserRefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	response, apiErr := u.service.Refresh(c.Context(), req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

// Delete menghapus user
// @Summary delete user by ID
// @Description menghapus user berdasarkan userID
// @ID user-delete
// @Accept json
// @Produce json
// @Tags Access
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} payload.RespMsgExample
// @Failure 400 {object} payload.RespWrap{error=payload.ErrorExample400}
// @Failure 500 {object} payload.RespWrap{error=payload.ErrorExample500}
// @Router /api/v1/users/{id} [delete]
func (u *UserHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	userID := c.Params("id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if claims.Identity == userIDInt {
		apiErr := rest_err.NewBadRequestError("Tidak dapat menghapus akun terkait (diri sendiri)!")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	apiErr := u.service.DeleteUser(c.Context(), userIDInt)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("user %d berhasil dihapus", userIDInt)})
}

// Get menampilkan user berdasarkan username
// @Summary get user by ID
// @Description menampilkan user berdasarkan userID
// @ID user-get
// @Accept json
// @Produce json
// @Tags Access
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} payload.RespWrap{data=dto.User}
// @Failure 400 {object} payload.RespWrap{error=payload.ErrorExample400}
// @Failure 500 {object} payload.RespWrap{error=payload.ErrorExample500}
// @Router /api/v1/users/{id} [get]
func (u *UserHandler) Get(c *fiber.Ctx) error {
	userID := c.Params("id")

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	user, apiErr := u.service.GetUser(c.Context(), userIDInt)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(payload.RespWrap{
		Data:  user,
		Error: nil,
	})
}

// GetProfile mengembalikan user yang sedang login
// @Summary get current profile
// @Description menampilkan profile berdasarkan user yang login saat ini
// @ID user-profile
// @Accept json
// @Produce json
// @Tags Access
// @Security BearerAuth
// @Success 200 {object} payload.RespWrap{data=dto.User}
// @Router /api/v1/profile [get]
func (u *UserHandler) GetProfile(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	user, apiErr := u.service.GetUser(c.Context(), claims.Identity)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": user})
}

// Find menampilkan list user
// @Summary find user
// @Description menampilkan daftar user
// @ID user-find
// @Accept json
// @Produce json
// @Tags Access
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param last_id query int false "Last ID sebagai cursor untuk page selanjutnya"
// @Param search query string false "Search apabila di isi akan melakukan pencarian string include"
// @Success 200 {object} payload.RespWrap{data=[]dto.User}
// @Failure 400 {object} payload.RespWrap{error=payload.ErrorExample400}
// @Failure 500 {object} payload.RespWrap{error=payload.ErrorExample500}
// @Router /api/v1/users [get]
func (u *UserHandler) Find(c *fiber.Ctx) error {
	limit := sfunc.StrToInt(c.Query("limit"), 10)
	cursor := sfunc.StrToInt(c.Query("last_id"), 0)
	search := c.Query("search")

	userList, apiErr := u.service.FindUsers(c.Context(), search, limit, cursor)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if userList == nil {
		userList = []dto.User{}
	}
	return c.JSON(fiber.Map{"error": nil, "data": userList})
}
