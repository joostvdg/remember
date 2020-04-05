package oauth

import (
	gocontext "context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/joostvdg/remember/pkg/context"
	"github.com/joostvdg/remember/pkg/remember"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//import "os"

//ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
//ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "https://remember.fly.dev/auth/google/callback",
	ClientID:     "883463116122-bkmkui3jg5bn3u65q2pkbmr4jgq0scng.apps.googleusercontent.com",
	ClientSecret: "G61xsYSiWqQ7bLFB7Ah0r3ZB",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

type GoogleOauthResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func OauthGoogleLogin(c echo.Context) error { //w http.ResponseWriter, r *http.Request) {

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(c.Response())

	/*
	   AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
	   validate that it matches the the state query parameter on your redirect callback.
	*/
	u := googleOauthConfig.AuthCodeURL(oauthState)
	//http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	return c.Redirect(307, u)
}

func OauthGoogleCallback(c echo.Context) error { //w http.ResponseWriter, r *http.Request) {
	cc := c.(*context.CustomContext)

	oauthState, _ := c.Cookie("oauthstate")
	if oauthState != nil && c.FormValue("state") != oauthState.Value {
		cc.Log.Info("invalid oauth google state")
		c.Redirect(307, "/")
		return nil
	}

	data, err := getUserDataFromGoogle(c.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		c.Redirect(307, "/")
		return err
	}

	var response GoogleOauthResponse
	// unmarshall it
	err1 := json.Unmarshal([]byte(data), &response)
	if err1 != nil {
		cc.Log.Warnf("error:", err1)
	}

	userIsFound := false
	var foundUser remember.User
	for _, user := range cc.MemoryStore.Users {
		if response.ID == user.Id {
			userIsFound = true
			foundUser = *user
		}
	}
	if !userIsFound {
		cc.Log.Warn("No user found")
	} else {
		cc.Log.Infof("Found user with %v lists", len(foundUser.Lists))
	}

	cc.Log.Infof("Authenticated user with ID: '%v'", response.ID)
	responseMessage := fmt.Sprintf("UserInfo: %s\n", data)
	return c.String(http.StatusAccepted, responseMessage)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.

	token, err := googleOauthConfig.Exchange(gocontext.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
