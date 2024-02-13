package handlers

import (
	"net/http"

	"github.com/Nuxnuxx/showcase_go/internal/services"
	authviews "github.com/Nuxnuxx/showcase_go/internal/views/auth_views"
	"github.com/Nuxnuxx/showcase_go/internal/views/errors_pages"
	"github.com/labstack/echo/v4"
)

type AuthServices interface {
	CreateUser(user services.User) error
}

func NewAuthHandler(as AuthServices) *AuthHandler {

	return &AuthHandler{
		AuthServices: as,
	}
}

type AuthHandler struct {
	AuthServices AuthServices
}

func (au *AuthHandler) Register(c echo.Context) error {
	if c.Request().Method == "POST" {
		user := services.User{
			Email:    c.FormValue("email"),
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}

		if err := c.Validate(user); err != nil {
			humanErrors := services.CreateHumanErrors(err)

			return renderView(c, authviews.Register(humanErrors))
		}

		err := au.AuthServices.CreateUser(user)

		if err != nil {
			return renderView(c, errors_pages.Error500Index())
		}

		return c.Redirect(http.StatusSeeOther, "/")
	}

	return renderView(c, authviews.RegisterIndex())
}

// func (au *AuthHandler) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		if err := next(c); err != nil {
// 			c.Error(err)
// 		}
//
// 		token, ok := c.Get("user").(*jwt.Token)
//
// 		if !ok {
// 			return c.Redirect(http.StatusNetworkAuthenticationRequired, "/login")
// 		}
//
// 		claims, err := au.AuthServices.VerifyToken(token.Raw)
//
// 		if err != nil {
// 			return c.Redirect(http.StatusNetworkAuthenticationRequired, "/login")
// 		}
//
// 		c.Set("claims", claims)
//
// 		return next(c)
// 	}
// }
