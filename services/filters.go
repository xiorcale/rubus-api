package services

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/astaxie/beego"
	beegoCtx "github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/kjuvi/rubus-api/models"
)

// FilterAdmin checks if the `User` has an administrator `Role`. If not, return
// an Unauthorized `JSONError`.
func FilterAdmin(c *beego.Controller) {
	claims := c.Ctx.Request.Context().Value("claims").(*models.Claims)
	if claims.Role != models.EnumRoleAdmin {
		c.Data["error"] = models.NewUnauthorizedError()
		c.Abort("JSONError")
	}
}

// FilterMeOrAdmin checks if the `User` is the same as the given `uid` or and admin.
// If not, return an Unauthorized `JSONError`.
func FilterMeOrAdmin(c *beego.Controller, uid int64) {
	claims := c.Ctx.Request.Context().Value("claims").(*models.Claims)
	if claims.UserID != uid && claims.Role != models.EnumRoleAdmin {
		c.Data["error"] = models.NewUnauthorizedError()
		c.Abort("JSONError")
	}
}

// FilterOwnerOrAdmin checks id the `User` is the same as the given `uid`
func FilterOwnerOrAdmin(c *beego.Controller, owner int64) {
	claims := c.Ctx.Request.Context().Value("claims").(*models.Claims)
	if owner != claims.UserID && claims.Role != models.EnumRoleAdmin {
		c.Data["error"] = models.NewUnauthorizedError()
		c.Abort("JSONError")
	}
}

// FilterUser verifies if the `User` is authenticated
var FilterUser = func(ctx *beegoCtx.Context) {
	unprotectedRoutes := []map[string]string{
		{"method": http.MethodGet, "url": "/v1/user/login"},
		// {"method": http.MethodPost, "url": "/v1/user/"},
	}

	for _, route := range unprotectedRoutes {
		if strings.HasPrefix(ctx.Input.URL(), route["url"]) && route["method"] == ctx.Request.Method {
			return
		}
	}

	// check if the authorization header is under the form: "Bearer token"
	authorization := ctx.Input.Header("Authorization")
	re := regexp.MustCompilePOSIX(`^Bearer (.+)$`)
	if !re.MatchString(authorization) {
		ctx.Output.Status = http.StatusUnauthorized
		ctx.Output.JSON(map[string]string{"message": "Unauthorized"}, false, false)
		return
	}

	// extract the token string
	tokenString := strings.Split(authorization, " ")[1]

	// parse and validate
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(beego.AppConfig.String("jwtsecret")), nil
	})

	if err != nil {
		ctx.Output.Status = http.StatusUnauthorized
		ctx.Output.JSON(map[string]string{"message": "Unauthorized"}, false, false)
		return
	}

	if !token.Valid {
		ctx.Output.Status = http.StatusUnauthorized
		ctx.Output.JSON(map[string]string{"message": "Unauthorized"}, false, false)
		return
	}

	// add the claims into the context
	ctxWithUser := context.WithValue(ctx.Request.Context(), interface{}("claims"), claims)
	ctx.Request = ctx.Request.WithContext(ctxWithUser)
}
