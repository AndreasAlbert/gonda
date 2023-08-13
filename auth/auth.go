/*
Package auth implements authentication handling
for the gonda web server.
*/
package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// An OAuthHandler provides all functionality required
// to authorize users against an OAuth provider.
type OAuthHandler struct {
	Name    string
	Config  oauth2.Config
	UserURL string
}

// randToken is a helper function that generates a random state token
// as a secret for the OAuth exchange
func randToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Fatal().Msg(fmt.Sprintf("Failed to read rand: %v", err))
	}
	return base64.StdEncoding.EncodeToString(b)
}

// HandleLogin handles the request that initiates the OAuth exchange.
// It redirects the user to the OAuth provider
func (h OAuthHandler) HandleLogin(ctx *gin.Context) {
	// Generate random state
	state := randToken()
	session := sessions.Default(ctx)
	session.Set("state", state)
	session.Save()

	// Redirect to identity provier
	url := h.Config.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleCallback handles the callback from the OAuth provider
// It extracts the short-lived authentication token,
// uses it to determine information about the user
// and finally attaches the user information to the session.
func (h OAuthHandler) HandleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	user_data, err := h.getUserData(code)
	if err != nil {
		ctx.JSON(401, "Failed to get user data")
		return
	}
	session := sessions.Default(ctx)
	session.Set("user_name", user_data["name"])
	session.Set("user_provider", user_data["provider"])
	session.Save()

	ctx.JSON(200, gin.H{"data": user_data})
}

// getUserData retrieves information about the user from the oauth provider
// code is the short-lived authentication token returned from the callback request
func (h OAuthHandler) getUserData(code string) (map[string]string, error) {

	// Exchange the short-term token for a longer-term token
	token, err := h.Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	// Ask the OAuth provider's API for information about the user
	r, err := http.NewRequest("GET", h.UserURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	r.Header.Set("Accept", "application/vnd.github+json")
	response, err := (&http.Client{}).Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	// Translate the reply into the format we need
	var response_data map[string]interface{}
	err = json.Unmarshal(content, &response_data)

	user_data := map[string]string{
		"name":     fmt.Sprintf("%v", response_data["login"]),
		"provider": h.Name,
	}

	return user_data, nil
}
