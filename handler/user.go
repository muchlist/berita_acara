package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/berita_acara/dto"
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

	response, apiErr := u.service.Login(login)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

// Register menambahkan user
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

	insertUserID, apiErr := u.service.InsertUser(dto.User{
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

// Edit mengedit user
func (u *UserHandler) Edit(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user dto.User
	var err error
	user.ID, err = strconv.Atoi(userID)
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := c.BodyParser(&user); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	userEdited, apiErr := u.service.EditUser(user)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": userEdited})
}

// RefreshToken
func (u *UserHandler) RefreshToken(c *fiber.Ctx) error {
	var payload dto.UserRefreshTokenRequest
	if err := c.BodyParser(&payload); err != nil {
		apiErr := rest_err.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	response, apiErr := u.service.Refresh(payload)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

// Delete menghapus user, idealnya melalui middleware is_admin
func (u *UserHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	userID := c.Params("user_id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if claims.Identity == userIDInt {
		apiErr := rest_err.NewBadRequestError("Tidak dapat menghapus akun terkait (diri sendiri)!")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	apiErr := u.service.DeleteUser(userIDInt)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("user %d berhasil dihapus", userIDInt)})
}

// Get menampilkan user berdasarkan username
func (u *UserHandler) Get(c *fiber.Ctx) error {
	userID := c.Params("id")

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		apiErr := rest_err.NewBadRequestError("kesalahan input, id harus berupa angka")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	user, apiErr := u.service.GetUser(userIDInt)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": user})
}

// GetProfile mengembalikan user yang sedang login
func (u *UserHandler) GetProfile(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	user, apiErr := u.service.GetUser(claims.Identity)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": user})
}

// Find menampilkan list user
func (u *UserHandler) Find(c *fiber.Ctx) error {
	limit := sfunc.StrToInt(c.Query("limit"), 10)
	cursor := sfunc.StrToInt(c.Query("last_id"), 0)
	search := c.Query("search")

	userList, apiErr := u.service.FindUsers(search, limit, cursor)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if userList == nil {
		userList = []dto.User{}
	}
	return c.JSON(fiber.Map{"error": nil, "data": userList})
}
