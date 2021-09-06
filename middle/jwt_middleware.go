package middle

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/berita_acara/utils/mjwt"
	"github.com/muchlist/berita_acara/utils/rest_err"
	"strings"
)

var (
	jwt = mjwt.NewJwt()
)

const (
	headerKey = "Authorization"
	bearerKey = "Bearer"
)

// NormalAuth memerlukan salah satu role inputan agar diloloskan ke proses berikutnya
// token tidak perlu fresh
func NormalAuth(rolesReq ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(headerKey)
		claims, err := authHaveRoleValidator(authHeader, false, rolesReq)
		if err != nil {
			return c.Status(err.Status()).JSON(fiber.Map{"error": err, "data": nil})
		}
		c.Locals(mjwt.CLAIMS, claims)
		return c.Next()
	}
}

// FreshAuth memerlukan salah satu role inputan agar diloloskan ke proses berikutnya
// token harus fresh (tidak hasil dari refresh token)
func FreshAuth(rolesReq ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(headerKey)
		claims, err := authHaveRoleValidator(authHeader, true, rolesReq)
		if err != nil {
			return c.Status(err.Status()).JSON(fiber.Map{"error": err, "data": nil})
		}

		c.Locals(mjwt.CLAIMS, claims)
		return c.Next()
	}
}

func authHaveRoleValidator(authHeader string, mustFresh bool, rolesAllowed []string) (*mjwt.CustomClaim, rest_err.APIError) {
	if !strings.Contains(authHeader, bearerKey) {
		apiErr := rest_err.NewUnauthorizedError("Unauthorized")
		return nil, apiErr
	}
	tokenString := strings.Split(authHeader, " ")
	if len(tokenString) != 2 {
		apiErr := rest_err.NewUnauthorizedError("Unauthorized")
		return nil, apiErr
	}
	token, apiErr := jwt.ValidateToken(tokenString[1])
	if apiErr != nil {
		return nil, apiErr
	}
	claims, apiErr := jwt.ReadToken(token)
	if apiErr != nil {
		return nil, apiErr
	}
	if mustFresh {
		if !claims.Fresh {
			apiErr := rest_err.NewUnauthorizedError("Memerlukan token yang baru untuk mengakses halaman ini")
			return nil, apiErr
		}
	}

	//if len(rolesAllowed) != 0 {
	//	if sfunc.InSlice(claims.Roles, rolesAllowed) {
	//		return claims, nil
	//	}
	//} else {
	//	return claims, nil
	//}

	apiErr = rest_err.NewUnauthorizedError(fmt.Sprintf("Unauthorized, memerlukan hak akses %s", rolesAllowed))
	return nil, apiErr
}
