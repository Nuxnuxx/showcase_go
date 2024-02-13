package services

import (
	"time"
	"github.com/Nuxnuxx/showcase_go/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func NewAuthServices(u User, uStore database.Store, secretKey string) *AuthService {
	return &AuthService{
		User:      u,
		UserStore: uStore,
		SecretKey: []byte(secretKey),
	}
}

type AuthService struct {
	User      User
	UserStore database.Store
	SecretKey []byte 
}

type User struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

func (as *AuthService) GetSecretKey() []byte {
	return as.SecretKey
}

func (as *AuthService) CreateUser(u User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users(email, password, username) VALUES($1, $2, $3)`

	_, err = as.UserStore.Db.Exec(
		stmt,
		u.Email,
		string(hashedPassword),
		u.Username,
	)

	return err
}

type JwtCustomClaims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (as *AuthService) GenerateToken(user User) (string, error) {
	claims := &JwtCustomClaims{
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(as.SecretKey)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (as *AuthService) CheckEmail(email string) (User, error) {

	query := `SELECT email, password, username FROM users
		WHERE email = ?`

	stmt, err := as.UserStore.Db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	as.User.Email = email
	err = stmt.QueryRow(
		as.User.Email,
	).Scan(
		&as.User.Email,
		&as.User.Password,
		&as.User.Username,
	)

	if err != nil {
		return User{}, err
	}

	return as.User, nil
}
