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
	"golang.org/x/crypto/bcrypt"
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

func (ah *AuthHandler) Logout(c echo.Context) error {
	cookie := http.Cookie{
		Name:     "user",
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	}

	c.SetCookie(&cookie)

	return c.Redirect(http.StatusSeeOther, "/")
}

func (ah *AuthHandler) Login(c echo.Context) error {
	if c.Request().Method == "POST" {
		user, err := ah.AuthServices.CheckEmail(c.FormValue("email"))

		// If the user is not found
		if err != nil {
			error := map[string]services.HumanErrors{
				"email": {
					Error: "Email not found",
					Value: c.FormValue("email"),
				},
			}

			return renderView(c, authviews.Login(error))
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(c.FormValue("password")),
		)

		// If the password is not correct
		if err != nil {
			error := map[string]services.HumanErrors{
				"internal": {
					Error: "Invalid Credentials",
					Value: "",
				},
			}

			return renderView(c, authviews.Login(error))
		}

		token, err := ah.AuthServices.GenerateToken(user)

		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return renderView(c, errors_pages.Error500Index())
		}

		cookie := http.Cookie{
			Name:     "user",
			Value:    token,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(24 * time.Hour),
		}

		c.SetCookie(&cookie)

		c.Response().Header().Set("HX-Redirect", "/")
		c.Response().WriteHeader(http.StatusOK)
		return nil
	}

	return renderView(c, authviews.LoginIndex())
}

func (ah *AuthHandler) Register(c echo.Context) error {
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

		userInDatabase, err := ah.AuthServices.CheckEmail(user.Email)

		if userInDatabase != (services.User{}) {
			humanErrors := map[string]services.HumanErrors{
				"email": {
					Error: "Email already exists",
					Value: user.Email,
				},
			}

			return renderView(c, authviews.Register(humanErrors))
		}

		err = ah.AuthServices.CreateUser(user)

		if err != nil {
			log.Errorf("Error Creating User: %v", err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return renderView(c, errors_pages.Error500Index())
		}

		token, err := ah.AuthServices.GenerateToken(user)

		if err != nil {
			log.Errorf("Error generating token: %v", err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return renderView(c, errors_pages.Error500Index())
		}

		cookie := http.Cookie{
			Name:     "user",
			Value:    token,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(24 * time.Hour),
		}

		c.SetCookie(&cookie)

		// INFO: To redirect when using HTMX you need to set the HX-Redirect header
		c.Response().Header().Set("HX-Redirect", "/")
		c.Response().WriteHeader(http.StatusOK)
		return nil
	}

	return renderView(c, authviews.RegisterIndex())
}

func (ah *AuthHandler) Profil(c echo.Context) error {
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

func (ah *AuthHandler) CheckLogged(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := c.Cookie("user")

		// If there is a error means he is not connected
		if err != nil {
			return next(c)
		}

		// The token is here means he is connected
		if token.Value != "" {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		return next(c)
	}
}

func (ah *AuthHandler) CheckNotLogged(c echo.Context, err error) error {
	token, ok := c.Get("user").(string)

	// Means the user is not connected
	if !ok {
		return c.Redirect(http.StatusSeeOther, "/auth/register")
	}

	// Means the user is not connected
	if token == "" {
		return c.Redirect(http.StatusSeeOther, "/auth/register")
	}

	return nil
}
