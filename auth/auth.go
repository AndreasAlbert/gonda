package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"golang.org/x/oauth2"
)

type OAuthHandler struct {
	Name   string
	Config oauth2.Config
}

func randToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		glog.Fatalf("[Gin-OAuth] Failed to read rand: %v\n", err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
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
func (h OAuthHandler) HandleCallback(ctx *gin.Context) {

	code := ctx.Query("code")
	data, _ := h.getUserData(code)
	session := sessions.Default(ctx)
	session.Set("user", data)
	session.Save()
	ctx.String(200, string(data))
}

func (h OAuthHandler) getUserData(code string) ([]byte, error) {
	token, err := h.Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	r, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		panic("")
	}
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	r.Header.Set("Accept", "application/vnd.github+json")
	glog.Info(token.AccessToken)
	response, err := (&http.Client{}).Do(r)

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
