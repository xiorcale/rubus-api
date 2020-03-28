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

// FilterUser verifies if the `User` is authenticated
var FilterUser = func(ctx *beegoCtx.Context) {
	unprotectedRoutes := []map[string]string{
		{"method": "GET", "url": "/v1/user/login"},
		{"method": "POST", "url": "/v1/user/"},
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
	tk := &models.Token{}
	token, err := jwt.ParseWithClaims(tokenString, tk, func(token *jwt.Token) (interface{}, error) {
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

	// add the user id into the context
	ctxWithUser := context.WithValue(ctx.Request.Context(), interface{}("user"), tk.UserID)
	ctx.Request = ctx.Request.WithContext(ctxWithUser)
}
