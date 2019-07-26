package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"kigo/model"
	"kigo/session"
)

// Scopes: OAuth 2.0 scopes provide a way to limit the amount of access that is granted to an access token.
var googleConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
	ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

const googleApiUrl = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

// Asks for google permissions.
func googleLogin(w http.ResponseWriter, r *http.Request) {
	// Starts session
	sess := session.MegaManager.SessionStart(w, r)
	// Retrieves session id
	state := sess.Get("sid")
	// Pass in state for later validation and redirect to url
	url := googleConfig.AuthCodeURL(state.(string))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Cleans up after google response.
func googleCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")
	// Handles invalid state
	if r.FormValue("state") != oauthState.Value {
		// Deletes session
		session.MegaManager.SessionDestroy(w, r)
		// Logs error
		log.Println("invalid oauth google state")
		// Redirects to homepage
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// Retrieves user data (name,email,etc)
	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// Creates user in database
	initOrCreateUser(data)
	// Redirects to protected endpoint
	http.Redirect(w, r, "/haiku", http.StatusTemporaryRedirect)
}

// Initializes or creates a user.
func initOrCreateUser(data []byte) {
	var user model.User
	// Stores req data in user
	json.Unmarshal(data, &user)
	// Creates new user or returns existing user (contingent on whether id is in DB)
	DB.FirstOrInit(&user, map[string]interface{}{"ID": user.ID})
}

// Exchanges code for google user's info.
func getUserDataFromGoogle(code string) ([]byte, error) {
	// Exchanges code for token
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	// Uses token to request user info
	res, err := http.Get(googleApiUrl + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer res.Body.Close()
	// Reads the response
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
