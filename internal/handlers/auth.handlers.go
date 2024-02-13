package handlers

import (
	"net/http"
	"time"

	"github.com/Nuxnuxx/showcase_go/internal/services"
	authviews "github.com/Nuxnuxx/showcase_go/internal/views/auth_views"
	"github.com/Nuxnuxx/showcase_go/internal/views/errors_pages"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type AuthServices interface {
	GetSecretKey() []byte
	CheckEmail(email string) (services.User, error)
	CreateUser(user services.User) error
	GenerateToken(user services.User) (string, error)
}

func NewAuthHandler(as AuthServices) *AuthHandler {

	return &AuthHandler{
		AuthServices: as,
	}
}

type AuthHandler struct {
	AuthServices AuthServices
}

func (au *AuthHandler) Login(c echo.Context) error {
	return renderView(c, authviews.LoginIndex())
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

		userInDatabase, err := au.AuthServices.CheckEmail(user.Email)

		if userInDatabase != (services.User{}) {
			humanErrors := map[string]services.HumanErrors{
				"email": {
					Error: "Email already exists",
					Value: user.Email,
				},
			}

			return renderView(c, authviews.Register(humanErrors))
		}

		err = au.AuthServices.CreateUser(user)

		if err != nil {
			log.Errorf("Error Creating User: %v", err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return renderView(c, errors_pages.Error500Index())
		}

		token, err := au.AuthServices.GenerateToken(user)

		if err != nil {
			log.Errorf("Error generating token: %v", err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return renderView(c, errors_pages.Error500Index())
		}

		cookie := http.Cookie{
			Name:    "user",
			Value:   token,
			Path:    "/",
			Secure:  true,
			Expires: time.Now().Add(24 * time.Hour),
		}

		c.SetCookie(&cookie)

		// INFO: To redirect when using HTMX you need to set the HX-Redirect header
		c.Response().Header().Set("HX-Redirect", "/")
		c.Response().WriteHeader(http.StatusOK)
		return nil
	}

	return renderView(c, authviews.RegisterIndex())
}

func (au *AuthHandler) Profil(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)

	if !ok {
		log.Errorf("Error getting claims from token: %v", token)
		return renderView(c, errors_pages.Error401Index())
	}

	claims, ok := token.Claims.(*services.JwtCustomClaims)

	if !ok {
		log.Errorf("Error getting claims from token: %v", token)
		return renderView(c, errors_pages.Error401Index())
	}

	user := services.User{
		Email:    claims.Email,
		Username: claims.Username,
	}

	return renderView(c, authviews.ProfilIndex(user))
}
