package helpers

import (
	"combustiblemon/keletron-tennis-be/database/models/UserModel"
	"combustiblemon/keletron-tennis-be/modules/logger"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	if strings.Contains(err.Error(), "Error:Field validation") {
		ctx.JSON(status, map[string]any{
			"errors": GenerateValidationError(err),
		})

		return
	}

	ctx.JSON(status, map[string]any{
		"errors": []string{err.Error()},
	})
}

func GetUser(ctx *gin.Context) (user *UserModel.User, exists bool) {
	userData, exists := ctx.Get("user")

	if exists {
		user, ok := userData.(*UserModel.User)

		if ok {
			return user, exists
		}
	}

	return nil, exists
}

func FormatDate(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute())
}

func ParseDate(date string) time.Time {
	loc, _ := time.LoadLocation("Europe/Athens")

	year, _ := strconv.Atoi(date[0:4])
	month, _ := strconv.Atoi(date[5:7])
	day, _ := strconv.Atoi(date[8:10])
	hour, _ := strconv.Atoi(date[11:13])
	minute, _ := strconv.Atoi(date[14:16])

	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, loc)
}

func firstToLower(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
		return s
	}
	lc := unicode.ToLower(r)
	if r == lc {
		return s
	}
	return string(lc) + s[size:]
}

type ValidationErrorInfo struct {
	Key   string
	Error string
	Info  string
}

func GenerateValidationError(err error) []ValidationErrorInfo {
	errParts := strings.Split(err.Error(), "\n")
	regKey, _ := regexp2.Compile("(?<=Key: ')\\b\\w+\\.\\w+\\b", 0)
	regTag, _ := regexp2.Compile("(?<=failed on the ')\\b\\w+\\b(?=' tag)", 0)

	errs := []ValidationErrorInfo{}

	for _, v := range errParts {
		if strings.Contains(v, "Error:Field validation") {
			m1, _ := regKey.FindStringMatch(v)
			m2, _ := regTag.FindStringMatch(v)

			key := firstToLower(strings.Replace(m1.String(), "ID", "_id", 1))

			tag := strings.ToLower(m2.String())

			errs = append(errs, ValidationErrorInfo{
				Key:   key,
				Error: tag,
				Info:  "",
			})
		}
	}

	return errs
}

func ObjectIDFromHex(hex string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(hex)

}
