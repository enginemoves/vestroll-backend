package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuth struct {
	Config *oauth2.Config
	State  string
}

// Creates a OAuth instance
func NewGoogleOAuth(clientID, clientSecret, redirectURL string) *GoogleOAuth {
	return &GoogleOAuth{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
		State: generateState(),
	}
}

// Generates a random state string for OAuth
func generateState() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate random state")
	}
	return base64.URLEncoding.EncodeToString(b)
}


func (g *GoogleOAuth) LoginURL(c *gin.Context) {
	// Optionally accept a custom state from POST body
	state := c.PostForm("state")
	if state == "" {
		state = g.State
	}
	url := g.Config.AuthCodeURL(state)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// HandleCallbackGin handles the OAuth2 callback and returns JSON via Gin
func (g *GoogleOAuth) HandleCallbackGin(c *gin.Context) {
	state := c.PostForm("state")
	if state == "" {
		state = c.Query("state")
	}
	if state != g.State {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid OAuth state"})
		return
	}

	code := c.PostForm("code")
	if code == "" {
		code = c.Query("code")
	}
	token, err := g.Config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code exchange failed", "details": err.Error()})
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed getting user info", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	c.DataFromReader(http.StatusOK, resp.ContentLength, "application/json", resp.Body, nil)
}
