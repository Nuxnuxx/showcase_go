package services

import (
	"time"

	"github.com/Nuxnuxx/showcase_go/internal/database"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func NewAuthServices(u User, uStore database.Store, secretKey string) *AuthService {

	return &AuthService{
		User:      u,
		UserStore: uStore,
		SecretKey: secretKey,
	}
}

type AuthService struct {
	User      User
	UserStore database.Store
	SecretKey string
}

type User struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=8,max=20"`
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

// func GenerateToken(user User) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"email":    user.Email,
// 		"username": user.Username,
// 		"password": user.Password,
// 		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Expiration time in seconds
// 		"iat":      time.Now().Unix(),                     // Issued at this time
// 	})
//
// 	signedToken, err := token.SignedString([]byte("secret"))
//
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return signedToken, nil
// }

// func VerifyToken(tokenString string) (jwt.MapClaims, error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return []byte("secret"), nil
// 	})
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if !token.Valid {
// 		return nil, err
// 	}
//
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		return claims, nil
// 	}
//
// 	return nil, err
// }

// func (as *AuthService) CheckEmail(email string) (User, error) {
//
// 	query := `SELECT id, email, password, username FROM users
// 		WHERE email = ?`
//
// 	stmt, err := as.UserStore.Db.Prepare(query)
// 	if err != nil {
// 		return User{}, err
// 	}
//
// 	defer stmt.Close()
//
// 	as.User.Email = email
// 	err = stmt.QueryRow(
// 		as.User.Email,
// 	).Scan(
// 		&as.User.ID,
// 		&as.User.Email,
// 		&as.User.Password,
// 		&as.User.Username,
// 	)
// 	if err != nil {
// 		return User{}, err
// 	}
//
// 	return as.User, nil
// }
