package providersGoogle

import (
	"combustiblemon/keletron-tennis-be/database/models/UserModel"
	helpers "combustiblemon/keletron-tennis-be/modules"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const COOKIE_MAX_AGE int = 3000

// Your credentials should be obtained from the Google
// Developer Console (https://console.developers.google.com).
var conf = &oauth2.Config{
	ClientID:     "",
	ClientSecret: "",
	RedirectURL:  "http://localhost:2000/auth/providers/google/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

func Init() {
	conf.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	conf.ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
}

func HandleLogin(w http.ResponseWriter, r *http.Request, oauthConf *oauth2.Config, oauthStateString string) error {
	loginURL, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		return err
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	loginURL.RawQuery = parameters.Encode()

	http.Redirect(w, r, loginURL.String(), http.StatusTemporaryRedirect)
	return nil
}

const oauthStateStringGl = "google_login_state"

func Start() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := HandleLogin(ctx.Writer, ctx.Request, conf, oauthStateStringGl)

		if err != nil {
			http.Redirect(ctx.Writer, ctx.Request, "/", http.StatusTemporaryRedirect)
		}
	}
}

type GoogleUserData struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Picture       string
}

/*
CallBackFromGoogle Function
*/
func CallBackFromGoogle(w http.ResponseWriter, r *http.Request) (*GoogleUserData, error) {
	state := r.FormValue("state")
	fmt.Println(state)
	if state != oauthStateStringGl {
		return nil, fmt.Errorf("invalid oauth state, expected " + oauthStateStringGl + ", got " + state + "\n")
	}

	code := r.FormValue("code")

	if code == "" {
		reason := r.FormValue("error_reason")
		return nil, fmt.Errorf("Code Not Found to provide AccessToken: %v", reason)
	}

	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("oauthConfGl.Exchange() failed with " + err.Error() + "\n")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil, fmt.Errorf("Get: " + err.Error() + "\n")
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil, fmt.Errorf("ReadAll: " + err.Error() + "\n")
	}

	// fmt.Println("parseResponseBody: " + string(response) + "\n")

	var d GoogleUserData
	err = json.Unmarshal(response, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func Callback() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := CallBackFromGoogle(ctx.Writer, ctx.Request)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		session, err := uuid.NewV7()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		token, err := helpers.CreateToken(UserModel.User{
			Name:      data.Name,
			Role:      "USER",
			FCMTokens: []string{},
			Session:   session.String(),
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		ctx.SetCookie("auth", token, COOKIE_MAX_AGE, "/", ctx.Request.URL.Host, true, true)
		ctx.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	}
}
