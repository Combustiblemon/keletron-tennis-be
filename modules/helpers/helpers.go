package helpers

import (
	"combustiblemon/keletron-tennis-be/database/models/UserModel"
	"combustiblemon/keletron-tennis-be/modules/logger"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func CreateToken(user UserModel.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email":   user.Email,
			"name":    user.Name,
			"role":    user.Role,
			"_id":     user.ID,
			"session": user.Session,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type UserToken struct {
	jwt.MapClaims
	ID      string `json:"_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Role    string `json:"role"`
	Session string `json:"session"`
	Expire  int    `json:"exp"`
}

func ParseToken(tokenString string) (*UserToken, error) {
	claims := UserToken{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(_token *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &claims, nil
}

type URL struct {
	Full   string
	Host   string
	Scheme string
	URI    string
}

func (url *URL) String() string {
	return url.Full
}

func GetURL(ctx *gin.Context) URL {
	scheme := "https"

	if ctx.Request.TLS == nil {
		scheme = "http"
	}

	return URL{
		scheme + "://" + ctx.Request.Host + ctx.Request.RequestURI,
		ctx.Request.Host,
		scheme,
		ctx.Request.RequestURI,
	}
}

const (
	COOKIE_MAX_AGE int    = 3000
	HOME_PATH      string = "/"
)

func SetAuthCookie(ctx *gin.Context, value string) {
	host := GetURL(ctx).Host

	if strings.Contains(host, "localhost") {
		host = ""
	}

	ctx.SetCookie("session", value, COOKIE_MAX_AGE, HOME_PATH, host, true, true)
}

func ClearAuthCookie(ctx *gin.Context) {
	host := GetURL(ctx).Host

	if strings.Contains(host, "localhost") {
		host = ""
	}

	ctx.SetCookie("session", "", COOKIE_MAX_AGE, HOME_PATH, host, true, true)
}

func SendError(ctx *gin.Context, status int, err error) {
	logger.Error(ctx, err.Error())

	ctx.JSON(status, map[string]string{
		"error": err.Error(),
	})
}

func GetUser(ctx *gin.Context) *UserModel.User {
	userData, exists := ctx.Get("user")

	if exists {
		user, ok := userData.(UserModel.User)

		if ok {
			return &user
		}
	}

	return nil
}
