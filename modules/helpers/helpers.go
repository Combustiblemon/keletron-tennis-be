package helpers

import (
	"combustiblemon/keletron-tennis-be/database/models/UserModel"
	"combustiblemon/keletron-tennis-be/modules/logger"
	"fmt"
	"os"
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

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(_token *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
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
	ctx.SetCookie("auth", value, COOKIE_MAX_AGE, HOME_PATH, GetURL(ctx).Host, true, true)
}

func ClearAuthCookie(ctx *gin.Context) {
	ctx.SetCookie("auth", "", COOKIE_MAX_AGE, HOME_PATH, GetURL(ctx).Host, true, true)
}

func SendError(ctx *gin.Context, status int, err error) {
	logger.Error(ctx, err.Error())

	ctx.JSON(status, map[string]string{
		"error": err.Error(),
	})
}
